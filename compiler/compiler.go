package compiler

import (
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