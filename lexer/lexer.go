package lexer

import (
	"github.com/px86/monkey/token"
	"io"
	"os"
)

type Lexer struct {
	source     string
	pos        int
	line       int
	column     int
	lastColumn int // to restore 'column' when a byte is retracted
}

func New(source string) *Lexer {
	return &Lexer{
		source: source,
		line:   1,
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

// return next character (byte) in the source code
func (lex *Lexer) next() byte {
	if lex.pos >= len(lex.source) {
		return 0
	}
	b := lex.source[lex.pos]
	lex.pos++
	if b == '\n' {
		lex.line++
		lex.lastColumn = lex.column
		lex.column = 0
	} else {
		lex.column++
	}
	return b
}

// move the input pointer one character (byte) backwards
func (lex *Lexer) retract() {
	if lex.pos > 0 {
		lex.pos--
	}
	b := lex.source[lex.pos]
	if b == '\n' {
		lex.line--
		lex.column = lex.lastColumn
	} else {
		lex.column--
	}
}

func (lex *Lexer) NextToken() token.Token {

	tok := token.Token{Line: lex.line, Column: lex.column}

	// this loop should break after each complete iteration (unless you call 'continue')
	// ensure there is a break statement at the bottom
	for {
		ch := lex.next()
		if ch == 0 {
			tok.Type = token.EOF
		}
		switch {
		case isWhitespace(ch):
			continue
		case ch == '*':
			tok.Type = token.ASTERISK
		case ch == ',':
			tok.Type = token.COMMA
		case ch == '-':
			tok.Type = token.MINUS
		case ch == '+':
			tok.Type = token.PLUS
		case ch == ';':
			tok.Type = token.SEMI_COLON
		case ch == '/':
			tok.Type = token.SLASH
		case ch == '(':
			tok.Type = token.LEFT_PAREN
		case ch == ')':
			tok.Type = token.RIGHT_PAREN
		case ch == '{':
			tok.Type = token.LEFT_BRACE
		case ch == '}':
			tok.Type = token.RIGHT_BRACE
		case ch == '[':
			tok.Type = token.LEFT_BRACKET
		case ch == ']':
			tok.Type = token.RIGHT_BRACKET
		case ch == '^':
			tok.Type = token.CARET
		case ch == '~':
			tok.Type = token.TILDE
		case ch == '=':
			if lex.next() == '=' {
				tok.Type = token.EQUAL_EQUAL
			} else {
				lex.retract()
				tok.Type = token.EQUAL
			}
		case ch == '!':
			if lex.next() == '=' {
				tok.Type = token.EXCLAMATION_EQUAL
			} else {
				lex.retract()
				tok.Type = token.EXCLAMATION
			}
		case ch == '>':
			if lex.next() == '=' {
				tok.Type = token.GREATER_THAN_EQUAL
			} else {
				lex.retract()
				tok.Type = token.GREATER_THAN
			}
		case ch == '<':
			if lex.next() == '=' {
				tok.Type = token.LESSER_THAN_EQUAL
			} else {
				lex.retract()
				tok.Type = token.LESSER_THAN
			}
		case ch == '&':
			if lex.next() == '&' {
				tok.Type = token.AMPERSAND_AMPERSAND
			} else {
				lex.retract()
				tok.Type = token.AMPERSAND
			}
		case ch == '|':
			if lex.next() == '|' {
				tok.Type = token.PIPE_PIPE
			} else {
				lex.retract()
				tok.Type = token.PIPE
			}
		// string literal
		case ch == '"':
			chars := []byte{}
		STRING:
			for c := lex.next(); c != 0; c = lex.next() {
				switch c {
				case 0:
					tok.Type = token.UNKNOWN
					break STRING
				case '"':
					break STRING
				case '\\':
					switch lex.next() {
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
						tok.Type = token.UNKNOWN
						break STRING
					}
				default:
					chars = append(chars, c)
				}
			}
			if tok.Type != token.UNKNOWN {
				tok.Type = token.STRING_LITERAL
				tok.Value = string(chars)
			}
		// number
		case isDigit(ch):
			// TODO: check floating point numbers
			value := int(ch) - int('0')
			for d := lex.next(); isDigit(d); d = lex.next() {
				value = value*10 + int(d) - int('0')
			}
			lex.retract()
			tok.Type = token.INTEGER
			tok.Value = value
		// identifier or keyword
		case isAlpha(ch):
			chars := []byte{ch}
			for c := lex.next(); isAlphaNumeric(c) || c == '_'; c = lex.next() {
				chars = append(chars, c)
			}
			lex.retract()
			value := string(chars)
			tok.Value = value
			if kw, ok := token.IsKeyword(value); ok {
				tok.Type = kw
			} else {
				tok.Type = token.IDENTIFIER
			}
		}
		break // break the infinite for loop
	}

	return tok
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

func isAlphaNumeric(c byte) bool {
	if 'a' <= c && c <= 'z' {
		return true
	}
	if 'A' <= c && c <= 'Z' {
		return true
	}
	return isDigit(c)
}
