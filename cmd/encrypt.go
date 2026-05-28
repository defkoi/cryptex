package cmd

import (
	"cryptex/internal/cryptex"
	"log"
	"os"
	"strings"
)

func Encrypt() {
	c, err := cryptex.New(iter)
	if err != nil {
		log.Fatal(err)
	}

	value := strings.TrimSpace(readUntilEnd("string"))

	c.Store(defaultKey, value)

	file, err := os.Create(cryptexFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	password, err := readPassword(true)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Encode(file, password); err != nil {
		log.Fatal(err)
	}
}
