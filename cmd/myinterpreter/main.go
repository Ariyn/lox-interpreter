package main

import (
	"fmt"
	lox "github.com/codecrafters-io/interpreter-starter-go"
	"log"
	"os"
	"strings"
)

var commandMap = map[string]bool{
	"tokenize":  true,
	"parse":     true,
	"interpret": true,
}

func main() {
	log.SetFlags(log.Lmsgprefix)
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if _, ok := commandMap[command]; !ok {
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", command)
		os.Exit(1)
	}

	// Uncomment this block to pass the first stage
	//
	filename := os.Args[2]
	fileContents, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	s := lox.Scanner{Source: string(fileContents)}
	switch command {
	case "tokenize":
		err := tokenize(s)
		if strings.Contains(err.Error(), "Unexpected character") {
			os.Exit(65)
		} else if strings.Contains(err.Error(), "Unterminated string") {
			os.Exit(65)
		}
		break
	case "parse":
		err := parse(s)
		if err != nil {
			log.Println(err.Error())
			os.Exit(65)
		}
	case "interpret":
		err := interpret(s)
		if err != nil {
			log.Println(err.Error())
			os.Exit(70)
		}
	}
}

func tokenize(scanner lox.Scanner) (err error) {
	tokens, err := scanner.ScanTokens()

	for _, t := range tokens {
		format := "%s %s %s"
		arguments := []any{strings.ToUpper(string(t.Type)), t.Lexeme}

		if t.Literal != nil {
			if t.Type == lox.STRING {
				arguments = append(arguments, t.Literal)
			} else {
				if t.Literal == float64(int(t.Literal.(float64))) {
					arguments = append(arguments, fmt.Sprintf("%.1f", t.Literal.(float64)))
				} else {
					arguments = append(arguments, fmt.Sprintf("%g", t.Literal.(float64)))
				}
			}
		} else {
			arguments = append(arguments, "null")
		}

		fmt.Printf(format+"\n", arguments...)
	}

	if err != nil {
		return err
	}

	return nil
}

func parse(scanner lox.Scanner) (err error) {
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return
	}

	parser := lox.NewParser(tokens)
	expr, err := parser.Parse()

	if err != nil {
		return
	}

	printer := lox.AstPrinter{}
	fmt.Println(printer.Print(expr))

	return nil
}

func interpret(scanner lox.Scanner) (err error) {
	tokens, err := scanner.ScanTokens()
	if err != nil {
		return
	}

	parser := lox.NewParser(tokens)
	expr, err := parser.Parse()

	if err != nil {
		return
	}

	interpreter := lox.NewInterpreter()
	value, err := interpreter.Interpret(expr)

	if err != nil {
		return
	}

	fmt.Println(value)

	return nil
}
