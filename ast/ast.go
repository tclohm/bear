package ast

import (
	"bear/token"
	"bytes"
)

type Node interface {
	TokenLiteral() string
	String() string
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
	if len(self.Statements) > 0 {
		return self.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (self *Program) String() string {
	var out bytes.Buffer

	for _, s := range self.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}

type LetStatement struct {
	Token 	token.Token // token.Let token
	Name 	*Identifier
	Value 	Expression
}

func (self *LetStatement) statementNode() {}
func (self *LetStatement) TokenLiteral() string { return self.Token.Literal }
func (self *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(self.TokenLiteral() + " ")
	out.WriteString(self.Name.String())
	out.WriteString(" = ")

	if self.Value != nil {
		out.WriteString(self.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	Token token.Token // token.IDENT token
	Value string
}

func (self *Identifier) expressionNode() {}
func (self *Identifier) TokenLiteral() string { return self.Token.Literal }
func (self *Identifier) String() string { return self.Value }

type ReturnStatement struct {
	Token 		token.Token // return token
	ReturnValue Expression
}

func (self *ReturnStatement) statementNode() {}
func (self *ReturnStatement) TokenLiteral() string { return self.Token.Literal }
func (self *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(self.TokenLiteral() + " ")

	if self.ReturnValue != nil {
		out.WriteString(self.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}

type ExpressionStatement struct {
	Token 		token.Token // first token of the expression
	Expression 	Expression
}

func (self *ExpressionStatement) statementNode() {}
func (self *ExpressionStatement) TokenLiteral() string { return self.Token.Literal }
func (self *ExpressionStatement) String() string {
	if self.Expression != nil {
		return self.Expression.String()
	}
	return ""
	
}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (self *IntegerLiteral) expressionNode() {}
func (self *IntegerLiteral) TokenLiteral() string { return self.Token.Literal }
func (self *IntegerLiteral) String() string { return self.Token.Literal }

type PrefixExpression struct {
	Token 		token.Token
	Operator 	string
	Right 		Expression
}

func (self *PrefixExpression) expressionNode() {}
func (self *PrefixExpression) TokenLiteral() string { return self.Token.Literal }
func (self *PrefixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(self.Operator)
	out.WriteString(self.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token 		token.Token
	Left 		Expression
	Operator 	string
	Right 		Expression
}

func (self *InfixExpression) expressionNode() {}
func (self *InfixExpression) TokenLiteral() string { return self.Token.Literal }
func (self *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(self.Left.String())
	out.WriteString(" " + self.Operator + " ")
	out.WriteString(self.Right.String())
	out.WriteString(")")

	return out.String()
}