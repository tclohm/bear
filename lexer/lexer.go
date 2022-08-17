package lexer

import "bear/token"

type Lexer struct {
	input			string
	position 		int // 	current position in input (current char)
	readPosition 	int // 	current reading position (after current char)
	ch 				byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// read next position
func (self *Lexer) readChar() {
	if self.readPosition >= len(self.input) {
		self.ch = 0
	} else {
		self.ch = self.input[self.readPosition]
	}
	self.position = self.readPosition
	self.readPosition += 1
}

func (self *Lexer) NextToken() token.Token {
	var tok token.Token

	self.skipWhiteSpace()

	switch self.ch {
	case '=':
		if self.peekChar() == '=' {
			ch := self.ch
			self.readChar()
			literal := string(ch) + string(self.ch)
			tok = token.Token{Type: token.EQ, Literal: literal}
		} else {
			tok = newToken(token.ASSIGN, self.ch)
		}
	case '+':
		tok = newToken(token.PLUS, self.ch)
	case '-':
		tok = newToken(token.MINUS, self.ch)
	case '!':
		if self.peekChar() == '=' {
			ch := self.ch
			self.readChar()
			literal := string(ch) + string(self.ch)
			tok = token.Token{Type: token.NOT_EQ, Literal: literal}
		} else {
			tok = newToken(token.BANG, self.ch)
		}
	case '/':
		tok = newToken(token.SLASH, self.ch)
	case '*':
		tok = newToken(token.ASTERISK, self.ch)
	case '<':
		tok = newToken(token.LT, self.ch)
	case '>':
		tok = newToken(token.GT, self.ch)
	case ';':
		tok = newToken(token.SEMICOLON, self.ch)
	case '(':
		tok = newToken(token.LPAREN, self.ch)
	case ')':
		tok = newToken(token.RPAREN, self.ch)
	case ',':
		tok = newToken(token.COMMA, self.ch)
	case '{':
		tok = newToken(token.LBRACE, self.ch)
	case '}':
		tok = newToken(token.RBRACE, self.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(self.ch) {
			tok.Literal = self.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(self.ch) {
			tok.Type = token.INT
			tok.Literal = self.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, self.ch)
		}
	}

	self.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func (self *Lexer) readIdentifier() string {
	position := self.position
	for isLetter(self.ch) {
		self.readChar()
	}
	return self.input[position:self.position]
}

func (self *Lexer) skipWhiteSpace() {
	for self.ch == ' ' || self.ch == '\t' || self.ch == '\n' || self.ch == '\r' {
		self.readChar()
	}
}

func (self *Lexer) readNumber() string {
	position := self.position
	for isDigit(self.ch) {
		self.readChar()
	}
	return self.input[position:self.position]
}

func (self *Lexer) peekChar() byte {
	if self.readPosition >= len(self.input) {
		return 0
	} else {
		return self.input[self.readPosition]
	}
}
