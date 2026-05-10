package cryptex

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"io"
)

func (c *Cryptex) Encode(w io.Writer, password string) error {
	salt := generateRand(saltSize)

	key, err := keyFromPassword(password, salt, iterations)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	iv := generateRand(ivSize)

	blockMode := cipher.NewCBCEncrypter(block, iv)

	buf := bytes.NewBuffer([]byte{})

	if err := gob.NewEncoder(buf).Encode(c.data); err != nil {
		return err
	}

	data := appendPadding(buf.Bytes(), aes.BlockSize)

	blockMode.CryptBlocks(data, data)

	encData := encryptedData{
		version: Version,
		salt:    salt,
		iv:      iv,
		data:    data,
	}

	if _, err := w.Write(encData.encode()); err != nil {
		return err
	}

	return nil
}
