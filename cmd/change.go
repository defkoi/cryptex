package cmd

import (
	"fmt"
	"os"
)

func Change(args []string) {
	if len(args) == 0 {
		fatal("specify what to change")
	}

	parseFlags(args[1:])

	switch args[0] {
	case "password":
		changePassword()
	case "iterations":
		changeIterations()
	default:
		fatalf("unknown subcommand '%s'", args[0])
	}
}

func changePassword() {
	r, err := openAndSaveBackup(cryptexFile)
	if err != nil {
		fatal(err)
	}

	oldPassword, err := getPassword(false)
	if err != nil {
		fatal(err)
	}

	c, err := decode(r, oldPassword)
	if err != nil {
		fatal(err)
	}

	file, err := os.Create(cryptexFile)
	if err != nil {
		fatal(err)
	}
	defer file.Close()

	newPassword, err := getNewPasswordWithRetry()
	if err != nil {
		fatal(err)
	}

	if err := c.Encode(file, newPassword); err != nil {
		fatal(err)
	}
}

func changeIterations() {
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

	file, err := os.Create(cryptexFile)
	if err != nil {
		fatal(err)
	}
	defer file.Close()

	old, new := yellowColor.Sprint(c.Iter()), yellowColor.Sprint(iter)

	if err := c.SetIter(iter); err != nil {
		fatal(err)
	}

	if err := c.Encode(file, password); err != nil {
		fatal(err)
	}

	fmt.Printf("changed: %s -> %s\n", old, new)
}
