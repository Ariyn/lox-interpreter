package codecrafters_interpreter_go
type StmtVisitor interface {
	VisitExpressionExpr(expr *Expression) (interface{}, error)
	VisitPrintExpr(expr *Print) (interface{}, error)
}

type Stmt interface {
	Accept(v StmtVisitor) (interface{}, error)
}
var _ Stmt = (*Expression)(nil)
type Expression struct {
	expression Expr
}

func NewExpression(expression Expr) *Expression {
	return &Expression{
		expression,
	}
}

func (e *Expression) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitExpressionExpr(e)
}

var _ Stmt = (*Print)(nil)
type Print struct {
	expression Expr
}

func NewPrint(expression Expr) *Print {
	return &Print{
		expression,
	}
}

func (e *Print) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitPrintExpr(e)
}

