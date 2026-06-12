package cryptex

import (
	"crypto/aes"
	"crypto/cipher"
	"io"
)

func (c *Cryptex) Encode(w io.Writer, password string) error {
	salt := generateRand(saltSize)

	data := encodeMap(c.data)

	data = appendPadding(data, aes.BlockSize)

	key, err := keyFromPassword(password, salt, c.iter)
	if err != nil {
		return err
	}

	data, err = encrypt(data, key)
	if err != nil {
		return err
	}

	encryptedData := encryptedData{
		iter: c.iter,
		meta: make([]byte, metaSize),
		salt: salt,
		data: data,
	}

	if _, err := w.Write(encryptedData.encode()); err != nil {
		return err
	}

	return nil
}

func encrypt(data []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCMWithRandomNonce(block)
	if err != nil {
		return nil, err
	}

	return gcm.Seal(nil, nil, data, nil), nil
}
