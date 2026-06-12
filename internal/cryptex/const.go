package cryptex

import (
	"math"
)

const (
	iterSize = 0x04
	metaSize = 0x10 - iterSize
	saltSize = 0x10

	keySize = 0x20

	MaxIter = math.MaxInt32 // is limited for using on 32-bit systems
)
