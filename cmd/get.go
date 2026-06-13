package cmd

import (
	"os"
)

func Get(args []string) {
	parseFlags(args)

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

	key, err := getKey()
	if err != nil {
		fatal(err)
	}

	if c.Has(key) {
		yellowColor.Println(c.Load(key))
	} else {
		grayColor.Println("empty")
	}
}
