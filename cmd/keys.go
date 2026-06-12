package cmd

import (
	"os"
)

func Keys() {
	file, err := os.Open(cryptexFile)
	if err != nil {
		fatal(err)
	}
	defer file.Close()

	password, err := getPassword(false)
	if err != nil {
		fatal(err)
	}

	c, err := decode(file, password)
	if err != nil {
		fatal(err)
	}

	isEmpty := true
	for k := range c.Keys() {
		isEmpty = false
		yellowColor.Fprintln(os.Stdout, k)
	}

	if isEmpty {
		grayColor.Fprintln(os.Stdout, "empty")
	}
}
