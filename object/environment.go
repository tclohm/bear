package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func (self *Environment) Get(name string) (Object, bool) {
	obj, ok := self.store[name]
	if !ok && self.outer != nil {
		obj, ok = self.outer.Get(name)
	}
	return obj, ok
}

func (self *Environment) Set(name string, value Object) Object {
	self.store[name] = value
	return value
}
