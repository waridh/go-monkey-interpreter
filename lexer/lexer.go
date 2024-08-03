package lexer

import (
	"github.com/waridh/go-monkey-interpreter/token"
)

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           byte // This is the character currently being pointed to
}

// New is the base constructor for the Lexer struct
func New(input string) *Lexer {
	l := &Lexer{
		input:        input,
		position:     0,
		readPosition: 0,
	}
	l.readChar() // setup the struct
	return l
}

// readChar mutates the internal state, and updates the character currently
// being pointed to
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// Swap to sentinel value
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l Lexer) peakAhead() byte {
	var ret byte
	if l.readPosition >= len(l.input) {
		ret = 0
	} else {
		ret = l.input[l.readPosition]
	}

	return ret
}

func (l *Lexer) readIdentifier() string {
	position := l.position

	for isLetter(l.ch) {
		l.readChar() // This will increment the position
	}

	l.readPosition -= 1

	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position

	for isDigit(l.ch) {
		l.readChar()
	}
	l.readPosition -= 1

	return l.input[position:l.position]
}

func (l *Lexer) skipWhiteSpace() {
	for isWhiteSpace(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) NextToken() token.Token {
	// Handle the single charcters first
	var tok token.Token
	l.skipWhiteSpace()

	switch l.ch {
	case 0:
		tok = newToken(token.EOF, "")
	case '=':
		if l.peakAhead() == '=' {
			tok = newToken(token.EQ, "==")
			l.readChar() // This is done to keep position consistent
		} else {
			tok = newToken(token.ASSIGN, "=")
		}
	case '+':
		tok = newToken(token.PLUS, "+")
	case '(':
		tok = newToken(token.LPAREN, "(")
	case ')':
		tok = newToken(token.RPAREN, ")")
	case '{':
		tok = newToken(token.LBRACE, "{")
	case '}':
		tok = newToken(token.RBRACE, "}")
	case ',':
		tok = newToken(token.COMMA, ",")
	case ';':
		tok = newToken(token.SEMICOLON, ";")
	case '!':
		if l.peakAhead() == '=' {
			tok = newToken(token.NOT_EQ, "!=")
			l.readChar() // This is done to keep position consistent
		} else {
			tok = newToken(token.BANG, "!")
		}
	case '-':
		tok = newToken(token.MINUS, "-")
	case '/':
		tok = newToken(token.SLASH, "/")
	case '*':
		tok = newToken(token.ASTERISK, "*")
	case '<':
		tok = newToken(token.LT, "<")
	case '>':
		tok = newToken(token.GT, ">")
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
		} else if isDigit(l.ch) {
			tok = newToken(token.INT, l.readNumber())
		} else {
			tok = newToken(token.ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{Type: tokenType, Literal: literal}
}

func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z') || (ch == '_')
}

func isDigit(ch byte) bool {
	return ('0' <= ch) && (ch <= '9')
}

func isWhiteSpace(ch byte) bool {
	return (ch == ' ') || (ch == '\t') || (ch == '\n') || (ch == '\r')
}
