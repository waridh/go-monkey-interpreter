package evaluator

import (
	"github.com/waridh/go-monkey-interpreter/ast"
	"github.com/waridh/go-monkey-interpreter/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return &object.Boolean{Value: node.Value}
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

func booleanObjectOfNativeBool(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}
