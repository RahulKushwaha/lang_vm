package main

import (
	"fmt"
	"lang_vm/lexer"
	"lang_vm/parser"
)

func main() {
	input := `
		if (5 + 10) {
			4 + 5;
		} else {
			99 + 99;
			55 - 8
		}
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	for _, statement := range program.Statements {
		fmt.Printf("%v\n", statement.String())
	}
}
