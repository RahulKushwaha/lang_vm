package main

import (
	"fmt"
	"lang_vm/lexer"
	"lang_vm/parser"
)

func main() {
	input := "3 + 5"
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	for _, statement := range program.Statements {
		fmt.Printf("%v\n", statement.String())
	}
}
