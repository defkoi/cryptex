package cmd

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	defaultKey = ""
)

var (
	cryptexFile string
)

func ParseFlags(args []string) {
	set := flag.NewFlagSet("main", flag.ContinueOnError)

	set.StringVar(&cryptexFile, "f", ".cryptex", "cryptex file")

	set.Parse(args)
}

func readPassword(confirm bool) (password string, err error) {
	fd := int(os.Stdin.Fd())

	state, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer func() { err = term.Restore(fd, state) }()

	fmt.Print("password: ")
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

		fmt.Print("confirm: ")
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

func readLine(prompt string) string {
	fmt.Printf("%s: ", prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
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
