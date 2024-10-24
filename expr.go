package codecrafters_interpreter_go

type ExprVisitor interface {
	VisitAssignExpr(expr *Assign) (interface{}, error)
	VisitLogicalExpr(expr *Logical) (interface{}, error)
	VisitTernaryExpr(expr *Ternary) (interface{}, error)
	VisitBinaryExpr(expr *Binary) (interface{}, error)
	VisitGroupingExpr(expr *Grouping) (interface{}, error)
	VisitLiteralExpr(expr *Literal) (interface{}, error)
	VisitUnaryExpr(expr *Unary) (interface{}, error)
	VisitCallExpr(expr *Call) (interface{}, error)
	VisitGetExpr(expr *Get) (interface{}, error)
	VisitSetExpr(expr *Set) (interface{}, error)
	VisitVariableExpr(expr *Variable) (interface{}, error)
	VisitThisExpr(expr *This) (interface{}, error)
	VisitSuperExpr(expr *Super) (interface{}, error)
}

type Expr interface {
	Accept(v ExprVisitor) (interface{}, error)
}

var _ Expr = (*Assign)(nil)

type Assign struct {
	name  Token
	value Expr
}

func NewAssign(name Token, value Expr) *Assign {
	return &Assign{
		name,
		value,
	}
}

func (e *Assign) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitAssignExpr(e)
}

var _ Expr = (*Logical)(nil)

type Logical struct {
	left     Expr
	operator Token
	right    Expr
}

func NewLogical(left Expr, operator Token, right Expr) *Logical {
	return &Logical{
		left,
		operator,
		right,
	}
}

func (e *Logical) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitLogicalExpr(e)
}

var _ Expr = (*Ternary)(nil)

type Ternary struct {
	condition Expr
	question  Token
	left      Expr
	colon     Token
	right     Expr
}

func NewTernary(condition Expr, question Token, left Expr, colon Token, right Expr) *Ternary {
	return &Ternary{
		condition,
		question,
		left,
		colon,
		right,
	}
}

func (e *Ternary) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitTernaryExpr(e)
}

var _ Expr = (*Binary)(nil)

type Binary struct {
	left     Expr
	operator Token
	right    Expr
}

func NewBinary(left Expr, operator Token, right Expr) *Binary {
	return &Binary{
		left,
		operator,
		right,
	}
}

func (e *Binary) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitBinaryExpr(e)
}

var _ Expr = (*Grouping)(nil)

type Grouping struct {
	expression Expr
}

func NewGrouping(expression Expr) *Grouping {
	return &Grouping{
		expression,
	}
}

func (e *Grouping) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitGroupingExpr(e)
}

var _ Expr = (*Literal)(nil)

type Literal struct {
	value interface{}
}

func NewLiteral(value interface{}) *Literal {
	return &Literal{
		value,
	}
}

func (e *Literal) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitLiteralExpr(e)
}

var _ Expr = (*Unary)(nil)

type Unary struct {
	operator Token
	right    Expr
}

func NewUnary(operator Token, right Expr) *Unary {
	return &Unary{
		operator,
		right,
	}
}

func (e *Unary) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitUnaryExpr(e)
}

var _ Expr = (*Call)(nil)

type Call struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func NewCall(callee Expr, paren Token, arguments []Expr) *Call {
	return &Call{
		callee,
		paren,
		arguments,
	}
}

func (e *Call) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitCallExpr(e)
}

var _ Expr = (*Get)(nil)

type Get struct {
	object Expr
	name   Token
}

func NewGet(object Expr, name Token) *Get {
	return &Get{
		object,
		name,
	}
}

func (e *Get) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitGetExpr(e)
}

var _ Expr = (*Set)(nil)

type Set struct {
	object Expr
	name   Token
	value  Expr
}

func NewSet(object Expr, name Token, value Expr) *Set {
	return &Set{
		object,
		name,
		value,
	}
}

func (e *Set) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitSetExpr(e)
}

var _ Expr = (*Variable)(nil)

type Variable struct {
	name Token
}

func NewVariable(name Token) *Variable {
	return &Variable{
		name,
	}
}

func (e *Variable) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitVariableExpr(e)
}

var _ Expr = (*This)(nil)

type This struct {
	keyword Token
}

func NewThis(keyword Token) *This {
	return &This{
		keyword,
	}
}

func (e *This) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitThisExpr(e)
}

var _ Expr = (*Super)(nil)

type Super struct {
	keyword Token
	method  Token
}

func NewSuper(keyword Token, method Token) *Super {
	return &Super{
		keyword,
		method,
	}
}

func (e *Super) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitSuperExpr(e)
}
