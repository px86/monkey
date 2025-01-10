package token

type TokenType int

const (
	UNKNOWN TokenType = iota + 1
	EOF

	ASTERISK            // *
	COMMA               // ,
	MINUS               // -
	PLUS                // +
	SEMI_COLON          // ;
	SLASH               // /
	LEFT_PAREN          // (
	RIGHT_PAREN         // )
	LEFT_BRACE          // {
	RIGHT_BRACE         // }
	LEFT_BRACKET        // [
	RIGHT_BRACKET       // ]
	EQUAL               // =
	EQUAL_EQUAL         // ==
	EXCLAMATION         // !
	EXCLAMATION_EQUAL   // !=
	GREATER_THAN        // >
	GREATER_THAN_EQUAL  // >=
	LESSER_THAN         // <
	LESSER_THAN_EQUAL   // <=
	TILDE               // ~
	AMPERSAND           // &
	AMPERSAND_AMPERSAND // &&
	PIPE                // |
	PIPE_PIPE           // ||
	CARET               // ^

	INTEGER
	FLOAT
	STRING_LITERAL

	IDENTIFIER
	KW_LET      // let
	KW_IF       // if
	KW_ELSE     // else
	KW_FUNCTION // fn
	KW_RETURN   // return
	KW_TRUE     // true
	KW_FALSE    // false
)

var kwMap = map[string]TokenType{
	"fn":     KW_FUNCTION,
	"let":    KW_LET,
	"if":     KW_IF,
	"else":   KW_ELSE,
	"return": KW_RETURN,
	"true":   KW_TRUE,
	"false":  KW_FALSE,
}

type Token struct {
	Type   TokenType
	Value  any
	Line   int // line on which token starts
	Column int // column on which token starts
}

func IsKeyword(s string) (kw TokenType, ok bool) {
	kw, ok = kwMap[s]
	return kw, ok
}
