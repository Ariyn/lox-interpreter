package lox_interpreter

func isAlphaNumeric(c string) bool {
	return isAlphabet(c) || isDigit(c)
}

func isAlphabet(c string) bool {
	return ('a' <= c[0] && c[0] <= 'z') || ('A' <= c[0] && c[0] <= 'Z') || c[0] == '_'
}

func isDigit(c string) bool {
	return '0' <= c[0] && c[0] <= '9'
}
