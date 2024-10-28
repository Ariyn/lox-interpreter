package main

import (
	lox "github.com/ariyn/lox_interpreter"
	"math/rand"
	"strings"
)

var _ (lox.Callable) = (*RandFunc)(nil)

type RandFunc struct {
}

func (g RandFunc) Call(_ *lox.Interpreter, arguments []interface{}) (v interface{}, err error) {
	x := strings.Split("abcdefghijklmnopqrstuvwxyz", "")
	rand.Shuffle(len(x), func(i, j int) {
		x[i], x[j] = x[j], x[i]
	})
	return strings.Join(x, "")[:6], nil
}

func (g RandFunc) Arity() int {
	return 0
}

func (g RandFunc) ToString() string {
	return "random word"
}

func main() {
	script := `print clock();
var x = rand();
print "RESULT " + x;
`
	scanner := lox.NewScanner(script)
	tokens, _ := scanner.ScanTokens()

	parser := lox.NewParser(tokens)
	statements, _ := parser.Parse()

	env := lox.NewEnvironment(nil)
	env.Define("rand", &RandFunc{})
	interpreter := lox.NewInterpreter(env)

	resolver := lox.NewResolver(interpreter)
	err := resolver.Resolve(statements...)
	if err != nil {
		panic(err)
	}

	_, err = interpreter.Interpret(statements)
	if err != nil {
		panic(err)
	}
}
