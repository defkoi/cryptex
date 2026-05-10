package cmd

import (
	"cryptex/internal/cryptex"
	"log"
	"os"
	"strings"
)

func Encrypt() {
	value := strings.TrimSpace(readLine("string"))

	password, err := readPassword(true)
	if err != nil {
		log.Fatal(err)
	}

	c := cryptex.New()
	c.Store(defaultKey, value)

	file, err := os.Create(cryptexFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := c.Encode(file, password); err != nil {
		log.Fatal(err)
	}
}
