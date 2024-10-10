package codecrafters_interpreter_go

type StmtVisitor interface {
	VisitVarStmt(expr *Var) (interface{}, error)
	VisitExpressionStmt(expr *Expression) (interface{}, error)
	VisitIfStmt(expr *If) (interface{}, error)
	VisitPrintStmt(expr *Print) (interface{}, error)
	VisitWhileStmt(expr *While) (interface{}, error)
	VisitBreakStmt(expr *Break) (interface{}, error)
	VisitBlockStmt(expr *Block) (interface{}, error)
}

type Stmt interface {
	Accept(v StmtVisitor) (interface{}, error)
}

var _ Stmt = (*Var)(nil)

type Var struct {
	name        Token
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

var _ Stmt = (*If)(nil)

type If struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func NewIf(condition Expr, thenBranch Stmt, elseBranch Stmt) *If {
	return &If{
		condition,
		thenBranch,
		elseBranch,
	}
}

func (e *If) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitIfStmt(e)
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

var _ Stmt = (*While)(nil)

type While struct {
	condition Expr
	body      Stmt
}

func NewWhile(condition Expr, body Stmt) *While {
	return &While{
		condition,
		body,
	}
}

func (e *While) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitWhileStmt(e)
}

var _ Stmt = (*Break)(nil)

type Break struct {
	keyword Token
}

func NewBreak(keyword Token) *Break {
	return &Break{
		keyword,
	}
}

func (e *Break) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitBreakStmt(e)
}

var _ Stmt = (*Block)(nil)

type Block struct {
	statements []Stmt
}

func NewBlock(statements []Stmt) *Block {
	return &Block{
		statements,
	}
}

func (e *Block) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitBlockStmt(e)
}
