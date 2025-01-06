package ast

import (
	"github.com/px86/monkey/token"
	"testing"
)

func TestString(t *testing.T) {
	program := &Program{
		Statements: []Statement{
			&LetStatement{
				Token: token.Token{Type: token.LET, Value: "let", Line: 1, Column: 0},
				Name: &Identifier{
					Token: token.Token{
						Type:   token.IDENTIFIER,
						Value:  "myVar",
						Line:   1,
						Column: 4,
					},
					Value: "myVar",
				},
				Value: &Identifier{
					Token: token.Token{
						Type:   token.IDENTIFIER,
						Value:  "anotherVar",
						Line:   1,
						Column: 12,
					},
					Value: "anotherVar",
				},
			},
		},
	}

	if program.String() != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", program.String())
	}
}
