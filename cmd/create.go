package cmd

import (
	"os"
)

func Create(args []string) {
	parseFlags(args)

	c, err := createCryptex()
	if err != nil {
		fatal(err)
	}

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
