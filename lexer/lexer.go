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

func (l *Lexer) readString() string {
	l.readChar()
	position := l.position
	for l.ch != '"' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) NextToken() token.Token {
	// Handle the single charcters first
	var tok token.Token
	l.skipWhiteSpace()

	switch l.ch {
	case 0:
		tok = newToken(token.EOF, 0)
	case '=':
		if l.peakAhead() == '=' {
			tok = token.Token{Type: token.EQ, Literal: "=="}
			l.readChar() // This is done to keep position consistent
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '"':
		tok = token.Token{Type: token.STRING, Literal: l.readString()}

	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '!':
		if l.peakAhead() == '=' {
			tok = token.Token{Type: token.NOT_EQ, Literal: "!="}
			l.readChar() // This is done to keep position consistent
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
		} else if isDigit(l.ch) {
			tok = token.Token{Type: token.INT, Literal: l.readNumber()}
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, literal byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(literal)}
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
