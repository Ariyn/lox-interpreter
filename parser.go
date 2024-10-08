package codecrafters_interpreter_go

import "log"

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

func (p *Parser) Parse() Expr {
	if recover() != nil {
		return nil
	}
	return p.Expression()
}

func (p *Parser) Expression() Expr {
	return p.equality()
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(BANG_EQUAL, EQUAL_EQUAL) {
		token := p.previous()
		right := p.comparison()
		expr = NewBinary(expr, token, right)
	}

	return expr
}

func (p *Parser) match(types ...TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t TokenType) bool {
	if p.tokens[p.current].Type == EOF {
		return false
	}
	return p.tokens[p.current].Type == t
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
	return len(p.tokens) <= p.current
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) comparison() Expr {
	expr := p.term()

	for p.match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL) {
		token := p.previous()
		right := p.term()
		expr = NewBinary(expr, token, right)
	}

	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()

	for p.match(MINUS, PLUS) {
		token := p.previous()
		right := p.factor()
		log.Println("here", expr, token, right)
		expr = NewBinary(expr, token, right)
	}

	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()

	for p.match(SLASH, STAR) {
		token := p.previous()
		right := p.unary()
		expr = NewBinary(expr, token, right)
	}

	return expr
}

func (p *Parser) unary() Expr {
	if p.match(BANG, MINUS) {
		token := p.previous()
		right := p.unary()
		return NewUnary(token, right)
	}

	return p.primary()
}

func (p *Parser) primary() Expr {
	if p.match(FALSE) {
		return NewLiteral(false)
	}
	if p.match(TRUE) {
		return NewLiteral(true)
	}
	if p.match(NIL) {
		return NewLiteral(nil)
	}

	if p.match(NUMBER, STRING) {
		log.Println("here match literal", p.previous())
		return NewLiteral(p.previous().Literal)
	}

	if p.match(LEFT_PAREN) {
		expr := p.Expression()
		p.consume(RIGHT_PAREN, "Expect ')' after expression.")
		return NewGrouping(expr)
	}

	panic("Expect expression.")
}

func (p *Parser) consume(t TokenType, message string) {
	if p.check(t) {
		p.advance()
	}
	panic(message)
}
