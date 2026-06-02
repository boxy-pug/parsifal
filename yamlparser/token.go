// Package yamlparser implements a recursive-decent YAML parser,
// handwritten as a learning project. Implements a subset of YAML.
package yamlparser

type tokenType string

type token struct {
	Type    tokenType
	Literal string
	Line    int
}

// The root object should be a map.

const (
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	LBRACE   = "{"
	RBRACE   = "}"
	LBRACKET = "["
	RBRACKET = "]"

	COMMA = ","
	COLON = ":"

	// DASH preceded by whitespace at beginning of line is
	// start or continuation of sequence
	DASH = "DASH"

	NEWLINE = "NEWLINE"

	// indentation spaces not tabs, has to be consistent within a block
	// but you could have various number of spaces across blocks.
	// what matters is "more indented than parent block" and "same indent as sibling block"
	INDENT = "INDENT" // col indent should be stored in literal
	DEDENT = "DEDENT"

	// Just lex any scalar type in here, and we'll figure out the type
	// when parsing
	SCALAR = "SCALAR"
	// Scalar types
	// STRING = "STRING" // no quotes, or double quotes with escape codes, and single quotes with double single quotes as only escape ''
	// NUMBER = "NUMBER" // int float, scientific notation(1e+12), hex 0x and octal (leading 0)
	// BOOL   = "BOOL"   // case insensitive: true, yes, y, on – false, no, off, n
	// NULL   = ""       // case insensitive null, empty value, empty string
)

// newToken is a little helper to make a token based on
// tokentype and literal
func newToken(tokenType tokenType, literal string, line int) token {
	return token{Type: tokenType, Literal: literal, Line: line}
}
