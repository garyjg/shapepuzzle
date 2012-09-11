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

func ComputeMask(grid [][]int) MaskBits {

	mbits := MaskBits(0)
	for r := 0; r < len(grid); r += 1 {
		bit := FirstMaskBit() >> uint(r * 8)
		for c := 0; c < len(grid[0]); c += 1 {
			if grid[r][c] != 0 {
				mbits = mbits | bit
			}
			bit = bit >> 1
		}
	}
	return mbits
}


func (mask MaskBits) Translate(row int, col int) MaskBits {
	 return (mask >> uint(row * 8)) >> uint(col)
}

