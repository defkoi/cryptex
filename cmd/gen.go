package cmd

import (
	"fmt"
	"os"

	"github.com/defkoi/passgen"
)

func Gen() {
	r, err := openAndSaveBackup(cryptexFile)
	if err != nil {
		fatal(err)
	}

	value, err := passgen.GeneratePassword(genLen)
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
