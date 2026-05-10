package cryptex

import (
	"crypto/aes"
)

const (
	encryptedDataHeaderSize = 1 + saltSize + ivSize
)

type encryptedData struct {
	version byte
	salt    []byte
	iv      []byte
	data    []byte
}

func decodeEncryptedData(data []byte) (encryptedData, error) {
	const minSize = encryptedDataHeaderSize + aes.BlockSize

	if len(data) < minSize ||
		(len(data)-encryptedDataHeaderSize)%aes.BlockSize != 0 {
		return encryptedData{}, ErrInvalidSize
	}

	return encryptedData{
		version: data[0],
		salt:    data[1 : 1+saltSize],
		iv:      data[1+saltSize : 1+saltSize+ivSize],
		data:    data[1+saltSize+ivSize:],
	}, nil
}

func (d *encryptedData) encode() []byte {
	b := make([]byte, 0, encryptedDataHeaderSize+len(d.data))
	b = append(b, d.version)
	b = append(b, d.salt...)
	b = append(b, d.iv...)
	b = append(b, d.data...)
	return b
}
