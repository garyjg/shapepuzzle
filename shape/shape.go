// -*- tab-width: 4; -*-

package shape

import (
    "fmt"
	"shapepuzzle/mask"
)

type Shape struct {
	id    int
	shape [][]int
	mask  mask.MaskBits
	gaps  mask.MaskBits
	row	  int
	col	  int
}

func NewShape(id int, grid [][]int) Shape {
	 s := Shape{ id, grid, 0, 0, 0, 0 }
	 (&s).ComputeMask()
	 return s
}

func (s Shape) NumRows() int {
	return len(s.shape)
}

func (s Shape) NumCols() int {
	return len(s.shape[0])
}


func (s Shape) ID() int {
	return s.id
}


func (s Shape) String() string {
	buf := fmt.Sprintf("Shape #%3d, mask:%v\n", s.id, s.mask)
	for r := 0; r < s.NumRows(); r += 1 {
		buf += fmt.Sprintf("%v\n", s.shape[r])
	}
	return buf
}


func (s *Shape) ComputeMask() mask.MaskBits {

	s.mask, s.gaps = mask.ComputeMask(s.shape)
	return s.mask
}


func (s Shape) Mask() mask.MaskBits {
	return s.mask
}


func (s Shape) GapMask() mask.MaskBits {
	return s.gaps
}


func (s Shape) OutlineMask() mask.MaskBits {
    return s.mask & (^ s.gaps)
}


func (s Shape) Translate(r int, c int) (Shape) {
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
	for r := 0; r < nrow; r += 1 {
		grid[r] = make([]int, ncol)
		for c := 0; c < ncol; c += 1 {
			grid[r][c] = s.shape[ncol-c-1][r]
		}
	}
	return NewShape(s.id, grid)
}

func (s Shape) Equals(b Shape) bool {
	nrow := s.NumRows()
	ncol := s.NumCols()
	result := (nrow == b.NumRows() && ncol == b.NumCols())
	if result {
		for r := 0; r < nrow; r += 1 {
			for c := 0; c < ncol; c += 1 {
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
	for r := 0; r < nrow; r += 1 {
		grid[r] = make([]int, ncol)
		for c := 0; c < ncol; c += 1 {
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

func (s Shape) Permutations() []Shape {
	shapes := []Shape{}
	for i := 0; i < 8; i += 1 {
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


