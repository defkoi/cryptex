package cryptex

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

func (c *Cryptex) Encode(w io.Writer, password string) error {
	salt := generateRand(saltSize)
	iv := generateRand(ivSize)

	data := encodeMap(c.data)
	data = appendPadding(data, aes.BlockSize)
	if err := encrypt(data, password, salt, iv); err != nil {
		return err
	}

	encryptedData := encryptedData{
		version: Version,
		salt:    salt,
		iv:      iv,
		data:    data,
	}

	if _, err := w.Write(encryptedData.encode()); err != nil {
		return err
	}

	return nil
}

func encrypt(data []byte, password string, salt []byte, iv []byte) error {
	key, err := keyFromPassword(password, salt, iterations)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	cipher.NewCBCEncrypter(block, iv).CryptBlocks(data, data)

	return nil
}
