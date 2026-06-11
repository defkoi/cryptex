package cryptex

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
	"io"
)

func (c *Cryptex) Encode(w io.Writer, password string) error {
	salt := generateRand(saltSize)
	iv := generateRand(ivSize)

	data := encodeMap(c.data)
	data = appendPadding(data, aes.BlockSize)

	key, err := keyFromPassword(password, salt, c.iter)
	if err != nil {
		return err
	}

	data, err = encrypt(data, key, iv, c.mode)
	if err != nil {
		return err
	}

	encryptedData := encryptedData{
		mode: c.mode,
		iter: c.iter,
		salt: salt,
		iv:   iv,
		data: data,
	}

	if _, err := w.Write(encryptedData.encode()); err != nil {
		return err
	}

	return nil
}

func encrypt(data []byte, key []byte, iv []byte, mode Mode) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	switch mode {
	case ModeCBC:
		ret := make([]byte, len(data))
		cipher.NewCBCEncrypter(block, iv).CryptBlocks(ret, data)
		return ret, nil
	case ModeGCM:
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		nonce := iv[:gcm.NonceSize()]
		return gcm.Seal(nil, nonce, data, nil), nil
	default:
		return nil, errors.New("unsupported mode")
	}
}
