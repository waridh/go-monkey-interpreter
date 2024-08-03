package parser

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/waridh/go-monkey-interpreter/ast"
	"github.com/waridh/go-monkey-interpreter/lexer"
)

type Test struct {
	expectedIdentifier string
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

func TestLetStatements(t *testing.T) {
	input := `
  let x = 5;
  let y = 10;
  let foobar = 838383;
  `

	expected := []Test{
		{"x"},
		{"y"},
		{"foobar"},
	}

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		var err bytes.Buffer
		for i := range len(program.Statements) {
			fmt.Fprintf(&err, "%s, Type: %T, TokenLiteral: %s, ", program.Statements[i].String(), program.Statements[i], program.Statements[i].TokenLiteral())
		}
		t.Fatalf("ParseProgram did not return expected number of items. Expected 3, got %d: %s", len(program.Statements), err.String())
	}

	for i, tt := range expected {
		stmt := program.Statements[i]
		if !testLetStatements(t, tt.expectedIdentifier, stmt) {
			return
		}
	}
}

func testLetStatements(t *testing.T, expected string, s ast.Statement) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("TokenLiteral of let statement != let")
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not *ast.LetStatement. got=%T", s)
		return false
	}

	if letStmt.Name.Value != expected {
		t.Errorf("letStmt.Name.Value != %s, got %s", expected, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != expected {
		t.Errorf("s.Name != %s, got %s", expected, letStmt.Name)
		return false
	}

	return true
}

// TestReturnStatements does the basic checks for the functionalities of the
// return statement parsing
func TestReturnStatements(t *testing.T) {
	input := `
  return 5;
  return 10;
  return add(15);
  `

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("ParseProgram did not return expected number of items. Expected 3, got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Expected ReturnStatement, but got %T", stmt)
			continue
		}

		toklit := returnStmt.TokenLiteral()
		if toklit != "return" {
			t.Errorf("Expected TokenLiteral to be return, instead got %s", toklit)
		}

	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := `5;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	expectedLen := 1
	if len(program.Statements) != expectedLen {
		t.Fatalf("ParseProgram did not return expected number of items. Expected %d, got %d", expectedLen, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Errorf("Could not get ExpressionStatement, got %T", program.Statements[0])
	}

	intLit, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("Could not get %s, got %T", "IntegerLiteral", program.Statements[0])
	}

	if intLit.Value != 5 {
		t.Errorf("ident.Value expected %d, got %d", 5, intLit.Value)
	}

	if intLit.TokenLiteral() != "5" {
		t.Errorf("TokentLiteral expected %s, got %q", "5", intLit.TokenLiteral())
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	expectedLen := 1
	if len(program.Statements) != expectedLen {
		t.Fatalf("ParseProgram did not return expected number of items. Expected %d, got %d", expectedLen, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Errorf("Could not get ExpressionStatement, got %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Errorf("Could not get %s, got %T", "Identifier", program.Statements[0])
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value expected foobar, got %q", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("TokentLiteral expected %s, got %q", "foobar", ident.TokenLiteral())
	}
}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, test := range prefixTests {
		l := lexer.New(test.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Errorf("Expected program.Statement to have %d elements, got %d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)

		if !ok {
			t.Errorf("Expected ast.ExpressionStatement but got %T", program.Statements[0])
		}

		expr, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("Expected ast.PrefixExpression, got %T", stmt.Expression)
		}

		if expr.Operator != test.operator {
			t.Fatalf("Expected %s, but got %s", test.operator, expr.Operator)
		}

		if !testIntegerLiteralExpression(t, expr.Right, test.integerValue) {
			return
		}
	}
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
