package compiler

import (
	"fmt"
	"bear/ast"
	"bear/code"
	"bear/object"
	"sort"
)


type EmittedInstruction struct {
	Opcode 		code.Opcode
	Position 	int
}

type CompilationScope struct {
	instructions 		code.Instructions
	lastInstruction 	EmittedInstruction
	previousInstruction EmittedInstruction
}

type Compiler struct {
	constants  	 		[]object.Object

	symbolTable 		*SymbolTable

	scopes 				[]CompilationScope
	scopeIndex 			int
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions: 			code.Instructions{},
		lastInstruction: 		EmittedInstruction{},
		previousInstruction: 	EmittedInstruction{},
	}

	return &Compiler{
		constants:		[]object.Object{},
		symbolTable: 	NewSymbolTable(),
		scopes: 		[]CompilationScope{mainScope},
		scopeIndex: 	0,
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

		if self.lastInstructionIs(code.OpPop) {
			self.removeLastPop()
		}

		jumpPos := self.emit(code.OpJump, 9999)

		afterConsequencePos := len(self.currentInstructions())
		self.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			self.emit(code.OpNull)
		} else {
			err := self.Compile(node.Alternative)
			if err != nil {
				return err
			}
			
			if self.lastInstructionIs(code.OpPop) {
				self.removeLastPop()
			}
		}

		afterAlternativePos := len(self.currentInstructions())
		self.changeOperand(jumpPos, afterAlternativePos)

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := self.Compile(s)
			if err != nil {
				return err
			}
		}

	case *ast.LetStatement:
		err := self.Compile(node.Value)
		if err != nil {
			return err
		}

		symbol := self.symbolTable.Define(node.Name.Value)
		if symbol.Scope == GlobalScope {
			self.emit(code.OpSetGlobal, symbol.Index)
		} else {
			self.emit(code.OpSetLocal, symbol.Index)
		}
		

	case *ast.Identifier:
		symbol, ok := self.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		if symbol.Scope == GlobalScope {
			self.emit(code.OpGetGlobal, symbol.Index)
		} else {
			self.emit(code.OpGetLocal, symbol.Index)
		}

		

	case *ast.StringLiteral:
		str := &object.String{Value: node.Value}
		self.emit(code.OpConstant, self.addConstant(str))

	case *ast.ArrayLiteral:
		for _, el := range node.Elements {
			err := self.Compile(el)
			if err != nil {
				return err
			}
		}

		self.emit(code.OpArray, len(node.Elements))

	case *ast.HashLiteral:
		keys := []ast.Expression{}
		for k := range node.Pairs {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i].String() < keys[j].String()
		})

		for _, k := range keys {
			err := self.Compile(k)
			if err != nil {
				return err
			}
			err = self.Compile(node.Pairs[k])
			if err != nil {
				return err
			}
		}

		self.emit(code.OpHash, len(node.Pairs)*2)

	case *ast.IndexExpression:
		err := self.Compile(node.Left)

		if err != nil {
			return err
		}

		err = self.Compile(node.Index)
		if err != nil {
			return err
		}

		self.emit(code.OpIndex)

	case *ast.FunctionLiteral:
		self.enterScope()

		err := self.Compile(node.Body)
		if err != nil {
			return err
		}

		if self.lastInstructionIs(code.OpPop) {
			self.replaceLastPopWithReturn()
		}

		if !self.lastInstructionIs(code.OpReturnValue) {
			self.emit(code.OpReturn)
		}

		instructions := self.leaveScope()

		compiledFn := &object.CompiledFunction{Instructions: instructions}
		self.emit(code.OpConstant, self.addConstant(compiledFn))

	case *ast.ReturnStatement:
		err := self.Compile(node.ReturnValue)
		if err != nil {
			return err
		}

		self.emit(code.OpReturnValue)

	case *ast.CallExpression:
		err := self.Compile(node.Function)
		if err != nil {
			return err
		}

		self.emit(code.OpCall)
	}

	return nil
}

func (self *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: self.currentInstructions(),
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
	posNewInstruction := len(self.currentInstructions())
	updatedInstructions := append(self.currentInstructions(), ins...)

	self.scopes[self.scopeIndex].instructions = updatedInstructions

	return posNewInstruction
}

func (self *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := self.scopes[self.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	self.scopes[self.scopeIndex].previousInstruction = previous
	self.scopes[self.scopeIndex].lastInstruction = last
}

func (self *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(self.currentInstructions()) == 0 {
		return false
	}

	return self.scopes[self.scopeIndex].lastInstruction.Opcode == op
}

func (self *Compiler) removeLastPop() {
	last := self.scopes[self.scopeIndex].lastInstruction
	previous := self.scopes[self.scopeIndex].previousInstruction

	old := self.currentInstructions()
	new := old[:last.Position]

	self.scopes[self.scopeIndex].instructions = new
	self.scopes[self.scopeIndex].lastInstruction = previous
}

func (self *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	ins := self.currentInstructions()

	for i := 0 ; i < len(newInstruction) ; i++ {
		ins[pos + i] = newInstruction[i]
	}
}

func (self *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(self.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)
	self.replaceInstruction(opPos, newInstruction)
}

func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
}

func (self *Compiler) currentInstructions() code.Instructions {
	return self.scopes[self.scopeIndex].instructions
}

func (self *Compiler) enterScope() {
	scope := CompilationScope{
		instructions: code.Instructions{},
		lastInstruction: EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	self.scopes = append(self.scopes, scope)
	self.scopeIndex++

	self.symbolTable = NewEnclosedSymbolTable(self.symbolTable)
}

func (self *Compiler) leaveScope() code.Instructions {
	instructions := self.currentInstructions()

	self.scopes = self.scopes[:len(self.scopes)-1]
	self.scopeIndex--

	self.symbolTable = self.symbolTable.Outer

	return instructions
}

func (self *Compiler) replaceLastPopWithReturn() {
	lastPos := self.scopes[self.scopeIndex].lastInstruction.Position
	self.replaceInstruction(lastPos, code.Make(code.OpReturnValue))

	self.scopes[self.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}