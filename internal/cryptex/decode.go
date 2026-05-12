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

	salt := decodedData.salt
	iv := decodedData.iv
	data := decodedData.data

	if err := decrypt(data, password, salt, iv); err != nil {
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

	return &Cryptex{data: dataMap}, nil
}

func decrypt(data []byte, password string, salt []byte, iv []byte) error {
	key, err := keyFromPassword(password, salt, iterations)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	cipher.NewCBCDecrypter(block, iv).CryptBlocks(data, data)

	return nil
}
