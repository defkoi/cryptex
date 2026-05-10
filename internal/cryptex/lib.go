package cryptex

import (
	"bytes"
	"crypto/aes"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha512"
	"io"
)

func keyFromPassword(password string, salt []byte, iter int) ([]byte, error) {
	key, err := pbkdf2.Key(sha512.New, password, salt, iter, keySize)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func appendPadding(data []byte, blockSize int) []byte {
	padSize := blockSize - len(data)%blockSize
	return append(data, bytes.Repeat([]byte{byte(padSize)}, padSize)...)
}

func removePadding(data []byte) ([]byte, error) {
	if len(data) < aes.BlockSize {
		return nil, ErrInvalidPadding
	}

	padSize := data[len(data)-1]
	if padSize > aes.BlockSize {
		return nil, ErrInvalidPadding
	}

	rem := len(data) - int(padSize)

	noPad, pad := data[:rem], data[rem:]

	for _, pad := range pad {
		if pad != padSize {
			return nil, ErrInvalidPadding
		}
	}

	return noPad, nil
}

func generateRand(size int) []byte {
	return fillRand(make([]byte, size))
}

func fillRand(data []byte) []byte {
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		panic(err)
	}
	return data
}
