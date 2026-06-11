package cmd

import (
	"cryptex/internal/cryptex"
	"os"
)

func Set() {
	rFile, err := os.Open(cryptexFile)
	if err != nil {
		if os.IsNotExist(err) {
			encryptNew()
			return
		}
		fatal(err)
	}
	defer rFile.Close()

	password, err := getPassword(false)
	if err != nil {
		fatal(err)
	}

	c, err := decode(rFile, password)
	if err != nil {
		fatal(err)
	}

	key, err := getKey()
	if err != nil {
		fatal(err)
	}

	value := getString()

	c.Store(key, value)

	wFile, err := os.Create(cryptexFile)
	if err != nil {
		fatal(err)
	}
	defer wFile.Close()

	if err := c.Encode(wFile, password); err != nil {
		fatal(err)
	}
}

func encryptNew() {
	mode, err := getMode()
	if err != nil {
		fatal(err)
	}

	c, err := cryptex.New(iter, mode)
	if err != nil {
		fatal(err)
	}

	password, err := getPassword(true)
	if err != nil {
		fatal(err)
	}

	key, err := getKey()
	if err != nil {
		fatal(err)
	}

	value := getString()

	c.Store(key, value)

	file, err := os.Create(cryptexFile)
	if err != nil {
		fatal(err)
	}
	defer file.Close()

	if err := c.Encode(file, password); err != nil {
		fatal(err)
	}
}
