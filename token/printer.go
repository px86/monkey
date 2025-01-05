package token

import "fmt"

func PrintToken(tok Token) {

	var tokTypeStr string

	switch tok.Type {

	case UNKNOWN:
		tokTypeStr = "UNKNOWN"
	case EOF:
		tokTypeStr = "EOF"
	case ASTERISK:
		tokTypeStr = "ASTERISK"
	case COMMA:
		tokTypeStr = "COMMA"
	case MINUS:
		tokTypeStr = "MINUS"
	case PLUS:
		tokTypeStr = "PLUS"
	case SEMI_COLON:
		tokTypeStr = "SEMI_COLON"
	case SLASH:
		tokTypeStr = "SLASH"
	case LEFT_PAREN:
		tokTypeStr = "LEFT_PAREN"
	case RIGHT_PAREN:
		tokTypeStr = "RIGHT_PAREN"
	case LEFT_BRACE:
		tokTypeStr = "LEFT_BRACE"
	case RIGHT_BRACE:
		tokTypeStr = "RIGHT_BRACE"
	case LEFT_BRACKET:
		tokTypeStr = "LEFT_BRACKET"
	case RIGHT_BRACKET:
		tokTypeStr = "RIGHT_BRACKET"
	case EQUAL:
		tokTypeStr = "EQUAL"
	case EQUAL_EQUAL:
		tokTypeStr = "EQUAL_EQUAL"
	case EXCLAMATION:
		tokTypeStr = "EXCLAMATION"
	case NOT_EQUAL:
		tokTypeStr = "NOT_EQUAL"
	case GREATER_THAN:
		tokTypeStr = "GREATER_THAN"
	case GREATER_THAN_EQUAL:
		tokTypeStr = "GREATER_THAN_EQUAL"
	case LESS_THAN:
		tokTypeStr = "LESS_THAN"
	case LESS_THAN_EQUAL:
		tokTypeStr = "LESS_THAN_EQUAL"
	case BITWISE_AND:
		tokTypeStr = "BITWISE_AND"
	case LOGICAL_AND:
		tokTypeStr = "LOGICAL_AND"
	case BITWISE_OR:
		tokTypeStr = "BITWISE_OR"
	case LOGICAL_OR:
		tokTypeStr = "LOGICAL_OR"
	case XOR:
		tokTypeStr = "XOR"
	case INTEGER:
		tokTypeStr = "INTEGER"
	case FLOAT:
		tokTypeStr = "FLOAT"
	case STRING:
		tokTypeStr = "STRING"
	case IDENTIFIER:
		tokTypeStr = "IDENTIFIER"
	case LET:
		tokTypeStr = "LET"
	case IF:
		tokTypeStr = "IF"
	case ELSE:
		tokTypeStr = "ELSE"
	case FUNCTION:
		tokTypeStr = "FUNCTION"
	case RETURN:
		tokTypeStr = "RETURN"

	default:
		tokTypeStr = "Oopsie...don't know about this type!"
	}

	fmt.Printf("{Type:%v, Value:%v, Line:%d, Column:%d}\n", tokTypeStr, tok.Value, tok.Line, tok.Column)
}
