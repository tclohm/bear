package vm

import (
	"fmt"
	"bear/code"
	"bear/compiler"
	"bear/object"
)

const StackSize = 2048
var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}

type VM struct {
	constants 	 []object.Object
	instructions code.Instructions

	stack 		 []object.Object
	sp 			 int // point the next value
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants: 	  bytecode.Constants,

		stack: 		  make([]object.Object, StackSize),
		sp: 		  0,
	}
}

func (self *VM) StackTop() object.Object {
	if self.sp == 0 {
		return nil
	}
	return self.stack[self.sp - 1]
}

func (self *VM) Run() error {
	for ip := 0 ; ip < len(self.instructions) ; ip++ {
		op := code.Opcode(self.instructions[ip])

		switch op {
		case code.OpConstant:
			constIndex := code.ReadUint16(self.instructions[ip+1:])
			// MARK: - ip point to opcode instead of an operand
			ip += 2
			err := self.push(self.constants[constIndex])
			if err != nil {
				return err
			}
		case code.OpAdd, code.OpSub, code.OpMul, code.OpDiv:
			err := self.executeBinaryOperation(op)
			if err != nil {
				return err
			}
		case code.OpPop:
			self.pop()
		case code.OpTrue:
			err := self.push(True)
			if err != nil {
				return err
			}
		case code.OpFalse:
			err := self.push(False)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (self *VM) push(o object.Object) error {
	if self.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	self.stack[self.sp] = o
	self.sp++

	return nil
}

func (self *VM) pop() object.Object {
	o := self.stack[self.sp - 1]
	self.sp--
	return o
}

func (self *VM) LastPoppedStackElem() object.Object {
	return self.stack[self.sp]
}

func (self *VM) executeBinaryOperation(op code.Opcode) error {
	right := self.pop()
	left := self.pop()

	leftType := left.Type()
	rightType := right.Type()

	if leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ {
		return self.executeBinaryIntegerOperation(op, left, right)
	}
	
	return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
}

func (self *VM) executeBinaryIntegerOperation(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	var result int64

	switch op {
	case code.OpAdd:
		result = leftValue + rightValue
	case code.OpSub:
		result = leftValue - rightValue
	case code.OpMul:
		result = leftValue * rightValue
	case code.OpDiv:
		result = leftValue / rightValue
	default:
		return fmt.Errorf("unknown integer operator: %d", op)
	}

	return self.push(&object.Integer{Value: result})
}