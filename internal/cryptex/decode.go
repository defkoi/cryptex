package cryptex

import (
	"crypto/aes"
	"crypto/cipher"
	"errors"
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

	mode := decodedData.mode
	iter := decodedData.iter
	salt := decodedData.salt
	iv := decodedData.iv
	data := decodedData.data

	key, err := keyFromPassword(password, salt, iter)
	if err != nil {
		return nil, err
	}

	data, err = decrypt(data, key, iv, mode)
	if err != nil {
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
		mode: mode,
	}, nil
}

func decrypt(data []byte, key []byte, iv []byte, mode Mode) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	switch mode {
	case ModeCBC:
		ret := make([]byte, len(data))
		cipher.NewCBCDecrypter(block, iv).CryptBlocks(ret, data)
		return ret, nil
	case ModeGCM:
		gcm, err := cipher.NewGCM(block)
		if err != nil {
			return nil, err
		}
		nonce := iv[:gcm.NonceSize()]
		return gcm.Open(nil, nonce, data, nil)
	default:
		return nil, errors.New("unsupported mode")
	}
}
