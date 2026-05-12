package cryptex

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"crypto/x509"
	"io"
)

/* << idea >>
 * cryptex asym create // create .key and .pub files
 * cryptex asym encrypt -k .pub
 * cryptex asym decrypt -k .key
 */

const (
	privateKeySize = 0x800

	label = "cryptex"
)

type Asym struct {
	data map[string]string
}

func NewAsym() *Asym {
	return &Asym{
		data: make(map[string]string),
	}
}

func (a *Asym) Encode(w io.Writer, publicKey *rsa.PublicKey) error {
	data, err := rsa.EncryptOAEP(
		sha512.New(),
		rand.Reader,
		publicKey,
		encodeMap(a.data),
		[]byte(label),
	)
	if err != nil {
		return err
	}

	if _, err := w.Write(data); err != nil {
		return err
	}

	return nil
}

func DecodeAsym(r io.Reader, privateKey *rsa.PrivateKey) (*Asym, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	decryptedData, err := rsa.DecryptOAEP(
		sha512.New(),
		rand.Reader,
		privateKey,
		data,
		[]byte(label),
	)
	if err != nil {
		return nil, err
	}

	dataMap, err := decodeMap(decryptedData)
	if err != nil {
		return nil, err
	}

	return &Asym{data: dataMap}, nil
}

func generateAsymKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, privateKeySize)
	if err != nil {
		panic(err)
	}
	return key
}

func keyToBytes(key *rsa.PrivateKey) (
	private []byte,
	public []byte,
) {
	private = x509.MarshalPKCS1PrivateKey(key)
	public = x509.MarshalPKCS1PublicKey(&key.PublicKey)
	return
}
