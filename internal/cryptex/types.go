package cryptex

import (
	"crypto/aes"
)

const encrypredDataMinLength = 1 + saltSize + ivSize + aes.BlockSize

type encryptedData struct {
	version byte
	salt    []byte
	iv      []byte
	data    []byte
}

func decodeEncryptedData(data []byte) (encryptedData, error) {
	if len(data) < encrypredDataMinLength ||
		len(data)-encrypredDataMinLength%aes.BlockSize != 0 {
		return encryptedData{}, ErrInvalidSize
	}

	return encryptedData{
		version: data[0],
		salt:    data[1:saltSize],
		iv:      data[1+saltSize : 1+saltSize+ivSize],
		data:    data[1+saltSize+ivSize:],
	}, nil
}

func (d *encryptedData) encode() []byte {
	b := make([]byte, 0, 1+saltSize+ivSize+len(d.data))
	b = append(b, d.version)
	b = append(b, d.salt...)
	b = append(b, d.iv...)
	b = append(b, d.data...)
	return b
}
