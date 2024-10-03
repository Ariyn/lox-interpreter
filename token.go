package codecrafters_interpreter_go

type TokenType string

const (
	// 단일 단어 토큰
	LEFT_PAREN  TokenType = "left paren"
	RIGHT_PAREN TokenType = "right paren"
	LEFT_BRACE  TokenType = "left brace"
	RIGHT_BRACE TokenType = "right brace"
	COMMA       TokenType = "comma"
	DOT         TokenType = "dot"
	MINUS       TokenType = "minus"
	PLUS        TokenType = "plus"
	SEMICOLON   TokenType = "semicolon"
	SLASH       TokenType = "slash"
	STAR        TokenType = "star"

	// 1~2 글자 토큰
	BANG          TokenType = "bang"
	BANG_EQUAL    TokenType = "bang equal"
	EQUAL         TokenType = "equal"
	EQUAL_EQUAL   TokenType = "equal equal"
	GREATER       TokenType = "greater"
	GREATER_EQUAL TokenType = "greater equal"
	LESS          TokenType = "less"
	LESS_EQUAL    TokenType = "less equal"

	// 리터럴
	IDENTIFIER TokenType = "identifier"
	STRING     TokenType = "string"
	NUMBER     TokenType = "number"

	// 키워드
	AND    TokenType = "add"
	CLASS  TokenType = "class"
	ELSE   TokenType = "else"
	FALSE  TokenType = "false"
	FUN    TokenType = "fun"
	FOR    TokenType = "for"
	IF     TokenType = "if"
	NIL    TokenType = "nil"
	OR     TokenType = "or"
	PRINT  TokenType = "print"
	RETURN TokenType = "return"
	SUPER  TokenType = "super"
	THIS   TokenType = "this"
	TRUE   TokenType = "true"
	VAR    TokenType = "var"
	WHILE  TokenType = "while"

	EOF TokenType = "eof"
)

var KeywordsMap = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"print":  PRINT,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
}

type Token struct {
	Type       TokenType
	Lexeme     string
	Literal    any
	LineNumber int
}

func (t Token) String() string {
	return string(t.Type) + " " + t.Lexeme + " " // + string(t.Literal)
}
