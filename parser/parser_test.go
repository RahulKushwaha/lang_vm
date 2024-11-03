package parser

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"lang_vm/lexer"
	"lang_vm/lexer/mocks"
	"lang_vm/token"
	"testing"
)

func TestExpressionParsing(t *testing.T) {
	type testCase struct {
		l           lexer.ILexer
		expectedOut string
	}

	testCases := map[string]testCase{}

	{
		l := &mocks.ILexer{}
		l.On("NextToken").Return(token.Token{Type: token.Int, Literal: "2"}).Once()
		l.On("NextToken").Return(token.Token{Type: token.Plus, Literal: "+"}).Once()
		l.On("NextToken").Return(token.Token{Type: token.Int, Literal: "5"}).Once()
		l.On("NextToken").Return(token.Token{Type: token.EOF, Literal: ""}).Once()

		testCases["simple_plus_expression"] = testCase{
			l:           l,
			expectedOut: "(2 + 5)",
		}

	}

	{
		l := &mocks.ILexer{}
		l.On("NextToken").Return(token.Token{Type: token.Int, Literal: "1"}).Once()
		l.On("NextToken").Return(token.Token{Type: token.Plus, Literal: "+"}).Once()
		l.On("NextToken").Return(token.Token{Type: token.Int, Literal: "2"}).Once()
		l.On("NextToken").Return(token.Token{Type: token.Plus, Literal: "+"}).Once()
		l.On("NextToken").Return(token.Token{Type: token.Int, Literal: "3"}).Once()
		l.On("NextToken").Return(token.Token{Type: token.Plus, Literal: "+"}).Once()
		l.On("NextToken").Return(token.Token{Type: token.Int, Literal: "4"}).Once()
		l.On("NextToken").Return(token.Token{Type: token.Plus, Literal: "+"}).Once()
		l.On("NextToken").Return(token.Token{Type: token.Int, Literal: "5"}).Once()
		l.On("NextToken").Return(token.Token{Type: token.EOF, Literal: ""}).Once()

		testCases["multiple_plus_expression"] = testCase{
			l:           l,
			expectedOut: "((((1 + 2) + 3) + 4) + 5)",
		}

	}

	for name, tc := range testCases {
		parser := New(tc.l)

		program := parser.ParseProgram()
		assert.Equal(t, tc.expectedOut, fmt.Sprintf("%v", program), name)
	}
}
