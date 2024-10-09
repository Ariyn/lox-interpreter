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
               | statement l

varDecl        → "var" IDENTIFIER ( "=" expression )? ";" ;

statement      → exprStmt
               | ifStmt
               | printStmt
               | block ;

exprStmt       → expression ";" ;
ifStmt         → "if" "(" expression ")" statement ( "else" statement )? ;
printStmt      → "print" expression ";" ;
block          → "{" declaration* "}" ;

expression     → assignment ;
assignment     → IDENTIFIER "=" assignment
               | ternary ;
ternary        → comma ( "?" comma ":" comma )* ;
comma          → equality ( "," comma )*
equality       → comparison ( ( "!=" | "==" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "-" | "+" ) factor )* ;
factor         → unary ( ( "/" | "*" ) unary )* ;
unary          → ( "!" | "-" ) unary
               | primary ;
primary        → NUMBER | STRING | "true" | "false" | "nil"
               | "(" expression ")"
               | IDENTIFIER ;
*/

type Parser struct {
	tokens  []Token
	current int
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

	stmt, err := p.Statement()
	if err != nil {
		p.synchronize()
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

	return p.expressionStatement()
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
	expr, err := p.ternary()
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

	return p.primary()
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
		}

		p.advance()
	}
}
