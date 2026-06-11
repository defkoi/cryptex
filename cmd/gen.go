package cmd

import (
	"fmt"
	"os"

	"github.com/defkoi/passgen"
)

func Gen() {
	rFile, err := os.Open(cryptexFile)
	if err != nil {
		fatal(err)
	}
	defer rFile.Close()

	value, err := passgen.GeneratePassword(genLen)
	if err != nil {
		fatal(err)
	}

	password, err := getPassword(false)
	if err != nil {
		fatal(err)
	}

	c, err := decode(rFile, password)
	if err != nil {
		fatal(err)
	}

	key := getKey()

	c.Store(key, value)

	wFile, err := os.Create(cryptexFile)
	if err != nil {
		fatal(err)
	}
	defer wFile.Close()

	if err := c.Encode(wFile, password); err != nil {
		fatal(err)
	}

	fmt.Printf("generated: %s\n", yellowColor.Sprint(value))
}
