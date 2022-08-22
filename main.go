package main

import (
	"fmt"
	"os"
	"os/user"
	"bear/repl"
)

const BEAR_TEXT = `
██████╗ ███████╗ █████╗ ██████╗ 
██╔══██╗██╔════╝██╔══██╗██╔══██╗
██████╔╝█████╗  ███████║██████╔╝
██╔══██╗██╔══╝  ██╔══██║██╔══██╗
██████╔╝███████╗██║  ██║██║  ██║
╚═════╝ ╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝
`

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf(BEAR_TEXT)
	fmt.Printf("Welcome to the Bear Programming language, %s", user.Username)
	fmt.Printf("\nType in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}