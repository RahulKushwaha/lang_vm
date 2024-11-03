package lexer

import (
	"bytes"
	"fmt"
	"lang_vm/token"
)

//go:generate mockery --name ILexer
type ILexer interface {
	NextToken() token.Token
}

type Lexer struct {
	input        string
	position     int
	nextPosition int
	ch           byte
}

func New(input string) *Lexer {
	if len(input) > 0 {
		l := &Lexer{input: input, position: 0, nextPosition: 0, ch: 0}
		return l
	}

	return &Lexer{input: input, position: 0, nextPosition: 0, ch: 0}
}

func (l *Lexer) getAllTokens() []token.Token {
	tokens := make([]token.Token, 0)
	for {
		t := l.NextToken()
		if t.Type == token.EOF {
			break
		}

		fmt.Println(t.Literal)
		tokens = append(tokens, t)
	}

	return tokens
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	if l.position >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.position]
	}

	skipWhiteSpaces(l)

	switch {
	// Digit
	case '0' <= l.ch && l.ch <= '9':
		tok = token.Token{Type: token.Int, Literal: getNumber(l)}

	case l.ch == '+':
		tok = token.Token{Type: token.Plus, Literal: string(l.ch)}
		l.position++

	case l.ch == '-':
		tok = token.Token{Type: token.Minus, Literal: string(l.ch)}
		l.position++

	case l.ch == '*':
		tok = token.Token{Type: token.Asterisk, Literal: string(l.ch)}
		l.position++

	case l.ch == '/':
		tok = token.Token{Type: token.Slash, Literal: string(l.ch)}
		l.position++

	case l.ch == '(':
		tok = token.Token{Type: token.LeftParen, Literal: string(l.ch)}
		l.position++

	case l.ch == ')':
		tok = token.Token{Type: token.RightParen, Literal: string(l.ch)}
		l.position++

	case l.ch == 0:
		tok = token.Token{Type: token.EOF}
	}

	return tok
}

func getNumber(l *Lexer) string {
	var number bytes.Buffer
	for ; l.position < len(l.input) && isDigit(l.input[l.position]); l.position++ {
		number.WriteString(string(l.input[l.position]))
	}

	return number.String()
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func skipWhiteSpaces(l *Lexer) {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' {
		l.position++
		l.ch = l.input[l.position]
	}
}
