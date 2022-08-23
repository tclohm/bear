package object

import "fmt"

const (
	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	NULL_OBJ 	= "NULL"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (self *Integer) Type() ObjectType { return INTEGER_OBJ }
func (self *Integer) Inspect() string { return fmt.Sprintf("%d", self.Value) }

type Boolean struct {
	Value bool
}

func (self *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (self *Boolean) Inspect() string { return fmt.Sprintf("%d", self.Value) }

type Null struct {}

func (self *Null) Type() ObjectType { return NULL_OBJ }
func (self *Null) Inspect() string { return "null" }