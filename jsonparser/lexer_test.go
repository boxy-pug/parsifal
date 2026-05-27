package jsonparser

import (
	"strings"
	"testing"
)

type wantToken struct {
	type_ tokenType
	lit   string
	line  int
}

func assertTokens(t *testing.T, input string, wants []wantToken) {
	t.Helper()

	l := NewLexer(strings.NewReader(input))

	for i, want := range wants {
		tok := l.NextToken()

		if tok.Type != want.type_ {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, want.type_, tok.Type)
		}
		if tok.Literal != want.lit {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, want.lit, tok.Literal)
		}
		if tok.Line != want.line {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d", i, want.line, tok.Line)
		}

		if tok.Type == ILLEGAL {
			break
		}
	}
}

func TestJsonLexer_BasicTokensAndLines(t *testing.T) {
	input := `{
	  "name": "Alice",
	  "age": 30,
	  "deleted": null
}`

	assertTokens(t, input, []wantToken{
		{type_: LBRACE, lit: "{", line: 1},
		{type_: STRING, lit: "name", line: 2},
		{type_: COLON, lit: ":", line: 2},
		{type_: STRING, lit: "Alice", line: 2},
		{type_: COMMA, lit: ",", line: 2},
		{type_: STRING, lit: "age", line: 3},
		{type_: COLON, lit: ":", line: 3},
		{type_: NUMBER, lit: "30", line: 3},
		{type_: COMMA, lit: ",", line: 3},
		{type_: STRING, lit: "deleted", line: 4},
		{type_: COLON, lit: ":", line: 4},
		{type_: NULL, lit: "", line: 4},
		{type_: RBRACE, lit: "}", line: 5},
		{type_: EOF, lit: "", line: 5},
	})
}

func TestJsonLexer_IllegalToken(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []wantToken
	}{
		{
			name:  "at sign",
			input: `@`,
			want: []wantToken{
				{type_: ILLEGAL, lit: "", line: 1},
			},
		},
		{
			name:  "question mark",
			input: `?`,
			want: []wantToken{
				{type_: ILLEGAL, lit: "", line: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.want)
		})
	}
}

func TestJsonLexer_StringEscapes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []wantToken
	}{
		{
			name:  "escaped quote",
			input: `{"value": "quote\"here"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "value", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: `quote"here`, line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "backslash",
			input: `{"value": "back\\slash"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "value", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: `back\slash`, line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "slash",
			input: `{"value": "slash\/here"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "value", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "slash/here", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "newline",
			input: "{\"value\": \"line1\\nline2\"}",
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "value", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "line1\nline2", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "tab",
			input: "{\"value\": \"tab\\there\"}",
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "value", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "tab\there", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "carriage return",
			input: "{\"value\": \"ret\\rurn\"}",
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "value", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "ret\rurn", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "backspace",
			input: "{\"value\": \"back\\bpace\"}",
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "value", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "back\bpace", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "form feed",
			input: "{\"value\": \"form\\ffeed\"}",
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "value", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "form\ffeed", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.want)
		})
	}
}

func TestJsonLexer_ValidNumbers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []wantToken
	}{
		{name: "integer", input: `42`, want: []wantToken{{type_: NUMBER, lit: "42", line: 1}, {type_: EOF, lit: "", line: 1}}},
		{name: "negative", input: `-17`, want: []wantToken{{type_: NUMBER, lit: "-17", line: 1}, {type_: EOF, lit: "", line: 1}}},
		{name: "zero", input: `0`, want: []wantToken{{type_: NUMBER, lit: "0", line: 1}, {type_: EOF, lit: "", line: 1}}},
		{name: "float", input: `3.14159`, want: []wantToken{{type_: NUMBER, lit: "3.14159", line: 1}, {type_: EOF, lit: "", line: 1}}},
		{name: "negative float", input: `-0.5`, want: []wantToken{{type_: NUMBER, lit: "-0.5", line: 1}, {type_: EOF, lit: "", line: 1}}},
		{name: "small decimal", input: `0.001`, want: []wantToken{{type_: NUMBER, lit: "0.001", line: 1}, {type_: EOF, lit: "", line: 1}}},
		{name: "scientific positive", input: `1.5e10`, want: []wantToken{{type_: NUMBER, lit: "1.5e10", line: 1}, {type_: EOF, lit: "", line: 1}}},
		{name: "scientific negative", input: `2.99792458e-8`, want: []wantToken{{type_: NUMBER, lit: "2.99792458e-8", line: 1}, {type_: EOF, lit: "", line: 1}}},
		{name: "scientific upper e", input: `6.022E23`, want: []wantToken{{type_: NUMBER, lit: "6.022E23", line: 1}, {type_: EOF, lit: "", line: 1}}},
		{name: "large number", input: `9007199254740991`, want: []wantToken{{type_: NUMBER, lit: "9007199254740991", line: 1}, {type_: EOF, lit: "", line: 1}}},
		{name: "negative zero", input: `-0`, want: []wantToken{{type_: NUMBER, lit: "-0", line: 1}, {type_: EOF, lit: "", line: 1}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.want)
		})
	}
}

func TestJsonLexer_InvalidNumbers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []wantToken
	}{
		{name: "leading zero", input: `042`, want: []wantToken{{type_: ILLEGAL, lit: "", line: 1}}},
		{name: "leading decimal", input: `.5`, want: []wantToken{{type_: ILLEGAL, lit: "", line: 1}}},
		{name: "trailing decimal", input: `5.`, want: []wantToken{{type_: ILLEGAL, lit: "", line: 1}}},
		{name: "plus sign", input: `+42`, want: []wantToken{{type_: ILLEGAL, lit: "", line: 1}}},
		{name: "hex", input: `0xFF`, want: []wantToken{{type_: ILLEGAL, lit: "", line: 1}}},
		{name: "octal", input: `0o77`, want: []wantToken{{type_: ILLEGAL, lit: "", line: 1}}},
		{name: "binary", input: `0b1010`, want: []wantToken{{type_: ILLEGAL, lit: "", line: 1}}},
		{name: "nan", input: `NaN`, want: []wantToken{{type_: ILLEGAL, lit: "", line: 1}}},
		{name: "infinity", input: `Infinity`, want: []wantToken{{type_: ILLEGAL, lit: "", line: 1}}},
		{name: "negative infinity", input: `-Infinity`, want: []wantToken{{type_: ILLEGAL, lit: "", line: 1}}},
		{name: "quoted number", input: `"42"`, want: []wantToken{{type_: STRING, lit: "42", line: 1}, {type_: EOF, lit: "", line: 1}}},
		{name: "underscores", input: `1_000_000`, want: []wantToken{{type_: ILLEGAL, lit: "", line: 1}}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.want)
		})
	}
}

func TestJsonLexer_InvalidStrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []wantToken
	}{
		{
			name:  "unknown escape",
			input: `{"value": "bad\qescape"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "value", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: ILLEGAL, lit: "", line: 1},
			},
		},
		{
			name: "raw newline",
			input: `{"value": "line1
line2"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "value", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: ILLEGAL, lit: "", line: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.want)
		})
	}
}

// Unicode Escape Sequences
//
// The \uXXXX syntax lets you include any Unicode character using its four-digit hexadecimal code point:
//
// {
//   "cafe": "Caf\u00E9",
//   "degree": "72\u00B0F",
//   "copyright": "\u00A9 2025",
//   "yen": "\u00A5500",
//   "greek": "\u03B1\u03B2\u03B3",
//   "checkmark": "\u2713 Complete"
// }
//
// Characters outside the Basic Multilingual Plane (like many emoji) require a surrogate pair — two \uXXXX sequences. For example, the grinning face emoji (U+1F600) is encoded as \uD83D\uDE00.
// Characters That Don't Need Escaping
//
// Most printable Unicode characters can appear directly in JSON strings without escaping. You can include accented letters, CJK characters, and even emoji directly:
//
// {
//   "french": "Crème brûlée",
//   "japanese": "東京タワー",
//   "emoji": "Hello 👋 World 🌍",
//   "math": "π ≈ 3.14159"
// }

func TestJsonLexer_UnicodeStrings(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []wantToken
	}{
		{
			name:  "unicode cafe",
			input: `{"cafe": "Caf\u00E9"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "cafe", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "Café", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "unicode degree",
			input: `{"degree": "72\u00B0F"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "degree", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "72°F", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "unicode copyright",
			input: `{"copyright": "\u00A9 2025"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "copyright", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "© 2025", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "unicode greek",
			input: `{"greek": "\u03B1\u03B2\u03B3"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "greek", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "αβγ", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "surrogate pair emoji",
			input: `{"emoji": "Hello \uD83D\uDE00"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "emoji", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "Hello 😀", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "direct unicode",
			input: `{"french": "Crème brûlée", "japanese": "東京タワー"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "french", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "Crème brûlée", line: 1},
				{type_: COMMA, lit: ",", line: 1},
				{type_: STRING, lit: "japanese", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "東京タワー", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "direct emoji",
			input: `{"emoji": "Hello 👋 World 🌍"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "emoji", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "Hello 👋 World 🌍", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.want)
		})
	}
}

func TestJsonLexer_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []wantToken
	}{
		{
			name:  "empty object",
			input: `{}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "empty array",
			input: `[]`,
			want: []wantToken{
				{type_: LBRACKET, lit: "[", line: 1},
				{type_: RBRACKET, lit: "]", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "empty string",
			input: `{"key": ""}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "key", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "no whitespace at all",
			input: `{"a":1,"b":true,"c":null}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "a", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: NUMBER, lit: "1", line: 1},
				{type_: COMMA, lit: ",", line: 1},
				{type_: STRING, lit: "b", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: BOOL, lit: "true", line: 1},
				{type_: COMMA, lit: ",", line: 1},
				{type_: STRING, lit: "c", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: NULL, lit: "", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "excessive whitespace",
			input: "  \t  {  \n\n  \"x\"  \t :  \n  42  \n  }  \t\t",
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "x", line: 3},
				{type_: COLON, lit: ":", line: 3},
				{type_: NUMBER, lit: "42", line: 4},
				{type_: RBRACE, lit: "}", line: 5},
				{type_: EOF, lit: "", line: 5},
			},
		},
		{
			name:  "deeply nested",
			input: `{"a":{"b":{"c":[1,[2,{"d":3}]]}}}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "a", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "b", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "c", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: LBRACKET, lit: "[", line: 1},
				{type_: NUMBER, lit: "1", line: 1},
				{type_: COMMA, lit: ",", line: 1},
				{type_: LBRACKET, lit: "[", line: 1},
				{type_: NUMBER, lit: "2", line: 1},
				{type_: COMMA, lit: ",", line: 1},
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "d", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: NUMBER, lit: "3", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: RBRACKET, lit: "]", line: 1},
				{type_: RBRACKET, lit: "]", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "mixed whitespace types",
			input: "{\r\n\t\"key\"\t:\r\n\t\"value\"\r\n}",
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "key", line: 2},
				{type_: COLON, lit: ":", line: 2},
				{type_: STRING, lit: "value", line: 3},
				{type_: RBRACE, lit: "}", line: 4},
				{type_: EOF, lit: "", line: 4},
			},
		},
		{
			name:  "just whitespace",
			input: "   \t\n  ",
			want: []wantToken{
				{type_: EOF, lit: "", line: 2},
			},
		},
		{
			name:  "number zero",
			input: `0`,
			want: []wantToken{
				{type_: NUMBER, lit: "0", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "negative number in array",
			input: `[-1, -2, -3]`,
			want: []wantToken{
				{type_: LBRACKET, lit: "[", line: 1},
				{type_: NUMBER, lit: "-1", line: 1},
				{type_: COMMA, lit: ",", line: 1},
				{type_: NUMBER, lit: "-2", line: 1},
				{type_: COMMA, lit: ",", line: 1},
				{type_: NUMBER, lit: "-3", line: 1},
				{type_: RBRACKET, lit: "]", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "string with all escapes",
			input: `{"x": "\"\\/\b\f\n\r\t"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "x", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "\"\\/\b\f\n\r\t", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "trailing comma (invalid)",
			input: `{"a": 1,}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "a", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: NUMBER, lit: "1", line: 1},
				{type_: COMMA, lit: ",", line: 1},
				{type_: ILLEGAL, lit: "", line: 1},
			},
		},
		{
			name:  "bare number",
			input: `42`,
			want: []wantToken{
				{type_: NUMBER, lit: "42", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "bare true",
			input: `true`,
			want: []wantToken{
				{type_: BOOL, lit: "true", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "bare false",
			input: `false`,
			want: []wantToken{
				{type_: BOOL, lit: "false", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "bare null",
			input: `null`,
			want: []wantToken{
				{type_: NULL, lit: "", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "bare string",
			input: `"hello"`,
			want: []wantToken{
				{type_: STRING, lit: "hello", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "consecutive braces and brackets",
			input: `{}[]{}[]`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: LBRACKET, lit: "[", line: 1},
				{type_: RBRACKET, lit: "]", line: 1},
				{type_: LBRACE, lit: "{", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: LBRACKET, lit: "[", line: 1},
				{type_: RBRACKET, lit: "]", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "array of mixed types",
			input: `[1, "two", true, null, {"three": 3}]`,
			want: []wantToken{
				{type_: LBRACKET, lit: "[", line: 1},
				{type_: NUMBER, lit: "1", line: 1},
				{type_: COMMA, lit: ",", line: 1},
				{type_: STRING, lit: "two", line: 1},
				{type_: COMMA, lit: ",", line: 1},
				{type_: BOOL, lit: "true", line: 1},
				{type_: COMMA, lit: ",", line: 1},
				{type_: NULL, lit: "", line: 1},
				{type_: COMMA, lit: ",", line: 1},
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "three", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: NUMBER, lit: "3", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: RBRACKET, lit: "]", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "unicode in key",
			input: `{"café": "latte"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "café", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "latte", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "emoji in string value",
			input: `{"status": "👍"}`,
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "status", line: 1},
				{type_: COLON, lit: ":", line: 1},
				{type_: STRING, lit: "👍", line: 1},
				{type_: RBRACE, lit: "}", line: 1},
				{type_: EOF, lit: "", line: 1},
			},
		},
		{
			name:  "multiple newlines between tokens",
			input: "{\n\n\n\"x\"\n\n:\n\n1\n\n}",
			want: []wantToken{
				{type_: LBRACE, lit: "{", line: 1},
				{type_: STRING, lit: "x", line: 4},
				{type_: COLON, lit: ":", line: 6},
				{type_: NUMBER, lit: "1", line: 8},
				{type_: RBRACE, lit: "}", line: 10},
				{type_: EOF, lit: "", line: 10},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertTokens(t, tt.input, tt.want)
		})
	}
}
