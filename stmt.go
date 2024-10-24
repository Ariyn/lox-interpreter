package codecrafters_interpreter_go

type StmtVisitor interface {
	VisitVarStmt(expr *Var) (interface{}, error)
	VisitFunStmt(expr *Fun) (interface{}, error)
	VisitExpressionStmt(expr *Expression) (interface{}, error)
	VisitIfStmt(expr *If) (interface{}, error)
	VisitPrintStmt(expr *Print) (interface{}, error)
	VisitWhileStmt(expr *While) (interface{}, error)
	VisitBreakStmt(expr *Break) (interface{}, error)
	VisitReturnStmt(expr *Return) (interface{}, error)
	VisitBlockStmt(expr *Block) (interface{}, error)
	VisitClassStmt(expr *Class) (interface{}, error)
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

var _ Stmt = (*Fun)(nil)

type Fun struct {
	name   Token
	params []Token
	body   Stmt
}

func NewFun(name Token, params []Token, body Stmt) *Fun {
	return &Fun{
		name,
		params,
		body,
	}
}

func (e *Fun) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitFunStmt(e)
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

var _ Stmt = (*Return)(nil)

type Return struct {
	keyword Token
	value   Expr
}

func NewReturn(keyword Token, value Expr) *Return {
	return &Return{
		keyword,
		value,
	}
}

func (e *Return) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitReturnStmt(e)
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

var _ Stmt = (*Class)(nil)

type Class struct {
	name       Token
	superClass *Variable
	methods    []*Fun
}

func NewClass(name Token, superClass *Variable, methods []*Fun) *Class {
	return &Class{
		name,
		superClass,
		methods,
	}
}

func (e *Class) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitClassStmt(e)
}
