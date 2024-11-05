package lox_interpreter

type StmtVisitor interface {
	VisitVarStmt(expr *VarStmt) (interface{}, error)
	VisitFunStmt(expr *FunStmt) (interface{}, error)
	VisitExpressionStmt(expr *ExpressionStmt) (interface{}, error)
	VisitIfStmt(expr *IfStmt) (interface{}, error)
	VisitPrintStmt(expr *PrintStmt) (interface{}, error)
	VisitWhileStmt(expr *WhileStmt) (interface{}, error)
	VisitBreakStmt(expr *BreakStmt) (interface{}, error)
	VisitReturnStmt(expr *ReturnStmt) (interface{}, error)
	VisitBlockStmt(expr *BlockStmt) (interface{}, error)
	VisitClassStmt(expr *ClassStmt) (interface{}, error)
}
type Stmt interface {
	Accept(v StmtVisitor) (interface{}, error)
}

var _ Stmt = (*VarStmt)(nil)

type VarStmt struct {
	name        Token
	initializer Expr
}

func NewVarStmt(name Token, initializer Expr) *VarStmt {
	return &VarStmt{
		name,
		initializer,
	}
}

func (e *VarStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitVarStmt(e)
}

var _ Stmt = (*FunStmt)(nil)

type FunStmt struct {
	name   Token
	params []Token
	body   []Stmt
}

func NewFunStmt(name Token, params []Token, body []Stmt) *FunStmt {
	return &FunStmt{
		name,
		params,
		body,
	}
}

func (e *FunStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitFunStmt(e)
}

var _ Stmt = (*ExpressionStmt)(nil)

type ExpressionStmt struct {
	expression Expr
}

func NewExpressionStmt(expression Expr) *ExpressionStmt {
	return &ExpressionStmt{
		expression,
	}
}

func (e *ExpressionStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitExpressionStmt(e)
}

var _ Stmt = (*IfStmt)(nil)

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

func NewIfStmt(condition Expr, thenBranch Stmt, elseBranch Stmt) *IfStmt {
	return &IfStmt{
		condition,
		thenBranch,
		elseBranch,
	}
}

func (e *IfStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitIfStmt(e)
}

var _ Stmt = (*PrintStmt)(nil)

type PrintStmt struct {
	expression Expr
}

func NewPrintStmt(expression Expr) *PrintStmt {
	return &PrintStmt{
		expression,
	}
}

func (e *PrintStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitPrintStmt(e)
}

var _ Stmt = (*WhileStmt)(nil)

type WhileStmt struct {
	condition Expr
	body      Stmt
}

func NewWhileStmt(condition Expr, body Stmt) *WhileStmt {
	return &WhileStmt{
		condition,
		body,
	}
}

func (e *WhileStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitWhileStmt(e)
}

var _ Stmt = (*BreakStmt)(nil)

type BreakStmt struct {
	keyword Token
}

func NewBreakStmt(keyword Token) *BreakStmt {
	return &BreakStmt{
		keyword,
	}
}

func (e *BreakStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitBreakStmt(e)
}

var _ Stmt = (*ReturnStmt)(nil)

type ReturnStmt struct {
	keyword Token
	value   Expr
}

func NewReturnStmt(keyword Token, value Expr) *ReturnStmt {
	return &ReturnStmt{
		keyword,
		value,
	}
}

func (e *ReturnStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitReturnStmt(e)
}

var _ Stmt = (*BlockStmt)(nil)

type BlockStmt struct {
	statements []Stmt
}

func NewBlockStmt(statements []Stmt) *BlockStmt {
	return &BlockStmt{
		statements,
	}
}

func (e *BlockStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitBlockStmt(e)
}

var _ Stmt = (*ClassStmt)(nil)

type ClassStmt struct {
	name       Token
	superClass *VariableExpr
	methods    []*FunStmt
}

func NewClassStmt(name Token, superClass *VariableExpr, methods []*FunStmt) *ClassStmt {
	return &ClassStmt{
		name,
		superClass,
		methods,
	}
}

func (e *ClassStmt) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitClassStmt(e)
}
