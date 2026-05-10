package cryptex

import "crypto/aes"

const (
	keySize  = 0x20
	saltSize = 0x0f
	ivSize   = aes.BlockSize

	iterations = 0x100000
)
