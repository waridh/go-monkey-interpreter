package evaluator

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/waridh/go-monkey-interpreter/ast"
	"github.com/waridh/go-monkey-interpreter/functools"
	"github.com/waridh/go-monkey-interpreter/object"
)

// For static values, we can just refer to the same objects
var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return evalLetStatement(node, val, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return booleanObjectOfNativeBool(node.Value)

	case *ast.Identifier:
		val, ok := env.Get(node.Value)
		if !ok {
			return newError("identity not found: %s", node.Value)
		}
		return val

	case *ast.FunctionLiteral:
		params := node.Parameter
		body := node.Body
		return &object.Function{Parameter: params, Body: body, Env: env}
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixOperator(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixOperator(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.CallExpression:
		function := Eval(node.Function, env)
		if isError(function) {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(function, args)
	}

	return nil
}

func applyFunction(function object.Object, args []object.Object) object.Object {
	fn, ok := function.(*object.Function)
	if !ok {
		return newError("not a function: %s", function.Type())
	}
	newEnv := object.NewEnclosedEnvironment(fn.Env)
	if len(args) != len(fn.Parameter) {
		var out bytes.Buffer
		expected := functools.Map(fn.Parameter, func(x *ast.Identifier) string { return x.String() })
		got := functools.Map(args, func(x object.Object) string { return x.Inspect() })
		out.WriteString("missing parameters:\n")
		out.WriteString("\texpected: ")
		out.WriteString(strings.Join(expected, ", "))
		out.WriteString("\n\tgot: ")
		out.WriteString(strings.Join(got, ", "))
		return newError(out.String())
	}
	for idx, arg := range args {
		newEnv.Set(fn.Parameter[idx].Value, arg)
	}
	evaluated := Eval(fn.Body, newEnv)
	return unwrapReturnValue(evaluated)
}

// Takes an array of ast.Expressions and return an array of the evaluated
// expressions. If an error was encountered, return an array holding just
// the encountered error.
func evalExpressions(exprs []ast.Expression, env *object.Environment) []object.Object {
	evals := []object.Object{}
	for _, expr := range exprs {
		eval := Eval(expr, env)
		if isError(eval) {
			return []object.Object{eval}
		}
		evals = append(evals, eval)
	}
	return evals
}

func evalProgram(node []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range node {
		result = Eval(stmt, env)
		switch res := result.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}
	}

	return result
}

func evalBlockStatements(node []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, stmt := range node {
		result = Eval(stmt, env)
		if result.Type() == object.RETURN_VALUE_OBJ || result.Type() == object.ERROR_OBJ {
			return result
		}
	}

	return result
}

func evalLetStatement(node *ast.LetStatement, val object.Object, env *object.Environment) object.Object {
	return env.Set(node.Name.Value, val)
}

func unwrapReturnValue(evaluated object.Object) object.Object {
	if evaluated == nil {
		return NULL
	} else if evaluated.Type() == object.RETURN_VALUE_OBJ {
		return evaluated.(*object.ReturnValue).Value
	} else {
		return evaluated
	}
}

func evalReturnStatement(stmt *ast.ReturnStatement, env *object.Environment) object.Object {
	var val object.Object
	if stmt.ReturnValue == nil {
		val = NULL
	} else {
		val = Eval(stmt.ReturnValue, env)
	}
	if isError(val) {
		return val
	} else {
		return &object.ReturnValue{Value: val}
	}
}

func evalIfExpression(node *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(node.Condition, env)
	if isError(cond) {
		return cond
	}

	if isTruthy(cond) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
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

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
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
			return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
		}
	default:
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalPrefixOperator(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalPrefixBang(right)
	case "-":
		return evalPrefixMinus(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
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
		return newError("unknown operator: %s %s %s", object.INTEGER_OBJ, operator, object.INTEGER_OBJ)
	}
}

func evalInfixBooleanExpression(operator string, left object.Object, right object.Object) object.Object {
	switch operator {
	case "==":
		return booleanObjectOfNativeBool(left == right)
	case "!=":
		return booleanObjectOfNativeBool(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
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
			return newError("unknown operator: %s%s", "-", right.Type())
		}
		return &object.Integer{Value: -v.Value}
	default:
		return newError("unknown operator: %s%s", "-", right.Type())
	}
}

func newError(format string, a ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

func booleanObjectOfNativeBool(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}
}
