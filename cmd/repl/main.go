package main

import (
	"bufio"
	"fmt"
	lox "github.com/ariyn/lox_interpreter"
	"os"
	"strings"
)

func main() {
	interpreter := lox.NewInterpreter()

	printHelp()

	reader := bufio.NewReader(os.Stdin)

	for true {
		fmt.Print("> ")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		input = strings.TrimSpace(input)
		if input == ".exit" {
			break
		} else if input == "help" {
			printHelp()
		} else {
			scanner := lox.NewScanner(input)
			singleLineTokens, err := scanner.ScanTokens()
			if err != nil {
				fmt.Println(err)
				continue
			}

			parser := lox.NewParser(singleLineTokens)
			statements, err := parser.Parse()
			if err != nil {
				fmt.Println(err)
				continue
			}

			v, err := interpreter.Interpret(statements)
			if err != nil {
				fmt.Println(err)
				continue
			}

			if v != nil {
				fmt.Println(lox.Stringify(v))
			}
		}
	}
}

func printHelp() {
	fmt.Println(`Welcome to the Lox REPL!
Type in Lox code and press enter to run it.
Type '.exit' to exit the REPL.
Type 'help' to see this message again.`)
}
