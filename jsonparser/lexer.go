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
	lastRune     rune
}

// NewLexer wraps input in bufio.Reader for rune level reading
// Calls readChar once to "load the lexer" with first rune
func NewLexer(input io.Reader) *Lexer {
	l := &Lexer{
		reader: bufio.NewReader(input),
	}
	l.readRune()
	l.line = 1
	return l
}

// NextToken is the main entrypoint, inspects current rune and dispatches
// to correct tokentype
func (l *Lexer) NextToken() token {
	var tok token

	// skip/ignore whitespace. calls ReadRune helper which increments line count if ch == '\n'
	l.skipWhitespace()

	switch l.ch {
	case '{':
		tok = newToken(LBRACE, string(l.ch), l.line)
	case '}':
		if l.lastRune == ',' {
			tok = newToken(ILLEGAL, "", l.line)
			break
		}
		tok = newToken(RBRACE, string(l.ch), l.line)
	case '[':
		tok = newToken(LBRACKET, string(l.ch), l.line)
	case ']':
		if l.lastRune == ',' {
			tok = newToken(ILLEGAL, "", l.line)
			break
		}
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
	l.lastRune = l.ch
	l.readRune()
	return tok
}

// readRune reads and advances to next rune
// Store rune in ch, add width to readPosition and set position to
// the old readposition on EOF or invalid set ch to 0
func (l *Lexer) readRune() {
	newRune, byteWidth, err := l.reader.ReadRune()
	if err != nil {
		l.ch = 0
		return
	}
	l.ch = newRune
	if l.ch == '\n' {
		l.line++
	}
	l.position = l.readPosition
	l.readPosition += byteWidth
}

// readString reads a string until '"' with no preceding \ escape
func (l *Lexer) readString() (string, error) {
	var res []rune
	l.readRune() // advance to first char in string

	for l.ch != '"' {
		if l.ch == '\n' {
			return "", fmt.Errorf("real newline in string")
		}
		// handle escape codes
		if l.ch == '\\' {
			l.readRune() // read one more to see what is escaped
			switch l.ch {
			case '"': // escaped double quote
				res = append(res, l.ch)
			case '\\': // escaped backslash
				res = append(res, l.ch)
			case '/': // escaped forward slash
				res = append(res, l.ch)
			case 'n': // escpaed newline
				res = append(res, '\n')
			case 't': // tab
				res = append(res, '\t')
			case 'r': // carriage return
				res = append(res, '\r')
			case 'b': // backspace
				res = append(res, '\b')
			case 'f':
				res = append(res, '\f')
			case 'u': // unicode utf16 code point
				code, err := l.readFourHexDigits()
				if err != nil {
					return "", fmt.Errorf("invalid unicode string")
				}
				// check for high surrogate, in that case check for '\u' read low surrogate and combine them
				if isHighSurrogate(code) {
					code, err = l.combineHighAndLowSurrogate(code)
					if err != nil {
						return "", fmt.Errorf("invalid low surrogate unicode string")
					}
				}
				res = append(res, code)
			default:
				return "", fmt.Errorf("unknown escape sequence")
			}
		} else {
			res = append(res, l.ch)
		}
		l.readRune()
	}
	return string(res), nil
}

func (l *Lexer) readNumber() string {
	var res []rune

	// if first ch is - that's ok, a negative number
	if l.ch == '-' {
		res = append(res, l.ch)
		l.readRune()
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
			l.readRune()
		case unicode.ToLower(l.ch) == 'e':
			res = append(res, l.ch)
			l.readRune()
			// '-' is allowed if it follows and e exponent
			if l.ch == '-' {
				res = append(res, l.ch)
				l.readRune()
			}
			// These are the only ways to end a number, i guess? not sure...
			// maybe move the peek eof to it's own case?
		case l.ch == '}' || l.ch == ']' || l.ch == ',' || l.ch == ' ' || l.ch == '\n' || l.peekRune() == 0:
			// number is not allowed to end with a '.'
			if res[len(res)-1] == '.' {
				return ""
			}
			// when exiting loop its bcs l.ch is something like
			// }, ] or , or eof – so we need to unread this rune to handle it properly later
			l.unreadRune()
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
		l.readRune()
	}
	// exiting loop when l.ch is not lowercase char
	// we need to unread one rune to handle it properly
	l.unreadRune()
	// TODO will this be a bug somehow?
	return string(res)
}

// readFourHexDigits() reads 4 unicode hex digits
// TODO need to check for high surrogate to support more emojis and chars
func (l *Lexer) readFourHexDigits() (rune, error) {
	var res []rune

	for range 4 {
		l.readRune() // advance from 'u' to first of four hex digits
		ch := unicode.ToLower(l.ch)
		if isDigit(ch) || (ch >= 'a' && ch <= 'f') {
			res = append(res, l.ch)
		} else {
			return 0, fmt.Errorf("invalid unicode string")
		}
	}

	code, err := strconv.ParseUint(string(res), 16, 32)
	if err != nil {
		return 0, fmt.Errorf("invalid unicode string")
	}
	return rune(code), nil
}

// Helper funcs:

func (l *Lexer) unreadRune() {
	if l.ch == '\n' {
		l.line--
	}
	l.reader.UnreadRune()

}

func isHighSurrogate(r rune) bool {
	return r >= 0xD800 && r <= 0xDBFF
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readRune()
	}
}

func isDigit(ch rune) bool {
	if ch >= '0' && ch <= '9' {
		return true
	}
	return false
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

func (l *Lexer) combineHighAndLowSurrogate(code rune) (rune, error) {

	l.readRune() // advance to possible '\'
	if l.ch == '\\' && l.peekRune() == 'u' {
		l.readRune() // advance to 'u'
		low, err := l.readFourHexDigits()
		if err != nil {
			return 0, fmt.Errorf("invalid unicode string")
		}
		// combine into a single rune
		code = 0x10000 + ((code - 0xD800) << 10) + (low - 0xDC00)
	}
	return code, nil
}
