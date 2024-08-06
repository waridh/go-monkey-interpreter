package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/waridh/go-monkey-interpreter/ast"
	"github.com/waridh/go-monkey-interpreter/functools"
)

type (
	ObjectType      string
	BuiltinFunction func(args ...Object) Object
)

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	STRING_OBJ       = "STRING"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
)

type Hashable interface {
	HashKey() HashKey
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: INTEGER_OBJ, Value: uint64(i.Value)}
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) HashKey() HashKey {
	var value uint64

	if b.Value {
		value = 1
	} else {
		value = 0
	}

	return HashKey{Type: BOOLEAN_OBJ, Value: value}
}

type Null struct{}

func (null *Null) Inspect() string  { return "null" }
func (null *Null) Type() ObjectType { return NULL_OBJ }

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

type Error struct {
	Message string
}

func (er *Error) Inspect() string  { return "ERROR: " + er.Message }
func (er *Error) Type() ObjectType { return ERROR_OBJ }

type Function struct {
	Parameter []*ast.Identifier
	Body      *ast.BlockStatement
	Env       *Environment
}

func (fn *Function) Inspect() string {
	var out bytes.Buffer

	params := functools.Map(fn.Parameter, func(x *ast.Identifier) string { return x.String() })

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(fn.Body.String())
	out.WriteString("}\n")

	return out.String()
}
func (fn *Function) Type() ObjectType { return FUNCTION_OBJ }

type String struct {
	Value string
}

func (str *String) Inspect() string  { return str.Value }
func (str *String) Type() ObjectType { return STRING_OBJ }
func (str *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(str.Value))
	return HashKey{Type: STRING_OBJ, Value: h.Sum64()}
}

type Array struct {
	Elements []Object
}

func (arr *Array) Type() ObjectType { return ARRAY_OBJ }
func (arr *Array) Inspect() string {
	var out bytes.Buffer

	ele := functools.Map(arr.Elements, func(x Object) string { return x.Inspect() })
	out.WriteString("[")
	out.WriteString(strings.Join(ele, ", "))
	out.WriteString("]")

	return out.String()
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (hash *Hash) Type() ObjectType { return HASH_OBJ }
func (hash *Hash) Inspect() string {
	var out bytes.Buffer

	pairs := make([]string, len(hash.Pairs))

	for _, value := range hash.Pairs {
		pairs = append(pairs, value.Key.Inspect()+": "+value.Key.Inspect())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}

type Builtin struct {
	Fn BuiltinFunction
}

func (bi *Builtin) Inspect() string  { return "builtin function" }
func (bi *Builtin) Type() ObjectType { return BUILTIN_OBJ }

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object), outer: nil}
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok := e.outer.Get(name)
		return obj, ok
	}
	return obj, ok
}

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}
