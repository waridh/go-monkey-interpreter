package parser

import (
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
		t.Fatalf("ParseProgram did not return expected number of items. Expected 3, got %d", len(program.Statements))
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
