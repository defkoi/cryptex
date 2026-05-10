package main

import (
	"cryptex/cmd"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) == 1 {
		showUsage()
		return
	}

	if command, ok := commands[os.Args[1]]; ok {
		cmd.ParseFlags(os.Args[2:])
		command()
	} else {
		fmt.Printf("unknown command '%s'\n", os.Args[1])
		showCommands()
	}
}

var commands = map[string]func(){
	"encrypt": cmd.Encrypt,
	"decrypt": cmd.Decrypt,
}

func showUsage() {
	fmt.Println("cryptex [command] [...flags]")
	fmt.Println()

	showCommands()
}

func showCommands() {
	fmt.Println("commands:")

	fmt.Println("\tencrypt")
	fmt.Println("\tdecrypt")
}
