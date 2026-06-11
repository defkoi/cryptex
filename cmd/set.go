package cmd

import (
	"os"
)

func Set() {
	r, err := openAndSaveBackup(cryptexFile)
	if err != nil {
		if os.IsNotExist(err) {
			setNew()
			return
		}
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

func setNew() {
	c, err := createCryptex()
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

	password, err := getPasswordWithRetry()
	if err != nil {
		fatal(err)
	}

	if err := c.Encode(file, password); err != nil {
		fatal(err)
	}
}
