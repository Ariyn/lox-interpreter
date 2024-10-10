package codecrafters_interpreter_go

import (
	"fmt"
)

type ParseError struct {
	Token   Token
	Message string
}

func (p *ParseError) Error() string {
	if p.Token.Type == EOF {
		return fmt.Sprintf("%d at end: %s", p.Token.LineNumber, p.Message)
	}
	return fmt.Sprintf("%d at '%s': %s", p.Token.LineNumber, p.Token.Lexeme, p.Message)
}

func newParseError(token Token, message string) error {
	return &ParseError{token, message}
}

/*
program        → declaration* EOF ;

declaration    → varDecl
               | funDecl
               | statement ;

varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;
funDecl        → "fun" IDENTIFIER "(" parameters? ")" block ;
parameters     → IDENTIFIER ( "," IDENTIFIER )* ;

statement      → exprStmt
               | ifStmt
               | printStmt
               | whileStmt
               | forStmt
               | jumpStmt
               | block ;

jumpStmt       → breakStmt
               | returnStmt ;

exprStmt       → expression ";" ;
ifStmt         → "if" "(" expression ")" statement ( "else" statement )? ;
printStmt      → "print" expression ";" ;
whileStmt      → "while" "(" expression ")" loopStatement ;
forStmt        → "for" "(" ( varDecl | exprStmt | ";" )
						   expression? ";"
						   expression? ")" loopStatement ;
breakStmt      → "break" ";" ;
returnStmt     → "return" expression? ";" ;
block          → "{" declaration* "}" ;

expression     → assignment ;
assignment     → IDENTIFIER "=" assignment
               | logic_and ;

logic_or       → logic_and ( "or" logic_and )* ;
logic_and      → ternary ( "and" ternary )* ;

ternary        → comma ( "?" comma ":" comma )* ;
comma          → equality ( "," comma )*
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary | call ;
call           → primary ( "(" arguments? ")" )? ;
arguments      → expression ( "," expression )* ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")"
               | IDENTIFIER ;
*/

type Parser struct {
	tokens   []Token
	current  int
	isInLoop bool
	isInFun  []bool
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() ([]Stmt, error) {
	var statements []Stmt
	for !p.isAtEnd() {
		stmt, err := p.Declaration()
		if err != nil {
			return nil, err
		}
		statements = append(statements, stmt)
	}

	return statements, nil
}

func (p *Parser) Declaration() (Stmt, error) {
	if p.match(VAR) {
		stmt, err := p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil, err
		}

		return stmt, nil
	}

	if p.match(FUN) {
		stmt, err := p.funDeclaration()
		if err != nil {
			p.synchronize()
			return nil, err
		}

		return stmt, nil
	}

	stmt, err := p.Statement()
	if err != nil {
		return nil, err
	}

	return stmt, nil
}

func (p *Parser) varDeclaration() (Stmt, error) {
	identifier, err := p.identifier()
	if err != nil {
		return nil, err
	}

	var initializer Expr
	if p.match(EQUAL) {
		initializer, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}

	err = p.consume(SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil, err
	}

	return NewVar(identifier, initializer), nil
}

func (p *Parser) funDeclaration() (Stmt, error) {
	p.isInFun = append(p.isInFun, true)
	defer func() { p.isInFun = p.isInFun[:len(p.isInFun)-1] }()

	identifier, err := p.identifier()
	if err != nil {
		return nil, err
	}

	err = p.consume(LEFT_PAREN, "Expect '(' after function name.")
	if err != nil {
		return nil, err
	}

	parameters, err := p.parameters()
	if err != nil {
		return nil, err
	}

	err = p.consume(RIGHT_PAREN, "Expect ')' after parameters.")
	if err != nil {
		return nil, err
	}

	err = p.consume(LEFT_BRACE, "Expect '{' before function body.")
	if err != nil {
		return nil, err
	}

	block, err := p.blockStatement()
	if err != nil {
		return nil, err
	}

	return NewFun(identifier, parameters, block), nil
}

func (p *Parser) parameters() ([]Token, error) {
	var parameters []Token
	for {
		if len(parameters) >= 255 {
			return nil, newParseError(p.peek(), "Cannot have more than 255 parameters.")
		}

		if p.check(RIGHT_PAREN) {
			break
		}

		parameter, err := p.identifier()
		if err != nil {
			return nil, err
		}

		parameters = append(parameters, parameter)

		if !p.match(COMMA) {
			break
		}
	}

	return parameters, nil
}

func (p *Parser) identifier() (Token, error) {
	if p.check(IDENTIFIER) {
		return p.advance(), nil
	}

	return Token{}, newParseError(p.peek(), "Expect identifier.")
}

func (p *Parser) Statement() (Stmt, error) {
	if p.match(PRINT) {
		return p.printStatement()
	}
	if p.match(LEFT_BRACE) {
		return p.blockStatement()
	}
	if p.match(IF) {
		return p.ifStatement()
	}
	if p.match(WHILE) {
		return p.whileStatement()
	}
	if p.match(FOR) {
		return p.forStatement()
	}
	if p.match(BREAK) {
		if !p.isInLoop {
			return nil, newParseError(p.previous(), "Expect break statement inside loop.")
		}

		return p.breakStatement()
	}
	if p.match(RETURN) {
		if len(p.isInFun) == 0 {
			return nil, newParseError(p.previous(), "Expect return statement inside function.")
		}

		return p.returnStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) whileStatement() (Stmt, error) {
	p.isInLoop = true
	defer func() {
		p.isInLoop = false
	}()

	err := p.consume(LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.Expression()
	if err != nil {
		return nil, err
	}

	err = p.consume(RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil, err
	}

	body, err := p.Statement()
	if err != nil {
		return nil, err
	}

	return NewWhile(condition, body), nil
}

/*
`for(var i=0; i<10; i=i+1) foo();` equals to `var i = 0; while(i < 10) {foo(); i=i+1}`
*/
func (p *Parser) forStatement() (Stmt, error) {
	p.isInLoop = true
	defer func() {
		p.isInLoop = false
	}()

	err := p.consume(LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil, err
	}

	var initializer Stmt
	if p.match(VAR) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	} else if p.match(SEMICOLON) {
		initializer = nil
	} else {
		initializer, err = p.expressionStatement()
		if err != nil {
			return nil, err
		}
	}

	var condition Expr
	if p.check(SEMICOLON) {
		condition = NewLiteral(true)
	} else {
		condition, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}

	err = p.consume(SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil, err
	}

	var increment Expr
	if p.check(RIGHT_PAREN) {
		increment = nil
	} else {
		increment, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}

	err = p.consume(RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil, err
	}

	body, err := p.Statement()
	if err != nil {
		return nil, err
	}

	whileStatement := NewWhile(condition, NewBlock([]Stmt{body, NewExpression(increment)}))
	if initializer != nil {
		return NewBlock([]Stmt{initializer, whileStatement}), nil
	}

	return whileStatement, nil
}

func (p *Parser) breakStatement() (Stmt, error) {
	breakToken := p.previous()

	err := p.consume(SEMICOLON, "Expect ';' after 'break'.")
	if err != nil {
		return nil, err
	}

	return NewBreak(breakToken), nil
}

func (p *Parser) returnStatement() (stmt Stmt, err error) {
	returnToken := p.previous()

	var value Expr
	if !p.check(SEMICOLON) {
		value, err = p.Expression()
		if err != nil {
			return nil, err
		}
	}

	err = p.consume(SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil, err
	}

	return NewReturn(returnToken, value), nil
}

func (p *Parser) printStatement() (Stmt, error) {
	expr, err := p.Expression()
	if err != nil {
		return nil, err
	}

	err = p.consume(SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil, err
	}

	return NewPrint(expr), nil
}

func (p *Parser) ifStatement() (Stmt, error) {
	err := p.consume(LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil, err
	}

	condition, err := p.Expression()
	if err != nil {
		return nil, err
	}

	err = p.consume(RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil, err
	}

	thenBranch, err := p.Statement()
	if err != nil {
		return nil, err
	}

	var elseBranch Stmt
	if p.match(ELSE) {
		elseBranch, err = p.Statement()
		if err != nil {
			return nil, err
		}
	}

	return NewIf(condition, thenBranch, elseBranch), nil
}

func (p *Parser) blockStatement() (Stmt, error) {
	var statements []Stmt
	for !p.check(RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.Declaration()
		if err != nil {
			return nil, err
		}

		statements = append(statements, stmt)
	}

	err := p.consume(RIGHT_BRACE, "Expect '}' after block.")
	return NewBlock(statements), err
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.Expression()
	if err != nil {
		return nil, err
	}

	err = p.consume(SEMICOLON, "Expect ';' after expression.")
	if err != nil {
		return nil, err
	}

	return NewExpression(expr), nil
}

func (p *Parser) Expression() (Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (Expr, error) {
	expr, err := p.or()
	if err != nil {
		return nil, err
	}

	if p.match(EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		if variable, ok := expr.(*Variable); ok {
			return NewAssign(variable.name, value), nil
		}

		return nil, newParseError(equals, "Invalid assignment target.")
	}

	return expr, nil
}

func (p *Parser) or() (Expr, error) {
	expr, err := p.and()
	if err != nil {
		return nil, err
	}

	for p.match(OR) {
		operator := p.previous()
		right, err := p.and()
		if err != nil {
			return nil, err
		}

		expr = NewLogical(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) and() (Expr, error) {
	expr, err := p.ternary()
	if err != nil {
		return nil, err
	}

	for p.match(AND) {
		operator := p.previous()
		right, err := p.ternary()
		if err != nil {
			return nil, err
		}

		expr = NewLogical(expr, operator, right)
	}

	return expr, nil
}

func (p *Parser) ternary() (Expr, error) {
	expr, err := p.comma()
	if err != nil {
		return nil, err
	}

	if p.match(QUESTION) {
		question := p.previous()
		trueExpr, err := p.comma()
		if err != nil {
			return nil, err
		}

		err = p.consume(COLON, "Expect ':' after true expression.")
		if err != nil {
			return nil, err
		}

		colon := p.previous()
		falseExpr, err := p.comma()
		if err != nil {
			return nil, err
		}

		expr = NewTernary(expr, question, trueExpr, colon, falseExpr)
	}

	return expr, nil
}

func (p *Parser) comma() (Expr, error) {
	if p.check(COMMA) {
		return nil, newParseError(p.peek(), "Expect Left-hand side of comma operator.")
	}

	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	for p.match(COMMA) {
		token := p.previous()
		right, err := p.comma()
		if err != nil {
			return nil, err
		}

		expr = NewBinary(expr, token, right)
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	if p.check(BANG_EQUAL, EQUAL_EQUAL) {
		return nil, newParseError(p.peek(), "Expect Left-hand side of equality operator.")
	}

	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		token := p.previous()
		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = NewBinary(expr, token, right)
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	if p.check(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		return nil, newParseError(p.peek(), "Expect Left-hand side of comparison operator.")
	}

	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		token := p.previous()
		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = NewBinary(expr, token, right)
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	if p.check(PLUS) {
		return nil, newParseError(p.peek(), "Expect Left-hand side of term operator.")
	}

	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(MINUS, PLUS) {
		token := p.previous()
		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = NewBinary(expr, token, right)
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	if p.check(SLASH, STAR) {
		return nil, newParseError(p.peek(), "Expect Left-hand side of factor operator.")
	}

	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(SLASH, STAR) {
		token := p.previous()
		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = NewBinary(expr, token, right)
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(BANG, MINUS) {
		token := p.previous()
		right, err := p.unary()

		if err != nil {
			return nil, err
		}

		return NewUnary(token, right), nil
	}

	return p.call()
}

func (p *Parser) call() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	if p.match(LEFT_PAREN) {
		var arguments []Expr

		arguments, err = p.arguments()
		if err != nil {
			return nil, err
		}

		err = p.consume(RIGHT_PAREN, "Expect ')' after arguments.")
		if err != nil {
			return nil, err
		}
		expr = NewCall(expr, p.previous(), arguments)
	}

	return expr, nil
}

func (p *Parser) arguments() (arguments []Expr, err error) {
	for {
		if len(arguments) >= 255 {
			return nil, newParseError(p.peek(), "Cannot have more than 255 arguments.")
		}

		if p.check(RIGHT_PAREN) {
			break
		}

		arg, err := p.Expression()
		if err != nil {
			return nil, err
		}

		arguments = append(arguments, arg)

		if !p.match(COMMA) {
			break
		}
	}

	return arguments, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(FALSE) {
		return NewLiteral(false), nil
	}
	if p.match(TRUE) {
		return NewLiteral(true), nil
	}
	if p.match(NIL) {
		return NewLiteral(nil), nil
	}

	if p.match(NUMBER, STRING) {
		return NewLiteral(p.previous().Literal), nil
	}

	if p.match(LEFT_PAREN) {
		expr, err := p.Expression()
		if err != nil {
			return nil, err
		}

		err = p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil, err
		}
		return NewGrouping(expr), nil
	}

	if p.match(IDENTIFIER) {
		return NewVariable(p.previous()), nil
	}

	return nil, newParseError(p.peek(), "Expect expression.")
}

func (p *Parser) consume(t TokenType, message string) (err error) {
	if p.check(t) {
		p.advance()
		return
	}

	return newParseError(p.peek(), message)
}

func (p *Parser) match(types ...TokenType) bool {
	if p.check(types...) {
		p.advance()
		return true
	}

	return false
}

func (p *Parser) check(types ...TokenType) bool {
	if p.tokens[p.current].Type == EOF {
		return false
	}

	for _, t := range types {
		if p.tokens[p.current].Type == t {
			return true
		}
	}

	return false
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return len(p.tokens) <= p.current+1
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == SEMICOLON {
			return
		}

		switch p.peek().Type {
		case CLASS:
			continue
		case FUN:
			continue
		case VAR:
			continue
		case FOR:
			continue
		case IF:
			continue
		case WHILE:
			continue
		case PRINT:
			continue
		case RETURN:
			return
		case BREAK:
			return
		}

		p.advance()
	}
}
