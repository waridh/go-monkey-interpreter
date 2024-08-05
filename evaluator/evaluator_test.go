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

func TestStringExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			`"Hello World!";`,
			"Hello World!",
		},
		{
			`"Hello " + "World!";`,
			"Hello World!",
		},
		{
			`"Hello" + " " + "World!";`,
			"Hello World!",
		},
		{
			`"Hello" == "World!";`,
			false,
		},
		{
			`"Hello" != "World!";`,
			true,
		},
		{
			`"Hello" == "Hello";`,
			true,
		},
		{
			`"Hello" != "Hello";`,
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

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			"if (true) { 10 }",
			10,
		},
		{
			"if (false) { 10 } else { 11 }",
			11,
		},
		{
			"if (false) { 10 }",
			nil,
		},
		{
			"if (true) { 10 } else { 11 }",
			10,
		},
		{
			"if (1) { 10 }",
			10,
		},
		{
			"if (1==1) { 10 }",
			10,
		},
		{
			"if (1 == 2) { 10 }",
			nil,
		},
		{
			"if (1 < 4) { 10 }",
			10,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testLiteralObject(t, evaluated, tt.expected)
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			"return 10;",
			10,
		},
		{
			"return 10; 9;",
			10,
		},
		{
			"return 2 * 5; 9;",
			10,
		},
		{
			"9; return 2 * 5; 9;",
			10,
		},
		{
			"return;",
			nil,
		},
		{
			"return; 70;",
			nil,
		},
		{
			"return 1 == 1;",
			true,
		},
		{
			"return 1 < 1;",
			false,
		},
		{
			`
      if (10 > 1) {
      if (10 > 1) {
      return 10;
      }
      return 1;
      }
      `,
			10,
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testLiteralObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"5 + true;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"type mismatch: INTEGER + BOOLEAN",
		},
		{
			"-true; 6;",
			"unknown operator: -BOOLEAN",
		},
		{
			"true + false;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5;",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) {true + false};",
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			`
    if (10 > 1) {
      if (10 > 1) {
        return true + false;
      }
      return 1;
    };`,
			"unknown operator: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"identity not found: foobar",
		},
		{
			`"Hello" - "World"`,
			"unknown operator: STRING - STRING",
		},
	}
	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testErrorObject(t, evaluated, tt.expected)
	}
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			"let a = 5; a;",
			5,
		},
		{
			"let a = 5 < 10; a;",
			true,
		},
		{
			"let a = 7 - 12; a;",
			-5,
		},
		{
			"let a = 5; let b = a * 3; b;",
			15,
		},
		{
			"let a = 5; let b = a * 2; let c = a + b; c;",
			15,
		},
	}
	for _, tt := range tests {
		testLiteralObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := `fn(x) {x + 2;};`

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)

	if !ok {
		t.Errorf("Expected %s, got %T. (%+v)", "object.Function", evaluated, evaluated)
	}

	numParams := 1
	if len(fn.Parameter) != numParams {
		t.Errorf("Expected %d parameter, but got %d. (%+v)", numParams, len(fn.Parameter), fn)
	}

	paramName := "x"
	if fn.Parameter[0].String() != paramName {
		t.Errorf("Expected %q as parameter, but got %q. (%+v)", paramName, fn.Parameter[0].String(), fn)
	}

	expectedBody := "(x + 2)"
	if fn.Body.String() != expectedBody {
		t.Errorf("Expected body to be %s, got %s. (%+v)", expectedBody, fn.Body.String(), fn)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			"let identity = fn(x) {x;}; identity(5);",
			5,
		},
		{
			"let identity = fn(x) {return x;}; identity(5);",
			5,
		},
		{
			"let double = fn(x) {return x * 2;}; double(5);",
			10,
		},
		{
			"let add = fn(x, y) {return x + y;}; add(2, 3);",
			5,
		},
		{
			"let add = fn(x, y) {return x + y;}; add(2, add(2,1));",
			5,
		},
		{
			"fn(x, y) {return x * y;}(2, 3);",
			6,
		},
		{
			"let not = fn(x) {!x;}; not(5);",
			false,
		},
		{
			"let not = fn(x) {!x;}; not(false);",
			true,
		},
		{
			"let pos = fn(x) { return x > 0;}; pos(5);",
			true,
		},
	}
	for _, tt := range tests {
		testLiteralObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			`
      let newAdder = fn(x) {
      return fn(y) {x + y;};
      };

      let addThree = newAdder(3);
      addThree(2);
      `,
			5,
		},
		{
			`
      let newAdder = fn(x) {
      return fn(y) {x + y;};
      };

      let apply = fn(x, y) {
      return x(y);
      }
      apply(apply(newAdder, 2), 3);
      `,
			5,
		},
	}
	for _, tt := range tests {
		testLiteralObject(t, testEval(tt.input), tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testLiteralObject(t *testing.T, input object.Object, expected any) bool {
	switch v := expected.(type) {
	case bool:
		return testBooleanObject(t, input, v)
	case int:
		return testIntegerObject(t, input, int64(v))
	case int64:
		return testIntegerObject(t, input, v)
	case string:
		return testStringObject(t, input, v)
	case nil:
		return testNullObject(t, input)
	}

	t.Errorf("Unhandled expected type for testLiteralObject. Got %T", expected)
	return false
}

func testErrorObject(t *testing.T, input object.Object, expected string) bool {
	err, ok := input.(*object.Error)
	if !ok {
		t.Errorf("Was not able to cast output into error object. Got %T. (%+v)", input, input)
		return false
	}

	if err.Message != expected {
		t.Errorf("Expected: %q, got %q. (%+v)", expected, err.Message, err)
		return false
	}

	return true
}

func testNullObject(t *testing.T, input object.Object) bool {
	if input != NULL {
		t.Errorf("Expected NULL, but got %T (%+v)", input, input)
		return false
	} else {
		return true
	}
}

func testStringObject(t *testing.T, input object.Object, expected string) bool {
	// Type check
	if input == nil {
		t.Errorf("Expected String, got nil")
		return false
	}
	objType := input.Type()
	if objType != object.STRING_OBJ {
		t.Errorf("Expected %q, got %q for input: %q", object.INTEGER_OBJ, input.Type(), input.Inspect())
		return false
	}

	if input.Inspect() != expected {
		t.Errorf("Expected %q, got %q", expected, input.Inspect())
		return false
	}

	str, ok := input.(*object.String)

	if !ok {
		t.Errorf("Failed to cast Object into %s. Got %T", "String", input)
		return false
	}

	if str.Value != expected {
		t.Errorf("Expected representation to be %q, got %q", expected, str.Value)
		return false
	}
	return true
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
