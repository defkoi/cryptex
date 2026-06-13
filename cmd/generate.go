package cmd

import (
	"fmt"
	"os"

	"github.com/defkoi/passgen"
)

func Generate(args []string) {
	parseFlags(args)

	r, err := openAndSaveBackup(cryptexFile)
	if err != nil {
		fatal(err)
	}

	if genCharSet != "" {
		if cs, err := passgen.CharSetFromModifier(genCharSet); err != nil {
			fatal(err)
		} else {
			passgen.SetCharSet(cs)
		}
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

	file, err := os.Create(cryptexFile)
	if err != nil {
		fatal(err)
	}
	defer file.Close()

	if err := c.Encode(file, password); err != nil {
		fatal(err)
	}

	fmt.Printf("generated: %s\n", yellowColor.Sprint(value))
}
