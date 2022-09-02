package compiler

import (
	"fmt"
	"bear/ast"
	"bear/code"
	"bear/object"
)

type Compiler struct {
	instructions code.Instructions
	constants  	 []object.Object
}

func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants: 	  []object.Object{},
	}
}

func (this *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := this.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := this.Compile(node.Expression)
		if err != nil {
			return err
		}
		this.emit(code.OpPop)
	case *ast.InfixExpression:
		err := this.Compile(node.Left)
		if err != nil {
			return err
		}

		err = this.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			this.emit(code.OpAdd)
		case "-":
			this.emit(code.OpSub)
		case "*":
			this.emit(code.OpMul)
		case "/":
			this.emit(code.OpDiv)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		this.emit(code.OpConstant, this.addConstant(integer))
	}

	return nil
}

func (this *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: this.instructions,
		Constants: 	  this.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants 	 []object.Object
}

func (this *Compiler) addConstant(obj object.Object) int {
	this.constants = append(this.constants, obj)
	return len(this.constants) - 1
}

func (this *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := this.addInstruction(ins)
	return pos
}

func (this *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(this.instructions)
	this.instructions = append(this.instructions, ins...)
	return posNewInstruction
}