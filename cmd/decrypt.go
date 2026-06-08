package cmd

import (
	"cryptex/internal/cryptex"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func Decrypt() {
	file, err := os.Open(cryptexFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	password, err := readPassword(false)
	if err != nil {
		log.Fatal(err)
	}

	key := strings.TrimSpace(readLine("key"))

	c, err := cryptex.Decode(file, password)
	if err != nil {
		if errors.Is(err, cryptex.ErrInvalidPadding) {
			log.Fatal("invalid password")
		} else {
			log.Fatal(err)
		}
	}

	if c.Has(key) {
		fmt.Fprintln(os.Stdout, c.Load(key))
	} else {
		fmt.Fprintln(os.Stderr, "empty")
	}
}
