package cmd

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"

	"golang.org/x/term"
)

const (
	defaultFile = ".cryptex"

	defaultKey = ""
)

var (
	cryptexFile string
)

func ParseFlags() {
	f := flag.String("f", defaultFile, "cryptex file")
	flag.Parse()

	cryptexFile = *f
}

func readPassword(confirm bool) (password string, err error) {
	fd := int(os.Stdin.Fd())

	state, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer func() { err = term.Restore(fd, state) }()

	fmt.Print("password > ")
	passBytes, err := term.ReadPassword(fd)
	if err != nil {
		return "", err
	}

	if len(passBytes) == 0 {
		return "", errors.New("no password provided")
	}

	password = string(passBytes)

	if confirm {
		clearLine()

		fmt.Print("confirm password > ")
		confirmBytes, err := term.ReadPassword(fd)
		if err != nil {
			return "", err
		}

		if password != string(confirmBytes) {
			return "", errors.New("passwords don't match")
		}
	}

	clearLine()

	return password, nil
}

func readInput() string {
	fmt.Println("ctrl+z for input")

	in, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	return string(in)
}

func clearLine() {
	fmt.Print("\r\x1b[K")
}

func rewrite(file *os.File) error {
	if err := file.Truncate(0); err != nil {
		return err
	}
	if _, err := file.Seek(0, 0); err != nil {
		return err
	}
	return nil
}
