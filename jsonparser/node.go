package jsonparser

// interface{} basically means that any go type satisfies this automatically
type Node interface{}

// Json types that satisfy Node interface

type Object struct {
	Pairs map[string]Node
}

type Array struct {
	Elements []Node
}

type String struct {
	Value string
}

type Number struct {
	Value float64
}

type Bool struct {
	Value bool
}

type Null struct{}
