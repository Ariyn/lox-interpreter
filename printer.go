package codecrafters_interpreter_go

import (
	"fmt"
)

type AstPrinter struct{}

func (ap *AstPrinter) Print(expr Expr) string {
	return expr.Accept(ap).(string)
}

func (ap *AstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	return ap.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (ap *AstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return ap.parenthesize("group", expr.expression)
}

func (ap *AstPrinter) VisitLiteralExpr(expr *Literal) interface{} {
	if expr.value == nil {
		return "nil"
	}
	return toString(expr.value)
}

func (ap *AstPrinter) VisitUnaryExpr(expr *Unary) interface{} {
	return ap.parenthesize(expr.operator.Lexeme, expr.right)
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var builder string
	builder += "(" + name

	for _, expr := range exprs {
		builder += " "
		d := expr.Accept(ap)
		builder += toString(d)
	}

	builder += ")"

	return builder
}

type RPNAstPrinter struct {
	AstPrinter
}

func (ap *RPNAstPrinter) Print(expr Expr) string {
	return expr.Accept(ap).(string)
}

func (ap *RPNAstPrinter) VisitBinaryExpr(expr *Binary) interface{} {
	l := expr.left.Accept(ap)
	r := expr.right.Accept(ap)

	return fmt.Sprintf("%s %s %s", toString(l), toString(r), expr.operator.Lexeme)
}

func (ap *RPNAstPrinter) VisitGroupingExpr(expr *Grouping) interface{} {
	return expr.expression.Accept(ap)
}

func toString(d interface{}) string {
	switch d.(type) {
	case string:
		return d.(string)
	case float64:
		if d.(float64) == float64(int(d.(float64))) {
			return fmt.Sprintf("%.1f", d)
		}
		return fmt.Sprintf("%g", d.(float64))
	case int:
		return fmt.Sprintf("%d", d)
	case int64:
		return fmt.Sprintf("%d", d)
	case bool:
		return fmt.Sprintf("%t", d)
	default:
		return "nil"
	}
}
