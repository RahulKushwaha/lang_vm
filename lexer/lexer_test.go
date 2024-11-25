package lexer

import (
	"github.com/google/go-cmp/cmp"
	"lang_vm/token"
	"testing"
)

func TestLexer(t *testing.T) {
	type input struct {
		input          string
		expectedTokens []token.Token
	}

	tests := map[string]input{
		"empty": {
			"",
			[]token.Token{},
		},
		"addition": {
			"2 + 5",
			[]token.Token{{token.Int, "2"}, {token.Plus, "+"},
				{token.Int, "5"}},
		},
		"minus": {
			"223 - 512312",
			[]token.Token{{token.Int, "223"}, {token.Minus, "-"},
				{token.Int, "512312"}},
		},
		"paren": {
			"(5 + 7)",
			[]token.Token{{token.LeftParen, "("}, {token.Int, "5"},
				{token.Plus, "+"}, {token.Int, "7"}, {token.RightParen, ")"}},
		},
		"let": {
			"let a = b;",
			[]token.Token{{token.Let, "let"}, {token.Identifier, "a"},
				{token.Assign, "="}, {token.Identifier, "b"}},
		},
		"multilet": {
			"let a = b; let c = d",
			[]token.Token{{token.Let, "let"}, {token.Identifier, "a"},
				{token.Assign, "="}, {token.Identifier, "b"},
				{token.Let, "let"}, {token.Identifier, "c"},
				{token.Assign, "="}, {token.Identifier, "d"}},
		},
		"if": {
			`
				if (5 + 10) {
					4 + 5
				} else {
					99 + 99
				}`,
			[]token.Token{{token.If, "if"}, {token.LeftParen, "("},
				{token.Int, "5"}, {token.Plus, "+"}, {token.Int, "10"},
				{token.RightParen, ")"}, {token.LeftBrace, "{"},
				{token.Int, "4"}, {token.Plus, "+"}, {token.Int, "5"},
				{token.RightBrace, "}"}, {token.Else, "else"}, {token.LeftBrace, "{"},
				{token.Int, "99"}, {token.Plus, "+"}, {token.Int, "99"},
				{token.RightBrace, "}"}},
		},
		"complex": {
			"2*(5+5*2)/3+(6/2+8)",
			[]token.Token{{token.Int, "2"}, {token.Asterisk, "*"},
				{token.LeftParen, "("}, {token.Int, "5"}, {token.Plus, "+"},
				{token.Int, "5"}, {token.Asterisk, "*"}, {token.Int, "2"},
				{token.RightParen, ")"}, {token.Slash, "/"}, {token.Int, "3"},
				{token.Plus, "+"},
				{token.LeftParen, "("}, {token.Int, "6"}, {token.Slash, "/"},
				{token.Int, "2"}, {token.Plus, "+"}, {token.Int, "8"},
				{token.RightParen, ")"}},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			l := New(test.input)
			got := l.GetAllTokens()

			if diff := cmp.Diff(got, test.expectedTokens); diff != "" {
				t.Errorf("expected: %v, got: %v", test.expectedTokens, got)
			}
		})
	}
}
