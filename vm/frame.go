package vm

import (
	"bear/code"
	"bear/object"
)

type Frame struct {
	fn *object.CompiledFunction
	ip int
}

func NewFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{fn: fn, ip: -1}
}

func (self *Frame) Instructions() code.Instructions {
	return self.fn.Instructions
}