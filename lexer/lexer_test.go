package lexer

import (
	token "github.com/px86/monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `let str = "This is a string.";
let x = 7;
let y = 13;

let add = fn(a, b) {
   return a + b;
};

let result = add(x, y);`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral any
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "str"},
		{token.EQUAL, nil},
		{token.STRING, "This is a string."},
		{token.SEMI_COLON, nil},
		{token.LET, "let"},
		{token.IDENTIFIER, "x"},
		{token.EQUAL, nil},
		{token.INTEGER, 7},
		{token.SEMI_COLON, nil},
		{token.LET, "let"},
		{token.IDENTIFIER, "y"},
		{token.EQUAL, nil},
		{token.INTEGER, 13},
		{token.SEMI_COLON, nil},
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.EQUAL, nil},
		{token.FUNCTION, "fn"},
		{token.LEFT_PAREN, nil},
		{token.IDENTIFIER, "a"},
		{token.COMMA, nil},
		{token.IDENTIFIER, "b"},
		{token.RIGHT_PAREN, nil},
		{token.LEFT_BRACE, nil},
		{token.RETURN, "return"},
		{token.IDENTIFIER, "a"},
		{token.PLUS, nil},
		{token.IDENTIFIER, "b"},
		{token.SEMI_COLON, nil},
		{token.RIGHT_BRACE, nil},
		{token.SEMI_COLON, nil},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.EQUAL, nil},
		{token.IDENTIFIER, "add"},
		{token.LEFT_PAREN, nil},
		{token.IDENTIFIER, "x"},
		{token.COMMA, nil},
		{token.IDENTIFIER, "y"},
		{token.RIGHT_PAREN, nil},
		{token.SEMI_COLON, nil},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, int(tt.expectedType), int(tok.Type))
		}
		if tok.Value != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Value)
		}
	}
}

func TestStringLiterals(t *testing.T) {
	input := `"This stiring has\nnewlines. And \t tabs."`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral any
	}{
		{token.STRING, "This stiring has\nnewlines. And \t tabs."},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, int(tt.expectedType), int(tok.Type))
		}
		if tok.Value != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Value)
		}
	}
}

func TestLineAndColumn(t *testing.T) {
	input := `let x = 1;
let foo = "bar";`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral any
		expectedLine    int
		expectedColumn  int
	}{
		{token.LET, "let", 1, 0},
		{token.IDENTIFIER, "x", 1, 4},
		{token.EQUAL, nil, 1, 6},
		{token.INTEGER, 1, 1, 8},
		{token.SEMI_COLON, nil, 1, 9},
		{token.LET, "let", 2, 0},
		{token.IDENTIFIER, "foo", 2, 4},
		{token.EQUAL, nil, 2, 8},
		{token.STRING, "bar", 2, 10},
		{token.SEMI_COLON, nil, 2, 15},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
		if tok.Value != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Value)
		}
		if tok.Value != tt.expectedLiteral {
			t.Fatalf("tests[%d] - line wrong. expected=%q, got=%q",
				i, tt.expectedLine, tok.Line)
		}
		if tok.Value != tt.expectedLiteral {
			t.Fatalf("tests[%d] - column wrong. expected=%q, got=%q",
				i, tt.expectedColumn, tok.Column)
		}
	}
}

func TestOperators(t *testing.T) {
	input := "= == === ! != > >= < <= & && | || ^"

	tests := []struct {
		expectedType token.TokenType
	}{
		{token.EQUAL},
		{token.EQUAL_EQUAL},
		{token.EQUAL_EQUAL},
		{token.EQUAL},
		{token.EXCLAMATION},
		{token.NOT_EQUAL},
		{token.GREATER_THAN},
		{token.GREATER_THAN_EQUAL},
		{token.LESS_THAN},
		{token.LESS_THAN_EQUAL},
		{token.BITWISE_AND},
		{token.LOGICAL_AND},
		{token.BITWISE_OR},
		{token.LOGICAL_OR},
		{token.XOR},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}
	}
}
