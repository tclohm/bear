package main

import (
	"fmt"
	"os"
	"os/user"
	"bear/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Welcome to the Bear Programming language, %s", user.Username)
	fmt.Printf("\nType in commands\n")
	repl.Start(os.Stdin, os.Stdout)
}