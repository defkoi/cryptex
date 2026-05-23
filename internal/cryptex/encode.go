package cryptex

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

func (c *Cryptex) Encode(w io.Writer, password string) error {
	iter := uint32(c.iter)
	salt := generateRand(saltSize)
	iv := generateRand(ivSize)

	data := encodeMap(c.data)
	data = appendPadding(data, aes.BlockSize)

	key, err := keyFromPassword(password, salt, iter)
	if err != nil {
		return err
	}

	if err := encrypt(data, key, iv); err != nil {
		return err
	}

	encryptedData := encryptedData{
		ver:  Version,
		iter: iter,
		salt: salt,
		iv:   iv,
		data: data,
	}

	if _, err := w.Write(encryptedData.encode()); err != nil {
		return err
	}

	return nil
}

func encrypt(data []byte, key []byte, iv []byte) error {
	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	cipher.NewCBCEncrypter(block, iv).CryptBlocks(data, data)

	return nil
}
