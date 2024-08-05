package parser

import (
	"fmt"
	"strconv"

	"github.com/waridh/go-monkey-interpreter/ast"
	"github.com/waridh/go-monkey-interpreter/lexer"
	"github.com/waridh/go-monkey-interpreter/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -X or !X
	CALL        // myFunction(x)
)

var precedences = map[token.TokenType]int{
	token.EQ:       EQUALS,
	token.NOT_EQ:   EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.PLUS:     SUM,
	token.MINUS:    SUM,
	token.ASTERISK: PRODUCT,
	token.SLASH:    PRODUCT,
	token.LPAREN:   CALL,
}

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	errors    []string
	peekToken token.Token

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBooleanLiteral)
	p.registerPrefix(token.FALSE, p.parseBooleanLiteral)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)

	// Step the lexer properly, and fill up the parser
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
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

// Token Check

func (p *Parser) isCurToken(s token.TokenType) bool {
	return p.curToken.Type == s
}

func (p *Parser) isPeekToken(s token.TokenType) bool {
	return p.peekToken.Type == s
}

// Returns false if unable to step forward
func (p *Parser) peekStep(s token.TokenType) bool {
	r := p.isPeekToken(s)
	if r {
		p.nextToken()
	} else {
		p.peekError(s)
	}
	return r
}

// Expression functionalities

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) writeError(msg string) {
	p.errors = append(p.errors, msg)
}

func (p *Parser) peekError(s token.TokenType) {
	msg := fmt.Sprintf("Parser expected %s but got %s", s, p.peekToken.Type)
	p.writeError(msg)
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
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

	p.nextToken()

	stmt.Value = p.parseExpression(LOWEST)
	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	stmt.ReturnValue = p.parseExpression(LOWEST)
	// TODO: Implement the expression node type
	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression(LOWEST)

	if p.isPeekToken(token.SEMICOLON) {
		p.nextToken()
	}
	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.isPeekToken(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]

		// Check if the next token is a valid infix operator
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}
	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("Could not parse %q as integer", err)
		p.writeError(msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (p *Parser) parseStringLiteral() ast.Expression {
	lit := &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	return lit
}

func (p *Parser) parseBooleanLiteral() ast.Expression {
	lit := &ast.Boolean{Token: p.curToken, Value: p.isCurToken(token.TRUE)}
	return lit
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	pe := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	pe.Right = p.parseExpression(PREFIX)

	return pe
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	ie := &ast.InfixExpression{Left: left, Operator: p.curToken.Literal, Token: p.curToken}

	precedence := p.curPrecedence()
	p.nextToken()
	ie.Right = p.parseExpression(precedence)

	return ie
}

func (p *Parser) parseIfExpression() ast.Expression {
	expr := &ast.IfExpression{Token: p.curToken}

	if !p.peekStep(token.LPAREN) {
		return nil
	} else {
		p.nextToken()
	}

	expr.Condition = p.parseExpression(LOWEST)

	if !p.peekStep(token.RPAREN) {
		return nil
	}

	if !p.peekStep(token.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStatement()

	if !p.isPeekToken(token.ELSE) {
		expr.Alternative = nil
	} else {
		p.nextToken()
		if !p.peekStep(token.LBRACE) {
			return nil
		}
		expr.Alternative = p.parseBlockStatement()
	}

	return expr
}

func (p *Parser) parseFunctionLiteral() ast.Expression {
	expr := &ast.FunctionLiteral{Token: p.curToken}

	if !p.peekStep(token.LPAREN) {
		return nil
	}

	expr.Parameter = p.parseFunctionParameters()

	if !p.peekStep(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStatement()

	return expr
}

// Collects the expressions that are separated by commas
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	// Handle the case where there is no params
	if p.isPeekToken(token.RPAREN) {
		p.nextToken()
	} else {
		p.nextToken()
		iden := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, iden)

		for p.isPeekToken(token.COMMA) {
			p.nextToken()
			p.nextToken()

			iden := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			identifiers = append(identifiers, iden)
		}

		if !p.peekStep(token.RPAREN) {
			return nil
		}
	}

	return identifiers
}

func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blkstmt := &ast.BlockStatement{Token: p.curToken, Statements: []ast.Statement{}}

	p.nextToken()

	for !p.isCurToken(token.RBRACE) && !p.isCurToken(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			blkstmt.Statements = append(blkstmt.Statements, stmt)
		}
		p.nextToken()
	}

	return blkstmt
}

func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	ce := &ast.CallExpression{Token: p.curToken, Function: fn}
	ce.Arguments = p.parseCallArguments()

	return ce
}

func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if p.isPeekToken(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpression(LOWEST))

	for p.isPeekToken(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.peekStep(token.RPAREN) {
		return nil
	}

	return args
}

// Error handling
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.writeError(msg)
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	expr := p.parseExpression(LOWEST)

	if !p.peekStep(token.RPAREN) {
		return nil
	}

	return expr
}
