package codecrafters_interpreter_go

import "fmt"

type FunctionType string

const (
	NONE        FunctionType = "NONE"
	FUNCTION    FunctionType = "FUNCTION"
	METHOD      FunctionType = "METHOD"
	INITIALIZER FunctionType = "INITIALIZER"
)

type ClassType string

const (
	NONE_CLASS ClassType = "NONE_CLASS"
	CLS        ClassType = "CLASS"
	SUBCLASS   ClassType = "SUBCLASS"
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
	currentClass     ClassType
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

func (r *Resolver) VisitVarStmt(stmt *VarStmt) (_ interface{}, err error) {
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

func (r *Resolver) VisitFunStmt(stmt *FunStmt) (_ interface{}, err error) {
	err = r.declare(stmt.name)
	if err != nil {
		return
	}

	r.define(stmt.name)

	err = r.resolveFunction(stmt, FUNCTION)
	return
}

func (r *Resolver) VisitClassStmt(expr *ClassStmt) (_ interface{}, err error) {
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

	if expr.superClass != nil {
		r.currentClass = SUBCLASS
		if expr.name.Lexeme == expr.superClass.name.Lexeme {
			return nil, NewCompileError(expr.superClass.name, "A class cannot inherit from itself.")
		}

		err = r.ResolveExpression(expr.superClass)
		if err != nil {
			return
		}
	}

	if expr.superClass != nil {
		r.beginScope()
		r.scope[len(r.scope)-1]["super"] = true
	}
	r.beginScope()
	defer r.endScope()
	defer func() {
		if expr.superClass != nil {
			r.endScope()
		}
	}()

	r.scope[len(r.scope)-1]["this"] = true

	for _, method := range expr.methods {
		functionType := METHOD
		if method.name.Lexeme == "init" {
			functionType = INITIALIZER
		}

		err = r.resolveFunction(method, functionType)
		if err != nil {
			return
		}
	}

	return nil, nil
}

func (r *Resolver) VisitExpressionStmt(expr *ExpressionStmt) (interface{}, error) {
	return nil, r.ResolveExpression(expr.expression)
}

func (r *Resolver) VisitIfStmt(expr *IfStmt) (_ interface{}, err error) {
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

func (r *Resolver) VisitPrintStmt(expr *PrintStmt) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.expression)
	return
}

func (r *Resolver) VisitWhileStmt(expr *WhileStmt) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.condition)
	if err != nil {
		return
	}

	err = r.ResolveStatements(expr.body)
	return
}

func (r *Resolver) VisitBreakStmt(expr *BreakStmt) (_ interface{}, err error) {
	return
}

func (r *Resolver) VisitReturnStmt(expr *ReturnStmt) (_ interface{}, err error) {
	if r.currentFunction == NONE {
		return nil, NewCompileError(expr.keyword, "Cannot return from top-level code.")
	}
	if r.currentFunction == INITIALIZER {
		return nil, NewCompileError(expr.keyword, "Cannot return a value from an initializer.")
	}

	if expr.value != nil {
		err = r.ResolveExpression(expr.value)
	}

	return
}

func (r *Resolver) VisitBlockStmt(stmt *BlockStmt) (_ interface{}, err error) {
	r.beginScope()
	defer r.endScope()

	err = r.ResolveStatements(stmt.statements...)
	return nil, err
}

func (r *Resolver) VisitAssignExpr(expr *AssignExpr) (_ interface{}, err error) {
	err = r.ResolveExpression(expr)
	if err != nil {
		return
	}

	err = r.resolveLocal(expr, expr.name)
	return
}

func (r *Resolver) VisitLogicalExpr(expr *LogicalExpr) (interface{}, error) {
	err := r.ResolveExpression(expr.left)
	if err != nil {
		return nil, err
	}

	err = r.ResolveExpression(expr.right)
	return nil, err
}

func (r *Resolver) VisitTernaryExpr(expr *TernaryExpr) (_ interface{}, err error) {
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

func (r *Resolver) VisitBinaryExpr(expr *BinaryExpr) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.left)
	if err != nil {
		return
	}

	err = r.ResolveExpression(expr.right)
	return
}

func (r *Resolver) VisitGroupingExpr(expr *GroupingExpr) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.expression)
	return
}

func (r *Resolver) VisitLiteralExpr(expr *LiteralExpr) (_ interface{}, err error) {
	return
}

func (r *Resolver) VisitUnaryExpr(expr *UnaryExpr) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.right)
	return
}

func (r *Resolver) VisitCallExpr(expr *CallExpr) (_ interface{}, err error) {
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

func (r *Resolver) VisitGetExpr(expr *GetExpr) (_ interface{}, err error) {
	return nil, r.ResolveExpression(expr.object)
}

func (r *Resolver) VisitSetExpr(expr *SetExpr) (_ interface{}, err error) {
	err = r.ResolveExpression(expr.value)
	if err != nil {
		return
	}

	return nil, r.ResolveExpression(expr.object)
}

func (r *Resolver) VisitVariableExpr(expr *VariableExpr) (interface{}, error) {
	if v, ok := r.scope[len(r.scope)-1][expr.name.Lexeme]; ok && !v {
		return nil, NewCompileError(expr.name, "Cannot read local variable in its own initializer.")
	}

	return nil, r.resolveLocal(expr, expr.name)
}

func (r *Resolver) VisitThisExpr(expr *ThisExpr) (_ interface{}, err error) {
	if !r.isCurrentlyClass {
		return nil, NewCompileError(expr.keyword, "Cannot use 'this' outside of a class.")
	}

	err = r.resolveLocal(expr, expr.keyword)
	return
}

func (r *Resolver) VisitSuperExpr(expr *SuperExpr) (interface{}, error) {
	if r.currentClass == NONE_CLASS {
		return nil, NewCompileError(expr.keyword, "Cannot use 'super' outside of a class.")
	} else if r.currentClass != SUBCLASS {
		return nil, NewCompileError(expr.keyword, "Cannot use 'super' in a class with no superclass.")
	}
	err := r.resolveLocal(expr, expr.keyword)
	if err != nil {
		return nil, err
	}

	return nil, nil
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

func (r *Resolver) resolveFunction(stmt *FunStmt, functionType FunctionType) (err error) {
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
