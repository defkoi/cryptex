package main

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/pbkdf2"
	"crypto/rand"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"golang.org/x/term"
)

const (
	keySize  = 0x20
	saltSize = 0x10

	defaultFile = ".cryptex"
	defaultIter = 0x100000
)

type cryptex struct {
	String string `json:"string"`
	IV     string `json:"initial_vector,omitempty"`

	Salt string `json:"salt,omitempty"`
	Iter int    `json:"iterations,omitempty"`

	Password string `json:"password,omitempty"` // optional
}

func (c *cryptex) encrypt() {
	c.beforeEncrypt()

	salt := generateRand(saltSize)

	block, err := aes.NewCipher(keyFromPassword(c.Password, salt, c.Iter))
	if err != nil {
		log.Fatal(err)
	}

	iv := generateRand(aes.BlockSize)

	blockMode := cipher.NewCBCEncrypter(block, iv)

	text := appendPadding([]byte(c.String), aes.BlockSize)

	blockMode.CryptBlocks(text, text)

	c.String = hex.EncodeToString(text)
	c.IV = hex.EncodeToString(iv)
	c.Salt = hex.EncodeToString(salt)
	c.Password = ""
}

func (c *cryptex) beforeEncrypt() {
	if len(c.String) < 1 {
		log.Fatal("empty 'string'")
	}

	if c.Iter < 1 {
		c.Iter = defaultIter
	}

	if c.Password == "" {
		c.Password = readPassword()
	}
}

func (c *cryptex) decrypt() {
	c.beforeDecrypt()

	salt, err := hex.DecodeString(c.Salt)
	if err != nil {
		log.Fatal(err)
	}

	block, err := aes.NewCipher(keyFromPassword(c.Password, salt, c.Iter))
	if err != nil {
		log.Fatal(err)
	}

	iv, err := hex.DecodeString(c.IV)
	if err != nil {
		log.Fatal(err)
	}

	blockMode := cipher.NewCBCDecrypter(block, iv)

	text, err := hex.DecodeString(c.String)
	if err != nil {
		log.Fatal(err)
	}

	blockMode.CryptBlocks(text, text)

	text, err = removePadding(text)
	if err != nil {
		log.Fatal(err)
	}

	c.String = string(text)
	c.IV = ""
	c.Salt = ""
	c.Iter = 0
	c.Password = ""
}

func (c *cryptex) beforeDecrypt() {
	if len(c.String) < 1 {
		log.Fatal("empty 'string'")
	}

	if len(c.String)%aes.BlockSize != 0 {
		log.Fatal("invalid 'string'")
	}

	if len(c.IV) != aes.BlockSize*2 {
		log.Fatal("invalid 'initial_vector'")
	}

	if len(c.Salt) != saltSize*2 {
		log.Fatal("invalid 'salt'")
	}

	if c.Iter < 1 {
		c.Iter = defaultIter
	}

	if c.Password == "" {
		c.Password = readPassword()
	}
}

func keyFromPassword(password string, salt []byte, iter int) []byte {
	key, err := pbkdf2.Key(sha512.New, password, salt, iter, keySize)
	if err != nil {
		log.Fatal(err)
	}
	return key
}

func generateRand(size int) []byte {
	return fillRand(make([]byte, size))
}

func appendPadding(data []byte, blockSize int) []byte {
	padSize := blockSize - len(data)%blockSize
	return append(data, bytes.Repeat([]byte{byte(padSize)}, padSize)...)
}

func removePadding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, errors.New("invalid padding")
	}
	remSize := int(data[len(data)-1])
	if remSize > len(data) {
		return nil, errors.New("invalid padding")
	}
	return data[:len(data)-remSize], nil
}

func fillRand(data []byte) []byte {
	if _, err := io.ReadFull(rand.Reader, data); err != nil {
		log.Fatal(err)
	}
	return data
}

type flags struct {
	file string
}

func parseFlags() flags {
	file := flag.String("f", defaultFile, "")
	flag.Parse()

	return flags{file: *file}
}

func main() {
	if len(os.Args) == 1 {
		usage()
		return
	}

	flags := parseFlags()

	switch os.Args[1] {
	case "encrypt":
		cmdWrapper(flags, cmdEncrypt)
	case "decrypt":
		cmdWrapper(flags, cmdDecrypt)
	default:
		fmt.Printf("unknown command '%s'\n", os.Args[1])
		commands()
	}
}

func usage() {
	fmt.Println("cryptex [command] [...flags]")
	fmt.Println()

	commands()
	fmt.Println()

	fmt.Println("{")
	fmt.Println("\tstring: string")
	fmt.Println("\tinitial_vector: string")
	fmt.Println("\tsalt: string")
	fmt.Println("\titerations: integer")
	fmt.Println("}")
}

func commands() {
	fmt.Println("commands:")
	fmt.Println("\tencrypt")
	fmt.Println("\tdecrypt")
}

func readPassword() string {
	fmt.Print("password: ")
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	if len(password) == 0 {
		log.Fatal("empty password")
	}
	fmt.Print( /* clear line */ "\r\x1b[K")
	return string(password)
}

func cmdWrapper(flags flags, f func(*cryptex)) {
	c := cryptex{}

	file, err := os.OpenFile(flags.file, os.O_RDWR, 0)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&c); err != nil {
		log.Fatal(err)
	}

	f(&c)

	rewrite(file)

	if err := json.NewEncoder(file).Encode(c); err != nil {
		log.Fatal(err)
	}
}

func rewrite(file *os.File) {
	if err := file.Truncate(0); err != nil {
		log.Fatal(err)
	}
	if _, err := file.Seek(0, 0); err != nil {
		log.Fatal(err)
	}
}

func cmdEncrypt(c *cryptex) {
	c.encrypt()
}

func cmdDecrypt(c *cryptex) {
	c.decrypt()
}
