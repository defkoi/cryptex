package cmd

import (
	"os"
)

func Del() {
	r, err := openAndSaveBackup(cryptexFile)
	if err != nil {
		fatal(err)
	}

	password, err := getPassword(false)
	if err != nil {
		fatal(err)
	}

	c, err := decode(r, password)
	if err != nil {
		fatal(err)
	}

	key, err := getKey()
	if err != nil {
		fatal(err)
	}

	c.Store(key, "")

	wFile, err := os.Create(cryptexFile)
	if err != nil {
		fatal(err)
	}
	defer wFile.Close()

	if err := c.Encode(wFile, password); err != nil {
		fatal(err)
	}
}
