package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/px86/monkey/lexer"
	"github.com/px86/monkey/token"
)

func main() {
	if len(os.Args) < 2 {
		exe := filepath.Base(os.Args[0])
		fmt.Fprintf(os.Stderr, "Usage: %s FILE...\n", exe)
		os.Exit(1)
	}
	start := time.Now()

	lexer, err := lexer.New(os.Args[1])
	if err != nil {
		panic(err)
	}
	for tok := lexer.NextToken(); tok.Type != token.EOF; tok = lexer.NextToken() {
		fmt.Println(tok)
	}
	fmt.Printf("\ntook %v\n", time.Since(start))
}
