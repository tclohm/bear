package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
	Name string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	store 			map[string]Symbol
	numDefinition 	int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func (self *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: self.numDefinition, Scope: GlobalScope}
	self.store[name] = symbol
	self.numDefinition++
	return symbol
}

func (self *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := self.store[name]
	return obj, ok
}