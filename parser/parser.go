package parser

import (
	"fmt"

	"github.com/waridh/go-monkey-interpreter/ast"
	"github.com/waridh/go-monkey-interpreter/lexer"
	"github.com/waridh/go-monkey-interpreter/token"
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	errors    []string
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Step the lexer properly, and fill up the parser
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram is the method that will return an AST from the input lexer
//
// TODO: Iteration occuring here is not very clean
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) isCurToken(s token.TokenType) bool {
	return p.curToken.Type == s
}

func (p *Parser) isPeekToken(s token.TokenType) bool {
	return p.peekToken.Type == s
}

func (p *Parser) peekStep(s token.TokenType) bool {
	r := p.isPeekToken(s)
	if r {
		p.nextToken()
	} else {
		msg := fmt.Sprintf("Parser expected %s but got %s", s, p.peekToken.Type)
		p.errors = append(p.errors, msg)
	}
	return r
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	default:
		return nil
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.peekStep(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.peekStep(token.ASSIGN) {
		return nil
	}

	// This is where expression parsing would happen
	// TODO: Implement the next expression ast node
	for !p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}
