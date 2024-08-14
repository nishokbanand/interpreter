package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/nishokbanand/interpreter/lexer"
	"github.com/nishokbanand/interpreter/token"
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
		for tok := lexer.NextToken(); tok.Type != token.EOF; tok = lexer.NextToken() {
			fmt.Printf("%#v\n", tok)
		}
	}
}
