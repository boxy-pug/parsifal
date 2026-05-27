package jsonparser

type tokenType string

type token struct {
	Type    tokenType
	Literal string
	Line    int
}

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	COMMA = ","
	COLON = ":"

	STRING = "STRING"
	NUMBER = "NUMBER"
	BOOL   = "BOOL"
	NULL   = ""
)

// newToken is a little helper to make a token based on
// tokentype and literal
func newToken(tokenType tokenType, literal string, line int) token {
	return token{Type: tokenType, Literal: literal, Line: line}
}
