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
		command(os.Args[2:])
	} else {
		fmt.Printf("unknown command '%s'\n\n", os.Args[1])
		showCommands()
	}
}

var commands = map[string]func([]string){
	"create":   cmd.Create,
	"set":      cmd.Set,
	"generate": cmd.Generate,
	"get":      cmd.Get,
	"keys":     cmd.Keys,
	"delete":   cmd.Delete,
	"change":   cmd.Change,
}

func showUsage() {
	fmt.Println("cryptex [command] [...flags]")
	fmt.Println()
	showCommands()
}

func showCommands() {
	fmt.Println("commands:")

	fmt.Println("\tcreate")
	fmt.Println("\tset")
	fmt.Println("\tgenerate")
	fmt.Println("\tget")
	fmt.Println("\tkeys")
	fmt.Println("\tdelete")
	fmt.Println("\tchange [password|iterations]")
}
