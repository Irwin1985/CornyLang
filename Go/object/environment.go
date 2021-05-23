package object

type Environment struct {
	symbolTable map[string]Object
	Parent      *Environment
}

func NewEnvironment() *Environment {
	var env = &Environment{symbolTable: make(map[string]Object), Parent: nil}
	return env
}

func (env *Environment) Set(key string, value Object) {
	env.symbolTable[key] = value
}

func (env *Environment) Get(key string) (Object, bool) {
	value, ok := env.symbolTable[key]
	if !ok && env.Parent != nil {
		value, ok = env.Parent.Get(key)
	}
	return value, ok
}
