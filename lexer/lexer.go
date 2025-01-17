package lexer

import (
	"fmt"
	"io"
	"os"

	"github.com/px86/monkey/token"
)

type Lexer struct {
	source string
	pos    int
	line   int
	column int
	eof    bool
}

func New(source string) *Lexer {
	return &Lexer{
		source: source,
		line:   1,
		column: 0,
	}
}

func FromFilePath(path string) (*Lexer, error) {
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
		lex.eof = true
		return 0
	}
	return lex.source[lex.pos]
}

func (lex *Lexer) consume() {
	if lex.pos >= len(lex.source) {
		lex.eof = true
		return
	}
	if lex.source[lex.pos] == '\n' {
		lex.line++
		lex.column = 0
	} else {
		lex.column++
	}
	lex.pos++
}

func (lex *Lexer) peekN(n int) string {
	if lex.pos+n < len(lex.source) {
		return lex.source[lex.pos : lex.pos+n]
	}
	return lex.source[lex.pos:]
}

func (lex *Lexer) singleCharToken(toktype token.TokenType) token.Token {
	tok := token.Token{Type: toktype, Value: nil, Line: lex.line, Column: lex.column}
	lex.consume()
	return tok
}

func (lex *Lexer) doubleCharToken(toktype token.TokenType) token.Token {
	tok := token.Token{Type: toktype, Value: nil, Line: lex.line, Column: lex.column}
	lex.consume()
	lex.consume()
	return tok
}

func (lex *Lexer) numberLiteralToken() token.Token {
	// TODO: check floating point numbers
	tok := token.Token{Line: lex.line, Column: lex.column}

	var value int64
	for d := lex.peek(); isDigit(d); d = lex.peek() {
		value = 10*value + int64(int(d)-int('0'))
		lex.consume()
	}
	tok.Type = token.INTEGER
	tok.Value = value
	return tok
}

func (lex *Lexer) stringLiteralToken() token.Token {

	if lex.peek() != '"' {
		panic("*Lexer.stringLiteralToken(): lex.peek() is not a double quote character")
	}

	tok := token.Token{Line: lex.line, Column: lex.column}
	lex.consume()

	chars := []byte{}
	for c := lex.peek(); c != 0; c = lex.peek() {

		// reached EOF, unterminated string literal
		if c == 0 {
			panic(fmt.Sprintf("%d:%d unterminated string literal", tok.Line, tok.Column))
		}
		// end of string literal
		if c == '"' {
			lex.consume()
			tok.Value = string(chars)
			tok.Type = token.STRING_LITERAL
			break
		}
		// escaped characters
		if c == '\\' {
			lex.consume() // consume the \ character
			echar := lex.peek()
			lex.consume()
			switch echar {
			case 'a':
				chars = append(chars, '\a')
			case 'n':
				chars = append(chars, '\n')
			case 't':
				chars = append(chars, '\t')
			case 'r':
				chars = append(chars, '\r')
			case 'v':
				chars = append(chars, '\v')
			case 'f':
				chars = append(chars, '\f')
			case '\\':
				chars = append(chars, '\\')
			case '"':
				chars = append(chars, '"')
			default:
				panic(fmt.Sprintf("unknown escaped character \\%c", echar))
			}
		} else {
			chars = append(chars, c)
			lex.consume()
		}
	}
	return tok
}

func (lex *Lexer) identifierOrKeywordToken() token.Token {
	tok := token.Token{Line: lex.line, Column: lex.column}
	chars := []byte{}
	for c := lex.peek(); isAlphaNumeric(c) || c == '_'; c = lex.peek() {
		chars = append(chars, c)
		lex.consume()
	}
	value := string(chars)
	tok.Value = value
	if kw, ok := token.IsKeyword(value); ok {
		tok.Type = kw
	} else {
		tok.Type = token.IDENTIFIER
	}
	return tok
}

func (lex *Lexer) NextToken() token.Token {

	ch := lex.peek()
	for isWhitespace(ch) {
		lex.consume()
		ch = lex.peek()
	}

	if ch == 0 {
		return token.Token{Type: token.EOF, Value: nil, Line: lex.line, Column: lex.column}
	}

	switch {
	case ch == '*':
		return lex.singleCharToken(token.ASTERISK)
	case ch == ',':
		return lex.singleCharToken(token.COMMA)
	case ch == '-':
		return lex.singleCharToken(token.MINUS)
	case ch == '+':
		return lex.singleCharToken(token.PLUS)
	case ch == ';':
		return lex.singleCharToken(token.SEMI_COLON)
	case ch == '/':
		return lex.singleCharToken(token.SLASH)
	case ch == '(':
		return lex.singleCharToken(token.LEFT_PAREN)
	case ch == ')':
		return lex.singleCharToken(token.RIGHT_PAREN)
	case ch == '{':
		return lex.singleCharToken(token.LEFT_BRACE)
	case ch == '}':
		return lex.singleCharToken(token.RIGHT_BRACE)
	case ch == '[':
		return lex.singleCharToken(token.LEFT_BRACKET)
	case ch == ']':
		return lex.singleCharToken(token.RIGHT_BRACKET)
	case ch == '^':
		return lex.singleCharToken(token.CARET)
	case ch == '~':
		return lex.singleCharToken(token.TILDE)
	case ch == '=':
		if lex.peekN(2) == "==" {
			return lex.doubleCharToken(token.EQUAL_EQUAL)
		}
		return lex.singleCharToken(token.EQUAL)
	case ch == '!':
		if lex.peekN(2) == "!=" {
			return lex.doubleCharToken(token.EXCLAMATION_EQUAL)
		}
		return lex.singleCharToken(token.EXCLAMATION)

	case ch == '>':
		if lex.peekN(2) == ">=" {
			return lex.doubleCharToken(token.GREATER_THAN_EQUAL)
		}
		return lex.singleCharToken(token.GREATER_THAN)

	case ch == '<':
		if lex.peekN(2) == "<=" {
			return lex.doubleCharToken(token.LESSER_THAN_EQUAL)
		}
		return lex.singleCharToken(token.LESSER_THAN)
	case ch == '&':
		if lex.peekN(2) == "&&" {
			return lex.doubleCharToken(token.AMPERSAND_AMPERSAND)
		}
		return lex.singleCharToken(token.AMPERSAND)
	case ch == '|':
		if lex.peekN(2) == "||" {
			return lex.doubleCharToken(token.PIPE_PIPE)
		}
		return lex.singleCharToken(token.PIPE)
	// string literal
	case ch == '"':
		return lex.stringLiteralToken()
	// number
	case isDigit(ch):
		return lex.numberLiteralToken()
	// identifier or keyword
	case isAlpha(ch):
		return lex.identifierOrKeywordToken()
	}
	panic("control should not reach here!")
}

func isWhitespace(c byte) bool {
	if c == ' ' || c == '\t' || c == '\n' || c == '\v' || c == '\r' {
		return true
	}
	return false
}

func isAlpha(c byte) bool {
	if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') {
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

func isAlphaNumeric(c byte) bool {
	return (isAlpha(c) || isDigit(c))
}
