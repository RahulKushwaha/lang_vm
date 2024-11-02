package ast

import (
	"bytes"
	"lang_vm/token"
)

type Expression interface {
	Node
}

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type BinaryExpression struct {
	Left     Expression
	Right    Expression
	Token    token.Token
	Operator string
}

func (b *BinaryExpression) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BinaryExpression) String() string {
	var out bytes.Buffer

	out.WriteString(token.LeftParen)
	out.WriteString(b.Left.String())
	out.WriteString(" " + b.Operator + " ")
	out.WriteString(b.Right.String())
	out.WriteString(token.RightParen)

	return out.String()
}
