package cryptex

import (
	"crypto/aes"
	"math"
)

const (
	keySize  = 0x20
	saltSize = 0x0f
	ivSize   = aes.BlockSize

	MaxIter = math.MaxInt32 // is limited for using on x32 systems
)
