package compiler

import (
	"fmt"
	"bear/ast"
	"bear/code"
	"bear/object"
)


type EmittedInstruction struct {
	Opcode 		code.Opcode
	Position 	int
}

type Compiler struct {
	instructions 		code.Instructions
	constants  	 		[]object.Object

	lastInstruction 	EmittedInstruction
	previousInstruction EmittedInstruction
}

func New() *Compiler {
	return &Compiler{
		instructions: 		 code.Instructions{},
		constants: 	  		 []object.Object{},
		lastInstruction: 	 EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
}

func (self *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := self.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := self.Compile(node.Expression)
		if err != nil {
			return err
		}
		self.emit(code.OpPop)
	case *ast.InfixExpression:
		if node.Operator == "<" {
			err := self.Compile(node.Right)
			if err != nil {
				return err
			}

			err = self.Compile(node.Left)
			if err != nil {
				return err
			}
			self.emit(code.OpGreaterThan)
			return nil
		}

		err := self.Compile(node.Left)
		if err != nil {
			return err
		}

		err = self.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "+":
			self.emit(code.OpAdd)
		case "-":
			self.emit(code.OpSub)
		case "*":
			self.emit(code.OpMul)
		case "/":
			self.emit(code.OpDiv)
		case ">":
			self.emit(code.OpGreaterThan)
		case "==":
			self.emit(code.OpEqual)
		case "!=":
			self.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		self.emit(code.OpConstant, self.addConstant(integer))
	case *ast.Boolean:
		if node.Value {
			self.emit(code.OpTrue)
		} else {
			self.emit(code.OpFalse)
		}
	case *ast.PrefixExpression:
		err := self.Compile(node.Right)
		if err != nil {
			return err
		}

		switch node.Operator {
		case "!":
			self.emit(code.OpBang)
		case "-":
			self.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IfExpression:
		err := self.Compile(node.Condition)
		if err != nil {
			return err
		}

		jumpNotTruthyPos := self.emit(code.OpJumpNotTruthy, 9999)

		err = self.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if self.lastInstructionIsPop() {
			self.removeLastPop()
		}

		afterConsequencePos := len(self.instructions)
		self.changeOperand(jumpNotTruthyPos, afterConsequencePos)
	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := self.Compile(s)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (self *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: self.instructions,
		Constants: 	  self.constants,
	}
}

type Bytecode struct {
	Instructions code.Instructions
	Constants 	 []object.Object
}

func (self *Compiler) addConstant(obj object.Object) int {
	self.constants = append(self.constants, obj)
	return len(self.constants) - 1
}

func (self *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := self.addInstruction(ins)
	self.setLastInstruction(op, pos)
	return pos
}

func (self *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(self.instructions)
	self.instructions = append(self.instructions, ins...)
	return posNewInstruction
}

func (self *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := self.lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}
	self.previousInstruction = previous
	self.lastInstruction = last
}

func (self *Compiler) lastInstructionIsPop() bool {
	return self.lastInstruction.Opcode == code.OpPop
}

func (self *Compiler) removeLastPop() {
	self.instructions = self.instructions[:self.lastInstruction.Position]
	self.lastInstruction = self.previousInstruction
}

func (self *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	for i := 0 ; i < len(newInstruction) ; i++ {
		self.instructions[pos + i] = newInstruction[i]
	}
}

func (self *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(self.instructions[opPos])
	newInstruction := code.Make(op, operand)
	self.replaceInstruction(opPos, newInstruction)
}