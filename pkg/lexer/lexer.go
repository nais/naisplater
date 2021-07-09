package lexer

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

const (
	TokenEnd = iota
	TokenCurlyLeft
	TokenCurlyRight
	TokenWhitespace
	TokenIdentifier
)

// runeClass returns the token class of an input character.
func runeClass(r rune) int {
	if unicode.IsSpace(r) {
		return TokenWhitespace
	}
	switch r {
	case eof:
		return TokenEnd
	case '{':
		return TokenCurlyLeft
	case '}':
		return TokenCurlyRight
	default:
		return TokenIdentifier
	}
}

// Scanner represents a lexical scanner.
type Scanner struct {
	r *bufio.Reader
}

// NewScanner returns a new instance of Scanner.
func NewScanner(r io.Reader) *Scanner {
	return &Scanner{r: bufio.NewReader(r)}
}

// Scan returns the next token and literal value.
func (s *Scanner) Scan() (class int, lit string) {
	ch := s.read()
	class = runeClass(ch)

	switch class {
	case TokenIdentifier:
		s.unread()
		lit = s.scanIdentifier()
	case TokenCurlyLeft, TokenCurlyRight, TokenWhitespace:
		s.unread()
		lit = s.scanAllOf(class)
	case TokenEnd:
		lit = ``
	default:
		lit = string(ch)
	}

	return
}

// ScanIgnoreWhitespace scans the next non-whitespace token.
func (s *Scanner) ScanIgnoreWhitespace() (tok int, lit string) {
	tok, lit = s.Scan()
	if tok == TokenWhitespace {
		tok, lit = s.Scan()
	}
	return
}

// scanWhitespace consumes the current rune and all contiguous whitespace.
func (s *Scanner) scanAllOf(tok int) string {
	var buf bytes.Buffer

	for {
		ch := s.read()
		class := runeClass(ch)
		if class == tok {
			buf.WriteRune(ch)
		} else if class == TokenEnd {
			return buf.String()
		} else {
			s.unread()
			return buf.String()
		}
	}
}

func (s *Scanner) scanIdentifier() string {
	var buf bytes.Buffer
	buf.WriteRune(s.read())
	escape := false

OUTER:
	for {
		ch := s.read()
		class := runeClass(ch)

		switch {
		case class == TokenEnd:
			break OUTER
		case escape || class == TokenIdentifier:
			buf.WriteRune(ch)
		default:
			s.unread()
			break OUTER
		}

		escape = false
	}

	return buf.String()
}

// read reads the next rune from the buffered reader.
// Returns the rune(0) if an error occurs (or io.EOF is returned).
func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

// unread places the previously read rune back on the reader.
func (s *Scanner) unread() { _ = s.r.UnreadRune() }

// eof represents a marker rune for the end of the reader.
var eof = rune(0)
