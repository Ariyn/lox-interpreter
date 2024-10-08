package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		os.Exit(64)
	}

	outputDir := args[0]

	defineAst(outputDir, "Expr", []string{
		"Ternary  : Expr condition, Token question, Expr left, Token colon, Expr right",
		"Binary   : Expr left, Token operator, Expr right",
		"Grouping : Expr expression",
		"Literal  : Object value",
		"Unary    : Token operator, Expr right",
	})
}

func defineAst(outputDir string, baseName string, types []string) (err error) {
	path := outputDir + "/" + strings.ToLower(baseName) + ".go"
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	fmt.Fprintln(f, "package codecrafters_interpreter_go")

	fmt.Fprintln(f, "type Visitor interface {")
	for _, t := range types {
		tokens := strings.Split(t, ":")

		className := strings.TrimSpace(tokens[0])
		fields := strings.TrimSpace(tokens[1])
		fields = convertFieldTypeOrder(fields)

		fmt.Fprintf(f, "	Visit%sExpr(expr *%s) interface{}\n", className, className)
	}
	fmt.Fprintln(f, "}\n")

	fmt.Fprintf(f, `type %s interface {
	Accept(v Visitor) interface{}
}
`, baseName)

	for _, t := range types {
		tokens := strings.Split(t, ":")
		className := strings.TrimSpace(tokens[0])
		fields := strings.TrimSpace(tokens[1])
		fields = convertFieldTypeOrder(fields)

		defineType(f, baseName, className, fields)
	}
	return nil
}

func convertFieldTypeOrder(fields string) string {
	fieldList := strings.Split(fields, ", ")
	for i, field := range fieldList {
		fieldTokens := strings.Split(field, " ")

		typ := strings.TrimSpace(fieldTokens[0])
		name := strings.TrimSpace(fieldTokens[1])

		if typ == "Object" {
			typ = "interface{}"
		}
		fieldList[i] = name + " " + typ
	}
	return strings.Join(fieldList, ", ")
}

func defineType(f *os.File, baseName string, className string, fieldList string) {
	//var _ = (*Scanner)(nil)
	fmt.Fprintf(f, "var _ %s = (*%s)(nil)\n", baseName, className)
	fmt.Fprintf(f, "type %s struct {\n", className)

	fields := strings.Split(fieldList, ", ")
	for _, field := range fields {
		fmt.Fprintf(f, "	%s\n", field)
	}
	fmt.Fprintln(f, "}\n")

	fmt.Fprintf(f, "func New%s(%s) *%s {\n", strings.ToUpper(className[:1])+className[1:], fieldList, className)
	fmt.Fprintf(f, "	return &%s{\n", className)
	for _, field := range fields {
		fieldTokens := strings.Split(field, " ")
		fmt.Fprintf(f, "		%s,\n", fieldTokens[0])
	}
	fmt.Fprintln(f, "	}")

	fmt.Fprintln(f, "}\n")

	fmt.Fprintf(f, "func (e *%s) Accept(v Visitor) interface{} {\n", className)
	fmt.Fprintf(f, "	return v.Visit%sExpr(e)\n", className)
	fmt.Fprintln(f, "}\n")
}
