package codecrafters_interpreter_go

import (
	"fmt"
	"log"
)

type RuntimeError struct {
	token   Token
	message string
}

func (r *RuntimeError) Error() string {
	return fmt.Sprintf("%d at '%s' %s", r.token.LineNumber, r.token.Lexeme, r.message)
}

func NewRuntimeError(token Token, message string) error {
	return &RuntimeError{token, message}
}

var _ StmtVisitor = (*Interpreter)(nil)
var _ ExprVisitor = (*Interpreter)(nil)

type Interpreter struct {
}

func NewInterpreter() *Interpreter {
	return &Interpreter{}
}

func (i *Interpreter) Interpret(expr []Stmt) (interface{}, error) {
	var err error
	for _, stmt := range expr {
		_err := i.execute(stmt)

		if err != nil {
			err = _err
			log.Printf("%s\n[line %d]", err.Error(), err.(*RuntimeError).token.LineNumber)
		}
	}

	return nil, err
}

func (i *Interpreter) execute(stmt Stmt) error {
	_, err := stmt.Accept(i)
	if err != nil {
		return err
	}

	return nil
}

func (i *Interpreter) VisitVarStmt(expr *Var) (interface{}, error) {
	//TODO implement me
	panic("implement me")
}

func (i *Interpreter) VisitExpressionStmt(expr *Expression) (interface{}, error) {
	_, err := i.evaluate(expr.expression)
	return nil, err
}

func (i *Interpreter) VisitPrintStmt(expr *Print) (interface{}, error) {
	value, err := i.evaluate(expr.expression)
	if err != nil {
		return nil, err
	}

	fmt.Println(Stringify(value))
	return nil, nil
}

func (i *Interpreter) VisitLiteralExpr(expr *Literal) (interface{}, error) {
	return expr.value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
	return i.evaluate(expr.expression)
}

func (i *Interpreter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) VisitTernaryExpr(expr *Ternary) (interface{}, error) {
	condition, err := i.evaluate(expr.condition)
	if err != nil {
		return nil, err
	}

	if i.isTruthy(condition) {
		return i.evaluate(expr.left)
	}

	return i.evaluate(expr.right)
}

func (i *Interpreter) VisitUnaryExpr(expr *Unary) (interface{}, error) {
	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.Type {
	case MINUS:
		isNumber := i.isAllNumber(right)
		if !isNumber {
			return nil, NewRuntimeError(expr.operator, "Operand must be a number.") // TODO: return error
		}

		return -right.(float64), nil
	case BANG:
		return !i.isTruthy(right), nil
	}

	return nil, nil // TODO: return error
}

func (i *Interpreter) isTruthy(value interface{}) bool {
	if value == nil {
		return false
	}

	if value == false {
		return false
	}

	return true
}

func (i *Interpreter) VisitBinaryExpr(expr *Binary) (interface{}, error) {
	left, err := i.evaluate(expr.left)
	if err != nil {
		return nil, err
	}

	right, err := i.evaluate(expr.right)
	if err != nil {
		return nil, err
	}

	switch expr.operator.Type {
	case MINUS:
		if !i.isAllNumber(left, right) {
			return nil, NewRuntimeError(expr.operator, "Operands must be two numbers or two strings.")
		}

		return left.(float64) - right.(float64), nil
	case SLASH:
		if !i.isAllNumber(left, right) {
			return nil, NewRuntimeError(expr.operator, "Operands must be two numbers or two strings.")
		}

		if right.(float64) == 0 {
			return nil, NewRuntimeError(expr.operator, "Division by zero.")
		}

		return left.(float64) / right.(float64), nil
	case STAR:
		if !i.isAllNumber(left, right) {
			return nil, NewRuntimeError(expr.operator, "Operands must be two numbers or two strings.")
		}

		return left.(float64) * right.(float64), nil
	case PLUS:
		if i.isAllNumber(left, right) {
			return left.(float64) + right.(float64), nil
		}
		if i.isAllString(left, right) {
			return left.(string) + right.(string), nil
		}

		if i.isAllStringOrNumber(left, right) {
			return Stringify(left) + Stringify(right), nil
		}

		return nil, NewRuntimeError(expr.operator, "Operands must be two numbers or two strings.")
	case GREATER:
		if i.isAllNumber(left, right) {
			return left.(float64) > right.(float64), nil
		} else if i.isAllString(left, right) {
			return left.(string) > right.(string), nil
		}

		return nil, NewRuntimeError(expr.operator, "Operands must be two numbers or two strings.")
	case GREATER_EQUAL:
		if i.isAllNumber(left, right) {
			return left.(float64) >= right.(float64), nil
		} else if i.isAllString(left, right) {
			return left.(string) >= right.(string), nil
		}

		return nil, NewRuntimeError(expr.operator, "Operands must be two numbers or two strings.")
	case LESS:
		if i.isAllNumber(left, right) {
			return left.(float64) < right.(float64), nil
		} else if i.isAllString(left, right) {
			return left.(string) < right.(string), nil
		}
	case LESS_EQUAL:
		if i.isAllNumber(left, right) {
			return left.(float64) <= right.(float64), nil
		} else if i.isAllString(left, right) {
			return left.(string) <= right.(string), nil
		}
	case EQUAL_EQUAL:
		return left == right, nil
	case BANG_EQUAL:
		return left != right, nil
	}

	return nil, nil // TODO: return error
}

func (i *Interpreter) isAllNumber(possibles ...interface{}) bool {
	for _, possible := range possibles {
		if _, ok := possible.(float64); !ok {
			return false
		}
	}

	return true
}

func (i *Interpreter) isAllString(possibles ...interface{}) bool {
	for _, possible := range possibles {
		if _, ok := possible.(string); !ok {
			return false
		}
	}

	return true
}

func (i *Interpreter) isAllStringOrNumber(possibles ...interface{}) bool {
	for _, possible := range possibles {
		_, isString := possible.(string)
		_, isNumber := possible.(float64)

		if !isString && !isNumber {
			return false
		}
	}

	return true
}

func Stringify(d interface{}) string {
	switch d.(type) {
	case float64:
		if d.(float64) == float64(int(d.(float64))) {
			return fmt.Sprintf("%.0f", d)
		}
		return fmt.Sprintf("%g", d.(float64))
	default:
		return toString(d)
	}
}
