package main

import (
	"fmt"
	"os"

	"github.com/nishokbanand/interpreter/repl"
)

func main() {
	fmt.Println("REPL starting")
	repl.Start(os.Stdin, os.Stdout)
}
