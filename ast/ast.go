package ast

import "github.com/waridh/go-monkey-interpreter/token"

// Node is the base interface for all the nodes of the abstract syntax tree
// in our system
type Node interface {
	TokenLiteral() string
}

// Statement is the interface that lets us identify if a struct is a statement
type Statement interface {
	Node
	statementNode()
}

// Expression interface identifies expression nodes
type Expression interface {
	Node
	expressionNode()
}

// Program is the top level struct that holds all the other nodes
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token token.Token // For the LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

type Identifier struct {
	Token token.Token
	Value string
}

func (id *Identifier) expressionNode() {}
func (id *Identifier) TokenLiteral() string {
	return id.Token.Literal
}
