package jsonparser

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"unicode"
)

type Lexer struct {
	reader       *bufio.Reader
	position     int  // byte offset of ch
	readPosition int  // points to the byte offset for next rune
	ch           rune // current rune
	line         int
}

// NewLexer wraps input in bufio.Reader for rune level reading
// Calls readChar once to "load the lexer" with first rune
func NewLexer(input io.Reader) *Lexer {
	l := &Lexer{
		reader: bufio.NewReader(input),
	}
	l.readChar()
	l.line = 1
	return l
}

// NextToken is the main entrypoint, inspects current rune and dispatches
// to correct tokentype
func (l *Lexer) NextToken() token {
	var tok token

	// skip/ignore whitespace, but increment linecount on '/n'
	l.skipWhitespace()

	switch l.ch {
	case '{':
		tok = newToken(LBRACE, string(l.ch), l.line)
	case '}':
		tok = newToken(RBRACE, string(l.ch), l.line)
	case '[':
		tok = newToken(LBRACKET, string(l.ch), l.line)
	case ']':
		tok = newToken(RBRACKET, string(l.ch), l.line)
	case ',':
		tok = newToken(COMMA, string(l.ch), l.line)
	case ':':
		tok = newToken(COLON, string(l.ch), l.line)
	case 0:
		tok = newToken(EOF, "", l.line)
	case '"': // handle strings
		literal, err := l.readString()
		if err != nil {
			tok = newToken(ILLEGAL, "", l.line)
			break
		}
		tok = newToken(STRING, literal, l.line)
	case 't', 'f', 'n':
		ident := l.readIdentifier()
		switch ident {
		case "true", "false":
			tok = newToken(BOOL, ident, l.line)
		case "null":
			tok = newToken(NULL, "", l.line)
		default:
			tok = newToken(ILLEGAL, "", l.line)
		}
	default:
		// Handle numbers
		if isDigit(l.ch) || l.ch == '-' {
			literal := l.readNumber()
			if literal == "" {
				tok = newToken(ILLEGAL, "", l.line)
				break
			}
			tok = newToken(NUMBER, literal, l.line)
		} else {
			tok = newToken(ILLEGAL, "", l.line)
		}
	}
	// TODO
	// fmt.Println(tok)
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
func newToken(tokenType tokenType, literal string, line int) token {
	return token{Type: tokenType, Literal: literal, Line: line}
}

// readString reads a string until '"' with no preceding \ escape
func (l *Lexer) readString() (string, error) {
	var res []rune
	l.readChar() // advance to first char in string

	for l.ch != '"' {
		if l.ch == '\n' {
			return "", fmt.Errorf("real newline in string")
		}
		if l.ch == '\\' {
			l.readChar() // read one more to see what is escaped
			switch l.ch {
			case '"':
				res = append(res, l.ch)
			case '\\':
				res = append(res, l.ch)
			case '/':
				res = append(res, l.ch)
			case 'n':
				res = append(res, '\n')
			case 't':
				res = append(res, '\t')
			case 'r':
				res = append(res, '\r')
			case 'b':
				res = append(res, '\b')
			case 'f':
				res = append(res, '\f')
			case 'u':
				literal, err := l.readUnicodeString()
				if err != nil {
					return "", fmt.Errorf("invalid unicode string")
				}
				res = append(res, literal)
				// TODO handle unicode code point stuff
			default:
				return "", fmt.Errorf("unknown escape sequence")
			}
		} else {
			res = append(res, l.ch)
		}
		l.readChar()
	}

	return string(res), nil
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		if l.ch == '\n' {
			l.line++
		}
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

	// if first ch is - that's ok, a negative number
	if l.ch == '-' {
		res = append(res, l.ch)
		l.readChar()
	}

	// if 1st ch is 0 and it's followed by a number , thats illegal
	// or if first char is '.'
	if (l.ch == '0' && isDigit(l.peekRune())) || l.ch == '.' {
		return ""
	}

	for {
		switch {
		case isDigit(l.ch), l.ch == '.':
			res = append(res, l.ch)
			l.readChar()
		case unicode.ToLower(l.ch) == 'e':
			res = append(res, l.ch)
			l.readChar()
			// '-' is allowed if it follows and e exponent
			if l.ch == '-' {
				res = append(res, l.ch)
				l.readChar()
			}
		case unicode.IsLetter(l.ch), l.ch == '_':
			return ""
		// These are the only ways to end a number, i guess? not sure...
		case l.ch == '}' || l.ch == ']' || l.ch == ',' || l.peekRune() == 0:
			if res[len(res)-1] == '.' {
				return ""
			}
			// when exiting loop its bcs l.ch is something like
			// }, ] or , or eof – so we need to unread this rune to handle it properly later
			l.reader.UnreadRune()
			return string(res)
		default:
			return ""

		}
	}
}

// for reading true, false or null
func (l *Lexer) readIdentifier() string {
	var res []rune

	for unicode.IsLower(l.ch) {
		res = append(res, l.ch)
		l.readChar()
	}
	// exiting loop when l.ch is not lowercase char
	// we need to unread one rune to handle it properly
	l.reader.UnreadRune()
	// TODO will this be a bug somehow?
	return string(res)
}

// reads a rune and immediately unreads it. returns the read rune
func (l *Lexer) peekRune() rune {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		return 0
	}
	l.reader.UnreadRune()
	return r
}

func (l *Lexer) readUnicodeString() (rune, error) {
	var res []rune
	var r rune

	for range 4 {
		l.readChar()
		ch := unicode.ToLower(l.ch)
		if isDigit(ch) || (ch >= 'a' && ch <= 'f') {
			res = append(res, ch)
		} else {
			return r, fmt.Errorf("invalid unicode string")
		}
	}

	code, err := strconv.ParseUint(string(res), 16, 32)
	if err != nil {
		return r, fmt.Errorf("invalid unicode string")
	}
	return rune(code), nil
}
