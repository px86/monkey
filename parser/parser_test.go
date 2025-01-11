package parser

import (
	"bytes"
	"github.com/px86/monkey/ast"
	"github.com/px86/monkey/lexer"
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

	input := []struct {
		stmt string
		tree string
	}{
		{"let a = 1;", "(let a 1)"},
		{"let b = foo(x, y);", "(let b (foo x y))"},
	}

	for i, testcase := range input {

		p := New(lexer.New(testcase.stmt))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("[TC %d] program does not contain 1 statement. got=%d",
				i, len(program.Statements))
		}

		lstmt, ok := program.Statements[0].(*ast.LetStatement)
		if !ok {
			t.Errorf("[TC %d] stmt not *ast.LetStatement. got=%T", i, program.Statements[0])
		}

		if lstmt.String() != testcase.tree {
			t.Errorf("[TC %d] AST string didn't match. expected=%q, got=%q",
				i, testcase.tree, lstmt.String())
		}
	}
}

func TestPrefixExpressions(t *testing.T) {

	input := []struct {
		expr string
		tree string
	}{
		{"-1;", "(- 1)"},
		{"-foo(x, y);", "(- (foo x y))"},
		{"!foo;", "(! foo)"},
	}

	for i, testcase := range input {

		p := New(lexer.New(testcase.expr))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("[TC %d] program does not contain 1 statement. got=%d",
				i, len(program.Statements))
		}

		estmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("[TC %d] statement not *ast.ExpressionStatement. got=%T", i, estmt)
		}

		prefix, ok := estmt.Expression.(*ast.PrefixExpr)
		if !ok {
			t.Errorf("[TC %d] expr not *ast.PrefixExpr. got=%T", i, estmt.Expression)
		}
		if prefix.String() != testcase.tree {
			t.Errorf("[TC %d] AST string didn't match. expected=%q, got=%q",
				i, testcase.tree, prefix.String())
		}
	}
}

func TestBinaryExpressions(t *testing.T) {

	input := []struct {
		expr string
		tree string
	}{
		{"1 + 2;", "(+ 1 2)"},
		{"1 * 2 + 3;", "(+ (* 1 2) 3)"},
		{"1 + 2 * 3 + 4;", "(+ (+ 1 (* 2 3)) 4)"},
		{"1 + 2 * 3 + 4/2 - 1;", "(- (+ (+ 1 (* 2 3)) (/ 4 2)) 1)"},
		{"bar() * foo + 3;", "(+ (* (bar) foo) 3)"},
		{"2 * (3 + 4);", "(* 2 (+ 3 4))"},
	}

	for i, testcase := range input {

		p := New(lexer.New(testcase.expr))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("[TC %d] program does not contain 1 statement. got=%d",
				i, len(program.Statements))
		}

		estmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Errorf("[TC %d] statement not *ast.ExpressionStatement. got=%T", i, estmt)
		}

		infix, ok := estmt.Expression.(*ast.InfixExpr)
		if !ok {
			t.Errorf("[TC %d] expr not *ast.InfixExpr. got=%T", i, estmt.Expression)
		}
		if infix.String() != testcase.tree {
			t.Errorf("[TC %d] AST string didn't match. expected=%q, got=%q",
				i, testcase.tree, infix.String())
		}
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
