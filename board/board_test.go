// -*- tab-width: 4; -*-

package board

import (
    "testing"
	"shapepuzzle/shape"
	"shapepuzzle/mask"
)


func TestRegionMask(t *testing.T) {

	if NewBoard(4, 4).RegionMask() != mask.MaskBits(0xf0f0f0f000000000) {
	    t.Errorf("Wrong region mask for 5x5")
	}
	if NewBoard(8, 8).RegionMask() != mask.MaskBits(0xffffffffffffffff) {
	    t.Errorf("Wrong region mask for 8x8")
	}
}


func testShapes() []shape.Shape {
    grids := [][][]int { { 
    	{1, 1, 0}, {1, 1, 1}}, {
		{1, 0, 1}, {1, 1, 1}}, {
		{1, 1, 1}, {0, 0, 1}} }
	return shape.MakeShapes(grids)
}


func checkReject(t *testing.T, b Board, rejects []shape.Shape, result bool) {
	reject := RejectBoard(b, rejects)
	if ! reject && result {
	    t.Errorf("Board should be rejected: %v", b)
	} else if reject && ! result {
	    t.Errorf("Board should NOT be rejected: %v", b)
	}
}


func TestNewBoard(t *testing.T) {

    b := NewBoard(8, 8)
    rejects := GapShapes(b)

	shapes := testShapes()
	
	checkReject(t, b.Place(shapes[1]), rejects, true)
	checkReject(t, b.Place(shapes[0].Translate(0, 5)), rejects, true)
	checkReject(t, b.Place(shapes[1].Translate(0, 1)), rejects, true)
	checkReject(t, b.Place(shapes[0].Translate(0, 4)), rejects, false)
	checkReject(t, b.Place(shapes[1].Translate(1, 1)), rejects, false)

	b = NewBoard(5, 5)
    rejects = GapShapes(b)
	checkReject(t, b.Place(shapes[2].Translate(3, 0)), rejects, true)

}
