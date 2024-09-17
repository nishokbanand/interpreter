package object

type Environment struct {
	store map[string]Object
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func (e *Environment) Get(name string) (Object, bool) {
	value, ok := e.store[name]
	return value, ok
}

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}
