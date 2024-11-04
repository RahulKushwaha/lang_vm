package lexer

import (
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
	case l.ch == '{':
		tok = token.Token{Type: token.LeftBrace, Literal: "{"}

	case l.ch == '}':
		tok = token.Token{Type: token.RightBrace, Literal: "}"}
	// Digit
	case '0' <= l.ch && l.ch <= '9':
		tok = token.Token{Type: token.Int, Literal: getNumber(l)}

	case l.ch == '+':
		tok = token.Token{Type: token.Plus, Literal: string(l.ch)}

	case l.ch == '-':
		tok = token.Token{Type: token.Minus, Literal: string(l.ch)}

	case l.ch == '*':
		tok = token.Token{Type: token.Asterisk, Literal: string(l.ch)}

	case l.ch == '/':
		tok = token.Token{Type: token.Slash, Literal: string(l.ch)}

	case l.ch == '(':
		tok = token.Token{Type: token.LeftParen, Literal: string(l.ch)}

	case l.ch == ')':
		tok = token.Token{Type: token.RightParen, Literal: string(l.ch)}

	case l.ch == '<':
		tok = token.Token{Type: token.LessThan, Literal: string(l.ch)}

	case isLetter(l.ch):
		identifier := l.readIdentifier()
		tok = token.Token{
			Type:    token.LookupIdentifierType(identifier),
			Literal: identifier,
		}

	case l.ch == 0:
		tok = token.Token{Type: token.EOF}
	}

	l.readChar()

	return tok
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func getNumber(l *Lexer) string {
	pos := l.position
	for isDigit(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func skipWhiteSpaces(l *Lexer) {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	pos := l.position
	for isLetter(l.ch) {
		l.readChar()
	}

	return l.input[pos:l.position]
}

func (l *Lexer) peekChar() byte {
	if l.nextPosition >= len(l.input) {
		return 0
	}

	return l.input[l.nextPosition]
}

func (l *Lexer) readChar() {
	l.ch = l.peekChar()
	l.position = l.nextPosition
	l.nextPosition++
}
