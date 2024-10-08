package main

import (
	"fmt"
	lex "github.com/codecrafters-io/interpreter-starter-go"
)

func main() {
	epression := lex.NewBinary(
		lex.NewUnary(
			lex.Token{lex.MINUS, "-", nil, 1},
			lex.NewLiteral(123),
		),
		lex.Token{lex.STAR, "*", nil, 1},
		lex.NewGrouping(lex.NewLiteral(45.67)))

	printer := lex.AstPrinter{}
	fmt.Println(printer.Print(epression))

	expression2 := lex.NewBinary(
		lex.NewGrouping(
			lex.NewBinary(
				lex.NewLiteral(1),
				lex.Token{lex.PLUS, "+", nil, 1},
				lex.NewLiteral(2),
			),
		),
		lex.Token{lex.STAR, "*", nil, 1},
		lex.NewGrouping(
			lex.NewBinary(
				lex.NewLiteral(4),
				lex.Token{lex.MINUS, "-", nil, 1},
				lex.NewLiteral(3),
			),
		),
	)
	rpnPrinter := lex.RPNAstPrinter{}
	fmt.Println(rpnPrinter.Print(expression2))
}
