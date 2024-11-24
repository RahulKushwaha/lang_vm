package compiler

import (
	"github.com/stretchr/testify/assert"
	"lang_vm/code"
	"lang_vm/lexer"
	"lang_vm/parser"
	"testing"
)

func TestCompilerByteCode(t *testing.T) {
	tests := map[string]struct {
		code     string
		byteCode code.Instructions
	}{
		"simple_addition": {
			code: "2 + 5",
			byteCode: code.NewBuilder().
				Add(code.OpConstant, 2).
				Add(code.OpConstant, 5).
				Add(code.OpAdd).
				Build(),
		},
		"multiple_statements": {
			code: `2 + 5;
					5 - 5; 
					45 / 54; `,
			byteCode: code.NewBuilder().
				Add(code.OpConstant, 2).
				Add(code.OpConstant, 5).
				Add(code.OpAdd).
				Add(code.OpConstant, 5).
				Add(code.OpConstant, 5).
				Add(code.OpSub).
				Add(code.OpConstant, 45).
				Add(code.OpConstant, 54).
				Add(code.OpDiv).
				Build(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			l := lexer.New(tc.code)
			p := parser.New(l)
			program := p.ParseProgram()

			c := NewCompiler()
			err := c.Compile(program)
			assert.NoError(t, err)
			assert.Equal(t, tc.byteCode, c.ByteCode().Instructions)
		})
	}
}
