// -*- tab-width: 4; -*-

// Package board contains the Board structure and methods for accumulating
// shape placements on a puzzle grid.
//
// Given a set of shapes, a board 8x8, and their permutations and
// translations on the board, find the combination which covers every spot
// on the board.  The board state is represented as a 64-bit mask, and
// shape placements can be converted to a mask on that board to quickly
// determine if that placement is valid, ie, it does not collide with any
// other placements.
//
// As soon as a board state is found on which a piece cannot be placed,
// that board state is discarded and not propagated.  Therefore not all
// pieces will need to be attempted before a placement sequence is
// recognized as a dead end.
//
// Because the first placement is symmetrical across the four quadrants of
// the board, it only needs to generate placements in one quadrant.  That
// at least reduces the search space by close to a fourth.
//
// This might be optimized by immediately pruning board states which cannot
// yield solutions because there are gaps which are too small for any shape.
// These gaps can be found by matching against a set of masks:
//
//  A square with one space inside, shifted to all positions on the board.
//  A space on the edge surrounded by filled positions.
//  A rectangle with two spaces inside.
//  A rectangle on the edge.
//  Corner masks which only fit one, two, three, or four spaces.
//
// [ [ 0, 1, 0 ], [ 1, 0, 1 ], [ 0, 1, 0 ] ]
//
// The corners do not need to be filled in, only the edges.  The point is
// that empty spaces are surrounded on all sides by either filled spaces or
// the border.
//
// Basically, given a rectangle with 1, 2, 3, or 4 empty spaces inside, AND
// that filled shape with the board, then OR the mask with the inside
// spaces removed.  If the result is not equal to the filled mask, then
// there are gaps which cannot fit a shape.
//
// Also, the masks for a shape placement do not change depending upon the
// incoming board, so really the placements can all be computed ahead of
// time and then the masks scanned for one that fits on the incoming board.
package board

import (
	"fmt"
	"log"

	"github.com/garyjg/shapepuzzle/mask"
	"github.com/garyjg/shapepuzzle/shape"
)

// Board is a number of rows and columns, a set of shape placements, and
// a current mask which provides a fast way to check if a new placement
// fits.
//
type Board struct {
	nrows      int
	ncols      int
	mask       mask.Bits
	placements []shape.Shape
}

// NewBoard initializes a new Board with nrows rows and ncols columns.
func NewBoard(nrows int, ncols int) Board {
	nb := Board{nrows: nrows, ncols: ncols, mask: 0}
	nb.placements = make([]shape.Shape, 0)
	return nb
}

// NumShapes returns the number of shapes placed on the board.
func (b Board) NumShapes() int {
	return len(b.placements)
}

// NumRows returns number of rows in the board.
func (b Board) NumRows() int {
	return b.nrows
}

// NumCols returns number of columns in the board.
func (b Board) NumCols() int {
	return b.ncols
}

// Mask returns the current mask for the board.
func (b Board) Mask() mask.Bits {
	return b.mask
}

// RegionMask creates a mask which matches all the points on the board.
// For example, the RegionMask for a 5x5 Board will have the first 5 bits
// set of the first 5 bytes, corresponding to the upper left 5x5 grid
// of the full 8x8 bit mask.
func (b Board) RegionMask() mask.Bits {
	cbits := mask.Bits(0)
	for c := 0; c < b.NumCols(); c++ {
		cbits = mask.FirstBit() | (cbits >> 1)
	}
	mbits := mask.Bits(0)
	for r := 0; r < b.NumRows(); r++ {
		mbits = (mbits >> 8) | cbits
	}
	return mbits
}

// Place adds the Shape to the Board and returns the new Board.  Note the
// original board is not changed, and there is no check that the Shape fits.
// This just adds the Shape to list of placements and updates the mask.
func (b Board) Place(p shape.Shape) Board {

	nb := b
	nb.placements = make([]shape.Shape, len(b.placements), len(b.placements)+1)
	copy(nb.placements, b.placements)
	nb.placements = append(nb.placements, p)
	nb.mask = nb.mask | p.Mask()
	return nb
}

func (b Board) String() string {

	// Fill the spots on the board with each shape's ID.
	nrow := b.NumRows()
	ncol := b.NumCols()
	buf := ""
	grid := make([][]int, nrow)
	for r := 0; r < nrow; r++ {
		grid[r] = make([]int, ncol)
		buf += "["
		for c := 0; c < ncol; c++ {
			var mbits mask.Bits
			mbits = mask.FirstBit().Translate(r, c)
			for _, p := range b.placements {
				if p.Mask()&mbits != 0 {
					grid[r][c] = p.ID()
					break
				}
			}
			buf += fmt.Sprintf(" %2d", grid[r][c])
		}
		buf += "]\n"
	}
	return buf
}

// Channel is a channel for passing Board states.
type Channel chan Board

// FirstPlacements generates every permutation of the shape at every position in
// the quadrant and push it to the channel, unless it matches one of the reject
// patterns.
func FirstPlacements(s shape.Shape, b Board, bc Channel) {

	rejects := GapShapes(b)
	perms := s.Permutations()
	ngen, nrej := 0, 0
	for _, p := range perms {
		height := p.NumRows()
		width := p.NumCols()
		for r := 0; r <= b.NumRows()/2 && r <= b.NumRows()-height; r++ {
			for c := 0; c <= b.NumCols()/2 && c <= b.NumCols()-width; c++ {
				place := p.Translate(r, c)
				nb := b.Place(place)
				if !RejectBoard(nb, rejects) {
					log.Printf("Generating first placement (S#%d):\n%v",
						place.ID(), nb)
					bc <- nb
					ngen++
				} else {
					log.Printf("Rejected first placement (S#%d):\n%v",
						place.ID(), nb)
					nrej++
				}
			}
		}
	}
	log.Printf("Total first placements (S#%d): %d generated, %d rejected.",
		s.ID(), ngen, nrej)
	close(bc)
}

// NextPlacements generates all possible board masks for placing the given shape
// on a blank Board base.  Then it tries to place each of those permutations on
// each Board on the boards channel.  Each board on which the shape can be
// placed successfully is passed to the moves Channel.
func NextPlacements(s shape.Shape, base Board, boards Channel,
	moves Channel) {

	// Generate all possible board masks for placing this shape.
	placements := make([]shape.Shape, 100)
	rejects := GapShapes(base)
	placements = placements[0:0]
	perms := s.Permutations()
	for i := 0; i < len(perms); i++ {
		s := &(perms[i])
		width := s.NumCols()
		height := s.NumRows()
		for r := 0; r <= base.NumRows()-height; r++ {
			for c := 0; c <= base.NumCols()-width; c++ {
				place := (*s).Translate(r, c)
				// If this shape at this place on a blank board would be
				// rejected, then reject it for any board.
				nb := base.Place(place)
				if !RejectBoard(nb, rejects) {
					placements = append(placements, place)
				} else {
					log.Printf("Rejected prepared placement (S#%d):\n%v",
						place.ID(), nb)
				}
			}
		}
	}
	// For each input board, find all the placements which fit, but reject the
	// ones known to not have room for future placements.
	for b := range boards {
		for _, place := range placements {
			if b.Mask()&place.Mask() == 0 {
				nb := b.Place(place)
				if !RejectBoard(nb, rejects) {
					log.Printf("Generating placement (S#%d):\n%v", place.ID(), nb)
					moves <- nb
				}
			}
		}
	}
	close(moves)
}

// Solve searches the solution space: each piece gets its own goroutine in which
// it generates all the available moves on the current board state.  It puts
// each valid move and the resulting board state onto its output channel.  The
// goroutines for each piece can be chained together so that if the final piece
// in the chain puts a new board state on its output channel, then that must be
// a solution to the puzzle.
//
func (b Board) Solve(shapes []shape.Shape) Channel {

	nshapes := len(shapes)

	// Set up a channel for each shape to be placed.
	channels := make([]Channel, nshapes)
	for i := 0; i < nshapes; i++ {
		channels[i] = make(Channel, 10000)
	}

	// Chain the channels.  Generate first placements for the first shape,
	// and tell it to put those new boards on its channel.
	go FirstPlacements(shapes[0], b, channels[0])

	for i := 1; i < nshapes; i++ {
		go NextPlacements(shapes[i], b, channels[i-1], channels[i])
	}

	// Finally listen for a solution (or not) to be pushed to the last
	// channel.

	return channels[nshapes-1]
}

// See if the gap mask defined in this shape indicates that this board
// state should be rejected as a possible solution.

func rejectGap(b Board, s shape.Shape) bool {

	// First see if the board matches the gap outline.
	if s.OutlineMask()&b.Mask() != s.OutlineMask() {
		return false
	}
	if s.GapMask()&b.Mask() == s.GapMask() {
		return false
	}
	return true
}

func searchGap(b Board, patterns []shape.Shape) int {
	for i, s := range patterns {
		if rejectGap(b, s) {
			return i
		}
	}
	return -1
}

func RejectBoard(b Board, patterns []shape.Shape) bool {
	return searchGap(b, patterns) >= 0
}

// GapShapes, given a board with a particular size, generates all the masks
// which if they match a board should cause the board to be rejected as a
// potential solution.
func GapShapes(b Board) []shape.Shape {

	grids := [][][]int{{
		{0, 1, 0}, {1, 2, 1}, {0, 1, 0}}, {
		{0, 1, 1, 0}, {1, 2, 2, 1}, {0, 1, 1, 0}}, {
		{0, 1, 1, 0}, {1, 2, 2, 1}, {1, 2, 2, 1}, {0, 1, 1, 0}}, {
		{0, 1, 1, 1, 0}, {1, 2, 2, 2, 1}, {0, 1, 1, 1, 0}}, {
		{0, 1, 1, 1, 1, 0}, {1, 2, 2, 2, 2, 1}, {0, 1, 1, 1, 1, 0}}, {
		{0, 1, 1, 0}, {1, 2, 2, 1}, {1, 2, 2, 1}, {0, 1, 1, 0}}, {
		{0, 1, 1, 0},
		{1, 2, 2, 1},
		{0, 1, 2, 1},
		{0, 0, 1, 0}}}

	region := b.RegionMask()
	shapes := []shape.Shape{}
	for id, g := range grids {
		s := shape.NewShape(id+100, g)
		perms := s.Permutations()
		for i := 0; i < len(perms); i++ {
			s := &(perms[i])
			width := s.NumCols()
			height := s.NumRows()
			for r := -1; r <= b.NumRows()-height+1; r++ {
				for c := -1; c <= b.NumCols()-width+1; c++ {
					shapes = append(shapes, s.Translate(r, c).Clip(region))
				}
			}
		}
	}
	return shapes
}
