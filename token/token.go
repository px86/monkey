package token

type TokenType int

const (
	UNKNOWN TokenType = iota + 1
	EOF

	ASTERISK           // *
	COMMA              // ,
	MINUS              // -
	PLUS               // +
	SEMI_COLON         // ;
	SLASH              // /
	LEFT_PAREN         // (
	RIGHT_PAREN        // )
	LEFT_BRACE         // {
	RIGHT_BRACE        // }
	LEFT_BRACKET       // [
	RIGHT_BRACKET      // ]
	EQUAL              // =
	EQUAL_EQUAL        // ==
	EXCLAMATION        // !
	NOT_EQUAL          // !=
	GREATER_THAN       // >
	GREATER_THAN_EQUAL // >=
	LESS_THAN          // <
	LESS_THAN_EQUAL    // <=
	BITWISE_AND        // &
	LOGICAL_AND        // &&
	BITWISE_OR         // |
	LOGICAL_OR         // ||
	XOR                // ^

	INTEGER
	FLOAT
	STRING

	IDENTIFIER
	LET
	IF
	ELSE
	FUNCTION
	RETURN
)

type Token struct {
	Type   TokenType
	Value  any
	Line   int // line on which token starts
	Column int // column on which token starts
}

var KeywordsMap = map[string]TokenType{
	"fn":     FUNCTION,
	"let":    LET,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
}
