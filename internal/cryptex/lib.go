package cryptex

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha256"
	"encoding/json"
	"io"
)

type Mode uint8

const (
	ModeCBC Mode = iota
	ModeGCM
)

func init() {
	block, err := aes.NewCipher(make([]byte, keySize))
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

	if gcm.NonceSize() > ivSize {
		panic("gcm: nonce > iv")
	}
}

func encodeMap(m map[string]string) []byte {
	buf := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buf).Encode(m); err != nil {
		panic(err)
	}
	return buf.Bytes()
}

func decodeMap(data []byte) (map[string]string, error) {
	var m map[string]string
	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&m); err != nil {
		return m, err
	}
	return m, nil
}

func keyFromPassword(
	password string,
	salt []byte,
	iter uint32,
) ([]byte, error) {
	key, err := pbkdf2.Key(sha256.New, password, salt, int(iter), keySize)
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
