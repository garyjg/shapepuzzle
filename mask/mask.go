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

