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
