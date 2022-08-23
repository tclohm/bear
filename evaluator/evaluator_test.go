package evaluator

import (
	"bear/lexer"
	"bear/object"
	"bear/parser"
	"testing"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input 		string
		expected 	int64
	}{
		{"5", 5},
		{"5", 5},
	}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

// MARK: -- helpers
func testEval(input string) object.Object {
	lex := lexer.New(input)
	p := parser.New(lex)
	program := p.ParseProgram()

	return Eval(program)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}