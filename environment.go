package codecrafters_interpreter_go

import "fmt"

type Environment struct {
	Enclosing *Environment
	Values    map[string]interface{}
}

func NewEnvironment(enclosing *Environment) *Environment {
	return &Environment{
		Enclosing: enclosing,
		Values:    make(map[string]interface{}),
	}
}

// FIXME: if distance is greater than the number of enclosing environments, this will return wrong one.
func (e *Environment) ancestor(distance int) *Environment {
	env := e
	for i := 0; i < distance; i++ {
		env = env.Enclosing
	}

	return env
}

func (e *Environment) GetAt(distance int, name string) (interface{}, error) {
	return e.ancestor(distance).Values[name], nil
}

func (e *Environment) Define(name string, value interface{}) {
	e.Values[name] = value
}

func (e *Environment) Assign(name Token, value interface{}) error {
	if _, ok := e.Values[name.Lexeme]; ok {
		e.Values[name.Lexeme] = value
		return nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Assign(name, value)
	}

	return NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'", name.Lexeme))
}

func (e *Environment) AssignAt(distance int, name Token, value interface{}) error {
	e.ancestor(distance).Values[name.Lexeme] = value
	return nil
}

func (e *Environment) Get(name Token) (interface{}, error) {
	if val, ok := e.Values[name.Lexeme]; ok {
		return val, nil
	}

	if e.Enclosing != nil {
		return e.Enclosing.Get(name)
	}

	return nil, NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'", name.Lexeme))
}
