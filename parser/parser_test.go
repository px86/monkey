package parser

import (
	"bytes"
	"github.com/px86/monkey/ast"
	"github.com/px86/monkey/lexer"
	"github.com/px86/monkey/token"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	for i, err := range p.Errors {
		t.Errorf("[%d] %s\n", i, err)
	}
}

func TestReturnStatements(t *testing.T) {
	input := "return foo(bar, baz);"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ReturnStatement)
	if !ok {
		t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
	}
	fcall, ok := stmt.ReturnValue.(*ast.FunctionCall)
	if !ok {
		t.Fatalf("stmt.ReturnValue not *ast.FunctionCall. got=%T", stmt.ReturnValue)
	}
	if fcall.Name.Value != "foo" {
		t.Fatalf("function name not %q. got=%q", "foo", fcall.Name.Value)
	}
}

func TestLetStatements(t *testing.T) {
	input := "let foo = bar(\"spam\", 10+20*30-5);"
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("stmt not *ast.LetStatement. got=%T", stmt)
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	// checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statements. got=%d", len(program.Statements))
	}

	expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("statement not *ast.ExpressionStatement. got=%T", expStmt)
	}

	ident, ok := expStmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", expStmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.String() != "foobar" {
		t.Errorf("ident.String not %s. got=%s", "foobar", ident.String())
	}
}

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

	be, ok := expStmt.Expression.(*ast.InfixExpr)
	if !ok {
		t.Fatalf("exp not *ast.BinaryExpression. got=%T", expStmt.Expression)
	}

	if be.Left == nil {
		t.Fatalf("be.Left is nil")
	}

	if be.Right == nil {
		t.Fatalf("be.Left is nil")
	}

	if be.String() != "(+ 10 30)" {
		t.Fatalf("String() not %q. got=%v", "(+ 10 30)", be.String())
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

	be, ok := expStmt.Expression.(*ast.InfixExpr)
	if !ok {
		t.Fatalf("exp not *ast.BinaryExpression. got=%T", expStmt.Expression)
	}

	if be.Operator.Type != token.PLUS {
		t.Fatalf("top operator not PLUS. got=%v\n.String()=%q", token.AsString(be.Operator.Type), be.String())
	}

}

func TestFunctionCall(t *testing.T) {

	cases := []struct {
		FunctionName string
		Args         []string
	}{
		{"foo", []string{}},
		{"bar", []string{"1", "2"}},
		{"baz", []string{"1+2*3-2", "100", "\"lorem ipsum.\""}},
		{"blahBlah", []string{"foo", "bar(10, 2+20*30-1)"}},
	}

	for i, tt := range cases {
		var buff bytes.Buffer
		buff.WriteString(tt.FunctionName + "(")
		for j, arg := range tt.Args {
			buff.WriteString(arg)
			if j != len(tt.Args)-1 {
				buff.WriteString(",")
			}
		}
		buff.WriteString(");")
		input := buff.String()

		l := lexer.New(input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("[%d] program.Statements does not contain 1 statement. got=%d", i, len(program.Statements))
		}

		expStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("[%d] statement not *ast.ExpressionStatement. got=%T", i, expStmt)
		}

		fc, ok := expStmt.Expression.(*ast.FunctionCall)
		if !ok {
			t.Fatalf("[%d] exp not *ast.Functioncall. got=%T", i, expStmt.Expression)
		}

		if fc.Name.Value != tt.FunctionName {
			t.Fatalf("[%d] function name not %q. got=%q", i, tt.FunctionName, fc.Name.Value)
		}

		if len(fc.Args) != len(tt.Args) {
			t.Fatalf("[%d] function args len not %d. got=%d", i, len(tt.Args), len(fc.Args))
		}

	}

}
