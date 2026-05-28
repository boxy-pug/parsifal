// Package ast contains shared ast node types for all the parsers in parsifal. Currently JSON, soon YAML and TOML too.
package ast

// Node is a value in a parsed tree. It's a "marker interface",
// it carries no methods.
type Node interface{}

// Object maps string keys to Node values. A Go map, order not preserved.
// Every node type has a Tag string for YAML explicit type annotation
// (like !!str ot !Custom). Empty when not Yaml or no tag present.
type Object struct {
	Tag   string
	Pairs map[string]Node
}

type Array struct {
	Tag      string
	Elements []Node
}

type String struct {
	Tag   string
	Value string
}

type Number struct {
	Tag   string
	Value float64
}

type Bool struct {
	Tag   string
	Value bool
}

type Null struct {
	Tag string
}

// TODO: Add YAML and TOML types
// YAML: timestamp, binary, set, ordered map
// TOML: Datetime, multiple types, localtime etc
// Ints with underscores and alternate bases
