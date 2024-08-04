package evaluator

import (
	"fmt"
	"testing"

	"github.com/waridh/go-monkey-interpreter/lexer"
	"github.com/waridh/go-monkey-interpreter/object"
	"github.com/waridh/go-monkey-interpreter/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{
			"5",
			5,
		},
		{
			"10",
			10,
		},
		{
			"10;",
			10,
		},
		{
			"-10;",
			-10,
		},
		{
			"-5",
			-5,
		},
		{
			"2 + 2",
			4,
		},
		{
			"5 + 5 + 5 + 5 + 5 - 10",
			15,
		},
		{
			"2 * 2 * 2 * 2 * 2 * 2",
			64,
		},
		{
			"-50 + 100 - 50",
			0,
		},
		{
			"5 * 2 + 10",
			20,
		},
		{
			"5 + 2 * 10",
			25,
		},
		{
			"50 / 2 * 2 + 10",
			60,
		},
		{
			"2 * (5 + 10)",
			30,
		},
		{
			"3 * 3 * 3 + 10",
			37,
		},
		{
			"3 * (3 * 3) + 10",
			37,
		},
		{
			"(5 + 10 * 2 + 15 / 3) * 2 + -10",
			50,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"true",
			true,
		},
		{
			"false",
			false,
		},
		{
			"false;",
			false,
		},
		{
			"true;",
			true,
		},
		{
			"5 == 5;",
			true,
		},
		{
			"6 == 5;",
			false,
		},
		{
			"5 != 5;",
			false,
		},
		{
			"6 != 5",
			true,
		},
		{
			"6 > 5;",
			true,
		},
		{
			"5 > 5;",
			false,
		},
		{
			"5 < 6",
			true,
		},
		{
			"8 < 6",
			false,
		},
		{
			"true==true",
			true,
		},
		{
			"true==false",
			false,
		},
		{
			"true!=false;",
			true,
		},
		{
			"true!=true",
			false,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testLiteralObject(t, evaluated, tt.expected)
	}
}

func TestBangPrefixExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{
			"!true",
			false,
		},
		{
			"!false",
			true,
		},
		{
			"!false;",
			true,
		},
		{
			"!true;",
			false,
		},
		{
			"!!true",
			true,
		},
		{
			"!!false",
			false,
		},
		{
			"!5",
			false,
		},
		{
			"!!5",
			true,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testLiteralObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}

func testLiteralObject(t *testing.T, input object.Object, expected any) bool {
	switch v := expected.(type) {
	case bool:
		return testBooleanObject(t, input, v)
	case int64:
		return testIntegerObject(t, input, v)
	}

	t.Errorf("Unhandled expected type for testLiteralObject. Got %T", expected)
	return false
}

func testIntegerObject(t *testing.T, input object.Object, expected int64) bool {
	// Type check
	if input == nil {
		t.Errorf("Expected Integer, got nil")
		return false
	}
	objType := input.Type()
	if objType != object.INTEGER_OBJ {
		t.Errorf("Expected %q, got %q for input: %q", object.INTEGER_OBJ, input.Type(), input.Inspect())
		return false
	}

	if input.Inspect() != fmt.Sprintf("%d", expected) {
		t.Errorf("Expected %d, got %q", expected, input.Inspect())
		return false
	}

	intObj, ok := input.(*object.Integer)

	if !ok {
		t.Errorf("Failed to cast Object into %s. Got %T", "Integer", input)
		return false
	}

	if intObj.Value != expected {
		t.Errorf("Expected representation to be %d, got %d", expected, intObj.Value)
		return false
	}
	return true
}

func testBooleanObject(t *testing.T, input object.Object, expected bool) bool {
	// Type check
	if input == nil {
		t.Errorf("Expected %s, got nil", "Boolean")
		return false
	}
	objType := input.Type()
	if objType != object.BOOLEAN_OBJ {
		t.Errorf("Expected %q, got %q", object.BOOLEAN_OBJ, input.Type())
		return false
	}

	if input.Inspect() != fmt.Sprintf("%t", expected) {
		t.Errorf("Expected %t, got %q", expected, input.Inspect())
		return false
	}

	boolObj, ok := input.(*object.Boolean)

	if !ok {
		t.Errorf("Failed to cast Object into %s. Got %T", "Integer", input)
		return false
	}

	if boolObj.Value != expected {
		t.Errorf("Expected representation to be %t, got %t", expected, boolObj.Value)
		return false
	}
	return true
}
