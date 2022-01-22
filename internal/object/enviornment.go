package object

type Environment struct {
	outer *Environment
	store map[string]Object
}

func NewEnv() *Environment {
	return &Environment{
		outer: nil,
		store: map[string]Object{},
	}
}

func NewEnclosedEnvironment(env *Environment) *Environment {
	e := NewEnv()
	e.outer = env
	return e
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}
