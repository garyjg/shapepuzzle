// -*- tab-width: 4; -*-

package mask

import(
    "testing"
)


func testTranslate(t *testing.T, start MaskBits, r int, c int, 
	 			   expect MaskBits, note string) {

	got := start.Translate(r, c)
	if got != expect {
	    t.Errorf("(%d,%d) translation got %v, expected %v: %v", 
				 r, c, got, expect, note)
	}
}


func TestMaskEqual(t *testing.T) {
	all := MaskBits(0xffffffffffffffff)
	testTranslate(t, all, 0, 0, all, "zero translate should be the same")
}

func TestMaskNegativeShift(t *testing.T) {
	testTranslate(t, MaskBits(0xffffffffffffffff), -2, -1, 
					 MaskBits(0xfefefefefefe0000), "")
}

func TestMaskPositiveShift(t *testing.T) {
	testTranslate(t, MaskBits(0xf0f0f0f000000000), 4, 4, 
					 MaskBits(0x000000000f0f0f0f), "move square to lower right")
}

func TestMaskShiftTooFar(t *testing.T) {
	testTranslate(t, MaskBits(0xf0f0f0f000000000), 4, 5, 
					 MaskBits(0x0000000007070707), "move past lower right")
}

func TestZero(t *testing.T) {
	testTranslate(t, MaskBits(0xf0f0f0f000000000), 8, 8, MaskBits(0),
					 "move off the board should be zero")
}
	
