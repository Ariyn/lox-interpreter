package codecrafters_interpreter_go

import "fmt"

type FunctionType string

const (
	NONE     FunctionType = "NONE"
	FUNCTION FunctionType = "FUNCTION"
	METHOD   FunctionType = "METHOD"
)

type CompileError struct {
	token   Token
	message string
}

func (r *CompileError) Error() string {
	return fmt.Sprintf("%d at '%s' %s", r.token.LineNumber, r.token.Lexeme, r.message)
}

func NewCompileError(token Token, message string) error {
	return &CompileError{token, message}
}

var _ ExprVisitor = (*Resolver)(nil)
var _ StmtVisitor = (*Resolver)(nil)

type Resolver struct {
	interpreter      *Interpreter
	scope            []map[string]bool
	currentFunction  FunctionType
	isCurrentlyClass bool
}

func NewResolver(interpreter *Interpreter) *Resolver {
	scope := make([]map[string]bool, 0)
	scope = append(scope, make(map[string]bool))

	return &Resolver{
		interpreter:      interpreter,
		scope:            scope,
		currentFunction:  NONE,
		isCurrentlyClass: false,
	}
}

func (r *Resolver) VisitVarStmt(stmt *Var) (_ interface{}, err error) {
	err = r.declare(stmt.name)
	if err != nil {
		return
	}

	if stmt.initializer != nil {
		err = r.ResolveExpression(stmt.initializer)
		if err != nil {
			return
		}
	}

	r.define(stmt.name)

	return
}

func (r *Resolver) declare(name Token) (err error) {
	scope := r.scope[len(r.scope)-1]
	if _, ok := scope[name.Lexeme]; ok {
		return NewCompileError(name, "Variable with this name already declared in this scope.")
	}
	scope[name.Lexeme] = false

	return nil
}

func (r *Resolver) define(name Token) {
	r.scope[len(r.scope)-1][name.Lexeme] = true
}

func (r *Resolver) VisitFunStmt(stmt *Fun) (_ interface{}, err error) {
	err = r.declare(stmt.name)
	if err != nil {
		return
	}

	r.define(stmt.name)

	err = r.resolveFunction(stmt, FUNCTION)
	return
}

func (r *Resolver) VisitClassStmt(expr *Class) (_ interface{}, err error) {
	isCurrentlyClass := r.isCurrentlyClass
	r.isCurrentlyClass = true
	defer func() {
		r.isCurrentlyClass = isCurrentlyClass
	}()

	err = r.declare(expr.name)
	if err != nil {
		return
	}

	r.define(expr.name)

	r.beginScope()
	defer r.endScope()

	r.scope[len(r.scope)-1]["this"] = true

	for _, method := range expr.methods {
		err = r.resolveFunction(method, METHOD)
		if err != nil {
			return
		}
	}

	return nil, nil
}

func (r *Resolver) VisitExpressionStmt(expr *Expression) (interface{}, error) {
	return nil, r.ResolveExpression(expr.expression)
}

func (r *Resolver) VisitIfStmt(expr *If) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.condition)
	if err != nil {
		return
	}

	err = r.ResolveStatements(expr.thenBranch)
	if err != nil {
		return
	}

	err = r.ResolveStatements(expr.elseBranch)
	return
}

func (r *Resolver) VisitPrintStmt(expr *Print) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.expression)
	return
}

func (r *Resolver) VisitWhileStmt(expr *While) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.condition)
	if err != nil {
		return
	}

	err = r.ResolveStatements(expr.body)
	return
}

func (r *Resolver) VisitBreakStmt(expr *Break) (_ interface{}, err error) {
	return
}

func (r *Resolver) VisitReturnStmt(expr *Return) (_ interface{}, err error) {
	if r.currentFunction == NONE {
		return nil, NewCompileError(expr.keyword, "Cannot return from top-level code.")
	}

	if expr.value != nil {
		err = r.ResolveExpression(expr.value)
	}

	return
}

func (r *Resolver) VisitBlockStmt(stmt *Block) (_ interface{}, err error) {
	r.beginScope()
	defer r.endScope()

	err = r.ResolveStatements(stmt.statements...)
	return nil, err
}

func (r *Resolver) VisitAssignExpr(expr *Assign) (_ interface{}, err error) {
	err = r.ResolveExpression(expr)
	if err != nil {
		return
	}

	err = r.resolveLocal(expr, expr.name)
	return
}

func (r *Resolver) VisitLogicalExpr(expr *Logical) (interface{}, error) {
	err := r.ResolveExpression(expr.left)
	if err != nil {
		return nil, err
	}

	err = r.ResolveExpression(expr.right)
	return nil, err
}

func (r *Resolver) VisitTernaryExpr(expr *Ternary) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.condition)
	if err != nil {
		return
	}

	err = r.ResolveExpression(expr.left)
	if err != nil {
		return
	}

	err = r.ResolveExpression(expr.right)
	return
}

func (r *Resolver) VisitBinaryExpr(expr *Binary) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.left)
	if err != nil {
		return
	}

	err = r.ResolveExpression(expr.right)
	return
}

func (r *Resolver) VisitGroupingExpr(expr *Grouping) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.expression)
	return
}

func (r *Resolver) VisitLiteralExpr(expr *Literal) (_ interface{}, err error) {
	return
}

func (r *Resolver) VisitUnaryExpr(expr *Unary) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.right)
	return
}

func (r *Resolver) VisitCallExpr(expr *Call) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.callee)
	if err != nil {
		return
	}

	for _, arg := range expr.arguments {
		err = r.ResolveExpression(arg)
		if err != nil {
			return
		}
	}

	return
}

func (r *Resolver) VisitGetExpr(expr *Get) (_ interface{}, err error) {
	return nil, r.ResolveExpression(expr.object)
}

func (r *Resolver) VisitSetExpr(expr *Set) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.value)
	if err != nil {
		return
	}

	return nil, r.ResolveExpression(expr.object)
}

func (r *Resolver) VisitVariableExpr(expr *Variable) (interface{}, error) {
	if v, ok := r.scope[len(r.scope)-1][expr.name.Lexeme]; ok && !v {
		return nil, NewCompileError(expr.name, "Cannot read local variable in its own initializer.")
	}

	return nil, r.resolveLocal(expr, expr.name)
}

func (r *Resolver) VisitThisExpr(expr *This) (_ interface{}, err error) {
	if !r.isCurrentlyClass {
		return nil, NewCompileError(expr.keyword, "Cannot use 'this' outside of a class.")
	}

	err = r.resolveLocal(expr, expr.keyword)
	return
}

func (r *Resolver) resolveLocal(expr Expr, name Token) (err error) {
	for i := len(r.scope) - 1; i >= 0; i-- {
		if _, ok := r.scope[i][name.Lexeme]; ok {
			_, err = r.interpreter.ResolveExpression(expr, len(r.scope)-1-i)
			return
		}
	}

	return NewCompileError(name, "Variable not found.")
}

func (r *Resolver) resolveFunction(stmt *Fun, functionType FunctionType) (err error) {
	enclosingFunction := r.currentFunction
	r.currentFunction = functionType
	defer func() {
		r.currentFunction = enclosingFunction
	}()

	// function require block as body. So don't need to begin new scope here.
	// TODO: check why example opens a scope at here. https://craftinginterpreters.com/resolving-and-binding.html
	//r.beginScope()
	//defer r.endScope()

	for _, param := range stmt.params {
		err = r.declare(param)
		if err != nil {
			return
		}

		r.define(param)
	}

	err = r.ResolveStatements(stmt.body)
	return
}

func (r *Resolver) beginScope() {
	r.scope = append(r.scope, make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scope = r.scope[:len(r.scope)-1]
}

func (r *Resolver) ResolveStatements(statements ...Stmt) (err error) {
	for _, stmt := range statements {
		_, err = stmt.Accept(r)
		if err != nil {
			return
		}
	}

	return nil
}

func (r *Resolver) ResolveExpression(expression Expr) (err error) {
	_, err = expression.Accept(r)
	return err
}
