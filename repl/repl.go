package repl

import (
	"bufio"
	"fmt"
	"io"
	"bear/lexer"
	"bear/parser"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprintf(out, ">> ")
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

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

const ERROR_FACE = `
ʕ⊙ᴥ⊙ʔ
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, ERROR_FACE)
	io.WriteString(out, "Whoop! An error occurred")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}