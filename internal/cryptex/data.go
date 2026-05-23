package cryptex

import (
	"crypto/aes"
	"encoding/binary"
)

const (
	encryptedDataHeaderSize = 1 + 4 + saltSize + ivSize
)

type encryptedData struct {
	ver  uint8
	iter uint32
	salt []byte
	iv   []byte
	data []byte
}

func decodeEncryptedData(data []byte) (encryptedData, error) {
	const minSize = encryptedDataHeaderSize + aes.BlockSize

	if len(data) < minSize ||
		(len(data)-encryptedDataHeaderSize)%aes.BlockSize != 0 {
		return encryptedData{}, ErrInvalidSize
	}

	return encryptedData{
		ver:  data[0],
		iter: binary.LittleEndian.Uint32(data[1:5]),
		salt: data[5 : 5+saltSize],
		iv:   data[5+saltSize : 5+saltSize+ivSize],
		data: data[5+saltSize+ivSize:],
	}, nil
}

func (d *encryptedData) encode() []byte {
	b := make([]byte, 0, encryptedDataHeaderSize+len(d.data))

	b = append(b, d.ver)
	b = binary.LittleEndian.AppendUint32(b, d.iter)
	b = append(b, d.salt...)
	b = append(b, d.iv...)
	b = append(b, d.data...)

	return b
}
