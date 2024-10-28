package lox_interpreter

import "time"

type Callable interface {
	Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error)
	Arity() int
	ToString() string
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

	if block, ok := f.declaration.body.(*BlockStmt); ok {
		value, err := interpreter.executeBlock(block.statements, env)
		if err != nil {
			return nil, err
		}

		if f.isInitializer {
			return f.closure.GetAtWithString(0, "this")
		}
		return value, nil
	} else {
		// FIXME: This should be an error.
		panic("FUNCTION BODY IS NOT BLOCK")
	}

	return nil, nil
}

func (f *LoxFunction) Arity() int {
	return len(f.declaration.params)
}

func (f *LoxFunction) ToString() string {
	return "<fn " + f.declaration.name.Lexeme + ">"
}

type Clock struct{}

func (c *Clock) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return float64(time.Now().UnixNano()), nil
}

func (c *Clock) Arity() int {
	return 0
}
