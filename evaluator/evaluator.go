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
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixOperator(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalStatements(node.Statements)
	case *ast.IfExpression:
		return evalIfExpression(node)
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

func evalIfExpression(node *ast.IfExpression) object.Object {
	cond := Eval(node.Condition)

	if isTruthy(cond) {
		return Eval(node.Consequence)
	} else if node.Alternative != nil {
		return Eval(node.Alternative)
	} else {
		return NULL
	}
}

func isTruthy(cond object.Object) bool {
	switch cond {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func evalInfixOperator(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == right.Type():
		switch {
		case left.Type() == object.INTEGER_OBJ:
			leftVal := left.(*object.Integer).Value
			rightVal := right.(*object.Integer).Value
			return evalInfixIntegerExpression(operator, leftVal, rightVal)
		case left.Type() == object.BOOLEAN_OBJ:
			return evalInfixBooleanExpression(operator, left, right)
		default:
			return NULL
		}
	default:
		return NULL
	}
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

func evalInfixIntegerExpression(operator string, left int64, right int64) object.Object {
	switch operator {
	case "+":
		return &object.Integer{Value: left + right}
	case "-":
		return &object.Integer{Value: left - right}
	case "*":
		return &object.Integer{Value: left * right}
	case "/":
		return &object.Integer{Value: left / right}
	case "<":
		return booleanObjectOfNativeBool(left < right)
	case ">":
		return booleanObjectOfNativeBool(left > right)
	case "==":
		return booleanObjectOfNativeBool(left == right)
	case "!=":
		return booleanObjectOfNativeBool(left != right)
	default:
		return NULL
	}
}

func evalInfixBooleanExpression(operator string, left object.Object, right object.Object) object.Object {
	switch operator {
	case "==":
		return booleanObjectOfNativeBool(left == right)
	case "!=":
		return booleanObjectOfNativeBool(left != right)
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
