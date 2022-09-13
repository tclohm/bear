package vm

import (
	"fmt"
	"bear/code"
	"bear/compiler"
	"bear/object"
)

const StackSize = 2048
const GlobalsSize = 65536
var True = &object.Boolean{Value: true}
var False = &object.Boolean{Value: false}
var Null = &object.Null{}

const MaxFrames = 1024

type VM struct {
	constants 	 []object.Object

	stack 		 []object.Object
	sp 			 int // point the next value

	globals 	 []object.Object

	frames 		 []*Frame
	framesIndex	 int
}

func (self *VM) currentFrame() *Frame {
	return self.frames[self.framesIndex - 1]
}

func (self *VM) pushFrame(f *Frame) {
	self.frames[self.framesIndex] = f
	self.framesIndex++
}

func (self *VM) popFrame() *Frame {
	self.framesIndex--
	return self.frames[self.framesIndex]
}

func New(bytecode *compiler.Bytecode) *VM {
	mainFn := &object.CompiledFunction{Instructions: bytecode.Instructions}
	mainFrame := NewFrame(mainFn, 0)

	frames := make([]*Frame, MaxFrames)
	frames[0] = mainFrame

	return &VM{
		constants: 	  bytecode.Constants,

		stack: 		  make([]object.Object, StackSize),
		sp: 		  0,

		globals: make([]object.Object, GlobalsSize),

		frames: frames,
		framesIndex: 1,
	}
}

func (self *VM) StackTop() object.Object {
	if self.sp == 0 {
		return nil
	}
	return self.stack[self.sp - 1]
}

func (self *VM) Run() error {
	var ip int
	var ins code.Instructions
	var op code.Opcode

	for self.currentFrame().ip < len(self.currentFrame().Instructions()) - 1 {
		self.currentFrame().ip++

		ip = self.currentFrame().ip
		ins = self.currentFrame().Instructions()

		op = code.Opcode(ins[ip])

		switch op {

		case code.OpConstant:
			constIndex := code.ReadUint16(ins[ip+1:])
			// MARK: - ip point to opcode instead of an operand
			self.currentFrame().ip += 2
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

		case code.OpEqual, code.OpNotEqual, code.OpGreaterThan:
			err := self.executeComparison(op)
			if err != nil {
				return err
			}

		case code.OpBang:
			err := self.executeBangOperator()
			if err != nil {
				return err
			}

		case code.OpMinus:
			err := self.executeMinusOperator()
			if err != nil {
				return err
			}

		case code.OpJump:
			pos := int(code.ReadUint16(ins[ip+1:]))
			self.currentFrame().ip = pos - 1

		case code.OpJumpNotTruthy:
			pos := int(code.ReadUint16(ins[ip+1:]))
			self.currentFrame().ip += 2

			condition := self.pop()
			if !isTruthy(condition) {
				self.currentFrame().ip = pos - 1
			}

		case code.OpNull:
			err := self.push(Null)
			if err != nil {
				return err
			}

		case code.OpSetGlobal:
			globalIndex := code.ReadUint16(ins[ip + 1:])
			self.currentFrame().ip += 2
			self.globals[globalIndex] = self.pop()

		case code.OpGetGlobal:
			globalIndex := code.ReadUint16(ins[ip + 1:])
			self.currentFrame().ip += 2

			err := self.push(self.globals[globalIndex])
			if err != nil {
				return err
			}

		case code.OpArray:
			numberElems := int(code.ReadUint16(ins[ip+1:]))
			self.currentFrame().ip += 2

			array := self.buildArray(self.sp-numberElems, self.sp)
			self.sp = self.sp - numberElems

			err := self.push(array)

			if err != nil {
				return err
			}

		case code.OpHash:
			numberElems := int(code.ReadUint16(ins[ip + 1:]))
			self.currentFrame().ip += 2

			hash, err := self.buildHash(self.sp - numberElems, self.sp)
			if err != nil {
				return err
			}

			self.sp = self.sp - numberElems

			err = self.push(hash)

			err = self.push(hash)
			if err != nil {
				return err
			}

		case code.OpIndex:
			index := self.pop()
			left := self.pop()

			err := self.executeIndexExpression(left, index)
			if err != nil {
				return err
			}

		case code.OpCall:
			fn, ok := self.stack[self.sp - 1].(*object.CompiledFunction)
			if !ok {
				return fmt.Errorf("calling non-function")
			}
			frame := NewFrame(fn, self.sp)
			self.pushFrame(frame)
			self.sp = frame.basePointer + fn.NumLocals

		case code.OpReturnValue:
			returnValue := self.pop()

			self.popFrame()
			self.pop()

			err := self.push(returnValue)
			if err != nil {
				return err
			}
			
		case code.OpReturn:
			self.popFrame()
			self.pop()

			err := self.push(Null)
			if err != nil {
				return err
			}
		
		case code.OpSetLocal:
			localIndex := code.ReadUint8(ins[ip + 1:])
			self.currentFrame().ip += 1

			frame := self.currentFrame()

			self.stack[frame.basePointer+int(localIndex)] = self.pop()

		case code.OpGetLocal:
			localIndex := code.ReadUint8(ins[ip + 1:])
			self.currentFrame().ip += 1

			frame := self.currentFrame()

			err := self.push(self.stack[frame.basePointer+int(localIndex)])
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

	switch {
	case leftType == object.INTEGER_OBJ && rightType == object.INTEGER_OBJ:
		return self.executeBinaryIntegerOperation(op, left, right)
	case leftType == object.STRING_OBJ && rightType == object.STRING_OBJ:
		return self.executeBinaryStringOperation(op, left, right)
	default:
		return fmt.Errorf("unsupported types for binary operation: %s %s", leftType, rightType)
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

func (self *VM) executeComparison(op code.Opcode) error {
	right := self.pop()
	left := self.pop()

	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return self.executeIntegerComparison(op, left, right)
	}

	switch op {
	case code.OpEqual:
		return self.push(nativeBoolToBooleanObject(right == left))
	case code.OpNotEqual:
		return self.push(nativeBoolToBooleanObject(right != left))
	default:
		return fmt.Errorf("unknown operator: %d (%s %s)", op, left.Type(), right.Type())
	}
}

func (self *VM) executeIntegerComparison(op code.Opcode, left, right object.Object) error {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch op {
	case code.OpEqual:
		return self.push(nativeBoolToBooleanObject(rightValue == leftValue))
	case code.OpNotEqual:
		return self.push(nativeBoolToBooleanObject(rightValue != leftValue))
	case code.OpGreaterThan:
		return self.push(nativeBoolToBooleanObject(leftValue > rightValue))
	default:
		return fmt.Errorf("unknown operator: %d", op)
	}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return True
	}
	return False
}

func (self *VM) executeBangOperator() error {
	operand := self.pop()

	switch operand {
	case True:
		return self.push(False)
	case False:
		return self.push(True)
	case Null:
		return self.push(True)
	default:
		return self.push(False)
	}
}

func (self *VM) executeMinusOperator() error {
	operand := self.pop()

	if operand.Type() != object.INTEGER_OBJ {
		return fmt.Errorf("unsupported type for negation: %s", operand.Type())
	}

	value := operand.(*object.Integer).Value
	return self.push(&object.Integer{Value: -value})
}

func isTruthy(obj object.Object) bool {
	switch obj := obj.(type) {
	case *object.Boolean:
		return obj.Value
	case *object.Null:
		return false
	default:
		return true
	}
}

func NewWithGlobalsStore(bytecode *compiler.Bytecode, s []object.Object) *VM {
	vm := New(bytecode)
	vm.globals = s
	return vm
}

func (self *VM) executeBinaryStringOperation(op code.Opcode, left, right object.Object) error {
	if op != code.OpAdd {
		return fmt.Errorf("unknown string operator: %d", op)
	}

	leftValue := left.(*object.String).Value
	rightValue := right.(*object.String).Value

	return self.push(&object.String{Value: leftValue + rightValue})
}

func (self *VM) buildArray(start, end int) object.Object {
	elements := make([]object.Object, end - start)

	for i := start ; i < end ; i++ {
		elements[i - start] = self.stack[i]
	}

	return &object.Array{Elements: elements}
}

func (self *VM) buildHash(start, end int) (object.Object, error) {
	hashedPairs := make(map[object.HashKey]object.HashPair)

	for i := start ; i < end ; i += 2 {
		key := self.stack[i]
		value := self.stack[i + 1]

		pair := object.HashPair{Key: key, Value: value}

		hashKey, ok := key.(object.Hashable)
		if !ok {
			return nil, fmt.Errorf("unusable as hash key: %s", key.Type())
		}

		hashedPairs[hashKey.HashKey()] = pair

	}

	return &object.Hash{Pairs: hashedPairs}, nil
}

func (self *VM) executeIndexExpression(left, index object.Object) error {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return self.executeArrayIndex(left, index)
	case left.Type() == object.HASH_OBJ:
		return self.executeHashIndex(left, index)
	default:
		return fmt.Errorf("index operator not supported: %s", left.Type())
	}
}


func (self *VM) executeArrayIndex(array, index object.Object) error {
	arrayObject := array.(*object.Array)
	i := index.(*object.Integer).Value
	max := int64(len(arrayObject.Elements) - 1)

	if i < 0 || i > max {
		return self.push(Null)
	}

	return self.push(arrayObject.Elements[i])
}

func (self *VM) executeHashIndex(hash, index object.Object) error {
	hashObject := hash.(*object.Hash)

	key, ok := index.(object.Hashable)

	if !ok {
		return fmt.Errorf("unusable as hash key: %s", index.Type())
	}

	pair, ok := hashObject.Pairs[key.HashKey()]
	if !ok {
		return self.push(Null)
	}

	return self.push(pair.Value)
}


