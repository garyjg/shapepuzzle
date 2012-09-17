// -*- tab-width: 4; -*-

package mask

import (
    "fmt"
)

type MaskBits uint64


func FirstMaskBit() MaskBits {
	return MaskBits(0x8000000000000000)
}


func (mask MaskBits) String() string {
	return fmt.Sprintf("0x%016x", uint64(mask))
}


// Turn a 2D array or slice into a mask on a nrow X ncol grid.  The
// most-significant bit in the mask is for the upper left corner of the
// grid, r=0 and c=0, indexed in row major order.  Each row always
// starts at a new byte, no matter how many columns.  That way the mask
// is independent of the board size and can be stashed with the shape.

func ComputeMask(grid [][]int) (MaskBits, MaskBits) {

	mbits := MaskBits(0)
	gapbits := MaskBits(0)
	for r := 0; r < len(grid); r += 1 {
		bit := FirstMaskBit() >> uint(r * 8)
		for c := 0; c < len(grid[0]); c += 1 {
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


func (mask MaskBits) Translate(row int, col int) MaskBits {

	// It is possible to translate a mask past the edges (positive and
	// negative), in which it should not wrap to the next row but be
	// truncated instead.

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
		for i := 0; i < 8; i += 1 {
		    blank = (blank << 8) + (0xff >> uint(-col))
		}
		mask = mask & MaskBits(blank)
		mask = mask << uint(-col)
	} else {
		for i := 0; i < 8; i += 1 {
		    blank = (blank << 8) + ((0xff << uint(col)) & 0xff)
		}
		mask = mask & MaskBits(blank)
		mask = mask >> uint(col)
	}   
	return mask
}

