package cmd

import (
	"cryptex/internal/cryptex"
	"errors"
	"fmt"
	"log"
	"os"
)

func Decrypt() {
	password, err := readPassword(false)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(cryptexFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	c, err := cryptex.Decode(file, password)
	if err != nil {
		if errors.Is(err, cryptex.ErrInvalidPadding) {
			log.Fatal("invalid password")
		} else {
			log.Fatal(err)
		}
	}

	if c.Has(defaultKey) {
		fmt.Fprintf(os.Stdout, "%s\n", c.Load(defaultKey))
	} else {
		fmt.Fprintf(os.Stderr, "empty\n")
	}
}
