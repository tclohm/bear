package parser

import (
	"bear/ast"
	"bear/lexer"
	"bear/token"
	"fmt"
)

type Parser struct {
	lex 		*lexer.Lexer
	errors 		[]string
	curToken 	token.Token
	peekToken 	token.Token
}

func New(self *lexer.Lexer) *Parser {
	p := &Parser{lex: self, errors: []string{}}
	// read two tokens, so curToken and peek are both set
	p.nextToken()

	p.nextToken()

	return p
}

func (self *Parser) Errors() []string {
	return self.errors
}

func (self *Parser) nextToken() {
	self.curToken = self.peekToken
	self.peekToken = self.lex.NextToken()
}

func (self *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for self.curToken.Type != token.EOF {
		stmt := self.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		self.nextToken()
	}

	return program
}

func (self *Parser) parseStatement() ast.Statement {
	switch self.curToken.Type {
	case token.LET:
		return self.parseLetStatement()
	case token.RETURN:
		return self.parseReturnStatement()
	default:
		return nil
	}
}

func (self *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: self.curToken}
	
	if !self.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: self.curToken, Value: self.curToken.Literal}

	if !self.expectPeek(token.ASSIGN) { return nil }

	// TODO -- skipping the expression until we encounter semicolon
	for !self.curTokenIs(token.SEMICOLON) { self.nextToken() }

	return stmt
}

func (self *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: self.curToken}
	
	self.nextToken()

	// TODO: -- sking the expressions until we encounter semicolon
	for !self.curTokenIs(token.SEMICOLON) { self.nextToken() }

	return stmt
}

func (self *Parser) curTokenIs(tt token.TokenType) bool {
	return self.curToken.Type == tt
}

func (self *Parser) peekTokenIs(tt token.TokenType) bool {
	return self.peekToken.Type == tt
}

func (self *Parser) expectPeek(tt token.TokenType) bool {
	if self.peekTokenIs(tt) {
		self.nextToken()
		return true
	} else {
		self.peekError(tt)
		return false
	}
}

func (self *Parser) peekError(tt token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead", tt, self.peekToken.Type)
	self.errors = append(self.errors, msg)
}