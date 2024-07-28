package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/waridh/go-monkey-interpreter/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Hello %s! This is the Monkey Interpreter", user.Username)

	repl.Start(os.Stdin, os.Stdout)
}
