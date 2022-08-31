package vm

import (
	"bear/code"
	"bear/compiler"
	"bear/object"
)

const StackSize = 2048

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
		}
	}

	return nil
}