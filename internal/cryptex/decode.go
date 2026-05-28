package cryptex

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

func Decode(r io.Reader, password string) (*Cryptex, error) {
	encData, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	decodedData, err := decodeEncryptedData(encData)
	if err != nil {
		return nil, err
	}

	iter := decodedData.iter
	salt := decodedData.salt
	iv := decodedData.iv
	data := decodedData.data

	key, err := keyFromPassword(password, salt, iter)
	if err != nil {
		return nil, err
	}

	if err := decrypt(data, key, iv); err != nil {
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

func decrypt(data []byte, key []byte, iv []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	cipher.NewCBCDecrypter(block, iv).CryptBlocks(data, data)

	return nil
}
