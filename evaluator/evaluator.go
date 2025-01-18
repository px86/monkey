package evaluator

import (
	"github.com/px86/monkey/ast"
	"github.com/px86/monkey/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {

	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	}

	return nil
}

// The evaluation of a list of statements returns the value of the
// evaluation of the last expression.
func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range stmts {
		result = Eval(stmt)
	}
	return result
}
