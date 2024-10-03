package main

import (
	"fmt"
	codecrafters_interpreter_go "github.com/codecrafters-io/interpreter-starter-go"
	"log"
	"os"
	"strings"
)

func main() {
	log.SetFlags(log.Lmsgprefix)
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Fprintln(os.Stderr, "Logs from your program will appear here!")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: ./your_program.sh tokenize <filename>")
		os.Exit(1)
	}

	command := os.Args[1]

	if command != "tokenize" {
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

	s := codecrafters_interpreter_go.Scanner{Source: string(fileContents)}
	tokens, err := s.ScanTokens()
	for _, t := range tokens {
		arguments := []any{strings.ToUpper(string(t.Type)), t.Lexeme}

		if t.Literal != nil {
			arguments = append(arguments, t.Literal)
		} else {
			arguments = append(arguments, "null")
		}

		fmt.Printf("%s %s %s\n", arguments...)
	}

	if err != nil {
		if strings.Contains(err.Error(), "Unexpected Character") {
			os.Exit(65)
		}
	}
}
