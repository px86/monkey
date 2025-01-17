package evaluator

import (
	"github.com/px86/monkey/ast"
	"github.com/px86/monkey/object"
)

func Eval(node ast.Node) object.Object {
	if il, ok := node.(*ast.IntegerLiteral); ok {
		return &object.Integer{Value: il.Value}
	}
	return nil
}
