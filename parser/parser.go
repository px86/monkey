package parser

import (
	"errors"
	"fmt"
	"github.com/px86/monkey/ast"
	"github.com/px86/monkey/lexer"
	"github.com/px86/monkey/token"
	"os"
)

/*******************************************************************************
* Resources On Parsing
*
* Robert Nystrom's article titled: 'Pratt Parsers: Expression Parsing Made Easy'
* https://journal.stuffwithstuff.com/2011/03/19/pratt-parsers-expression-parsing-made-easy/
*
*
* Jonathan Blow's video titled: 'Discussion with Casey Muratori about how easy precedence is...'
*  https://www.youtube.com/watch?v=fIPO4G42wYE
*
********************************************************************************/

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

/*
type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)
*/

type Parser struct {
	l      *lexer.Lexer
	Errors []error // TODO: implement parse error reporting.

	curToken  token.Token
	peekToken token.Token

	// TODO: implement the parsing yourself first. If it does not work, consult the interpreter book.
	// prefixParseFns map[token.TokenType]prefixParseFn
	// infixParseFns  map[token.TokenType]infixParseFn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Get two tokens from the lexer and populate peekToken and curToken.
	p.nextToken()
	p.nextToken()

	// p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	// p.registerPrefix(token.IDENTIFIER, p.parseIdentifier)

	return p
}

// Moves peekToken to curToken, and fills peekToken with the next token from lexer.
// If peekToken is EOF, the lexer is not called anymore. Further calls to p.nextToken
// will have no effect. It is up to the parsing logic to gracefully handle the EOF token.
func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	if p.peekToken.Type != token.EOF {
		p.peekToken = p.l.NextToken()
	}
}

/*
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}
*/

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
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
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}
	if !p.expectPeek(token.IDENTIFIER) {
		p.Errors = append(p.Errors,
			errors.New(fmt.Sprintf("expected IDENTIFIER, got=%v", token.TypeStr(p.peekToken.Type))))
		return nil
	}
	ident, ok := p.curToken.Value.(string)
	if !ok {
		p.Errors = append(p.Errors,
			errors.New(fmt.Sprintf("identifier value is not string. got=%v", p.curToken.Value)))
	}
	stmt.Name = &ast.Identifier{Token: p.curToken, Value: ident}

	if !p.expectPeek(token.EQUAL) {
		p.Errors = append(p.Errors,
			errors.New(fmt.Sprintf("expected EQUAL, got=%v", token.TypeStr(p.peekToken.Type))))
		return nil
	}
	for !p.curTokenIs(token.SEMI_COLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}
	p.nextToken()

	// TODO: parse the expression
	for !p.curTokenIs(token.SEMI_COLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(PREC_LOWEST)

	if p.peekTokenIs(token.SEMI_COLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseIdentifier() ast.Expression {
	s, ok := p.curToken.Value.(string)
	if !ok {
		p.Errors = append(p.Errors,
			errors.New(fmt.Sprintf("identifier value is not string. got=%v", p.curToken.Value)))
	}

	return &ast.Identifier{Token: p.curToken, Value: s}
}

// See Jonathan Blow's video link above.
func (p *Parser) parseIncreasingPrecedence(left ast.Expression, minPrec int) ast.Expression {
	next := p.curToken
	if !isBinaryOperator(next.Type) {
		return left
	}
	nextPrec := precOf(next.Type)
	if nextPrec <= minPrec {
		return left
	}
	p.nextToken()
	right := p.parseExpression(nextPrec)
	return &ast.BinaryExpression{left, next, right}
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
		value, ok := p.curToken.Value.(int)
		if !ok {
			p.Errors = append(p.Errors, errors.New(
				fmt.Sprintf("value of INTEGER is not int. got=%v", p.curToken.Value)))
		}
		leaf = &ast.Integer{Token: p.curToken, Value: value}
	case token.STRING:
		value, ok := p.curToken.Value.(string)
		if !ok {
			p.Errors = append(p.Errors, errors.New(
				fmt.Sprintf("value of STRING is not string. got=%v", p.curToken.Value)))
		}
		leaf = &ast.Str{Token: p.curToken, Value: value}
	}

	p.nextToken()
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
	case token.LESS_THAN:
		prec = PREC_LESSGREATER
	case token.LESS_THAN_EQUAL:
		prec = PREC_LESSGREATER
	case token.GREATER_THAN:
		prec = PREC_LESSGREATER
	case token.GREATER_THAN_EQUAL:
		prec = PREC_LESSGREATER
	case token.FUNCTION:
		prec = PREC_CALL
	}
	return prec
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// Match the expected token with peekToken.
// If matched, call p.nextToken().
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	} else {
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
