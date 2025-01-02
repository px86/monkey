package token

type Token struct {
	Type   TokenType
	Value  any
	Line   int
	Column int
}

type TokenType string

const (
	ASTERISK_SIGN  = "*"
	BANG_SIGN      = "!"
	COMMA          = ","
	ELSE           = "else"
	EOF            = ""
	EQUAL_SIGN     = "="
	FALSE          = "false"
	FUNCTION       = "fn"
	GREATER_THAN   = ">"
	IDENTIFIER     = "IDENTIFIER"
	IF             = "if"
	ILLEGAL        = "ILLEGAL"
	INTEGER        = "INT"
	LEFT_BRACE     = "{"
	LEFT_PAREN     = "("
	LESS_THAN      = "<"
	LET            = "let"
	MINUS_SIGN     = "-"
	PLUS_SIGN      = "+"
	RETURN         = "return"
	RIGHT_BRACE    = "}"
	RIGHT_PAREN    = ")"
	SEMI_COLON     = ";"
	SLASH          = "/"
	STRING_LITERAL = "STRING_LITERAL"
	TRUE           = "true"
)
