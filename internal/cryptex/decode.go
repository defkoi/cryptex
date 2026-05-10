package cryptex

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"io"
)

func Decode(r io.Reader, password string) (*Cryptex, error) {
	encData, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	parsedEncData, err := decodeEncryptedData(encData)
	if err != nil {
		return nil, err
	}

	if err := parsedEncData.validate(); err != nil {
		return nil, err
	}

	_, salt, iv, data :=
		parsedEncData.version,
		parsedEncData.salt,
		parsedEncData.iv,
		parsedEncData.data

	key, err := keyFromPassword(password, salt, iterations)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)

	blockMode.CryptBlocks(data, data)

	data, err = removePadding(data)
	if err != nil {
		return nil, err
	}

	dataMap := make(map[string]string)
	if err := gob.NewDecoder(
		bytes.NewReader(data),
	).Decode(&dataMap); err != nil {
		return nil, err
	}

	return &Cryptex{data: dataMap}, nil
}
