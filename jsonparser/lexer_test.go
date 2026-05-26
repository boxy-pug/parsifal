package jsonparser

import (
	"strings"
	"testing"
)

func TestJsonLexer_SingleCharTokens(t *testing.T) {
	input := strings.NewReader(`{[{}]}[]{:,:}`)

	l := NewLexer(input)

	tests := []struct {
		expectedType    tokenType
		expectedLiteral string
	}{
		{LBRACE, "{"},
		{LBRACKET, "["},
		{LBRACE, "{"},
		{RBRACE, "}"},
		{RBRACKET, "]"},
		{RBRACE, "}"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{LBRACE, "{"},
		{COLON, ":"},
		{COMMA, ","},
		{COLON, ":"},
		{RBRACE, "}"},
		{EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestJsonLexer(t *testing.T) {

	input := strings.NewReader(`{"name": "Alice", "age": 30, "active": true, "scores": [95, 82, 77], "address": {"city": "Portland", "zip": "97201"}, "nickname": "Ali\"ce", "balance": 99.95, "deleted": null}`)

	l := NewLexer(input)

	tests := []struct {
		expectedType    tokenType
		expectedLiteral string
	}{
		{LBRACE, "{"},
		{STRING, "name"},
		{COLON, ":"},
		{STRING, "Alice"},
		{COMMA, ","},
		{STRING, "age"},
		{COLON, ":"},
		{NUMBER, "30"},
		{COMMA, ","},
		{STRING, "active"},
		{COLON, ":"},
		{BOOL, "true"},
		{COMMA, ","},
		{STRING, "scores"},
		{COLON, ":"},
		{LBRACKET, "["},
		{NUMBER, "95"},
		{COMMA, ","},
		{NUMBER, "82"},
		{COMMA, ","},
		{NUMBER, "77"},
		{RBRACKET, "]"},
		{COMMA, ","},
		{STRING, "address"},
		{COLON, ":"},
		{LBRACE, "{"},
		{STRING, "city"},
		{COLON, ":"},
		{STRING, "Portland"},
		{COMMA, ","},
		{STRING, "zip"},
		{COLON, ":"},
		{STRING, "97201"},
		{RBRACE, "}"},
		{COMMA, ","},
		{STRING, "nickname"},
		{COLON, ":"},
		{STRING, "Ali\"ce"},
		{COMMA, ","},
		{STRING, "balance"},
		{COLON, ":"},
		{NUMBER, "99.95"},
		{COMMA, ","},
		{STRING, "deleted"},
		{COLON, ":"},
		{NULL, ""},
		{EOF, ""},
	}

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}
		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
