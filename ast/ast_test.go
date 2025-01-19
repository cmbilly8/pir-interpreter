package ast

import (
	"pir-interpreter/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&YarStatement{
				Token: token.Token{Type: token.YAR, Literal: "yar"},
				Name: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "variableName"},
					Value: "variableName",
				},
				Value: &Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVariable"},
					Value: "anotherVariable",
				},
			},
		},
	}
	if program.String() != "yar variableName be anotherVariable." {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
