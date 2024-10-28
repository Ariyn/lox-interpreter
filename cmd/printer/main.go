package main

import (
	"fmt"
	lex "github.com/ariyn/lox_interpreter"
)

func main() {
	epression := lex.NewBinaryExpr(
		lex.NewUnaryExpr(
			lex.Token{lex.MINUS, "-", nil, 1},
			lex.NewLiteralExpr(123),
		),
		lex.Token{lex.STAR, "*", nil, 1},
		lex.NewGroupingExpr(lex.NewLiteralExpr(45.67)))

	printer := lex.AstPrinter{}
	fmt.Println(printer.Print([]lex.Stmt{lex.NewExpressionStmt(epression)}))

	expression2 := lex.NewBinaryExpr(
		lex.NewGroupingExpr(
			lex.NewBinaryExpr(
				lex.NewLiteralExpr(1),
				lex.Token{lex.PLUS, "+", nil, 1},
				lex.NewLiteralExpr(2),
			),
		),
		lex.Token{lex.STAR, "*", nil, 1},
		lex.NewGroupingExpr(
			lex.NewBinaryExpr(
				lex.NewLiteralExpr(4),
				lex.Token{lex.MINUS, "-", nil, 1},
				lex.NewLiteralExpr(3),
			),
		),
	)
	rpnPrinter := lex.RPNAstPrinter{}
	fmt.Println(rpnPrinter.Print(expression2))
}
