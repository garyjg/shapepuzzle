// -*- tab-width: 4; -*-

package shape

import (
	"fmt"

	"github.com/garyjg/shapepuzzle/mask"
)

// Shape contains an id, a mask for its shape in the uppermost leftmost corner
// of a grid, and a row
type Shape struct {
	id    int
	shape [][]int
	mask  mask.Bits
	gaps  mask.Bits
	row   int
	col   int
}

// NewShape intializes a shape described by the 2D grid with a given id and
// grid position at the upper left (0, 0), then the masks are updated from
// the current position.
func NewShape(id int, grid [][]int) Shape {
	s := Shape{id, grid, 0, 0, 0, 0}
	(&s).ComputeMask()
	return s
}

// MakeShapes generates an array of Shapes from an array of 2D grids using
// NewShape. The index of each grid in the grids array is used as the id
// for that shape.  The grids for each shape do not need to have the same
// dimensions, so it is possible to initialize a set of shapes like so:
//
//     grids := [][][]int { {
// 	    {1, 1, 0}, {1, 1, 1}}, {
// 	    {1, 0, 1}, {1, 1, 1}}, {
// 	    {1, 0, 0, 0}, {1, 1, 1, 1}, {1, 0, 0, 0}}}
//     return shape.MakeShapes(grids)
//
// The above code returns an array of 3 shapes with id's 1, 2, and 3.
//
func MakeShapes(grids [][][]int) []Shape {
	shapes := []Shape{}
	for id, grid := range grids {
		s := NewShape(id+1, grid)
		shapes = append(shapes, s)
	}
	return shapes
}

// NumRows returns the number of rows in a shape, equivalent to the length
// of the first grid dimension.
func (s Shape) NumRows() int {
	return len(s.shape)
}

// NumCols returns the number of columns in a shape, in other words, the
// length of the second grid dimension.
func (s Shape) NumCols() int {
	return len(s.shape[0])
}

// ID returns the integer id of the Shape.
func (s Shape) ID() int {
	return s.id
}

// String formats a shape into text, one line for each row, and each column
// represented as a string of 0 and 1.
func (s Shape) String() string {
	buf := fmt.Sprintf("Shape #%3d, mask:%v\n", s.id, s.mask)
	for r := 0; r < s.NumRows(); r++ {
		buf += fmt.Sprintf("%v\n", s.shape[r])
	}
	return buf
}

// ComputeMask computes the bit mask from the Shape's current grid.
func (s *Shape) ComputeMask() mask.Bits {

	s.mask, s.gaps = mask.ComputeMask(s.shape)
	return s.mask
}

// Clip clears all bits in the Shape which are not within the given region
// mask.
func (s Shape) Clip(region mask.Bits) Shape {
	s.mask = s.mask & region
	s.gaps = s.gaps & region
	return s
}

// Mask returns the bit mask for the Shape.
func (s Shape) Mask() mask.Bits {
	return s.mask
}

// GapMask returns the bit mask for any gaps in the Shape.
func (s Shape) GapMask() mask.Bits {
	return s.gaps
}

// OutlineMask returns a mask with only the bits of Shape which are not part
// of the gap mask.
func (s Shape) OutlineMask() mask.Bits {
	return s.mask & (^s.gaps)
}

// Translate moves the Shape by rows r and columns c, then recomputes the
// mask and gaps for the new location.
func (s Shape) Translate(r int, c int) Shape {
	s.row = s.row + r
	s.col = s.col + c
	s.mask = s.mask.Translate(r, c)
	s.gaps = s.gaps.Translate(r, c)
	return s
}

func (s Shape) rotate() Shape {
	// Rotate 90 degrees clockwise
	nrow := s.NumCols()
	ncol := s.NumRows()
	grid := make([][]int, nrow)
	for r := 0; r < nrow; r++ {
		grid[r] = make([]int, ncol)
		for c := 0; c < ncol; c++ {
			grid[r][c] = s.shape[ncol-c-1][r]
		}
	}
	return NewShape(s.id, grid)
}

// Equals compares the Shape with the Shape b and returns true if the sizes
// of the Shapes are the same and all of their grid values match.  So both
// mask and gap grid points must be the same.
func (s Shape) Equals(b Shape) bool {
	nrow := s.NumRows()
	ncol := s.NumCols()
	result := (nrow == b.NumRows() && ncol == b.NumCols())
	if result {
		for r := 0; r < nrow; r++ {
			for c := 0; c < ncol; c++ {
				result = result && (s.shape[r][c] == b.shape[r][c])
			}
		}
	}
	// fmt.Println("comparing %v == %v ==> %v", s, b, result)
	return result
}

func (s Shape) flip() Shape {
	nrow := s.NumRows()
	ncol := s.NumCols()
	grid := make([][]int, nrow)
	for r := 0; r < nrow; r++ {
		grid[r] = make([]int, ncol)
		for c := 0; c < ncol; c++ {
			grid[r][c] = s.shape[nrow-r-1][c]
		}
	}
	return NewShape(s.id, grid)
}

func searchShapes(shapes []Shape, pred func(s Shape) bool) (bool, int) {
	for i, s := range shapes {
		if pred(s) {
			return true, i
		}
	}
	return false, 0
}

// Permutations returns all distinct Shapes generated from rotating Shape
// and flipping Shape all possible ways.
func (s Shape) Permutations() []Shape {
	shapes := []Shape{}
	for i := 0; i < 8; i++ {
		if i == 4 {
			s = s.flip()
		}
		found, _ := searchShapes(shapes, func(b Shape) bool { return b.Equals(s) })
		if !found {
			shapes = append(shapes, s)
		}
		s = s.rotate()
	}
	return shapes
}
