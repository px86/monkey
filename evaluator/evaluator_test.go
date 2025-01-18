package evaluator

import (
	"github.com/px86/monkey/lexer"
	"github.com/px86/monkey/object"
	"github.com/px86/monkey/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	testcases := []struct {
		expr     string
		expected int64
	}{
		{"5", 5},
		{"101", 101},
		{"1721", 1721},
	}

	for _, testcase := range testcases {
		obj := testEval(testcase.expr)
		testIntegerObject(t, obj, testcase.expected)
	}
}

func testEval(src string) object.Object {
	p := parser.New(lexer.New(src))
	prog := p.ParseProgram()
	return Eval(prog)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	iobj, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%v)", obj, obj)
		return false
	}
	if iobj.Value != expected {
		t.Errorf("object value does not match. got=%d. exprected=%d", iobj.Value, expected)
		return false
	}
	return true
}
