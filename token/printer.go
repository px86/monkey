package token

import "fmt"

// Return a string representing the token.
func String(tok Token) string {
	tokTypeStr := AsString(tok.Type)
	return fmt.Sprintf("{Type:'%v', Value:'%v', Line:'%d', Column:'%d'}", tokTypeStr, tok.Value, tok.Line, tok.Column)
}

// Print the string representation of token. Suited for debugging purpose.
func PrintToken(tok Token) {
	fmt.Println(String(tok))
}

func AsString(toktype TokenType) string {
	switch toktype {

	case UNKNOWN:
		return "UNKNOWN"
	case EOF:
		return "EOF"
	case ASTERISK:
		return "*"
	case COMMA:
		return ","
	case MINUS:
		return "-"
	case PLUS:
		return "+"
	case SEMI_COLON:
		return ";"
	case SLASH:
		return "/"
	case LEFT_PAREN:
		return "("
	case RIGHT_PAREN:
		return ")"
	case LEFT_BRACE:
		return "{"
	case RIGHT_BRACE:
		return "}"
	case LEFT_BRACKET:
		return "["
	case RIGHT_BRACKET:
		return "]"
	case EQUAL:
		return "="
	case EQUAL_EQUAL:
		return "=="
	case EXCLAMATION:
		return "!"
	case EXCLAMATION_EQUAL:
		return "!="
	case GREATER_THAN:
		return ">"
	case GREATER_THAN_EQUAL:
		return ">="
	case LESSER_THAN:
		return "<"
	case LESSER_THAN_EQUAL:
		return "<="
	case AMPERSAND:
		return "&"
	case AMPERSAND_AMPERSAND:
		return "&&"
	case PIPE:
		return "|"
	case PIPE_PIPE:
		return "||"
	case CARET:
		return "^"
	case INTEGER:
		return "INTEGER"
	case FLOAT:
		return "FLOAT"
	case STRING_LITERAL:
		return "STRING"
	case IDENTIFIER:
		return "IDENTIFIER"
	case KW_LET:
		return "let"
	case KW_IF:
		return "if"
	case KW_ELSE:
		return "else"
	case KW_FUNCTION:
		return "fn"
	case KW_RETURN:
		return "return"

	default:
		return ""
	}
}

// Return TokenType as string. Useful for error reporting.
func TypeStr2(toktype TokenType) string {
	switch toktype {

	case UNKNOWN:
		return "UNKNOWN"
	case EOF:
		return "EOF"
	case ASTERISK:
		return "ASTERISK"
	case COMMA:
		return "COMMA"
	case MINUS:
		return "MINUS"
	case PLUS:
		return "PLUS"
	case SEMI_COLON:
		return "SEMI_COLON"
	case SLASH:
		return "SLASH"
	case LEFT_PAREN:
		return "LEFT_PAREN"
	case RIGHT_PAREN:
		return "RIGHT_PAREN"
	case LEFT_BRACE:
		return "LEFT_BRACE"
	case RIGHT_BRACE:
		return "RIGHT_BRACE"
	case LEFT_BRACKET:
		return "LEFT_BRACKET"
	case RIGHT_BRACKET:
		return "RIGHT_BRACKET"
	case EQUAL:
		return "EQUAL"
	case EQUAL_EQUAL:
		return "EQUAL_EQUAL"
	case EXCLAMATION:
		return "EXCLAMATION"
	case EXCLAMATION_EQUAL:
		return "NOT_EQUAL"
	case GREATER_THAN:
		return "GREATER_THAN"
	case GREATER_THAN_EQUAL:
		return "GREATER_THAN_EQUAL"
	case LESSER_THAN:
		return "LESS_THAN"
	case LESSER_THAN_EQUAL:
		return "LESS_THAN_EQUAL"
	case AMPERSAND:
		return "BITWISE_AND"
	case AMPERSAND_AMPERSAND:
		return "LOGICAL_AND"
	case PIPE:
		return "BITWISE_OR"
	case PIPE_PIPE:
		return "LOGICAL_OR"
	case CARET:
		return "XOR"
	case INTEGER:
		return "INTEGER"
	case FLOAT:
		return "FLOAT"
	case STRING_LITERAL:
		return "STRING"
	case IDENTIFIER:
		return "IDENTIFIER"
	case KW_LET:
		return "LET"
	case KW_IF:
		return "IF"
	case KW_ELSE:
		return "ELSE"
	case KW_FUNCTION:
		return "FUNCTION"
	case KW_RETURN:
		return "RETURN"

	default:
		return ""
	}
}
