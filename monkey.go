package main

import (
	"github.com/px86/monkey/repl"
	"os"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
