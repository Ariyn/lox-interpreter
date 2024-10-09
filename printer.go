package codecrafters_interpreter_go

import (
	"fmt"
)

var _ StmtVisitor = (*AstPrinter)(nil)
var _ ExprVisitor = (*AstPrinter)(nil)

type AstPrinter struct{}

func (ap *AstPrinter) Print(stmts []Stmt) (string, error) {
	for _, stmt := range stmts {
		v, err := stmt.Accept(ap)
		if err != nil {
			return "", err
		}

		fmt.Println(v)
	}

	return "", nil
}

func (ap *AstPrinter) evaluate(expr Expr) (interface{}, error) {
	return expr.Accept(ap)
}

func (ap *AstPrinter) VisitVarStmt(stmt *Var) (interface{}, error) {
	statementString := fmt.Sprintf("var (%s", stmt.name)

	if stmt.initializer != nil {
		d, err := ap.parenthesize("=", stmt.initializer)
		if err != nil {
			return "", err
		}
		statementString += toString(d)
	}

	statementString += ")"
	return statementString, nil
}

func (ap *AstPrinter) VisitExpressionStmt(expr *Expression) (interface{}, error) {
	return expr.expression.Accept(ap)
}

func (ap *AstPrinter) VisitPrintStmt(expr *Print) (interface{}, error) {
	return ap.parenthesize("print", expr.expression)
}

func (ap *AstPrinter) VisitBlockStmt(expr *Block) (interface{}, error) {
	build := "{"
	for _, stmt := range expr.statements {
		d, err := stmt.Accept(ap)
		if err != nil {
			return "", err
		}

		build += toString(d)
	}
	build += "}"

	return build, nil
}

func (ap *AstPrinter) VisitAssignExpr(expr *Assign) (interface{}, error) {
	return ap.parenthesize("= "+expr.name.Lexeme, expr.value)
}

func (ap *AstPrinter) VisitTernaryExpr(expr *Ternary) (interface{}, error) {
	return ap.parenthesize("?:", expr.condition, expr.left, expr.right)
}

func (ap *AstPrinter) VisitBinaryExpr(expr *Binary) (interface{}, error) {
	return ap.parenthesize(expr.operator.Lexeme, expr.left, expr.right)
}

func (ap *AstPrinter) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
	return ap.parenthesize("group", expr.expression)
}

func (ap *AstPrinter) VisitLiteralExpr(expr *Literal) (interface{}, error) {
	if expr.value == nil {
		return "nil", nil
	}
	return toString(expr.value), nil
}

func (ap *AstPrinter) VisitUnaryExpr(expr *Unary) (interface{}, error) {
	return ap.parenthesize(expr.operator.Lexeme, expr.right)
}

func (ap *AstPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	var builder string
	builder += "(" + name

	for _, expr := range exprs {
		builder += " "
		d, err := expr.Accept(ap)
		if err != nil {
			return "", err
		}

		builder += toString(d)
	}

	builder += ")"

	return builder, nil
}

type RPNAstPrinter struct {
	AstPrinter
}

func (ap *RPNAstPrinter) Print(expr Expr) (string, error) {
	v, err := expr.Accept(ap)
	return v.(string), err
}

func (ap *RPNAstPrinter) VisitBinaryExpr(expr *Binary) (interface{}, error) {
	l, err := expr.left.Accept(ap)
	if err != nil {
		return "", err
	}

	r, err := expr.right.Accept(ap)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s %s %s", toString(l), toString(r), expr.operator.Lexeme), nil
}

func (ap *RPNAstPrinter) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
	return expr.expression.Accept(ap)
}

func (ap *AstPrinter) VisitVariableExpr(expr *Variable) (interface{}, error) {
	return expr.name.Lexeme, nil
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
