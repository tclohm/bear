package parser

import (
	"bear/ast"
	"bear/lexer"
	"bear/token"
	"fmt"
	"strconv"
)

const (
	_ int = iota
	LOWEST
	EQUALS			// ==
	LESSGREATER		// > or <
	SUM 			// +
	PRODUCT 		// *
	PREFIX 			// -X or !X
	CALL 			// myFunction(X)
)

var precedences = map[token.TokenType]int{
	token.EQ: 		EQUALS,
	token.NOT_EQ: 	EQUALS,
	token.LT:		LESSGREATER,
	token.GT:		LESSGREATER,
	token.PLUS: 	SUM,
	token.MINUS: 	SUM,
	token.SLASH: 	PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn func(ast.Expression) ast.Expression
)

type Parser struct {
	lex 			*lexer.Lexer
	errors 			[]string
	curToken 		token.Token
	peekToken 		token.Token

	prefixParseFns 	map[token.TokenType]prefixParseFn
	infixParseFns	map[token.TokenType]infixParseFn
}

func New(self *lexer.Lexer) *Parser {
	p := &Parser{lex: self, errors: []string{}}
	// read two tokens, so curToken and peek are both set
	p.nextToken()

	p.nextToken()

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.TRUE, p.parseBoolean)
	p.registerPrefix(token.FALSE, p.parseBoolean)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.IF, p.parseIfExpression)
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.PLUS, p.parseInfixExpression) 
	p.registerInfix(token.MINUS, p.parseInfixExpression) 
	p.registerInfix(token.SLASH, p.parseInfixExpression) 
	p.registerInfix(token.ASTERISK, p.parseInfixExpression) 
	p.registerInfix(token.EQ, p.parseInfixExpression) 
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression) 
	p.registerInfix(token.LT, p.parseInfixExpression) 
	p.registerInfix(token.GT, p.parseInfixExpression)

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
		return self.parseExpressionStatement()
	}
}

func (self *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: self.curToken}
	
	if !self.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: self.curToken, Value: self.curToken.Literal}

	if !self.expectPeek(token.ASSIGN) { return nil }

	self.nextToken()

	stmt.Value = self.parseExpression(LOWEST)

	if self.peekTokenIs(token.SEMICOLON) { self.nextToken() }

	return stmt
}

func (self *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: self.curToken}
	
	self.nextToken()

	stmt.ReturnValue = self.parseExpression(LOWEST)

	if self.peekTokenIs(token.SEMICOLON) { self.nextToken() }

	return stmt
}

func (self *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: self.curToken}

	stmt.Expression = self.parseExpression(LOWEST)

	if self.peekTokenIs(token.SEMICOLON) { self.nextToken() }

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

func (self *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	self.prefixParseFns[tokenType] = fn
}

func (self *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	self.infixParseFns[tokenType] = fn
}

func (self *Parser) parseExpression(precedence int) ast.Expression {
	prefix := self.prefixParseFns[self.curToken.Type]
	if prefix == nil { 
		self.noPrefixParseFnError(self.curToken.Type)
		return nil 
	}
	leftExp := prefix()

	for !self.peekTokenIs(token.SEMICOLON) && precedence < self.peekPrecedence() {
		infix := self.infixParseFns[self.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		self.nextToken()

		leftExp = infix(leftExp)
	}
	return leftExp
}

func (self *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: self.curToken, Value: self.curToken.Literal}
}

func (self *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: self.curToken}

	value, err := strconv.ParseInt(self.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as integer", self.curToken.Literal)
		self.errors = append(self.errors, msg)
		return nil
	}

	lit.Value = value

	return lit
}

func (self *Parser) noPrefixParseFnError(tokenType token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function for %s found", tokenType)
	self.errors = append(self.errors, msg)
}

func (self *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token: self.curToken,
		Operator: self.curToken.Literal,
	}

	self.nextToken()

	expression.Right = self.parseExpression(PREFIX)

	return expression
}

func (self *Parser) peekPrecedence() int {
	if peeked, ok := precedences[self.peekToken.Type]; ok {
		return peeked
	}
	return LOWEST
}

func (self *Parser) curPrecedence() int {
	if current, ok := precedences[self.curToken.Type]; ok {
		return current
	}
	return LOWEST
}

func (self *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token: 		self.curToken,
		Operator: 	self.curToken.Literal,
		Left: 		left,
	}

	precedence := self.curPrecedence()
	self.nextToken()
	expression.Right = self.parseExpression(precedence)

	return expression
}

func (self *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{Token: self.curToken, Value: self.curTokenIs(token.TRUE)}
}

func (self *Parser) parseGroupedExpression() ast.Expression {
	self.nextToken()

	expression := self.parseExpression(LOWEST)

	if !self.expectPeek(token.RPAREN) {
		return nil
	}

	return expression

}

func (self *Parser) parseIfExpression() ast.Expression {
	expression := &ast.IfExpression{Token: self.curToken}

	if !self.expectPeek(token.LPAREN) { return nil }

	self.nextToken()

	expression.Condition = self.parseExpression(LOWEST)

	if !self.expectPeek(token.RPAREN) { return nil }
	if !self.expectPeek(token.LBRACE) { return nil }

	expression.Consequence = self.parseBlockStatement()

	if self.peekTokenIs(token.ELSE) {
		self.nextToken()

		if !self.expectPeek(token.LBRACE) { return nil }

		expression.Alternative = self.parseBlockStatement()
	}

	return expression
}

func (self *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: self.curToken}
	block.Statements = []ast.Statement{}

	self.nextToken()

	for !self.curTokenIs(token.RBRACE) && !self.curTokenIs(token.EOF) {
		stmt := self.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		self.nextToken()
	}

	return block
}

func (self *Parser) parseFunctionLiteral() ast.Expression {
	literal := &ast.FunctionLiteral{Token: self.curToken}

	if !self.expectPeek(token.LPAREN) { return nil }

	literal.Parameters = self.parseFunctionParameters()

	if !self.expectPeek(token.LBRACE) { return nil }

	literal.Body = self.parseBlockStatement()

	return literal
}

func (self *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if self.peekTokenIs(token.RPAREN) {
		self.nextToken()
		return identifiers
	}

	self.nextToken()

	identifier := &ast.Identifier{Token: self.curToken, Value: self.curToken.Literal}
	identifiers = append(identifiers, identifier)

	for self.peekTokenIs(token.COMMA) {
		self.nextToken()
		self.nextToken()
		identifier := &ast.Identifier{Token: self.curToken, Value: self.curToken.Literal}
		identifiers = append(identifiers, identifier)
	}

	if !self.expectPeek(token.RPAREN) { return nil }

	return identifiers
}

func (self *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	expression := &ast.CallExpression{Token: self.curToken, Function: function}
	expression.Arguments = self.parseExpressionList(token.RPAREN)
	return expression
}

func (self *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	if self.peekTokenIs(token.RPAREN) {
		self.nextToken()
		return args
	}

	self.nextToken()
	args = append(args, self.parseExpression(LOWEST))

	for self.peekTokenIs(token.COMMA) {
		self.nextToken()
		self.nextToken()
		args = append(args, self.parseExpression(LOWEST))
	}

	if !self.expectPeek(token.RPAREN) { return nil }

	return args
}

func (self *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: self.curToken, Value: self.curToken.Literal}
}

func (self *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: self.curToken}
	array.Elements = self.parseExpressionList(token.RBRACKET)
	return array
}

func (self *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if self.peekTokenIs(end) {
		self.nextToken()
		return list
	}

	self.nextToken()
	list = append(list, self.parseExpression(LOWEST))

	for self.peekTokenIs(token.COMMA) {
		self.nextToken()
		self.nextToken()
		list = append(list, self.parseExpression(LOWEST))
	}

	if !self.expectPeek(end) {
		return nil
	}

	return list
}