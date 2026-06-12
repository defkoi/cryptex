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
		fmt.Printf("unknown command '%s'\n\n", os.Args[1])
		showCommands()
	}
}

var commands = map[string]func(){
	"create": cmd.Create,
	"set":    cmd.Set,
	"gen":    cmd.Gen,
	"get":    cmd.Get,
	"keys":   cmd.Keys,
	"del":    cmd.Del,
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
	fmt.Println("\tgen")
	fmt.Println("\tget")
	fmt.Println("\tkeys")
	fmt.Println("\tdel")
}
