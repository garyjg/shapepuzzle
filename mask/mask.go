// -*- tab-width: 4; -*-

// Package mask provides a bit mask type which represents cells on a 2D
// grid.  For example, the mask can keep track of which cells are occupied
// by shapes on a puzzle grid, and it can also represent a puzzle shape.
package mask

import (
    "fmt"
)

// Bits type is an 8-byte unsigned integer to efficiently represent a
// 2D mask on any grid up to 8x8.
type Bits uint64


// FirstBit returns a Bits mask initialized with only the very
// first bit set, corresponding to the upper left corner of a grid.
func FirstBit() Bits {
	return Bits(0x8000000000000000)
}


// String method of Bits formats the mask into a hexadecimal
// representation.
func (mask Bits) String() string {
	return fmt.Sprintf("0x%016x", uint64(mask))
}


// ComputeMask turns a 2D array or slice into a mask on a nrow X ncol grid.
// The most-significant bit in the mask is for the upper left corner of the
// grid, r=0 and c=0, indexed in row major order.  Each row always starts
// at a new byte, no matter how many columns.  That way the mask is
// independent of the board size and can be stashed with the shape.
func ComputeMask(grid [][]int) (Bits, Bits) {

	mbits := Bits(0)
	gapbits := Bits(0)
	for r := 0; r < len(grid); r++ {
		bit := FirstBit() >> uint(r * 8)
		for c := 0; c < len(grid[0]); c++ {
			if grid[r][c] != 0 {
				mbits = mbits | bit
			}
			if grid[r][c] == 2 {
			    gapbits = gapbits | bit
			}
			bit = bit >> 1
		}
	}
	return mbits, gapbits
}


// Translate adjusts mask as if the bit pattern on the grid were translated
// by row and col places.  It is possible to translate a mask past the
// edges (positive and negative), in which case it is truncated rather than
// wrapped.
func (mask Bits) Translate(row int, col int) Bits {

	if row < 0 {
	    mask = mask << uint(-row * 8)
		row = 0
	} else {
	    mask = mask >> uint(row * 8)
	}

	// Truncating columns means ANDing each byte before shifting it to the
	// right column.
    blank := uint64(0)
	if col < 0 {
		for i := 0; i < 8; i++ {
		    blank = (blank << 8) + (0xff >> uint(-col))
		}
		mask = mask & Bits(blank)
		mask = mask << uint(-col)
	} else {
		for i := 0; i < 8; i++ {
		    blank = (blank << 8) + ((0xff << uint(col)) & 0xff)
		}
		mask = mask & Bits(blank)
		mask = mask >> uint(col)
	}   
	return mask
}

