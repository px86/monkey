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
		{token.KW_LET, "let"},
		{token.IDENTIFIER, "str"},
		{token.EQUAL, nil},
		{token.STRING_LITERAL, "This is a string."},
		{token.SEMI_COLON, nil},
		{token.KW_LET, "let"},
		{token.IDENTIFIER, "x"},
		{token.EQUAL, nil},
		{token.INTEGER, int64(7)},
		{token.SEMI_COLON, nil},
		{token.KW_LET, "let"},
		{token.IDENTIFIER, "y"},
		{token.EQUAL, nil},
		{token.INTEGER, int64(13)},
		{token.SEMI_COLON, nil},
		{token.KW_LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.EQUAL, nil},
		{token.KW_FUNCTION, "fn"},
		{token.LEFT_PAREN, nil},
		{token.IDENTIFIER, "a"},
		{token.COMMA, nil},
		{token.IDENTIFIER, "b"},
		{token.RIGHT_PAREN, nil},
		{token.LEFT_BRACE, nil},
		{token.KW_RETURN, "return"},
		{token.IDENTIFIER, "a"},
		{token.PLUS, nil},
		{token.IDENTIFIER, "b"},
		{token.SEMI_COLON, nil},
		{token.RIGHT_BRACE, nil},
		{token.SEMI_COLON, nil},
		{token.KW_LET, "let"},
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
		if tok.Type == token.INTEGER {
			val, ok := tok.Value.(int64)
			if !ok {
				t.Fatalf("tests[%d] - value not int64. got=%T",
					i, tok.Value)
			}
			expectedVal, _ := tt.expectedLiteral.(int64)
			if val != expectedVal {
				t.Fatalf("tests[%d] - literal value wrong. expected=%v, got=%v",
					i, expectedVal, val)
			}
		}
	}
}

func TestStringLiterals(t *testing.T) {
	input := `"This string has\nnewlines. And \t tabs."`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral any
	}{
		{token.STRING_LITERAL, "This string has\nnewlines. And \t tabs."},
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

	stream := []struct {
		expectedType    token.TokenType
		expectedLiteral any
		expectedLine    int
		expectedColumn  int
	}{
		{token.KW_LET, "let", 1, 0},
		{token.IDENTIFIER, "x", 1, 4},
		{token.EQUAL, nil, 1, 6},
		{token.INTEGER, int64(1), 1, 8},
		{token.SEMI_COLON, nil, 1, 9},
		{token.KW_LET, "let", 2, 0},
		{token.IDENTIFIER, "foo", 2, 4},
		{token.EQUAL, nil, 2, 8},
		{token.STRING_LITERAL, "bar", 2, 10},
		{token.SEMI_COLON, nil, 2, 15},
	}

	lex := New(input)
	for i, tt := range stream {
		tok := lex.NextToken()
		if tok.Type != tt.expectedType {
			t.Errorf("token[%d] - TokenType did not match. expected=%q, got=%q",
				i, token.AsString(tt.expectedType), token.AsString(tok.Type))
		}
		if tok.Type == token.INTEGER {
			val, ok := tok.Value.(int64)
			if !ok {
				t.Fatalf("token[%d] - INTEGER value not of type int64. got=%T",
					i, tok.Value)
			}
			expectedVal, _ := tt.expectedLiteral.(int64)
			if val != expectedVal {
				t.Fatalf("token[%d] - literal value did not match. expected=%d, got=%d",
					i, expectedVal, val)
			}
		}
		if tok.Line != tt.expectedLine {
			t.Errorf("token[%d] - line number did not match. expected=%d, got=%d",
				i, tt.expectedLine, tok.Line)
		}
		if tok.Column != tt.expectedColumn {
			t.Errorf("token[%d] - column number did not match. expected=%d, got=%d",
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
		{token.EXCLAMATION_EQUAL},
		{token.GREATER_THAN},
		{token.GREATER_THAN_EQUAL},
		{token.LESSER_THAN},
		{token.LESSER_THAN_EQUAL},
		{token.AMPERSAND},
		{token.AMPERSAND_AMPERSAND},
		{token.PIPE},
		{token.PIPE_PIPE},
		{token.CARET},
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
