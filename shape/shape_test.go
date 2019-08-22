// -*- tab-width: 4; -*-

package shape

import (
	"reflect"
	"testing"

	"github.com/garyjg/shapepuzzle/mask"
)

func TestNewShape(t *testing.T) {
	type args struct {
		id   int
		grid [][]int
	}
	tests := []struct {
		name string
		args args
		want Shape
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewShape(tt.args.id, tt.args.grid); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewShape() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMakeShapes(t *testing.T) {
	type args struct {
		grids [][][]int
	}
	tests := []struct {
		name string
		args args
		want []Shape
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MakeShapes(tt.args.grids); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MakeShapes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShape_NumRows(t *testing.T) {
	type fields struct {
		id    int
		shape [][]int
		mask  mask.Bits
		gaps  mask.Bits
		row   int
		col   int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shape{
				id:    tt.fields.id,
				shape: tt.fields.shape,
				mask:  tt.fields.mask,
				gaps:  tt.fields.gaps,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			if got := s.NumRows(); got != tt.want {
				t.Errorf("Shape.NumRows() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShape_NumCols(t *testing.T) {
	type fields struct {
		id    int
		shape [][]int
		mask  mask.Bits
		gaps  mask.Bits
		row   int
		col   int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shape{
				id:    tt.fields.id,
				shape: tt.fields.shape,
				mask:  tt.fields.mask,
				gaps:  tt.fields.gaps,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			if got := s.NumCols(); got != tt.want {
				t.Errorf("Shape.NumCols() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShape_ID(t *testing.T) {
	type fields struct {
		id    int
		shape [][]int
		mask  mask.Bits
		gaps  mask.Bits
		row   int
		col   int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shape{
				id:    tt.fields.id,
				shape: tt.fields.shape,
				mask:  tt.fields.mask,
				gaps:  tt.fields.gaps,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			if got := s.ID(); got != tt.want {
				t.Errorf("Shape.ID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShape_String(t *testing.T) {
	type fields struct {
		id    int
		shape [][]int
		mask  mask.Bits
		gaps  mask.Bits
		row   int
		col   int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shape{
				id:    tt.fields.id,
				shape: tt.fields.shape,
				mask:  tt.fields.mask,
				gaps:  tt.fields.gaps,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			if got := s.String(); got != tt.want {
				t.Errorf("Shape.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShape_ComputeMask(t *testing.T) {
	type fields struct {
		id    int
		shape [][]int
		mask  mask.Bits
		gaps  mask.Bits
		row   int
		col   int
	}
	tests := []struct {
		name   string
		fields fields
		want   mask.Bits
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shape{
				id:    tt.fields.id,
				shape: tt.fields.shape,
				mask:  tt.fields.mask,
				gaps:  tt.fields.gaps,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			if got := s.ComputeMask(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Shape.ComputeMask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShape_Clip(t *testing.T) {
	type fields struct {
		id    int
		shape [][]int
		mask  mask.Bits
		gaps  mask.Bits
		row   int
		col   int
	}
	type args struct {
		region mask.Bits
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Shape
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shape{
				id:    tt.fields.id,
				shape: tt.fields.shape,
				mask:  tt.fields.mask,
				gaps:  tt.fields.gaps,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			if got := s.Clip(tt.args.region); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Shape.Clip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShape_Mask(t *testing.T) {
	type fields struct {
		id    int
		shape [][]int
		mask  mask.Bits
		gaps  mask.Bits
		row   int
		col   int
	}
	tests := []struct {
		name   string
		fields fields
		want   mask.Bits
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shape{
				id:    tt.fields.id,
				shape: tt.fields.shape,
				mask:  tt.fields.mask,
				gaps:  tt.fields.gaps,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			if got := s.Mask(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Shape.Mask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShape_GapMask(t *testing.T) {
	type fields struct {
		id    int
		shape [][]int
		mask  mask.Bits
		gaps  mask.Bits
		row   int
		col   int
	}
	tests := []struct {
		name   string
		fields fields
		want   mask.Bits
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shape{
				id:    tt.fields.id,
				shape: tt.fields.shape,
				mask:  tt.fields.mask,
				gaps:  tt.fields.gaps,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			if got := s.GapMask(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Shape.GapMask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShape_OutlineMask(t *testing.T) {
	type fields struct {
		id    int
		shape [][]int
		mask  mask.Bits
		gaps  mask.Bits
		row   int
		col   int
	}
	tests := []struct {
		name   string
		fields fields
		want   mask.Bits
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shape{
				id:    tt.fields.id,
				shape: tt.fields.shape,
				mask:  tt.fields.mask,
				gaps:  tt.fields.gaps,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			if got := s.OutlineMask(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Shape.OutlineMask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShape_Translate(t *testing.T) {
	type fields struct {
		id    int
		shape [][]int
		mask  mask.Bits
		gaps  mask.Bits
		row   int
		col   int
	}
	type args struct {
		r int
		c int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   Shape
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shape{
				id:    tt.fields.id,
				shape: tt.fields.shape,
				mask:  tt.fields.mask,
				gaps:  tt.fields.gaps,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			if got := s.Translate(tt.args.r, tt.args.c); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Shape.Translate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShape_rotate(t *testing.T) {
	type fields struct {
		id    int
		shape [][]int
		mask  mask.Bits
		gaps  mask.Bits
		row   int
		col   int
	}
	tests := []struct {
		name   string
		fields fields
		want   Shape
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shape{
				id:    tt.fields.id,
				shape: tt.fields.shape,
				mask:  tt.fields.mask,
				gaps:  tt.fields.gaps,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			if got := s.rotate(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Shape.rotate() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShape_Equals(t *testing.T) {
	// predefine some simple shape grids
	cross := [][]int{{0, 1, 0}, {1, 1, 1}, {0, 1, 0}}
	bracket := [][]int{{1, 1}, {1, 0}, {1, 1}}
	type fields struct {
		id    int
		shape [][]int
	}
	type args struct {
		b Shape
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{"simple cross equals itself", fields{1, cross}, args{NewShape(1, cross)}, true},
		{"equals itself but different id", fields{1, cross}, args{NewShape(2, cross)}, true},
		{"cross does not equal bracket", fields{1, cross}, args{NewShape(2, bracket)}, false},
		{"bracket equals bracket", fields{1, bracket}, args{NewShape(2, bracket)}, true},
		{"bracket == flipped bracket", fields{1, bracket}, args{NewShape(2, bracket).flip()}, true},
		{"bracket != rotated bracket", fields{1, bracket}, args{NewShape(2, bracket).rotate()}, false},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewShape(tt.fields.id, tt.fields.shape)
			if got := s.Equals(tt.args.b); got != tt.want {
				t.Errorf("Shape.Equals() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShape_flip(t *testing.T) {
	type fields struct {
		id    int
		shape [][]int
		mask  mask.Bits
		gaps  mask.Bits
		row   int
		col   int
	}
	tests := []struct {
		name   string
		fields fields
		want   Shape
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shape{
				id:    tt.fields.id,
				shape: tt.fields.shape,
				mask:  tt.fields.mask,
				gaps:  tt.fields.gaps,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			if got := s.flip(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Shape.flip() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchShapes(t *testing.T) {
	type args struct {
		shapes []Shape
		pred   func(s Shape) bool
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := searchShapes(tt.args.shapes, tt.args.pred)
			if got != tt.want {
				t.Errorf("searchShapes() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("searchShapes() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestShape_Permutations(t *testing.T) {
	type fields struct {
		id    int
		shape [][]int
		mask  mask.Bits
		gaps  mask.Bits
		row   int
		col   int
	}
	tests := []struct {
		name   string
		fields fields
		want   []Shape
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := Shape{
				id:    tt.fields.id,
				shape: tt.fields.shape,
				mask:  tt.fields.mask,
				gaps:  tt.fields.gaps,
				row:   tt.fields.row,
				col:   tt.fields.col,
			}
			if got := s.Permutations(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Shape.Permutations() = %v, want %v", got, tt.want)
			}
		})
	}
}
