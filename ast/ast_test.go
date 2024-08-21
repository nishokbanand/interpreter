package ast

import (
	"testing"

	"github.com/nishokbanand/interpreter/token"
)

func TestStringMethod(t *testing.T) {
	program := Program{
		Statements: []StatmentNode{
			&LetStatement{
				Token: token.Token{
					Type:    token.LET,
					Literal: "let",
				},
				Name: &Identifier{
					Token: token.Token{
						Type:    token.IDENT,
						Literal: "x",
					},
					Value: "x",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "y"},
					Value: "y",
				},
			},
		},
	}
	if program.String() != "let x = y;" {
		t.Errorf("program.String wrong got=%q", program.String())
	}
}
