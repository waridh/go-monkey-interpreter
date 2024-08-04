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
		return booleanObjectOfNativeBool(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right)
		return evalPrefixOperator(node.Operator, right)
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

func evalPrefixOperator(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalPrefixBang(right)
	case "-":
		return evalPrefixMinus(right)
	default:
		return NULL
	}
}

func evalPrefixBang(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalPrefixMinus(right object.Object) object.Object {
	switch v := right.(type) {
	case *object.Integer:
		if v.Type() != object.INTEGER_OBJ {
			return NULL
		}
		return &object.Integer{Value: -v.Value}
	default:
		return NULL
	}
}

func booleanObjectOfNativeBool(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}
