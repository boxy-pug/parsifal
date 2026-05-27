package jsonparser

import (
	"fmt"
	"strconv"
)

// Parser transforms a stream of tokens from the lexer to an AST.
// Holds the lexer and curretn token being examined.
type Parser struct {
	l   *Lexer
	cur token
}

// NewParser creates parser from existing Lexer, calls next immediately
// to load first token.
func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	p.next()
	return p
}

// Parse parses a single JSON value into a Node, then verifies that theres nothing left,
// no trailing tokens. Returns root node or error. The public API entry point.
func (p *Parser) Parse() (Node, error) {
	node, err := p.parseValue()
	if err != nil {
		return nil, err
	}
	if p.cur.Type != EOF {
		return nil, fmt.Errorf("unexpected token after root value")
	}
	return node, nil
}

// next pulls the next token from lexer and stores it in p.cur
// overwrites previous token
func (p *Parser) next() {
	p.cur = p.l.NextToken()
}

// expect checks if current token matches expected type. if yes, advance to next token
// and return true. if not, return false and no advance
func (p *Parser) expect(t tokenType) bool {
	if p.cur.Type == t {
		p.next()
		return true
	}
	return false
}

// parseValue is the dispatcher, it looks at p.cur.Type and routes to the right parse func.
// returns constructed node and error. The core of the recursive decent
func (p *Parser) parseValue() (Node, error) {

	//dispatch, switch statement based on p.
	switch p.cur.Type {
	case STRING:
		return p.parseString()
	case NUMBER:
		return p.parseNumber()
	case LBRACE:
		return p.parseObject()
	case LBRACKET:
		return p.parseArray()
	case BOOL:
		return p.parseBool()
	case NULL:
		return p.parseNull()

		// handle the different token types here
	default:
		return nil, fmt.Errorf("%w: %s", ErrUnexpectedToken, p.cur.Type)
	}
}

func (p *Parser) parseString() (Node, error) {
	node := String{Value: p.cur.Literal}
	p.next()
	return node, nil
}

func (p *Parser) parseNumber() (Node, error) {
	number, err := strconv.ParseFloat(p.cur.Literal, 64)
	if err != nil {
		return nil, err
	}
	node := Number{Value: number}
	p.next()
	return node, nil
}

func (p *Parser) parseObject() (Node, error) {
	var obj Object
	obj.Pairs = make(map[string]Node)

	p.next() // move past {
	if p.cur.Type == RBRACE {
		// if empty {}, move past ending brace and return empty
		p.next()
		return obj, nil
	}

	// if not empty {}, expect a string and a colon after, if not return error
	for {
		if p.cur.Type != STRING {
			return nil, ErrMissingKey
		}

		key := p.cur.Literal // Save the key

		p.next() // move past string to colon, hopefully
		if p.cur.Type != COLON {
			return nil, ErrMissingColon
		}

		p.next() // move past colon for the actual value, call parseValue recursively

		value, err := p.parseValue()
		if err != nil {
			return nil, ErrInvalidObjectValue
		}

		// add key and value to obj map
		obj.Pairs[key] = value

		// we're already at token after value here, parseValue has advanced it
		// rule of thumb: if you eat it, advance it
		switch p.cur.Type {
		case RBRACE: // ending brace, advance past it and return object
			p.next()
			return obj, nil
		case COMMA: // advance past it and continue the loop
			p.next()
		default: // anything else is error
			return nil, ErrInvalidObjectValue
		}
	}
}

func (p *Parser) parseArray() (Node, error) {
	var arr Array

	p.next() // move bast bracket
	if p.cur.Type == RBRACKET {
		p.next()
		return arr, nil
	}

	for {

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		arr.Elements = append(arr.Elements, value)

		switch p.cur.Type {
		case RBRACKET:
			p.next()
			return arr, nil
		case COMMA:
			p.next()
		default:
			return nil, ErrInvalidArray

		}
	}
}

func (p *Parser) parseBool() (Node, error) {
	var b Bool

	switch p.cur.Literal {
	case "true":
		b.Value = true
	case "false":
		b.Value = false
	default:
		return nil, ErrInvalidBool
	}

	p.next()
	return b, nil
}

func (p *Parser) parseNull() (Node, error) {
	p.next()
	return Null{}, nil
}
