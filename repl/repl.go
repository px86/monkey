package repl

import (
	"bufio"
	"fmt"
	"github.com/px86/monkey/evaluator"
	"github.com/px86/monkey/lexer"
	"github.com/px86/monkey/parser"
	"io"
	"os"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		prog := p.ParseProgram()
		if len(p.Errors) > 0 {
			for _, err := range p.Errors {
				fmt.Fprintf(os.Stderr, "%s\n", err)
			}
			return
		}
		result := evaluator.Eval(prog)
		if result != nil {
			io.WriteString(out, result.Inspect())
			io.WriteString(out, "\n")
		}
	}
}
