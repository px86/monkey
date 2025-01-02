package lexer

import (
	token "github.com/px86/monkey/token"
	"io"
	"os"
)

type Lexer struct {
	source string
	pos    int
	line   int
	column int
}

func New(path string) (*Lexer, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	bytes, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return &Lexer{
		source: string(bytes),
		line:   1,
	}, nil
}

func (lex *Lexer) peek() byte {
	if lex.pos >= len(lex.source) {
		return 0
	}
	return lex.source[lex.pos]
}

func (lex *Lexer) consume() byte {
	b := lex.source[lex.pos]
	lex.pos++
	if b == '\n' {
		lex.line++
		lex.column = 0
	} else {
		lex.column++
	}
	return b
}

func isWhitespace(c byte) bool {
	switch c {
	case ' ':
		return true
	case '\t':
		return true
	case '\n':
		return true
	case '\v':
		return true
	case '\r':
		return true
	default:
		return false
	}
}

func isAlpha(c byte) bool {
	if 'a' <= c && c <= 'z' {
		return true
	}
	if 'A' <= c && c <= 'Z' {
		return true
	}
	return false
}

func isDigit(c byte) bool {
	if '0' <= c && c <= '9' {
		return true
	}
	return false
}

func (lex *Lexer) NextToken() token.Token {
	var tok token.Token

	for {
		if lex.pos >= len(lex.source) {
			tok := token.Token{Type: token.EOF, Value: "", Line: lex.line, Column: lex.column}
			lex.pos = -1
			return tok
		}
		if lex.pos == -1 {
			panic("don't bother the lexer anymore...")
		}
		ch := lex.source[lex.pos]
		switch {
		case isWhitespace(ch):
			lex.consume()
			continue
		case ch == '/':
			tok = token.Token{Type: token.SLASH, Value: token.SLASH, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		case ch == '*':
			tok = token.Token{Type: token.ASTERISK_SIGN, Value: token.ASTERISK_SIGN, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		case ch == '!':
			tok = token.Token{Type: token.BANG_SIGN, Value: token.BANG_SIGN, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		case ch == '>':
			tok = token.Token{Type: token.GREATER_THAN, Value: token.GREATER_THAN, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		case ch == '<':
			tok = token.Token{Type: token.LESS_THAN, Value: token.LESS_THAN, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		case ch == ',':
			tok = token.Token{Type: token.COMMA, Value: token.COMMA, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		case ch == ';':
			tok = token.Token{Type: token.SEMI_COLON, Value: token.SEMI_COLON, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		case ch == '(':
			tok = token.Token{Type: token.LEFT_PAREN, Value: token.LEFT_PAREN, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		case ch == ')':
			tok = token.Token{Type: token.RIGHT_PAREN, Value: token.RIGHT_PAREN, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		case ch == '{':
			tok = token.Token{Type: token.LEFT_BRACE, Value: token.LEFT_BRACE, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		case ch == '}':
			tok = token.Token{Type: token.RIGHT_BRACE, Value: token.RIGHT_BRACE, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		case ch == '+':
			tok = token.Token{Type: token.PLUS_SIGN, Value: token.PLUS_SIGN, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		case ch == '-':
			tok = token.Token{Type: token.MINUS_SIGN, Value: token.MINUS_SIGN, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		case ch == '=':
			tok = token.Token{Type: token.EQUAL_SIGN, Value: token.EQUAL_SIGN, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok

		case ch == '"':
			start := lex.pos
			line := lex.line
			column := lex.column
			lex.consume() // cosume the " character
			for lex.peek() != 0 {
				c := lex.consume()
				if c == '"' {
					stringLiteral := lex.source[start+1 : lex.pos-1]
					tok = token.Token{Type: token.STRING_LITERAL, Value: stringLiteral, Line: line, Column: column}
					return tok
				}
			}
			val := lex.source[start-1:] // unterminated string literal
			tok = token.Token{Type: token.ILLEGAL, Value: val, Line: line, Column: column}
			return tok

		case isAlpha(ch):
			start := lex.pos
			column := lex.column
			for isAlpha(lex.peek()) || lex.peek() == '_' {
				lex.consume()
			}
			val := lex.source[start:lex.pos]
			switch val {
			case token.LET:
				tok = token.Token{Type: token.LET, Value: val, Line: lex.line, Column: column}
			case token.FUNCTION:
				tok = token.Token{Type: token.FUNCTION, Value: val, Line: lex.line, Column: column}
			case token.TRUE:
				tok = token.Token{Type: token.TRUE, Value: val, Line: lex.line, Column: column}
			case token.FALSE:
				tok = token.Token{Type: token.FALSE, Value: val, Line: lex.line, Column: column}
			case token.IF:
				tok = token.Token{Type: token.IF, Value: val, Line: lex.line, Column: column}
			case token.ELSE:
				tok = token.Token{Type: token.ELSE, Value: val, Line: lex.line, Column: column}
			case token.RETURN:
				tok = token.Token{Type: token.RETURN, Value: val, Line: lex.line, Column: column}
			default:
				tok = token.Token{Type: token.IDENTIFIER, Value: val, Line: lex.line, Column: column}
			}
			return tok

		case isDigit(ch):
			column := lex.column
			var value int
			for isDigit(lex.peek()) {
				value = value*10 + int(lex.consume()) - int('0')
			}
			tok = token.Token{Type: token.INTEGER, Value: value, Line: lex.line, Column: column}
			return tok

		default:
			tok = token.Token{Type: token.ILLEGAL, Value: token.ILLEGAL, Line: lex.line, Column: lex.column}
			lex.consume()
			return tok
		}
	}

}
