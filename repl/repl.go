package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/nishokbanand/interpreter/evaluate"
	"github.com/nishokbanand/interpreter/lexer"
	"github.com/nishokbanand/interpreter/parser"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(">>")
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		input := scanner.Text()
		lexer := lexer.New(input)
		parser := parser.New(lexer)
		program := parser.ParseProgram()
		if len(parser.Errors()) != 0 {
			printParseErrors(out, parser.Errors())
			continue
		}
		evaluated := evaluate.Eval(program)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
