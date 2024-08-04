package evaluator

import (
	"github.com/waridh/go-monkey-interpreter/ast"
	"github.com/waridh/go-monkey-interpreter/object"
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	case *ast.IntegerLiteral:
		return evalIntegerLiteral(node)
	}

	return nil
}

func evalStatements(node []ast.Statement) object.Object {
	var result object.Object
	for _, stmt := range node {
		result = Eval(stmt)
	}

	return result
}

func evalIntegerLiteral(node *ast.IntegerLiteral) *object.Integer {
	return &object.Integer{Value: node.Value}
}
