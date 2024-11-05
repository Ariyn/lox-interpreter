package lox_interpreter

import (
	"fmt"
)

var NO_RETURN_AT_ROOT = true

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
	for k := range interpreter.env.Values {
		scope[len(scope)-1][k] = true
	}

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
		err = r.ResolveExpressions(stmt.initializer)
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

		err = r.ResolveExpressions(expr.superClass)
		if err != nil {
			return
		}
	}

	if expr.superClass != nil {
		r.beginScope()
		defer r.endScope()
		r.scope[len(r.scope)-1]["super"] = true
	}
	r.beginScope()
	defer r.endScope()

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
	return nil, r.ResolveExpressions(expr.expression)
}

func (r *Resolver) VisitIfStmt(expr *IfStmt) (_ interface{}, err error) {
	err = r.ResolveExpressions(expr.condition)
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
	err = r.ResolveExpressions(expr.expression)
	return
}

func (r *Resolver) VisitWhileStmt(expr *WhileStmt) (_ interface{}, err error) {
	err = r.ResolveExpressions(expr.condition)
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
	if NO_RETURN_AT_ROOT && r.currentFunction == NONE {
		return nil, NewCompileError(expr.keyword, "Cannot return from top-level code.")
	}
	if r.currentFunction == INITIALIZER {
		return nil, NewCompileError(expr.keyword, "Cannot return a value from an initializer.")
	}

	if expr.value != nil {
		err = r.ResolveExpressions(expr.value)
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
	err = r.ResolveExpressions(expr.value)
	if err != nil {
		return
	}

	err = r.resolveLocal(expr, expr.name)
	return
}

func (r *Resolver) VisitLogicalExpr(expr *LogicalExpr) (interface{}, error) {
	err := r.ResolveExpressions(expr.left)
	if err != nil {
		return nil, err
	}

	err = r.ResolveExpressions(expr.right)
	return nil, err
}

func (r *Resolver) VisitTernaryExpr(expr *TernaryExpr) (_ interface{}, err error) {
	err = r.ResolveExpressions(expr.condition)
	if err != nil {
		return
	}

	err = r.ResolveExpressions(expr.left)
	if err != nil {
		return
	}

	err = r.ResolveExpressions(expr.right)
	return
}

func (r *Resolver) VisitBinaryExpr(expr *BinaryExpr) (_ interface{}, err error) {
	err = r.ResolveExpressions(expr.left)
	if err != nil {
		return
	}

	err = r.ResolveExpressions(expr.right)
	return
}

func (r *Resolver) VisitGroupingExpr(expr *GroupingExpr) (_ interface{}, err error) {
	err = r.ResolveExpressions(expr.expression)
	return
}

func (r *Resolver) VisitLiteralExpr(expr *LiteralExpr) (_ interface{}, err error) {
	return
}

func (r *Resolver) VisitUnaryExpr(expr *UnaryExpr) (_ interface{}, err error) {
	err = r.ResolveExpressions(expr.right)
	return
}

func (r *Resolver) VisitCallExpr(expr *CallExpr) (_ interface{}, err error) {
	err = r.ResolveExpressions(expr.callee)
	if err != nil {
		return
	}

	for _, arg := range expr.arguments {
		err = r.ResolveExpressions(arg)
		if err != nil {
			return
		}
	}

	return
}

func (r *Resolver) VisitGetExpr(expr *GetExpr) (_ interface{}, err error) {
	return nil, r.ResolveExpressions(expr.object)
}

func (r *Resolver) VisitSetExpr(expr *SetExpr) (_ interface{}, err error) {
	err = r.ResolveExpressions(expr.value)
	if err != nil {
		return
	}

	return nil, r.ResolveExpressions(expr.object)
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

func (r *Resolver) VisitDictionaryExpr(expr *DictionaryExpr) (interface{}, error) {
	dict := make(map[string]Expr)
	for k, v := range expr.mapExpr {
		err := r.ResolveExpressions(v)
		if err != nil {
			return nil, err
		}

		if _, ok := dict[k.Lexeme]; ok {
			return nil, NewCompileError(k, "Duplicate key in dictionary.")
		}
		dict[k.Lexeme] = v
	}

	return dict, nil
}

func (r *Resolver) VisitSelectExpr(expr *SelectExpr) (interface{}, error) {
	err := r.ResolveExpressions(expr.object)
	if err != nil {
		return nil, err
	}

	if _, ok := expr.object.(*VariableExpr); !ok {
		return nil, NewCompileError(Token{}, "Only variable can have properties.")
	}

	err = r.ResolveExpressions(expr.name)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (r *Resolver) VisitListExpr(expr *ListExpr) (interface{}, error) {
	return nil, r.ResolveExpressions(expr.values...)
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

	r.beginScope()
	defer r.endScope()

	for _, param := range stmt.params {
		err = r.declare(param)
		if err != nil {
			return
		}

		r.define(param)
	}

	err = r.ResolveStatements(stmt.body...)
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

func (r *Resolver) ResolveExpressions(expressions ...Expr) (err error) {
	for _, expr := range expressions {
		_, err = expr.Accept(r)
		if err != nil {
			return err
		}
	}
	return nil
}

// Resolve resolves the given statements. This is the only entrance for the resolver.
func (r *Resolver) Resolve(statements ...Stmt) (err error) {
	err = r.ResolveStatements(statements...)
	if err != nil {
		return
	}

	for _, s := range r.scope {
		for k := range s {
			if !s[k] {
				return NewCompileError(Token{}, "Local variable '"+k+"' is not used.")
			}
		}
	}

	return nil
}
