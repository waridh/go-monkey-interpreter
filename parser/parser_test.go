package parser

import (
	"fmt"
	"testing"

	"github.com/waridh/go-monkey-interpreter/ast"
	"github.com/waridh/go-monkey-interpreter/lexer"
)

type infixTest struct {
	left     any
	operator string
	right    any
}

type prefixTest struct {
	operator string
	right    any
}

type pairTest struct {
	key   any
	value any
}

func checkParserErrors(t *testing.T, p *Parser) {
	es := p.Errors()

	if len(es) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(es))
	for _, msg := range es {
		t.Errorf("parser error msg: %s", msg)
	}
	t.FailNow()
}

func getProgram(t *testing.T, input string) *ast.Program {
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	} else {
		return program
	}
	return nil
}

func TestLetStatements(t *testing.T) {
	expected := []struct {
		input              string
		expectedIdentifier string
		expectedValue      any
	}{
		{"let x = 5;", "x", 5},
		{"let y = 10;", "y", 10},
		{"let foobar = 838383;", "foobar", 838383},
	}

	for _, tt := range expected {
		program := getProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("ParseProgram did not return expected number of items. Expected 1, got %d", len(program.Statements))
		}

		stmt := program.Statements[0]

		if !testLetStatements(t, stmt, tt.expectedIdentifier, tt.expectedValue) {
			return
		}
	}
}

func testLetStatements(t *testing.T, s ast.Statement, iden string, value any) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("TokenLiteral of let statement != let")
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != iden {
		t.Errorf("letStmt.Name.Value != %s, got %s", iden, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != iden {
		t.Errorf("s.Name != %s, got %s", iden, letStmt.Name)
		return false
	}

	if !testExpression(t, letStmt.Value, value) {
		return false
	}

	return true
}

// TestReturnStatements does the basic checks for the functionalities of the
// return statement parsing
func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input string
		value any
	}{
		{
			"return 5;",
			5,
		},
		{
			"return 10;",
			10,
		},
		{
			"return foobar;",
			"foobar",
		},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("ParseProgram did not return expected number of items. Expected %d, got %d", 1, len(program.Statements))
		}

		returnStmt, ok := program.Statements[0].(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Expected ReturnStatement, but got %T", program.Statements[0])
		}

		toklit := returnStmt.TokenLiteral()
		if toklit != "return" {
			t.Errorf("Expected TokenLiteral to be return, instead got %s", toklit)
		}

		if !testExpression(t, returnStmt.ReturnValue, tt.value) {
			return
		}
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	program := getProgram(t, input)

	expectedLen := 1
	if len(program.Statements) != expectedLen {
		t.Fatalf("ParseProgram did not return expected number of items. Expected %d, got %d", expectedLen, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Errorf("Could not get ExpressionStatement, got %T", program.Statements[0])
	}

	if !testIntegerLiteralExpression(t, stmt.Expression, 5) {
		return
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := []struct {
		input    string
		expected string
	}{
		{
			"foobar;",
			"foobar",
		},
		{
			"barfoo",
			"barfoo",
		},
	}
	for _, tt := range input {
		program := getProgram(t, tt.input)

		expectedLen := 1
		if len(program.Statements) != expectedLen {
			t.Fatalf("ParseProgram did not return expected number of items. Expected %d, got %d", expectedLen, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("Could not get ExpressionStatement, got %T", program.Statements[0])
		}

		if !testExpression(t, stmt.Expression, tt.expected) {
			return
		}
	}
}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		right    any
	}{
		{"!5;", "!", 5},
		{"!foobar;", "!", "foobar"},
		{"-15;", "-", 15},
		{"-barfoo;", "-", "barfoo"},
		{"!true;", "!", true},
		{"!false;", "!", false},
		{"!false", "!", false},
	}

	for _, test := range prefixTests {

		program := getProgram(t, test.input)

		if len(program.Statements) != 1 {
			t.Errorf("Expected program.Statement to have %d elements, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("Expected ast.ExpressionStatement but got %T", program.Statements[0])
		}

		if !testPrefixExpression(t, stmt.Expression, test.operator, test.right) {
			return
		}
	}
}

func TestInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{
		{"5 + 5", 5, "+", 5},
		{"5 - 5", 5, "-", 5},
		{"5 * 5", 5, "*", 5},
		{"5 / 5", 5, "/", 5},
		{"5 > 5", 5, ">", 5},
		{"5 < 5", 5, "<", 5},
		{"5 == 5", 5, "==", 5},
		{"5 != 5", 5, "!=", 5},
		{"true != 5", true, "!=", 5}, // Gibberish
		{"false == false", false, "==", false},
		{"true != false", true, "!=", false},
	}

	for _, test := range infixTests {

		program := getProgram(t, test.input)

		if len(program.Statements) != 1 {
			t.Errorf("Expected program.Statement to have %d elements, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("Expected ast.ExpressionStatement but got %T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, test.leftValue, test.operator, test.rightValue) {
			return
		}

	}
}

func TestStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"testing";`, "testing"},
		{`"hello world!";`, "hello world!"},
	}

	for _, test := range tests {

		program := getProgram(t, test.input)

		if len(program.Statements) != 1 {
			t.Errorf("Expected program.Statement to have %d elements, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("Expected ast.ExpressionStatement but got %T, (%+v)", program.Statements[0], program.Statements)
		}

		if !testStringLiteral(t, stmt.Expression, test.expected) {
			return
		}
	}
}

func TestArrayLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected []any
	}{
		{`[1, 2 * 2, 3 + 3];`, []any{1, infixTest{2, "*", 2}, infixTest{3, "+", 3}}},
	}

	for _, test := range tests {

		program := getProgram(t, test.input)

		if len(program.Statements) != 1 {
			t.Errorf("Expected program.Statement to have %d elements, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("Expected ast.ExpressionStatement but got %T, (%+v)", program.Statements[0], program.Statements)
		}

		if !testExpression(t, stmt.Expression, test.expected) {
			return
		}
	}
}

func TestHashLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected map[any]any
	}{
		{`{foo: 1, 64: 2, true: 3};`, map[any]any{
			"foo": 1,
			64:    2,
			true:  3,
		}},
		{`{1: 1 + 10, 64: 2 * 2, true: 3 / 3};`, map[any]any{
			1:    infixTest{1, "+", 10},
			64:   infixTest{2, "*", 2},
			true: infixTest{3, "/", 3},
		}},
	}

	for _, test := range tests {

		program := getProgram(t, test.input)

		if len(program.Statements) != 1 {
			t.Errorf("Expected program.Statement to have %d elements, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("Expected ast.ExpressionStatement but got %T, (%+v)", program.Statements[0], program.Statements)
		}

		hash, ok := stmt.Expression.(*ast.HashLiteral)

		if !ok {
			t.Errorf("Unable to cast to ast.HashLiteral. Got %T", stmt.Expression)
		}
		for key, value := range hash.Pairs {
			switch k := key.(type) {
			case *ast.IntegerLiteral:
				testExpression(t, value, test.expected[int(k.Value)])
			case *ast.Boolean:
				testExpression(t, value, test.expected[k.Value])
			case *ast.Identifier:
				testExpression(t, value, test.expected[k.Value])
			default:
				t.Fatalf("Unexpected branch")
			}
		}
	}
}

func TestEmptyHashLiteral(t *testing.T) {
	input := "{}"

	program := getProgram(t, input)

	if len(program.Statements) != 1 {
		t.Errorf("Expected program.Statement to have %d elements, got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Errorf("Expected ast.ExpressionStatement but got %T, (%+v)", program.Statements[0], program.Statements)
	}

	hash, ok := stmt.Expression.(*ast.HashLiteral)

	if !ok {
		t.Errorf("Unable to cast to ast.HashLiteral. Got %T", stmt.Expression)
	}

	if len(hash.Pairs) != 0 {
		t.Errorf("Expected 0 elements, got %d", len(hash.Pairs))
	}
}

func testArrayLiteral(t *testing.T, exprs ast.Expression, expected []any) bool {
	array, ok := exprs.(*ast.ArrayLiteral)
	if !ok {
		t.Errorf("Expected %s, got %T", "ast.ArrayLiteral", exprs)
		return false
	}

	if len(array.Elements) != len(expected) {
		t.Errorf("Expected %d elements, got %d. (%+v)", len(expected), len(array.Elements), array.Elements)
		return false
	}

	for i, ele := range array.Elements {
		return testExpression(t, ele, expected[i])
	}

	return true
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"true;",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 2 == true",
			"((3 > 2) == true)",
		},
		{
			"3 < 2 == false",
			"((3 < 2) == false)",
		},
		{
			"(5 + 3) * 3",
			"((5 + 3) * 3)",
		},
		{
			"5 / (5 + 5)",
			"(5 / (5 + 5))",
		},
		{
			"5 + (5 + 5) + 5",
			"((5 + (5 + 5)) + 5)",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(a * c) + d",
			"((a + add((a * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * [1, 2, 3, 4][b * c] * d",
			"((a * ([1, 2, 3, 4][(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * [1, 2][1])",
			"add((a * (b[2])), (b[1]), (2 * ([1, 2][1])))",
		},
	}

	for _, test := range tests {

		program := getProgram(t, test.input)

		actual := program.String()
		if actual != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := []struct {
		input       string
		condition   infixTest
		consequence any
		alternative any
	}{
		{
			"if (x < y) { x }",
			infixTest{"x", "<", "y"},
			"x",
			nil,
		},
		{
			"if (x < y) { x } else { y }",
			infixTest{"x", "<", "y"},
			"x",
			"y",
		},
	}

	for _, tt := range input {
		program := getProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected length of %d, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Unable to cast to ast.ExpressionStatement, got %T", program.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.IfExpression)

		if !ok {
			t.Fatalf("Unable to cast to ast.IfExpression, got %T", stmt.Expression)
		}

		if !testInfixExpression(t, expr.Condition, tt.condition.left, tt.condition.operator, tt.condition.right) {
			return
		}

		consequence, ok := expr.Consequence.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Unable to cast to ast.ExpressionStatement, got %T", expr.Consequence.Statements[0])
		}

		if !testExpression(t, consequence.Expression, tt.consequence) {
			return
		}

		alternative := expr.Alternative
		if alternative == nil {
			if tt.alternative != nil {
				t.Fatalf("Expected Alternative to be %t, but got nil", tt.alternative)
			}
		} else {
			alt, ok := alternative.Statements[0].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("Unable to cast to ast.ExpressionStatement, got %T", alternative.Statements[0])
			}

			if !testExpression(t, alt.Expression, tt.alternative) {
				return
			}
		}

	}
}

func TestFunctionLiteralExpression(t *testing.T) {
	input := `fn(x, y) {x + y};`
	program := getProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected length of %d, got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Unable to cast to ast.ExpressionStatement, got %T", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("Unable to cast to ast.FunctionLiteral, got %T", stmt.Expression)
	}

	if len(expr.Parameter) != 2 {
		t.Fatalf("Expected %d parameters, but got %d", 2, len(expr.Parameter))
	}

	if !testExpression(t, expr.Parameter[0], "x") {
		return
	}
	if !testExpression(t, expr.Parameter[1], "y") {
		return
	}

	body, ok := expr.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Expected ast.ExpressionStatement, but got %T", expr.Body.Statements[0])
	}

	if !testInfixExpression(t, body.Expression, "x", "+", "y") {
		return
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{
			"fn(){};",
			[]string{},
		},
		{
			"fn(a){};",
			[]string{"a"},
		},
		{
			"fn(a, b){};",
			[]string{"a", "b"},
		},
		{
			"fn(a, b, c){};",
			[]string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.input)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected %q, got %q", "ExpressionStatement", program.Statements[0])
		}

		fn, ok := stmt.Expression.(*ast.FunctionLiteral)

		if !ok {
			t.Fatalf("Expected %q, got %q", "ast.FunctionLiteral", stmt.Expression)
		}

		if len(fn.Parameter) != len(tt.expectedParams) {
			t.Errorf("Expected %d parameters, got %d", len(tt.expectedParams), len(fn.Parameter))
		}

		for i, exp := range tt.expectedParams {
			testExpression(t, fn.Parameter[i], exp)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5);`
	program := getProgram(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("Expected length of %d, got %d", 1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Unable to cast to ast.ExpressionStatement, got %T", program.Statements[0])
	}

	expr, ok := stmt.Expression.(*ast.CallExpression)

	if !ok {
		t.Fatalf(
			"Unable to cast to %s, got %T",
			"ast.CallExpression",
			stmt.Expression,
		)
	}

	if !testIdentifierExpression(t, expr.Function, "add") {
		return
	}

	if len(expr.Arguments) != 3 {
		t.Fatalf("Expected %d arguments, but got %d", 3, len(expr.Arguments))
	}

	if !testExpression(t, expr.Arguments[0], 1) {
		return
	}

	if !testInfixExpression(t, expr.Arguments[1], 2, "*", 3) {
		return
	}

	if !testInfixExpression(t, expr.Arguments[2], 4, "+", 5) {
		return
	}
}

func TestCallExpressionArgumentParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{
			"add();",
			[]string{},
		},
		{
			"add(a);",
			[]string{"a"},
		},
		{
			"add(a, b);",
			[]string{"a", "b"},
		},
		{
			"add(a, b, c);",
			[]string{"a", "b", "c"},
		},
		{
			"add(1 + 1, b, c);",
			[]string{"(1 + 1)", "b", "c"},
		},
		{
			"add(1 + 1, a + b * c, c);",
			[]string{"(1 + 1)", "(a + (b * c))", "c"},
		},
	}

	for _, tt := range tests {
		program := getProgram(t, tt.input)

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Expected %q, got %q", "ExpressionStatement", program.Statements[0])
		}

		fn, ok := stmt.Expression.(*ast.CallExpression)

		if !ok {
			t.Fatalf("Expected %q, got %q", "ast.CallExpression", stmt.Expression)
		}

		if len(fn.Arguments) != len(tt.expectedParams) {
			t.Errorf("Expected %d parameters, got %d", len(tt.expectedParams), len(fn.Arguments))
		}

		for i, exp := range tt.expectedParams {
			testExpressionString(t, fn.Arguments[i], exp)
		}
	}
}

func TestIndexExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{
		{
			"myArray[1 + 1]",
			infixTest{
				left:     "myArray",
				operator: "[",
				right: infixTest{
					left:     1,
					operator: "+",
					right:    1,
				},
			},
		},
	}
	for _, tt := range tests {
		program := getProgram(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("Expected length of %d, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Fatalf("Unable to cast to ast.ExpressionStatement, got %T", program.Statements[0])
		}

		if !testExpression(t, stmt.Expression, tt.expected) {
			return
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	input := []struct {
		input    string
		expected bool
	}{
		{
			"true;",
			true,
		},
		{
			"false;",
			false,
		},
		{
			"true",
			true,
		},
		{
			"false",
			false,
		},
	}
	for _, tt := range input {
		program := getProgram(t, tt.input)

		expectedLen := 1
		if len(program.Statements) != expectedLen {
			t.Fatalf("ParseProgram did not return expected number of items. Expected %d, got %d", expectedLen, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("Could not get ExpressionStatement, got %T", program.Statements[0])
		}

		if !testExpression(t, stmt.Expression, tt.expected) {
			return
		}
	}
}

func testStringLiteral(t *testing.T, expr ast.Expression, expected string) bool {
	str, ok := expr.(*ast.StringLiteral)
	if !ok {
		t.Errorf("Unable to cast to %s. Got %T. (%+v)", "ast.StringLiteral", expr, expr)
		return false
	}

	if str.Value != expected {
		t.Errorf("Expected %q, got %q", expected, str.Value)
		return false
	}
	return true
}

func testNilExpression(t *testing.T, expr ast.Expression) bool {
	if expr == nil {
		return true
	} else {
		t.Errorf("Expected nil, got %T", expr)
		return false
	}
}

func testExpression(t *testing.T, expr ast.Expression, expected any) bool {
	if expected == nil {
		return testNilExpression(t, expr)
	}
	switch v := expected.(type) {
	case int:
		return testIntegerLiteralExpression(t, expr, int64(v))
	case int64:
		return testIntegerLiteralExpression(t, expr, v)
	case string:
		return testIdentifierExpression(t, expr, v)
	case bool:
		return testBooleanExpression(t, expr, v)
	case infixTest:
		return testInfixExpression(t, expr, v.left, v.operator, v.right)
	case prefixTest:
		return testPrefixExpression(t, expr, v.operator, v.right)
	case []any:
		switch e := expr.(type) {
		case *ast.ArrayLiteral:
			return testArrayLiteral(t, e, v)
		}
	}

	t.Errorf("Unsupported expected type, got expr: %T, and expected: %T", expr, expected)
	return false
}

func testInfixExpression(t *testing.T, expr ast.Expression, left any, operator string, right any) bool {
	switch op := expr.(type) {
	case *ast.InfixExpression:
		if !testExpression(t, op.Left, left) {
			return false
		}

		if op.Operator != operator {
			t.Errorf("Operator does not match. Expected %q, got %q", operator, op.Operator)
			return false
		}

		if !testExpression(t, op.Right, right) {
			return false
		}
	case *ast.IndexExpression:
		if !testExpression(t, op.Left, left) {
			return false
		}

		if operator != "[" {
			t.Errorf("Operator does not match. Expected %q, got %q", operator, "[")
			return false
		}

		if !testExpression(t, op.Index, right) {
			return false
		}
	default:
		t.Errorf("Unable to cast Expression to InfixExpression. Got %T", expr)
		return false
	}

	return true
}

func testPrefixExpression(t *testing.T, expr ast.Expression, operator string, right interface{}) bool {
	op, ok := expr.(*ast.PrefixExpression)

	if !ok {
		t.Errorf("Unable to cast Expression to %q. Got %T", "PrefixExpression", expr)
		return false
	}

	if op.Operator != operator {
		t.Errorf("Operator does not match. Expected %q, got %q", operator, op.Operator)
		return false
	}

	if !testExpression(t, op.Right, right) {
		return false
	}

	return true
}

func testIntegerLiteralExpression(t *testing.T, expr ast.Expression, expected int64) bool {
	lit, ok := expr.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("Expression not IntegerLiteral, got %T", expr)
		return false
	}
	if lit.Value != expected {
		t.Errorf("IntegerLiteral does not hold the expected value. Expected: %d, got %d", expected, lit.Value)
		return false
	}
	if lit.TokenLiteral() != fmt.Sprintf("%d", expected) {
		t.Errorf("TokenLiteral mismatch, expected %d, got %s", expected, lit.TokenLiteral())
	}

	return true
}

func testIdentifierExpression(t *testing.T, expr ast.Expression, expected string) bool {
	id, ok := expr.(*ast.Identifier)

	if !ok {
		t.Errorf("Expression not %q, got %T", "Identifier", expr)
		return false
	}

	if id.Value != expected {
		t.Errorf("Identifier did not have expected value. Expected: %q, got %q", expected, id.Value)
		return false
	}

	if id.TokenLiteral() != expected {
		t.Errorf("TokenLiteral mismatch, expected %q, got %q", expected, id.TokenLiteral())
		return false
	}
	return true
}

func testBooleanExpression(t *testing.T, expr ast.Expression, expected bool) bool {
	bo, ok := expr.(*ast.Boolean)

	if !ok {
		t.Errorf("Expression not %q, got %T", "Identifier", expr)
		return false
	}

	if bo.Value != expected {
		t.Errorf("Identifier did not have expected value. Expected: %t, got %t", expected, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", expected) {
		t.Errorf("TokenLiteral mismatch, expected %t, got %q", expected, bo.TokenLiteral())
		return false
	}
	return true
}

func testExpressionString(t *testing.T, expr ast.Expression, expected string) bool {
	if expr.String() != expected {
		t.Errorf("Expected %q, but got %q", expected, expr.String())
		return false
	}
	return true
}
