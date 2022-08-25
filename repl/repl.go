package repl

import (
	"bufio"
	"fmt"
	"io"
	"bear/lexer"
	"bear/parser"
	"bear/evaluator"
	"bear/object"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Fprintf(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		lex := lexer.New(line)
		par := parser.New(lex)

		program := par.ParseProgram()
		if len(par.Errors()) != 0 {
			printParserErrors(out, par.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

const ERROR_FACE = `
ʕ⊙ᴥ⊙ʔ
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, ERROR_FACE)
	io.WriteString(out, "Whoops! An error occurred")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "ℹ️  ")
		io.WriteString(out, msg+"\n\n")
	}
}