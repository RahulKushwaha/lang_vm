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
			got := l.getAllTokens()

			if diff := cmp.Diff(got, test.expectedTokens); diff != "" {
				t.Errorf("expected: %v, got: %v", test.expectedTokens, got)
			}
		})
	}
}
