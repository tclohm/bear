package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func (self *Environment) Get(name string) (Object, bool) {
	obj, ok := self.store[name]
	return obj, ok
}

func (self *Environment) Set(name string, value Object) Object {
	self.store[name] = value
	return value
}
