package sham

import (
	"bufio"
	"bytes"
	"fmt"
)

const eof = rune(0)

func isWhitespace(c rune) bool    { return c == ' ' || c == '\t' || c == '\r' || c == '\n' }
func isAlpha(c rune) bool         { return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') }
func isDigit(c rune) bool         { return '0' <= c && c <= '9' }
func isPositiveDigit(c rune) bool { return '1' <= c && c <= '9' }
func isAlphaNumeric(c rune) bool  { return isAlpha(c) || isDigit(c) }

type Scanner struct {
	r   *bufio.Reader
	buf *bytes.Buffer
}

func NewScanner(b []byte) *Scanner {
	return &Scanner{
		r:   bufio.NewReader(bytes.NewBuffer(b)),
		buf: bytes.NewBuffer(nil),
	}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return eof
	}
	return ch
}

func (s *Scanner) unread() {
	if err := s.r.UnreadRune(); err != nil {
		// This should never happen considering the bufio.Reader is not exported and
		// we aren't directly exporting any relevant functions.
		panic("sham: " + err.Error())
	}
}

func Tokenize(source []byte) ([]Token, error) {
	s := NewScanner(source)
	tokens := make([]Token, 0)

	for {
		t, lit := s.Scan()
		if t == TokWS {
			continue
		} else if t == TokEOF {
			return tokens, nil
		} else if t == TokInvalid {
			return nil, fmt.Errorf("unknown token: %q", lit)
		}

		tokens = append(tokens, newToken(t, lit))
	}
}

func (s *Scanner) Scan() (tok TokenType, lit string) {
	ch := s.read()

	if isWhitespace(ch) {
		s.unread()
		return s.scanWhitespace()
	} else if isAlpha(ch) {
		s.unread()
		return s.scanIdent()
	} else if isDigit(ch) || ch == '-' {
		s.unread()
		return s.scanNumber()
	} else if ch == '/' {
		return TokRegex, s.scanRegex()
	}

	switch ch {
	case eof:
		return TokEOF, ""
	case '{':
		return TokLBrace, string(ch)
	case '}':
		return TokRBrace, string(ch)
	case '[':
		return TokLBracket, string(ch)
	case ']':
		return TokRBracket, string(ch)
	case '(':
		return TokLParen, string(ch)
	case ')':
		return TokRParen, string(ch)
	case ':':
		return TokColon, string(ch)
	case ',':
		return TokComma, string(ch)
	case '"':
		return TokString, s.scanString(QuoteDouble)
	case '`':
		return TokFString, s.scanString(QuoteBacktick)
	}

	return TokInvalid, string(ch)
}

func (s *Scanner) scanIdent() (TokenType, string) {
	s.buf.Reset()
	s.buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isAlpha(ch) {
			s.unread()
			break
		} else {
			_, _ = s.buf.WriteRune(ch)
		}
	}

	if token, ok := keywordMap[s.buf.String()]; ok {
		return token, s.buf.String()
	}

	return TokIdent, s.buf.String()
}

func (s *Scanner) scanWhitespace() (tok TokenType, lit string) {
	s.buf.Reset()
	s.buf.WriteRune(s.read())

	for {
		if ch := s.read(); ch == eof {
			break
		} else if !isWhitespace(ch) {
			s.unread()
			break
		} else {
			s.buf.WriteRune(ch)
		}
	}

	return TokWS, s.buf.String()
}

func (s *Scanner) scanString(qt QuoteType) string {
	s.buf.Reset()

	for {
		if ch := s.read(); ch == eof {
			break
		} else if ch == rune(qt) {
			break
		} else {
			_, _ = s.buf.WriteRune(ch)
		}
	}

	return s.buf.String()
}

func (s *Scanner) scanNumber() (TokenType, string) {
	tokType := TokInteger

	s.buf.Reset()
	ch := s.read()
	if ch == '-' {
		s.buf.WriteRune(ch)
		ch = s.read()
	}

	if ch == '0' {
		ch = s.read()
		if isDigit(ch) {
			return TokInvalid, string(ch)
		}
	} else {
		if !isPositiveDigit(ch) {
			return TokInvalid, string(ch)
		}

		s.buf.WriteRune(ch)
		ch = s.read()

		for isDigit(ch) {
			s.buf.WriteRune(ch)
			ch = s.read()
		}
	}

	if ch == '.' {
		tokType = TokFloat

		s.buf.WriteRune(ch)
		ch = s.read()

		for isDigit(ch) {
			s.buf.WriteRune(ch)
			ch = s.read()
		}
	}

	if ch == 'e' || ch == 'E' {
		tokType = TokFloat

		s.buf.WriteRune(ch)
		ch = s.read()

		if ch == '-' || ch == '+' {
			s.buf.WriteRune(ch)
			ch = s.read()
		}

		if !isDigit(ch) {
			return TokInvalid, string(ch)
		}

		s.buf.WriteRune(ch)
		ch = s.read()

		for isDigit(ch) {
			s.buf.WriteRune(ch)
			ch = s.read()
		}
	}

	if ch != eof {
		s.unread()
	}

	return tokType, s.buf.String()
}

func (s *Scanner) scanRegex() string {
	s.buf.Reset()

	ch := '/'
	for {
		prev := ch
		ch = s.read()
		if ch == eof {
			break
		} else if ch == '/' && prev != '\\' {
			break
		} else {
			_, _ = s.buf.WriteRune(ch)
		}
	}

	return s.buf.String()
}
