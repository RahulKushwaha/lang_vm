package lexer

import "lang_vm/token"

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
	l := &Lexer{input: input}
	return l
}

func (l *Lexer) NextToken() token.Token {
	return token.Token{}
}
