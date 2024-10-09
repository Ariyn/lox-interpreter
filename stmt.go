package codecrafters_interpreter_go
type StmtVisitor interface {
	VisitVarStmt(expr *Var) (interface{}, error)
	VisitExpressionStmt(expr *Expression) (interface{}, error)
	VisitPrintStmt(expr *Print) (interface{}, error)
}

type Stmt interface {
	Accept(v StmtVisitor) (interface{}, error)
}
var _ Stmt = (*Var)(nil)
type Var struct {
	name Token
	initializer Expr
}

func NewVar(name Token, initializer Expr) *Var {
	return &Var{
		name,
		initializer,
	}
}

func (e *Var) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitVarStmt(e)
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

