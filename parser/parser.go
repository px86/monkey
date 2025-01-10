package parser

import (
	"errors"
	"fmt"
	"github.com/px86/monkey/ast"
	"github.com/px86/monkey/lexer"
	"github.com/px86/monkey/token"
	"os"
)

const (
	_ int = iota
	PREC_LOWEST
	PREC_EQUALS
	PREC_LESSGREATER
	PREC_SUM
	PREC_PRODUCT
	PREC_PREFIX
	PREC_CALL
)

// NOTE TO MYSELF: it the responsibility of each terminal parse function to put the
// parser.curToken to the very next token that should be read.
// For example, if once parseLetStatement returns, the parser.curToken should point
// at the token exactly after the semi colon.

type Parser struct {
	l      *lexer.Lexer
	Errors []error

	curToken  token.Token
	nextToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Get two tokens from the lexer and populate peekToken and curToken.
	p.advance()
	p.advance()

	// p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	// p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)

	return p
}

// Moves peekToken to curToken, and fills peekToken with the next token from lexer.
// If peekToken is EOF, the lexer is not called anymore. Further calls to p.advance
// will have no effect. It is up to the parsing logic to gracefully handle the EOF token.
func (p *Parser) advance() {
	p.curToken = p.nextToken
	if p.nextToken.Type != token.EOF {
		p.nextToken = p.l.NextToken()
	}
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}
	if len(p.Errors) > 0 {
		for i, err := range p.Errors {
			fmt.Fprintf(os.Stderr, "Error %2d: %s\n", i, err)
		}
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.KW_LET:
		return p.parseLetStatement()
	case token.KW_RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectNextThenAdvance(token.IDENTIFIER) {
		return nil
	}

	ident, ok := p.curToken.Value.(string)
	if !ok {
		p.Errors = append(p.Errors,
			errors.New(fmt.Sprintf("identifier value is not string. got=%v", p.curToken.Value)))
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: ident}

	if !p.expectNextThenAdvance(token.EQUAL) {
		return nil
	}
	p.advance() // move past =

	stmt.Value = p.parseExpression(PREC_LOWEST)

	if !p.expectCurrentThenAdvance(token.SEMI_COLON) {
		return nil
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken} // return keyword
	p.advance()
	stmt.ReturnValue = p.parseExpression(PREC_LOWEST)
	if !p.expectCurrentThenAdvance(token.SEMI_COLON) {
		return nil
	}
	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(PREC_LOWEST)

	if p.curTokenIs(token.SEMI_COLON) {
		p.advance()
	}
	return stmt
}

func (p *Parser) parseIntegerLiteral() *ast.IntegerLiteral {
	value, ok := p.curToken.Value.(int)
	if !ok {
		p.Errors = append(p.Errors,
			errors.New(fmt.Sprintf("at line:%d, column:%d, %s value not of type %s. got=%T",
				p.curToken.Line, p.curToken.Column,
				token.TypeStr(p.curToken.Type), "int", p.curToken.Value)))
		return nil
	}
	integer := &ast.IntegerLiteral{Token: p.curToken, Value: value}
	p.advance()
	return integer
}

func (p *Parser) parseStringLiteral() *ast.StringLiteral {
	value, ok := p.curToken.Value.(string)
	if !ok {
		p.Errors = append(p.Errors,
			errors.New(fmt.Sprintf("at line:%d, column:%d, %s value not of type %s. got=%T",
				p.curToken.Line, p.curToken.Column,
				token.TypeStr(p.curToken.Type), "string", p.curToken.Value)))
		return nil
	}
	s := &ast.StringLiteral{Token: p.curToken, Value: value}
	p.advance()
	return s
}

func (p *Parser) parseIdentifier() *ast.Identifier {
	id, ok := p.curToken.Value.(string)
	if !ok {
		p.Errors = append(p.Errors,
			errors.New(fmt.Sprintf("at line:%d, column:%d, %s value not of type %s. got=%T",
				p.curToken.Line, p.curToken.Column,
				token.TypeStr(p.curToken.Type), "string", p.curToken.Value)))
		return nil
	}
	identifier := &ast.Identifier{Token: p.curToken, Value: id}
	p.advance()
	return identifier
}

func (p *Parser) parseFunctionExpression() *ast.FunctionExpr {
	fexpr := &ast.FunctionExpr{Token: p.curToken} // fn keyword
	p.advance()
	if !p.expectCurrentThenAdvance(token.LEFT_PAREN) {
		return nil
	}
	// args
	for !p.curTokenIs(token.RIGHT_PAREN) && !p.curTokenIs(token.EOF) {
		fexpr.Args = append(fexpr.Args, p.parseIdentifier())
		if p.curTokenIs(token.COMMA) {
			p.advance()
		}
	}

	if !p.expectCurrentThenAdvance(token.RIGHT_PAREN) {
		return nil
	}

	if !p.expectCurrentThenAdvance(token.LEFT_BRACE) {
		return nil
	}
	// body

	for !p.curTokenIs(token.RIGHT_BRACE) && !p.curTokenIs(token.EOF) {
		fexpr.Body = append(fexpr.Body, p.parseStatement())
		// if p.curTokenIs(token.SEMI_COLON) {
		// 	p.advance()
		// }
	}

	if !p.expectCurrentThenAdvance(token.RIGHT_BRACE) {
		return nil
	}

	return fexpr
}

func (p *Parser) parseFunctionCall() *ast.FunctionCall {
	ident := p.parseIdentifier()
	if ident == nil {
		return nil
	}
	var fcall *ast.FunctionCall
	if p.expectCurrentThenAdvance(token.LEFT_PAREN) {
		fcall = &ast.FunctionCall{Name: ident}
		for !p.curTokenIs(token.RIGHT_PAREN) {
			fcall.Args = append(fcall.Args, p.parseExpression(PREC_LOWEST))
			if p.curTokenIs(token.COMMA) {
				p.advance()
			}
		}
		p.advance()
	}
	return fcall
}

// See Jonathan Blow's video link above.
func (p *Parser) parseIncreasingPrecedence(left ast.Expression, minPrec int) ast.Expression {
	operator := p.curToken
	if !isBinaryOperator(operator.Type) {
		return left
	}
	nextPrec := precOf(operator.Type)
	if nextPrec <= minPrec {
		return left
	}
	p.advance()
	right := p.parseExpression(nextPrec)
	return &ast.InfixExpr{left, operator, right}
}

func (p *Parser) parseExpression(minPrec int) ast.Expression {
	left := p.parseLeaf()
	for {
		node := p.parseIncreasingPrecedence(left, minPrec)
		if left == node {
			break
		}
		left = node
	}
	return left
}

func (p *Parser) parseLeaf() ast.Expression {
	var leaf ast.Expression

	switch p.curToken.Type {
	case token.INTEGER:
		leaf = p.parseIntegerLiteral()
	case token.STRING_LITERAL:
		leaf = p.parseStringLiteral()
	case token.KW_FUNCTION:
		leaf = p.parseFunctionExpression()
	case token.IDENTIFIER:
		if p.nextTokenIs(token.LEFT_PAREN) {
			leaf = p.parseFunctionCall()
		} else {
			leaf = p.parseIdentifier()
		}
	}

	// p.advance()
	return leaf
}

func precOf(toktype token.TokenType) int {
	prec := PREC_LOWEST
	switch toktype {
	case token.PLUS:
		prec = PREC_SUM
	case token.MINUS:
		prec = PREC_SUM
	case token.ASTERISK:
		prec = PREC_PRODUCT
	case token.SLASH:
		prec = PREC_PRODUCT
	case token.LESSER_THAN:
		prec = PREC_LESSGREATER
	case token.LESSER_THAN_EQUAL:
		prec = PREC_LESSGREATER
	case token.GREATER_THAN:
		prec = PREC_LESSGREATER
	case token.GREATER_THAN_EQUAL:
		prec = PREC_LESSGREATER
	case token.KW_FUNCTION:
		prec = PREC_CALL
	}
	return prec
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) nextTokenIs(t token.TokenType) bool {
	return p.nextToken.Type == t
}

func (p *Parser) expectNextThenAdvance(t token.TokenType) bool {
	if p.nextToken.Type == t {
		p.advance()
		return true
	} else {
		p.Errors = append(p.Errors,
			errors.New(fmt.Sprintf("at line:%d, column:%d, expected %s, got=%s",
				p.nextToken.Line, p.nextToken.Column,
				token.TypeStr(t), token.TypeStr(p.nextToken.Type))))
		return false
	}
}

func (p *Parser) expectCurrentThenAdvance(t token.TokenType) bool {
	if p.curToken.Type == t {
		p.advance()
		return true
	} else {
		p.Errors = append(p.Errors,
			errors.New(fmt.Sprintf("at line:%d, column:%d, expected %s, got=%s",
				p.nextToken.Line, p.nextToken.Column,
				token.TypeStr(t), token.TypeStr(p.nextToken.Type))))
		return false
	}
}

func isBinaryOperator(toktype token.TokenType) bool {
	switch toktype {
	case token.PLUS:
		return true
	case token.MINUS:
		return true
	case token.ASTERISK:
		return true
	case token.SLASH:
		return true
	default:
		return false
	}
}
