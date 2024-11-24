package main

import (
	"fmt"
	"lang_vm/compiler"
	"lang_vm/lexer"
	"lang_vm/parser"
	"log"
)

func main() {
	input := `
		if (5 + 10) {
			4 + 5
		} else {
			99 + 99
			55 - 8
		}
`

	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	c := compiler.NewCompiler()
	err := c.Compile(program)
	if err != nil {
		log.Fatalf("error encountered: %v", err)
	}

	fmt.Println(c.ByteCode().Instructions.String())

	//for _, statement := range program.Statements {
	//	fmt.Printf("%v\n", statement.String())
	//}
}
