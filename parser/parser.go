package parser

import (
	"bear/ast"
	"bear/lexer"
	"bear/token"
)

type Parser struct {
	lex 		*lexer.Lexer

	curToken 	token.Token
	peekToken 	token.Token
}

func New(self *lexer.Lexer) *Parser {
	p := &Parser{lex: self}
	// read two tokens, so curToken and peek are both set
	p.nextToken()

	p.nextToken()

	return p
}

func (self *Parser) nextToken() {
	self.curToken = peekToken
	self.peekToken = self.lex.NextToken()
}