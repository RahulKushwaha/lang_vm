package ast

import (
	"bytes"
	"lang_vm/token"
)

type Expression interface {
	Node
	expressionNode()
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

type ExpressionStatement struct {
	Token      token.Token // the first token of the expression
	Expression Expression
}

func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (es *ExpressionStatement) statementNode() {
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

func (b *BinaryExpression) expressionNode() {}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) expressionNode() {
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("\n{\n")
	for _, s := range b.Statements {
		out.WriteString("\n\t")
		out.WriteString(s.String())
	}
	out.WriteString("\n}\n")

	return out.String()
}

func (b *BlockStatement) statementNode() {
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(i.Condition.String())
	out.WriteString(" ")
	out.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

func (i *IfExpression) expressionNode() {}
