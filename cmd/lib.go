package cmd

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	defaultKey = ""

	defaultIter = 600_000
)

var (
	cryptexFile string
	iter        uint
)

func ParseFlags(args []string) {
	set := flag.NewFlagSet("cryptex", flag.ExitOnError)

	set.StringVar(&cryptexFile, "f", ".cryptex", "cryptex file")
	set.UintVar(&iter, "iter", defaultIter, "iterations")

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

func readUntilEnd(prompt string) string {
	const (
		endMark    = "\\end"
		notEndMark = "\\\\end"
	)

	fmt.Println("enter using '\\end'")
	fmt.Printf("%s: ", prompt)

	scanner := bufio.NewScanner(os.Stdin)
	buf := []byte{}
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.HasSuffix(line, []byte(endMark)) {
			if bytes.HasSuffix(line, []byte(notEndMark)) {
				line = append(line[:len(line)-len(notEndMark)], endMark...)
			} else {
				buf = append(buf, line[:len(line)-len(endMark)]...)
				break
			}
		}
		buf = append(buf, append(line, '\n')...)
	}

	return string(buf)
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
