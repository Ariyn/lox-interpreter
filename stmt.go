package codecrafters_interpreter_go
type StmtVisitor interface {
	VisitExpressionStmt(expr *Expression) (interface{}, error)
	VisitPrintStmt(expr *Print) (interface{}, error)
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
	return v.VisitExpressionStmt(e)
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
	return v.VisitPrintStmt(e)
}

