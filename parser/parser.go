package parser

import (
	"fmt"
	"lang_vm/ast"
	"lang_vm/lexer"
	"lang_vm/token"
	"strconv"
)

const (
	_ int = iota
	Lowest
	Equals        // =
	LessOrGreater // < or >
	Sum           // +
	Product       // *
	Prefix        // -X or !X
	Call          // myFunction(X)
	Index         // array[index]
)

var precedences = map[token.TokenType]int{
	token.Equal:       Equals,
	token.NotEqual:    Equals,
	token.LessThan:    LessOrGreater,
	token.GreaterThan: LessOrGreater,
	token.Plus:        Sum,
	token.Minus:       Sum,
	token.Slash:       Product,
	token.Asterisk:    Product,
	token.LeftParen:   Call,
	token.LeftBracket: Index,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l      lexer.ILexer
	errors []string

	currentToken token.Token
	peekToken    token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func New(l lexer.ILexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.prefixParseFns[token.Int] = p.parseIntegerLiteral

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.infixParseFns[token.Plus] = p.parseInfixExpression
	p.infixParseFns[token.Minus] = p.parseInfixExpression
	p.infixParseFns[token.Slash] = p.parseInfixExpression
	p.infixParseFns[token.Asterisk] = p.parseInfixExpression

	// Read two tokens, so currentToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Precedence

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return Lowest
}

func (p *Parser) currentPrecedence() int {
	if p, ok := precedences[p.currentToken.Type]; ok {
		return p
	}

	return Lowest
}

func (p *Parser) ParseProgram() *ast.Program {
	program := ast.Program{
		Statements: []ast.Statement{p.parseExpressionStatement()},
	}

	return &program
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	return &ast.ExpressionStatement{
		Token:      p.currentToken,
		Expression: p.parseExpression(Lowest),
	}
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.currentToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	leftExp := prefix()

	for !p.peekTokenIs(token.Semicolon) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.currentToken}

	value, err := strconv.ParseInt(p.currentToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.BinaryExpression{
		Token:    p.currentToken,
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.currentPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)
	return expression
}
