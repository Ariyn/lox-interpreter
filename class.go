package lox_interpreter

import "fmt"

var _ Callable = (*LoxClass)(nil)

type LoxClass struct {
	name       string
	superclass *LoxClass
	methods    map[string]Callable
}

func NewLoxClass(name string, superclass *LoxClass, methods map[string]Callable) *LoxClass {
	return &LoxClass{
		name,
		superclass,
		methods,
	}
}

func (l *LoxClass) Call(interpreter *Interpreter, arguments []interface{}) (_ interface{}, err error) {
	instance := NewLoxInstance(l)
	init := l.findMethod("init")
	if init != nil {
		_, err = init.Bind(instance).Call(interpreter, arguments)
	}

	return instance, err
}

func (l *LoxClass) Arity() int {
	if init := l.findMethod("init"); init != nil {
		return init.Arity()
	}

	return 0
}

func (l *LoxClass) findMethod(name string) Callable {
	if method, ok := l.methods[name]; ok {
		return method
	}

	return nil
}

func (l *LoxClass) ToString() string {
	return fmt.Sprintf("<cls %s>", l.name)
}

func (l *LoxClass) Bind(instance *LoxInstance) Callable {
	return l
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
		if ok := value.(*LiteralExpr); ok != nil {
			return value.(*LiteralExpr).value, nil
		}
		return value, nil
	}

	if method := l.class.findMethod(name.Lexeme); method != nil {
		if methodFunction, ok := method.(*LoxFunction); ok {
			return methodFunction.Bind(l), nil
		}

		return method, nil
	}

	if l.class.superclass != nil {
		return l.class.superclass.findMethod(name.Lexeme), nil
	}

	return nil, NewEnvironmentError(name, fmt.Sprintf("Undefined property '%s'.", name.Lexeme))
}

func (l *LoxInstance) Set(name Token, value interface{}) error {
	l.fields[name.Lexeme] = value
	return nil
}
