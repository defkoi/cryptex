package cmd

import (
	"bufio"
	"bytes"
	"cryptex/internal/cryptex"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
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
	errPasswordsMismatch  = errors.New("passwords mismatch")
	errNoPasswordProvided = errors.New("no password provided")
)

var (
	cryptexFile   string
	inputKey      string
	inputString   string
	inputPassword string
	iter          uint
	genLen        int
	genCharSet    string
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
		&genCharSet,
		"cm",
		"",
		"generation charset modifier",
	)

	set.Parse(args)
}

func readPassword(confirm bool) (password string, err error) {
	fd := int(os.Stdin.Fd())

	state, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer func() {
		if err != nil {
			term.Restore(fd, state)
		} else {
			err = term.Restore(fd, state)
		}
	}()

	fmt.Print("password: ")
	defer clearLine()
	passBytes, err := term.ReadPassword(fd)
	if err != nil {
		return "", err
	}

	if len(passBytes) == 0 {
		return "", errNoPasswordProvided
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
			return "", errPasswordsMismatch
		}
	}

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
		p, err := readPassword(confirm)
		return p, err
	}
	return inputPassword, nil
}

func decode(r io.Reader, p string) (*cryptex.Cryptex, error) {
	c, err := cryptex.Decode(r, p)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func fatal(a any) {
	redColor.Println(a)
	os.Exit(0)
}

func validateKey(key string) error {
	if key == "" {
		return errors.New("empty key")
	}

	var prev byte
	for _, ch := range []byte(key) {
		if !('a' <= ch && ch <= 'z' ||
			ch == '.' && prev != '.') {
			return errors.New("invalid key format")
		}
		prev = ch
	}

	return nil
}

func openAndSaveBackup(filepath string) (io.Reader, error) {
	rFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer rFile.Close()

	data, err := io.ReadAll(rFile)
	if err != nil {
		return nil, err
	}

	wFile, err := os.Create(addSuffix(filepath, "backup"))
	if err != nil {
		return nil, err
	}
	defer wFile.Close()

	if _, err := wFile.Write(data); err != nil {
		return nil, err
	}

	cp := make([]byte, len(data))
	copy(cp, data)
	return bytes.NewReader(cp), nil
}

func addSuffix(p string, suffix string) string {
	dir := path.Dir(p)
	ext := path.Ext(p)
	name, _ := strings.CutSuffix(path.Base(p), ext)
	if name != "" && !strings.HasSuffix(name, "_") {
		name += "_"
	}
	return path.Join(dir, name+suffix+ext)
}

func createCryptex() (*cryptex.Cryptex, error) {
	c, err := cryptex.New(iter)
	if err != nil {
		return nil, err
	}

	return c, nil
}

func getPasswordWithRetry() (string, error) {
	for {
		password, err := getPassword(true)
		if err != nil {
			if errors.Is(err, errNoPasswordProvided) ||
				errors.Is(err, errPasswordsMismatch) {
				fmt.Printf("%v: retry: ", err)
				continue
			}
			return "", err
		}
		return password, nil
	}
}
