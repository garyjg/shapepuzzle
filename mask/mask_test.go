// -*- tab-width: 4; -*-

package mask

import (
	"testing"
)

func testTranslate(t *testing.T, start Bits, r int, c int,
	expect Bits, note string) {

	got := start.Translate(r, c)
	if got != expect {
		t.Errorf("(%d,%d) translation got %v, expected %v: %v",
			r, c, got, expect, note)
	}
}

func TestString(t *testing.T) {
	type stringTest struct {
		bits Bits
		text string
	}
	tests := []stringTest{
		{Bits(0xffffffffffffffff), "0xffffffffffffffff"},
		{Bits(0x1), "0x0000000000000001"},
		{Bits(0x8000000000000000), "0x8000000000000000"},
	}
	for _, test := range tests {
		got := test.bits.String()
		if got != test.text {
			t.Errorf("got Bits string %s, expected %s", got, test.text)
		}
	}
}

func TestMaskEqual(t *testing.T) {
	all := Bits(0xffffffffffffffff)
	testTranslate(t, all, 0, 0, all, "zero translate should be the same")
}

func TestMaskNegativeShift(t *testing.T) {
	testTranslate(t, Bits(0xffffffffffffffff), -2, -1,
		Bits(0xfefefefefefe0000), "")
}

func TestMaskPositiveShift(t *testing.T) {
	testTranslate(t, Bits(0xf0f0f0f000000000), 4, 4,
		Bits(0x000000000f0f0f0f), "move square to lower right")
}

func TestMaskShiftTooFar(t *testing.T) {
	testTranslate(t, Bits(0xf0f0f0f000000000), 4, 5,
		Bits(0x0000000007070707), "move past lower right")
}

func TestZero(t *testing.T) {
	testTranslate(t, Bits(0xf0f0f0f000000000), 8, 8, Bits(0),
		"move off the board should be zero")
}
