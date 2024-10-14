package codecrafters_interpreter_go

import "time"

type Callable interface {
	Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error)
	Arity() int
	ToString() string
}

var _ Callable = (*Function)(nil)

type Function struct {
	declaration *Fun
	closure     *Environment
}

func NewFunction(declaration *Fun, closure *Environment) *Function {
	return &Function{
		declaration,
		closure,
	}
}

func (f *Function) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	env := NewEnvironment(f.closure)
	for i, param := range f.declaration.params {
		env.Define(param.Lexeme, arguments[i])
	}

	if block, ok := f.declaration.body.(*Block); ok {
		value, err := interpreter.executeBlock(block.statements, env)
		if err != nil {
			return nil, err
		}

		return value, nil
	} else {
		// FIXME: This should be an error.
		panic("FUNCTION BODY IS NOT BLOCK")
	}

	return nil, nil
}

func (f *Function) Arity() int {
	return len(f.declaration.params)
}

func (f *Function) ToString() string {
	return "<fn " + f.declaration.name.Lexeme + ">"
}

type Clock struct{}

func (c *Clock) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return float64(time.Now().UnixNano()), nil
}

func (c *Clock) Arity() int {
	return 0
}
