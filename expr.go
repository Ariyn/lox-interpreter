package codecrafters_interpreter_go

type Visitor interface {
	VisitBinaryExpr(expr *Binary) interface{}
	VisitGroupingExpr(expr *Grouping) interface{}
	VisitLiteralExpr(expr *Literal) interface{}
	VisitUnaryExpr(expr *Unary) interface{}
}

type Expr interface {
	Accept(v Visitor) interface{}
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

func (e *Binary) Accept(v Visitor) interface{} {
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

func (e *Grouping) Accept(v Visitor) interface{} {
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

func (e *Literal) Accept(v Visitor) interface{} {
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

func (e *Unary) Accept(v Visitor) interface{} {
	return v.VisitUnaryExpr(e)
}
