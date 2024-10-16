package codecrafters_interpreter_go

import "fmt"

var _ Callable = (*LoxClass)(nil)

type LoxClass struct {
	name    string
	methods map[string]Callable
}

func NewLoxClass(name string, methods map[string]Callable) *LoxClass {
	return &LoxClass{
		name,
		methods,
	}
}

func (l *LoxClass) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	return NewLoxInstance(l), nil
}

func (l *LoxClass) Arity() int {
	return 0
}

func (l *LoxClass) ToString() string {
	return fmt.Sprintf("<cls %s>", l.name)
}

type LoxInstance struct {
	class  *LoxClass
	fields map[string]interface{}
}

func NewLoxInstance(class *LoxClass) *LoxInstance {
	return &LoxInstance{
		class,
		make(map[string]interface{}),
	}
}

func (l *LoxInstance) ToString() string {
	return fmt.Sprintf("<inst %s>", l.class.name)
}

func (l *LoxInstance) Get(name Token) (interface{}, error) {
	if value, ok := l.fields[name.Lexeme]; ok {
		return value, nil
	}

	if method, ok := l.class.methods[name.Lexeme]; ok {
		return method, nil
	}

	return nil, NewEnvironmentError(name, fmt.Sprintf("Undefined property '%s'.", name.Lexeme))
}

func (l *LoxInstance) Set(name Token, value interface{}) error {
	l.fields[name.Lexeme] = value
	return nil
}
