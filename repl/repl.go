package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/waridh/go-monkey-interpreter/evaluator"
	"github.com/waridh/go-monkey-interpreter/lexer"
	"github.com/waridh/go-monkey-interpreter/object"
	"github.com/waridh/go-monkey-interpreter/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Print(PROMPT)

		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			printParseError(out, p.Errors())
			continue
		}

		env := object.NewEnvironment()
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}

	}
}

func printParseError(out io.Writer, errors []string) {
	io.WriteString(out, "Ran into some parser errors\nparser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
