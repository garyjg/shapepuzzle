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


type Placement struct {
	shape shape.Shape
	row int
	col int
	mask mask.MaskBits
}

func (p *Placement) ComputeMask() {
	 // Get the shape's mask, shift it to this row and col, and or it.
	 p.mask = p.shape.Mask().Translate(p.row, p.col)
	 log.Printf("Placement.ComputeMask(%v, %v, %v) ==> %v", 
	 	p.shape.Mask(), p.row, p.col, p.mask)
}


func (p Placement) Mask() mask.MaskBits {
	 return p.mask
}

func NewPlacement(s shape.Shape, row int, col int) Placement {
	 place := Placement{s, row, col, 0}
	 (&place).ComputeMask()
	 return place
}


type Board struct {
	nrows int
	ncols int
	mask mask.MaskBits
	placements []Placement
}

func NewBoard(nrows int, ncols int) Board {
	nb := Board{nrows:nrows, ncols:ncols, mask:0}
	nb.placements = make([]Placement, 0)
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


func (b Board) Place(p Placement) Board {

	log.Printf("Place(%v)\n", p)
	nb := b
	nb.placements = make([]Placement, len(b.placements), len(b.placements)+1)
	copy(nb.placements, b.placements)
	nb.placements = append(nb.placements, p)
	log.Printf("adding placement %d mask %v to board mask %v\n", 
			   len(nb.placements), p.Mask(), nb.Mask())
	nb.mask = nb.mask | p.Mask()
	log.Printf("result: %v\n", nb.Mask())
	return nb
}


func (b Board) String() string {

	// Fill the spots on the board with each shape's ID.
	nrow := b.NumRows()
	ncol := b.NumCols()
	buf := fmt.Sprintf("%d shapes.\n", len(b.placements))
	grid := make([][]int, nrow)
	for r := 0; r < nrow; r += 1 {
		grid[r] = make([]int, ncol)
		for c := 0; c < ncol; c += 1 {
			var mbits mask.MaskBits
			mbits = mask.FirstMaskBit().Translate(r, c)
			for _, p := range b.placements {
				if p.Mask() & mbits != 0 {
					grid[r][c] = p.shape.ID()
					break
				}
			}
		}
		buf = buf + fmt.Sprintf("%v\n", grid[r])
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
	// quadrant and push it to the channel.
	perms := s.Permutations()
	for _, p := range perms {
		height := p.NumRows()
		width := p.NumCols()
		for r := 0; r <= b.NumRows()/2 && r <= b.NumRows() - height; r += 1 {
			for c := 0; c <= b.NumCols()/2 && c <= b.NumCols() - width; c += 1 {
				place := NewPlacement(p, r, c)
				log.Printf("Placement mask => %v\n", place.Mask())
				nb := b.Place(place)
				log.Printf("Generating first placement:\n%v", nb)
				bc <- nb
			}
		}
	}
	close(bc)
}


func NextPlacements(s shape.Shape, base Board, boards BoardChannel, 
					moves BoardChannel) {

	// Generate all possible board masks for placing this shape.
	placements := make([]Placement, 100)
	placements = placements[0:0]
	perms := s.Permutations()
	for i := 0; i < len(perms); i += 1 {
		s := &(perms[i])
		width := s.NumCols()
		height := s.NumRows()
		for r := 0; r <= base.NumRows() - height; r += 1 {
			for c := 0; c <= base.NumCols() - width; c += 1 {
				place := NewPlacement(*s, r, c)
				placements = append(placements, place)
			}
		}
	}
	// For each input board, find all the placements which fit.
	for b := range boards {
		for _, place := range placements {
			if b.Mask() & place.Mask() == 0 {
				log.Printf("Received board:\n%v", b)
			    nb := b.Place(place)
				log.Printf("Found next placement:\n%v", nb)
				moves <- nb
			}
		}
	}
	close(moves)
}



// Given a board with a particular size, generate all the masks which if
// they match a board should cause the board to be rejected as a potential
// solution.

func gapMasks(b Board) []mask.MaskBits {

	shapes := [][][]int { { 
		{0,1,0}, {1,2,1}, {0,1,0} }, {
		{0,1,1,0}, {1,2,2,1}, {0,1,1,0} }, {
		{0,1,1,1,0}, {1,2,2,2,1}, {0,1,1,1,0} }, {
		{0,1,1,1,1,0}, {1,2,2,2,2,1}, {0,1,1,1,1,0} }, {
		{0,1,1,0}, {1,2,2,1}, {1,2,2,1}, {0,1,1,0} }, {
		{0,1,1,0}, {1,2,2,1}, {0,1,2,1}, {0,1,1,1} } }



}
