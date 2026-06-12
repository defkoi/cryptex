package cryptex

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

func Decode(r io.Reader, password string) (*Cryptex, error) {
	readed, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	decodedData, err := decodeEncryptedData(readed)
	if err != nil {
		return nil, err
	}

	iter := decodedData.iter
	_ = decodedData.meta
	salt := decodedData.salt
	data := decodedData.data

	key, err := keyFromPassword(password, salt, iter)
	if err != nil {
		return nil, err
	}

	data, err = decrypt(data, key)
	if err != nil {
		return nil, err
	}

	data, err = removePadding(data)
	if err != nil {
		return nil, err
	}

	dataMap, err := decodeMap(data)
	if err != nil {
		return nil, err
	}

	return &Cryptex{
		data: dataMap,
		iter: iter,
	}, nil
}

func decrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return nil, err
	}

	return gcm.Open(nil, nil, data, nil)
}
