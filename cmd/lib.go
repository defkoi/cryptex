package cmd

import (
	"bufio"
	"bytes"
	"cryptex/internal/cryptex"
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/term"
)

const (
	defaultIter   = 600_000
	defaultGenLen = 8
)

var (
	grayColor   = color.New(color.FgWhite, color.Faint)
	yellowColor = color.New(color.FgYellow)
	redColor    = color.New(color.FgRed)
)

var (
	cryptexFile   string
	inputKey      string
	inputString   string
	inputPassword string
	iter          uint
	genLen        int
	mode          string
)

func ParseFlags(args []string) {
	set := flag.NewFlagSet("cryptex", flag.ExitOnError)

	set.StringVar(
		&cryptexFile,
		"f",
		".cryptex",
		"cryptex file",
	)

	set.StringVar(
		&inputKey,
		"k",
		"",
		"key",
	)

	set.StringVar(
		&inputString,
		"s",
		"",
		"string",
	)

	set.StringVar(
		&inputPassword,
		"p",
		"",
		"password",
	)

	set.UintVar(
		&iter,
		"i",
		defaultIter,
		fmt.Sprintf("iterations (max %d)", cryptex.MaxIter),
	)

	set.IntVar(
		&genLen,
		"l",
		defaultGenLen,
		"generation length",
	)

	set.StringVar(
		&mode,
		"m",
		"cbc",
		"mode",
	)

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

func getKey() (string, error) {
	var key string
	if inputKey == "" {
		key = strings.TrimSpace(readLine("key"))
	} else {
		key = inputKey
	}

	if err := validateKey(key); err != nil {
		return "", err
	}

	return key, nil
}

func getString() string {
	if inputString == "" {
		return strings.TrimSpace(readUntilEnd("string"))
	}
	return inputString
}

func getPassword(confirm bool) (string, error) {
	if inputPassword == "" {
		return readPassword(confirm)
	}
	return inputPassword, nil
}

func getMode() (cryptex.Mode, error) {
	switch strings.ToLower(mode) {
	case "cbc":
		return cryptex.ModeCBC, nil
	case "gcm":
		return cryptex.ModeGCM, nil
	default:
		return 0, errors.New("unsupported mode")
	}
}

func decode(f *os.File, p string) (*cryptex.Cryptex, error) {
	c, err := cryptex.Decode(f, p)
	if err != nil {
		if errors.Is(err, cryptex.ErrInvalidPadding) {
			return nil, errors.New("invalid password")
		}
		return nil, err
	}
	return c, nil
}

func fatal(a any) {
	redColor.Println(a)
	os.Exit(0)
}

func validateKey(key string) error {
	buf := []byte(key)
	for _, ch := range buf {
		if !('0' <= ch && ch <= '9' ||
			'a' <= ch && ch <= 'z' ||
			ch == '.') {
			return errors.New("invalid key format")
		}
	}
	return nil
}
