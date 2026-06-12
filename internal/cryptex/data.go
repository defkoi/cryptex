package cryptex

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/binary"
)

const encryptedDataHeaderSize = iterSize + metaSize + saltSize

var overhead int

func init() {
	block, err := aes.NewCipher(make([]byte, keySize))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		panic(err)
	}

	overhead = gcm.Overhead()
}

type encryptedData struct {
	// header
	iter uint32
	meta []byte // unused
	salt []byte

	// data
	data []byte // include nonce
}

func decodeEncryptedData(data []byte) (encryptedData, error) {
	if len(data) < encryptedDataHeaderSize+aes.BlockSize+overhead ||
		(len(data)-encryptedDataHeaderSize-overhead)%aes.BlockSize != 0 {
		return encryptedData{}, ErrInvalidSize
	}

	return encryptedData{
		iter: binary.LittleEndian.Uint32(data[:iterSize]),
		meta: data[iterSize : iterSize+metaSize],
		salt: data[iterSize+metaSize : iterSize+metaSize+saltSize],
		data: data[iterSize+metaSize+saltSize:],
	}, nil
}

func (d *encryptedData) encode() []byte {
	b := make([]byte, 0, encryptedDataHeaderSize+len(d.data))

	b = binary.LittleEndian.AppendUint32(b, d.iter)
	b = append(b, d.meta[:metaSize]...)
	b = append(b, d.salt[:saltSize]...)
	b = append(b, d.data...)

	return b
}
