package cmd

import (
	"os"
)

func Delete(args []string) {
	parseFlags(args)

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

	file, err := os.Create(cryptexFile)
	if err != nil {
		fatal(err)
	}
	defer file.Close()

	if err := c.Encode(file, password); err != nil {
		fatal(err)
	}
}
