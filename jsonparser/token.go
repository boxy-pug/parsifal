package jsonparser

type tokenType string

type token struct {
	Type    tokenType
	Literal string
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
