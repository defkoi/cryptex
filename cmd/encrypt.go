package cmd

import (
	"cryptex/internal/cryptex"
	"log"
	"os"
	"strings"
)

func Encrypt() {
	rFile, err := os.Open(cryptexFile)
	if err != nil {
		if os.IsNotExist(err) {
			encryptNew()
			return
		}
		log.Fatal(err)
	}
	defer rFile.Close()

	password, err := readPassword(false)
	if err != nil {
		log.Fatal(err)
	}

	c, err := cryptex.Decode(rFile, password)
	if err != nil {
		log.Fatal(err)
	}

	key := strings.TrimSpace(readLine("key"))

	value := strings.TrimSpace(readUntilEnd("string"))

	c.Store(key, value)

	wFile, err := os.Create(cryptexFile)
	if err != nil {
		log.Fatal(err)
	}
	defer wFile.Close()

	if err := c.Encode(wFile, password); err != nil {
		log.Fatal(err)
	}
}

func encryptNew() {
	c, err := cryptex.New(iter)
	if err != nil {
		log.Fatal(err)
	}

	password, err := readPassword(true)
	if err != nil {
		log.Fatal(err)
	}

	key := strings.TrimSpace(readLine("key"))

	value := strings.TrimSpace(readUntilEnd("string"))

	c.Store(key, value)

	file, err := os.Create(cryptexFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	if err := c.Encode(file, password); err != nil {
		log.Fatal(err)
	}
}
