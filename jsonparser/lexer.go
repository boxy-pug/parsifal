package jsonparser

import (
	"bufio"
	"fmt"
	"io"
	"unicode"
)

type Lexer struct {
	reader       *bufio.Reader
	position     int  // byte offset of ch
	readPosition int  // points to the byte offset for next rune
	ch           rune // current rune
}

// NewLexer wraps input in bufio.Reader for rune level reading
// Calls readChar once to "load the lexer" with first rune
func NewLexer(input io.Reader) *Lexer {
	l := &Lexer{
		reader: bufio.NewReader(input),
	}
	l.readChar()
	return l
}

// NextToken is the main entrypoint, inspects current rune and dispatches
// to correct tokentype
func (l *Lexer) NextToken() token {
	var tok token

	l.skipWhitespace() // skip/ignore whitespace

	switch l.ch {
	case '{':
		tok = newToken(LBRACE, string(l.ch))
	case '}':
		tok = newToken(RBRACE, string(l.ch))
	case '[':
		tok = newToken(LBRACKET, string(l.ch))
	case ']':
		tok = newToken(RBRACKET, string(l.ch))
	case ',':
		tok = newToken(COMMA, string(l.ch))
	case ':':
		tok = newToken(COLON, string(l.ch))
	case 0:
		tok = newToken(EOF, "")
	case '"': // handle strings
		tok.Literal = l.readString()
		tok.Type = STRING
	case 't', 'f', 'n':
		ident := l.readIdentifier()
		switch ident {
		case "true", "false":
			tok = newToken(BOOL, ident)
		case "null":
			tok = newToken(NULL, "")
		default:
			tok = newToken(ILLEGAL, "")
		}
	default:
		// Handle numbers
		if isDigit(l.ch) || l.ch == '-' {
			tok.Type = NUMBER
			tok.Literal = l.readNumber()
		}
	}
	// TODO
	fmt.Println(tok)
	l.readChar()
	return tok
}

// readChar reads and advances to next rune
// Store rune in ch, add width to readPosition and set position to
// the old readposition on EOF or invalid set ch to 0
func (l *Lexer) readChar() {
	newRune, byteWidth, err := l.reader.ReadRune()
	if err != nil {
		l.ch = 0
		return
	}
	l.ch = newRune
	l.position = l.readPosition
	l.readPosition += byteWidth
}

// newToken is a little helper to make a token based on
// tokentype and literal
func newToken(tokenType tokenType, literal string) token {
	return token{Type: tokenType, Literal: literal}
}

// readString reads a string until '"' with no preceding \ escape
func (l *Lexer) readString() string {
	var res []rune
	l.readChar() // advance to first char in string

	for l.ch != '"' {
		if l.ch == '\\' {
			l.readChar() // read one more to see what is escaped
			switch l.ch {
			case '"':
				res = append(res, '"')
			}
		} else {
			res = append(res, l.ch)
		}
		l.readChar()
	}

	return string(res)
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

func isDigit(ch rune) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}
	return false
}

func (l *Lexer) readNumber() string {
	var res []rune

	for isDigit(l.ch) || l.ch == '.' {
		res = append(res, l.ch)
		next := l.peek()
		if next == ',' || next == ']' || next == '}' {
			break
		}
		l.readChar()
	}
	return string(res)
}

func (l *Lexer) readIdentifier() string {
	var res []rune

	for unicode.IsLower(l.ch) {
		res = append(res, l.ch)
		if l.peek() == ',' {
			break
		}
		l.readChar()
	}
	return string(res)
}

// reads a rune and immediately unreads it. returns the read rune
func (l *Lexer) peek() rune {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		return 0
	}
	l.reader.UnreadRune()
	return r
}
