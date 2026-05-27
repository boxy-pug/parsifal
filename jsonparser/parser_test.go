package jsonparser

import (
	"strings"
	"testing"
)

func assertParse(t *testing.T, input string) Node {
	t.Helper()
	l := NewLexer(strings.NewReader(input))
	p := NewParser(l)
	node, err := p.Parse()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	return node
}

func assertParseError(t *testing.T, input string) {
	t.Helper()
	l := NewLexer(strings.NewReader(input))
	p := NewParser(l)
	_, err := p.Parse()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func assertNodeType[T any](t *testing.T, node Node) T {
	t.Helper()
	v, ok := node.(T)
	if !ok {
		t.Fatalf("expected %T, got %T", v, node)
	}
	return v
}

func TestParse_SimpleObject(t *testing.T) {
	node := assertParse(t, `{"a": 1}`)
	obj := assertNodeType[Object](t, node)

	if len(obj.Pairs) != 1 {
		t.Fatalf("expected 1 pair, got %d", len(obj.Pairs))
	}

	num := assertNodeType[Number](t, obj.Pairs["a"])
	if num.Value != 1 {
		t.Fatalf("expected 1, got %f", num.Value)
	}
}

func TestParse_Arrays(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  int
	}{
		{"empty", `[]`, 0},
		{"single number", `[1]`, 1},
		{"multiple numbers", `[1, 2, 3]`, 3},
		{"mixed types", `[1, "hello", true, null]`, 4},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := assertParse(t, tt.input)
			arr := assertNodeType[Array](t, node)

			if len(arr.Elements) != tt.want {
				t.Fatalf("expected %d elements, got %d", tt.want, len(arr.Elements))
			}
		})
	}
}

func TestParse_Numbers(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  float64
	}{
		{"integer", `42`, 42},
		{"negative", `-17`, -17},
		{"zero", `0`, 0},
		{"float", `3.14`, 3.14},
		{"negative float", `-0.5`, -0.5},
		{"scientific positive", `1.5e10`, 1.5e10},
		{"scientific negative exp", `2.99e-8`, 2.99e-8},
		{"scientific upper E", `6.022E23`, 6.022e23},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := assertParse(t, tt.input)
			num := assertNodeType[Number](t, node)
			if num.Value != tt.want {
				t.Fatalf("expected %f, got %f", tt.want, num.Value)
			}
		})
	}
}

func TestParse_Strings(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"simple", `"hello"`, "hello"},
		{"empty", `""`, ""},
		{"escaped quote", `"say \"hi\""`, `say "hi"`},
		{"escaped backslash", `"a\\b"`, `a\b`},
		{"escaped newline", `"line1\nline2"`, "line1\nline2"},
		{"escaped tab", `"col\there"`, "col\there"},
		{"unicode", `"\u00E9"`, "é"},
		{"unicode surrogate pair", `"\uD83D\uDE00"`, "😀"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := assertParse(t, tt.input)
			str := assertNodeType[String](t, node)
			if str.Value != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, str.Value)
			}
		})
	}
}

func TestParse_BoolAndNull(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  Node
	}{
		{"true", `true`, Bool{Value: true}},
		{"false", `false`, Bool{Value: false}},
		{"null", `null`, Null{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := assertParse(t, tt.input)
			if node != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, node)
			}
		})
	}
}

func TestParse_NestedStructures(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, node Node)
	}{
		{
			name:  "nested object",
			input: `{"a": {"b": 1}}`,
			check: func(t *testing.T, node Node) {
				obj := assertNodeType[Object](t, node)
				inner := assertNodeType[Object](t, obj.Pairs["a"])
				num := assertNodeType[Number](t, inner.Pairs["b"])
				if num.Value != 1 {
					t.Fatalf("expected 1, got %f", num.Value)
				}
			},
		},
		{
			name:  "nested array",
			input: `[[1, 2], [3, 4]]`,
			check: func(t *testing.T, node Node) {
				arr := assertNodeType[Array](t, node)
				if len(arr.Elements) != 2 {
					t.Fatalf("expected 2 elements, got %d", len(arr.Elements))
				}
				inner := assertNodeType[Array](t, arr.Elements[0])
				if len(inner.Elements) != 2 {
					t.Fatalf("expected 2 inner elements, got %d", len(inner.Elements))
				}
			},
		},
		{
			name:  "array of objects",
			input: `[{"a": 1}, {"b": 2}]`,
			check: func(t *testing.T, node Node) {
				arr := assertNodeType[Array](t, node)
				if len(arr.Elements) != 2 {
					t.Fatalf("expected 2 elements, got %d", len(arr.Elements))
				}
				obj := assertNodeType[Object](t, arr.Elements[0])
				num := assertNodeType[Number](t, obj.Pairs["a"])
				if num.Value != 1 {
					t.Fatalf("expected 1, got %f", num.Value)
				}
			},
		},
		{
			name:  "object with array value",
			input: `{"nums": [1, 2, 3]}`,
			check: func(t *testing.T, node Node) {
				obj := assertNodeType[Object](t, node)
				arr := assertNodeType[Array](t, obj.Pairs["nums"])
				if len(arr.Elements) != 3 {
					t.Fatalf("expected 3 elements, got %d", len(arr.Elements))
				}
			},
		},
		{
			name:  "deeply nested",
			input: `{"a": {"b": {"c": [1, {"d": true}]}}}`,
			check: func(t *testing.T, node Node) {
				obj := assertNodeType[Object](t, node)
				obj2 := assertNodeType[Object](t, obj.Pairs["a"])
				obj3 := assertNodeType[Object](t, obj2.Pairs["b"])
				arr := assertNodeType[Array](t, obj3.Pairs["c"])
				num := assertNodeType[Number](t, arr.Elements[0])
				if num.Value != 1 {
					t.Fatalf("expected 1, got %f", num.Value)
				}
				obj4 := assertNodeType[Object](t, arr.Elements[1])
				b := assertNodeType[Bool](t, obj4.Pairs["d"])
				if b.Value != true {
					t.Fatalf("expected true, got %v", b.Value)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := assertParse(t, tt.input)
			tt.check(t, node)
		})
	}
}

func TestParse_EmptyStructures(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"empty object", `{}`},
		{"empty array", `[]`},
		{"empty string value", `{"k": ""}`},
		{"object with empty array", `{"a": []}`},
		{"array with empty object", `[{}]`},
		{"empty object in empty array", `[{}]`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertParse(t, tt.input)
		})
	}
}

func TestParse_MultipleObjectKeys(t *testing.T) {
	node := assertParse(t, `{"a": 1, "b": 2, "c": 3}`)
	obj := assertNodeType[Object](t, node)

	if len(obj.Pairs) != 3 {
		t.Fatalf("expected 3 pairs, got %d", len(obj.Pairs))
	}

	for _, key := range []string{"a", "b", "c"} {
		num := assertNodeType[Number](t, obj.Pairs[key])
		if num.Value == 0 {
			t.Fatalf("expected non-zero value for key %q", key)
		}
	}
}

func TestParse_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"trailing comma in object", `{"a": 1,}`},
		{"trailing comma in array", `[1, 2,]`},
		{"missing colon", `{"a" 1}`},
		{"missing key", `{1}`},
		{"missing value", `{"a":}`},
		{"trailing garbage", `1 2`},
		{"bare comma", `,,`},
		{"unclosed object", `{"a": 1`},
		{"unclosed array", `[1, 2`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertParseError(t, tt.input)
		})
	}
}
