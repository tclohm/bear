package ast

import (
	"bear/token"
	"bytes"
	"strings"
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

type Boolean struct {
	Token token.Token
	Value bool
}

func (self *Boolean) expressionNode() {}
func (self *Boolean) TokenLiteral() string { return self.Token.Literal }
func (self *Boolean) String() string { return self.Token.Literal }

type IfExpression struct {
	Token 		token.Token
	Condition 	Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (self *IfExpression) expressionNode() {}
func (self *IfExpression) TokenLiteral() string { return self.Token.Literal }
func (self *IfExpression) String() string { 
	var out bytes.Buffer

	out.WriteString("if")
	out.WriteString(self.Condition.String())
	out.WriteString(" ")
	out.WriteString(self.Consequence.String())

	if self.Alternative != nil {
		out.WriteString("else ")
		out.WriteString(self.Alternative.String())
	}

	return out.String()

}

type BlockStatement struct {
	Token token.Token
	Statements []Statement
}

func (self *BlockStatement) statementNode() {}
func (self *BlockStatement) TokenLiteral() string { return self.Token.Literal }
func (self *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range self.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}


type FunctionLiteral struct {
	Token token.Token
	Parameters []*Identifier
	Body *BlockStatement
}

func (self *FunctionLiteral) expressionNode() {}
func (self *FunctionLiteral) TokenLiteral() string { return self.Token.Literal }
func (self *FunctionLiteral) String() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range self.Parameters {
		params = append(params, p.String())
	}

	out.WriteString(self.TokenLiteral())
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") ")
	out.WriteString(self.Body.String())

	return out.String()
}

type CallExpression struct {
	Token token.Token
	Function Expression
	Arguments []Expression
}

func (self *CallExpression) expressionNode() {}
func (self *CallExpression) TokenLiteral() string { return self.Token.Literal }
func (self *CallExpression) String() string {
	var out bytes.Buffer

	args := []string{}
	for _, a := range self.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(self.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (self *StringLiteral) expressionNode() {}
func (self *StringLiteral) TokenLiteral() string { return self.Token.Literal }
func (self *StringLiteral) String() string { return self.Token.Literal }

type ArrayLiteral struct {
	Token token.Token
	Elements []Expression
}

func (self *ArrayLiteral) expressionNode() {}
func (self *ArrayLiteral) TokenLiteral() string { return self.Token.Literal }
func (self *ArrayLiteral) String() string {
	var out bytes.Buffer

	elements := []string{}
	for _, element := range self.Elements {
		elements = append(elements, element.String())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}