package lox_interpreter

type ExprVisitor interface {
	VisitAssignExpr(expr *AssignExpr) (interface{}, error)
	VisitLogicalExpr(expr *LogicalExpr) (interface{}, error)
	VisitTernaryExpr(expr *TernaryExpr) (interface{}, error)
	VisitBinaryExpr(expr *BinaryExpr) (interface{}, error)
	VisitGroupingExpr(expr *GroupingExpr) (interface{}, error)
	VisitLiteralExpr(expr *LiteralExpr) (interface{}, error)
	VisitUnaryExpr(expr *UnaryExpr) (interface{}, error)
	VisitCallExpr(expr *CallExpr) (interface{}, error)
	VisitGetExpr(expr *GetExpr) (interface{}, error)
	VisitSetExpr(expr *SetExpr) (interface{}, error)
	VisitVariableExpr(expr *VariableExpr) (interface{}, error)
	VisitThisExpr(expr *ThisExpr) (interface{}, error)
	VisitSuperExpr(expr *SuperExpr) (interface{}, error)
	VisitDictionaryExpr(expr *DictionaryExpr) (interface{}, error)
	VisitSelectExpr(expr *SelectExpr) (interface{}, error)
	VisitListExpr(expr *ListExpr) (interface{}, error)
}
type Expr interface {
	Accept(v ExprVisitor) (interface{}, error)
}

var _ Expr = (*AssignExpr)(nil)

type AssignExpr struct {
	name  Token
	value Expr
}

func NewAssignExpr(name Token, value Expr) *AssignExpr {
	return &AssignExpr{
		name,
		value,
	}
}

func (e *AssignExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitAssignExpr(e)
}

var _ Expr = (*LogicalExpr)(nil)

type LogicalExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func NewLogicalExpr(left Expr, operator Token, right Expr) *LogicalExpr {
	return &LogicalExpr{
		left,
		operator,
		right,
	}
}

func (e *LogicalExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitLogicalExpr(e)
}

var _ Expr = (*TernaryExpr)(nil)

type TernaryExpr struct {
	condition Expr
	question  Token
	left      Expr
	colon     Token
	right     Expr
}

func NewTernaryExpr(condition Expr, question Token, left Expr, colon Token, right Expr) *TernaryExpr {
	return &TernaryExpr{
		condition,
		question,
		left,
		colon,
		right,
	}
}

func (e *TernaryExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitTernaryExpr(e)
}

var _ Expr = (*BinaryExpr)(nil)

type BinaryExpr struct {
	left     Expr
	operator Token
	right    Expr
}

func NewBinaryExpr(left Expr, operator Token, right Expr) *BinaryExpr {
	return &BinaryExpr{
		left,
		operator,
		right,
	}
}

func (e *BinaryExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitBinaryExpr(e)
}

var _ Expr = (*GroupingExpr)(nil)

type GroupingExpr struct {
	expression Expr
}

func NewGroupingExpr(expression Expr) *GroupingExpr {
	return &GroupingExpr{
		expression,
	}
}

func (e *GroupingExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitGroupingExpr(e)
}

var _ Expr = (*LiteralExpr)(nil)

type LiteralExpr struct {
	value interface{}
}

func NewLiteralExpr(value interface{}) *LiteralExpr {
	return &LiteralExpr{
		value,
	}
}

func (e *LiteralExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitLiteralExpr(e)
}

var _ Expr = (*UnaryExpr)(nil)

type UnaryExpr struct {
	operator Token
	right    Expr
}

func NewUnaryExpr(operator Token, right Expr) *UnaryExpr {
	return &UnaryExpr{
		operator,
		right,
	}
}

func (e *UnaryExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitUnaryExpr(e)
}

var _ Expr = (*CallExpr)(nil)

type CallExpr struct {
	callee    Expr
	paren     Token
	arguments []Expr
}

func NewCallExpr(callee Expr, paren Token, arguments []Expr) *CallExpr {
	return &CallExpr{
		callee,
		paren,
		arguments,
	}
}

func (e *CallExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitCallExpr(e)
}

var _ Expr = (*GetExpr)(nil)

type GetExpr struct {
	object Expr
	name   Token
}

func NewGetExpr(object Expr, name Token) *GetExpr {
	return &GetExpr{
		object,
		name,
	}
}

func (e *GetExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitGetExpr(e)
}

var _ Expr = (*SetExpr)(nil)

type SetExpr struct {
	object Expr
	name   Token
	value  Expr
}

func NewSetExpr(object Expr, name Token, value Expr) *SetExpr {
	return &SetExpr{
		object,
		name,
		value,
	}
}

func (e *SetExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitSetExpr(e)
}

var _ Expr = (*VariableExpr)(nil)

type VariableExpr struct {
	name Token
}

func NewVariableExpr(name Token) *VariableExpr {
	return &VariableExpr{
		name,
	}
}

func (e *VariableExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitVariableExpr(e)
}

var _ Expr = (*ThisExpr)(nil)

type ThisExpr struct {
	keyword Token
}

func NewThisExpr(keyword Token) *ThisExpr {
	return &ThisExpr{
		keyword,
	}
}

func (e *ThisExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitThisExpr(e)
}

var _ Expr = (*SuperExpr)(nil)

type SuperExpr struct {
	keyword Token
	method  Token
}

func NewSuperExpr(keyword Token, method Token) *SuperExpr {
	return &SuperExpr{
		keyword,
		method,
	}
}

func (e *SuperExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitSuperExpr(e)
}

var _ Expr = (*DictionaryExpr)(nil)

type DictionaryExpr struct {
	mapExpr map[Token]Expr
}

func NewDictionaryExpr(mapExpr map[Token]Expr) *DictionaryExpr {
	return &DictionaryExpr{
		mapExpr,
	}
}

func (e *DictionaryExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitDictionaryExpr(e)
}

var _ Expr = (*SelectExpr)(nil)

type SelectExpr struct {
	object Expr
	name   Expr
}

func NewSelectExpr(object Expr, name Expr) *SelectExpr {
	return &SelectExpr{
		object,
		name,
	}
}

func (e *SelectExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitSelectExpr(e)
}

var _ Expr = (*ListExpr)(nil)

type ListExpr struct {
	values []Expr
}

func NewListExpr(values []Expr) *ListExpr {
	return &ListExpr{
		values,
	}
}

func (e *ListExpr) Accept(v ExprVisitor) (interface{}, error) {
	return v.VisitListExpr(e)
}
