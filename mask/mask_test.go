// -*- tab-width: 4; -*-

package mask

import(
    "testing"
)


func TestMaskEqual(t *testing.T) {
	bits := MaskBits(0xffffffffffffffff)
	if bits.Translate(0, 0) != bits {
	    t.Errorf("(0,0) translation should not change mask.")
	}

}

func TestMaskNegativeShift(t *testing.T) {
	bits := MaskBits(0xffffffffffffffff)
	if bits.Translate(-2, -1) != MaskBits(0xfefefefefefe0000) {
	    t.Errorf("(-2,-1) translation should have two blank rows and a blank column.")
	}
}

func TestMaskPositiveShift(t *testing.T) {
	bits := MaskBits(0xf0f0f0f000000000)
	if bits.Translate(4, 4) != MaskBits(0x000000000f0f0f0f) {
	    t.Errorf("(4,4) %v translation should move square to lower right, got %v.",
	        MaskBits(0xf0f0f0f000000000), bits.Translate(4, 4))
	}
}

func TestMaskShiftTooFar(t *testing.T) {
	bits := MaskBits(0xf0f0f0f000000000)
	if bits.Translate(4, 5) != MaskBits(0x0000000007070707) {
	    t.Errorf("(4,4) %v translation should move square to lower right, got %v.",
	        MaskBits(0xf0f0f0f000000000), bits.Translate(4, 5))
	}
}

