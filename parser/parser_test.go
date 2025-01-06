package parser

import (
	"github.com/px86/monkey/ast"
	"github.com/px86/monkey/lexer"
	"github.com/px86/monkey/token"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	for i, err := range p.Errors {
		t.Errorf("[%2d] %s\n", i, err)
	}
}

func TestReturnStatements(t *testing.T) {
	input := `
return 12;
return 17;
return 30071;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	// checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.returnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.String() != "return ;" {
			t.Errorf("returnStmt.String not 'return', got %q", returnStmt.String())
		}
	}
}

// func TestIdentifierExpression(t *testing.T) {
// 	input := "foobar;"

// 	l := lexer.New(input)
// 	p := New(l)

// 	program := p.ParseProgram()
// 	// checkParserErrors(t, p)

// 	if len(program.Statements) != 1 {
// 		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
// 	}

// 	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
// 	if !ok {
// 		t.Fatalf("statement not *ast.ExpressionStatement. got=%T", expStmt)
// 	}

// 	ident, ok := expStmt.Expression.(*ast.Identifier)
// 	if !ok {
// 		t.Fatalf("exp not *ast.Identifier. got=%T", expStmt.Expression)
// 	}

// 	if ident.Value != "foobar" {
// 		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
// 	}

// 	if ident.String() != "foobar" {
// 		t.Errorf("ident.String not %s. got=%s", "foobar", ident.String())
// 	}
// }

func TestIntegerExpression(t *testing.T) {
	input := "10 + 30;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement not *ast.ExpressionStatement. got=%T", expStmt)
	}

	be, ok := expStmt.Expression.(*ast.BinaryExpression)
	if !ok {
		t.Fatalf("exp not *ast.BinaryExpression. got=%T", expStmt.Expression)
	}

	if be.Left == nil {
		t.Fatalf("be.Left is nil")
	}

	if be.Right == nil {
		t.Fatalf("be.Left is nil")
	}

	if be.String() != "10 PLUS 30" {
		t.Fatalf("String() not %q. got=%v", "10 PLUS 30", be.String())
	}
}

func TestBinaryExpression(t *testing.T) {
	input := "1 + 2 * 3 + 4;"
	/*
	        +
	       / \
	      +   4
	     / \
	   1   *
	      / \
	     2   3
	*/

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement not *ast.ExpressionStatement. got=%T", expStmt)
	}

	be, ok := expStmt.Expression.(*ast.BinaryExpression)
	if !ok {
		t.Fatalf("exp not *ast.BinaryExpression. got=%T", expStmt.Expression)
	}

	if be.Operator.Type != token.PLUS {
		t.Fatalf("top operator not PLUS. got=%v\n.String()=%q", token.TypeStr(be.Operator.Type), be.String())
	}

}
