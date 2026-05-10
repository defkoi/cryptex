package cryptex

import (
	"crypto/aes"
)

const encDataMinLength = 1 + saltSize + ivSize

type encryptedData struct {
	version byte
	salt    []byte
	iv      []byte
	data    []byte
}

func decodeEncryptedData(data []byte) (*encryptedData, error) {
	if len(data) < encDataMinLength {
		return nil, ErrInvalidSize
	}

	return &encryptedData{
		version: data[0],
		salt:    data[1:saltSize],
		iv:      data[1+saltSize : 1+saltSize+ivSize],
		data:    data[1+saltSize+ivSize:],
	}, nil
}

func (d *encryptedData) validate() error {
	if len(d.data) == 0 || len(d.data)%aes.BlockSize != 0 {
		return ErrInvalidSize
	}
	return nil
}

func (d *encryptedData) encode() []byte {
	b := make([]byte, 0, encDataMinLength)
	b = append(b, d.version)
	b = append(b, d.salt...)
	b = append(b, d.iv...)
	b = append(b, d.data...)
	return b
}
