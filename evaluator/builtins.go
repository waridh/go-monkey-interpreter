package evaluator

import (
	"fmt"

	"github.com/waridh/go-monkey-interpreter/object"
)

func builtinLenCheck(funcName string, expected int, args []object.Object) object.Object {
	if len(args) != expected {
		return newError("wrong number of arguments for %s. got=%d, want=%d", funcName, len(args), expected)
	}
	return nil
}

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if err := builtinLenCheck("len", 1, args); err != nil {
				return err
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("argument to `%s` not supported, got=%s", "len", arg.Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if err := builtinLenCheck("len", 1, args); err != nil {
				return err
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[0]
			default:
				return newError("argument to `%s` not supported, got=%s", "first", arg.Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if err := builtinLenCheck("len", 1, args); err != nil {
				return err
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return NULL
				}
				return arg.Elements[len(arg.Elements)-1]
			default:
				return newError("argument to `%s` not supported, got=%s", "last", arg.Type())
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if err := builtinLenCheck("len", 1, args); err != nil {
				return err
			}
			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return &object.Array{Elements: []object.Object{}}
				}
				return &object.Array{Elements: arg.Elements[1:]}
			default:
				return newError("argument to `%s` not supported, got=%s", "rest", arg.Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if err := builtinLenCheck("len", 2, args); err != nil {
				return err
			}
			switch arg := args[0].(type) {
			case *object.Array:
				length := len(arg.Elements)
				newElements := make([]object.Object, length+1, length+1)
				copy(newElements, arg.Elements)
				newElements[length] = args[1]
				return &object.Array{Elements: newElements}
			default:
				return newError("argument to `%s` not supported, got=%s", "len", arg.Type())
			}
		},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
