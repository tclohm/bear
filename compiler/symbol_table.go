package compiler

type SymbolScope string

const (
	LocalScope SymbolScope = "LOCAL"
	GlobalScope SymbolScope = "GLOBAL"
)

type Symbol struct {
	Name string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer 			*SymbolTable

	store 			map[string]Symbol
	numDefinitions 	int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

func (self *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: self.numDefinitions}
	if self.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}
	self.store[name] = symbol
	self.numDefinitions++
	return symbol
}

func (self *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := self.store[name]
	if !ok && self.Outer != nil {
		obj, ok := self.Outer.Resolve(name)
		return obj, ok
	}
	return obj, ok
}