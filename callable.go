package lox_interpreter

import (
	"fmt"
	"time"
)

type Callable interface {
	Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error)
	Arity() int
	ToString() string
	Bind(instance *LoxInstance) Callable
}

var _ Callable = (*LoxFunction)(nil)

type LoxFunction struct {
	declaration   *FunStmt
	closure       *Environment
	isInitializer bool
}

func NewFunction(declaration *FunStmt, closure *Environment, isInitializer bool) *LoxFunction {
	return &LoxFunction{
		declaration,
		closure,
		isInitializer,
	}
}

func (f *LoxFunction) Bind(instance *LoxInstance) Callable {
	env := NewEnvironment(f.closure)
	env.Define("this", instance)

	return NewFunction(f.declaration, env, f.isInitializer)
}

func (f *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	env := NewEnvironment(f.closure)
	for i, param := range f.declaration.params {
		env.Define(param.Lexeme, arguments[i])
	}

	value, err := interpreter.executeBlock(f.declaration.body, env)
	if err != nil {
		return nil, err
	}

	if f.isInitializer {
		return f.closure.GetAtWithString(0, "this")
	}
	return value, nil
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.params)
}

func (f *LoxFunction) ToString() string {
	return "<fn " + f.declaration.name.Lexeme + ">"
}

var _ Callable = (*Clock)(nil)

type Clock struct{}

func (c *Clock) ToString() string {
	return "<native fn clock>"
}

func (c *Clock) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return float64(time.Now().UnixNano()), nil
}

func (c *Clock) Arity() int {
	return 0
}

func (c *Clock) Bind(instance *LoxInstance) Callable {
	return c
}

var _ Callable = (*Len)(nil)

type Len struct{}

func (l Len) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	switch arg := arguments[0].(type) {
	case string:
		return float64(len(arg)), nil
	case ListType:
		return float64(len(arg)), nil
	default:
		return nil, fmt.Errorf("Argument must be a string or an array.")
	}
}

func (l Len) Arity() int {
	return 1
}

func (l Len) ToString() string {
	return "<native fn len>"
}

func (l Len) Bind(instance *LoxInstance) Callable {
	return l
}
