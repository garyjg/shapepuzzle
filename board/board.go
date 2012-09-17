// -*- tab-width: 4; -*-

package board

// A Board is a number of rows and columns, a set of shape placements, and
// a current mask which provides a fast way to check if a new placement
// fits.


import (
	"shapepuzzle/shape"
	"shapepuzzle/mask"
	"fmt"
	"log"
)


type Board struct {
	nrows int
	ncols int
	mask mask.MaskBits
	placements []shape.Shape
}

func NewBoard(nrows int, ncols int) Board {
	nb := Board{nrows:nrows, ncols:ncols, mask:0}
	nb.placements = make([]shape.Shape, 0)
	return nb
}

func (b Board) NumShapes() int {
	return len(b.placements)
}

func (b Board) NumRows() int {
	return b.nrows
}

func (b Board) NumCols() int {
	return b.ncols
}

func (b Board) Mask() mask.MaskBits {
	return b.mask
}


func (b Board) RegionMask() mask.MaskBits {
	cbits := mask.MaskBits(0)
	for c := 0; c < b.NumCols(); c += 1 {
	    cbits = mask.FirstMaskBit() | (cbits >> 1)
	}
	mbits := mask.MaskBits(0)
	for r := 0; r < b.NumRows(); r += 1 {
		mbits = (mbits >> 8) | cbits
	}
	return mbits
}


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
	for r := 0; r < nrow; r += 1 {
		grid[r] = make([]int, ncol)
		buf += "["
		for c := 0; c < ncol; c += 1 {
			var mbits mask.MaskBits
			mbits = mask.FirstMaskBit().Translate(r, c)
			for _, p := range b.placements {
				if p.Mask() & mbits != 0 {
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


// Given a set of shapes, a board 8x8, and their permutations and
// translations on the board, find the combination which covers every spot
// on the board.  The board state is represented as a 64-bit mask, and
// shape placements can be converted to a mask on that board to quickly
// determine if that placement is valid, ie, it does not collied any other
// placements.
// 
// As soon as a board state is found on which a piece cannot be placed,
// that board state is discarded and not propagated.  Therefore not all
// pieces will need to be attempted before a placement sequence is
// recognized as a dead end.
// 
// Because the first placement is symmetrical across the four quadrants of
// the board, it only needs to generate placements in one quadrant.  That
// at least reduces the search space by close to a fourth.

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

// Also, the masks for a shape placement do not change depending upon the
// incoming board, so really the placements can all be computed ahead of
// time and then the masks scanned for one that fits on the incoming board.


type BoardChannel chan Board

func FirstPlacements(s shape.Shape, b Board, bc BoardChannel) {

	// Just generate every permutation of the shape at every position in the
	// quadrant and push it to the channel, unless it matches one of the reject
	// patterns.
	rejects := GapShapes(b)
	perms := s.Permutations()
	ngen, nrej := 0, 0
	for _, p := range perms {
		height := p.NumRows()
		width := p.NumCols()
		for r := 0; r <= b.NumRows()/2 && r <= b.NumRows() - height; r += 1 {
			for c := 0; c <= b.NumCols()/2 && c <= b.NumCols() - width; c += 1 {
				place := p.Translate(r, c)
				nb := b.Place(place)
				if !RejectBoard(nb, rejects) {
					log.Printf("Generating first placement (S#%d):\n%v", 
							   place.ID(), nb)
					bc <- nb
					ngen += 1
				} else {
				    log.Printf("Rejected first placement (S#%d):\n%v", 
							   place.ID(), nb)
					nrej += 1
				}
			}
		}
	}
	log.Printf("Total first placements (S#%d): %d generated, %d rejected.", 
			   s.ID(), ngen, nrej)
	close(bc)
}


func NextPlacements(s shape.Shape, base Board, boards BoardChannel, 
					moves BoardChannel) {

	// Generate all possible board masks for placing this shape.
	placements := make([]shape.Shape, 100)
	rejects := GapShapes(base)
	placements = placements[0:0]
	perms := s.Permutations()
	for i := 0; i < len(perms); i += 1 {
		s := &(perms[i])
		width := s.NumCols()
		height := s.NumRows()
		for r := 0; r <= base.NumRows() - height; r += 1 {
			for c := 0; c <= base.NumCols() - width; c += 1 {
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
			if b.Mask() & place.Mask() == 0 {
			    nb := b.Place(place)
				if ! RejectBoard(nb, rejects) {
					log.Printf("Generating placement (S#%d):\n%v", place.ID(), nb)
					moves <- nb
				}
			}
		}
	}
	close(moves)
}


// To search the solution space, each piece gets its own goroutine in which
// it generates all the available moves on the current board state.  It
// puts each valid move and the resulting board state onto its output
// channel.  The goroutines for each piece can be chained together so that
// if the final piece in the chain puts a new board state on its output
// channel, then that must be a solution to the puzzle.
// 


func (b Board) Solve(shapes []shape.Shape) (BoardChannel) {

	nshapes := len(shapes)

	// Set up a channel for each shape to be placed.
	channels := make ([]BoardChannel, nshapes)
	for i := 0; i < nshapes; i += 1 {
		channels[i] = make(BoardChannel, 10000)
	}
	
	// Chain the channels.  Generate first placements for the first shape,
	// and tell it to put those new boards on its channel.
	go FirstPlacements(shapes[0], b, channels[0])

	for i := 1; i < nshapes; i += 1 {
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
	if s.OutlineMask() & b.Mask() != s.OutlineMask() {
	    return false
	}
	if s.GapMask() & b.Mask() == s.GapMask() {
	    return false
	}
	return true
}


func searchGap(b Board, patterns []shape.Shape) (int) {
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


// Given a board with a particular size, generate all the masks which if
// they match a board should cause the board to be rejected as a potential
// solution.

func GapShapes(b Board) []shape.Shape {

	grids := [][][]int { { 
		{0,1,0}, {1,2,1}, {0,1,0} }, {
		{0,1,1,0}, {1,2,2,1}, {0,1,1,0} }, {
		{0,1,1,0}, {1,2,2,1}, {1,2,2,1}, {0,1,1,0} }, {
		{0,1,1,1,0}, {1,2,2,2,1}, {0,1,1,1,0} }, {
		{0,1,1,1,1,0}, {1,2,2,2,2,1}, {0,1,1,1,1,0} }, {
		{0,1,1,0}, {1,2,2,1}, {1,2,2,1}, {0,1,1,0} }, {
		{0,1,1,0},
		{1,2,2,1},
		{0,1,2,1},
		{0,0,1,0} } }

	region := b.RegionMask()
	shapes := []shape.Shape{}
	for id, g := range grids {
		s := shape.NewShape(id+100, g)
		perms := s.Permutations()
		for i := 0; i < len(perms); i += 1 {
			s := &(perms[i])
			width := s.NumCols()
			height := s.NumRows()
			for r := -1; r <= b.NumRows() - height + 1; r += 1 {
				for c := -1; c <= b.NumCols() - width + 1; c += 1 {
					shapes = append(shapes, s.Translate(r, c).Clip(region))
				}
			}
		}
	}
	return shapes
}

