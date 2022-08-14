package ast

import "bear/token"

type Node interface {
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (self *Program) TokenLiteral() string {
	if len(self.Statements[0].TokenLiteral()) {
		return self.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type LetStatement struct {
	Token 	token.Token // token.Let token
	Name 	*Identifier
	Value 	Expression
}

func (self *LetStatement) statementNode() {}
func (self *LetStatement) TokenLiteral() string { return self.Token.Literal }

type Identifier struct {
	Token token.Token // token.IDENT token
	Value string
}

func (self *Identifier) expressionNode() {}
func (self *Identifier) TokenLiteral() string { return self.Token.Literal }