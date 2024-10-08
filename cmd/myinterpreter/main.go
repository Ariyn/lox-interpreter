package main

import (
	"fmt"
	lox "github.com/codecrafters-io/interpreter-starter-go"
	"log"
	"os"
	"strings"
)

var commandMap = map[string]bool{
	"tokenize": true,
	"parse":    true,
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
	tokens, err := s.ScanTokens()

	if command == "tokenize" {
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
	} else if command == "parse" {
		parser := lox.NewParser(tokens)
		expr := parser.Parse()

		if expr == nil {
			os.Exit(65)
		}

		printer := lox.AstPrinter{}
		fmt.Println(printer.Print(expr))
	}

	if err != nil {
		if strings.Contains(err.Error(), "Unexpected character") {
			os.Exit(65)
		} else if strings.Contains(err.Error(), "Unterminated string") {
			os.Exit(65)
		}
	}
}
