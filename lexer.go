package sham

import (
	"errors"
	"fmt"
)

type TokenType int

const (
	TokInvalid TokenType = iota
	TokEOF
	// Structural tokens
	TokLBrace
	TokRBrace
	TokLBracket
	TokRBracket
	TokLParen
	TokRParen
	TokColon
	TokComma

	TokString
	TokInteger
	TokIdent
)

type Token struct {
	Type  TokenType
	Value string
}

func newToken(t TokenType, v string) Token {
	fmt.Printf("Adding token: %v: %s\n", t, v)
	return Token{Type: t, Value: v}
}

func isWhitespace(c byte) bool    { return c == ' ' || c == '\t' || c == '\r' || c == '\n' }
func isAlpha(c byte) bool         { return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') }
func isDigit(c byte) bool         { return '0' <= c && c <= '9' }
func isPositiveDigit(c byte) bool { return '1' <= c && c <= '9' }
func isAlphaNumeric(c byte) bool  { return isAlpha(c) || isDigit(c) }

func Tokenize(source string) ([]Token, error) {
	tokens := make([]Token, 0)

	i := 0
	for i < len(source) {
		c := source[i]

		if isWhitespace(c) {
			i++
			continue
		}

		var t Token
		if isAlpha(c) {
			v, err := scanIdent(source, i)
			if err != nil {
				return nil, err
			}
			t = newToken(TokIdent, v)
			i += len(v)
		} else if isDigit(c) {
			v, err := scanInteger(source, i)
			if err != nil {
				return nil, err
			}
			t = newToken(TokInteger, v)
			i += len(v) - 1
		} else {
			switch c {
			case '{':
				t = newToken(TokLBrace, "{")
			case '}':
				t = newToken(TokRBrace, "}")
			case '[':
				t = newToken(TokLBracket, "[")
			case ']':
				t = newToken(TokRBracket, "]")
			case '(':
				t = newToken(TokLParen, "(")
			case ')':
				t = newToken(TokRParen, ")")
			case ':':
				t = newToken(TokColon, ":")
			case ',':
				t = newToken(TokComma, ",")
			case '"':
				v, err := scanString(source, i)
				if err != nil {
					return nil, err
				}
				t = newToken(TokString, v)
				i += len(v) + 1
			default:
				return nil, fmt.Errorf("unknown token provided: %q", c)
			}
		}

		tokens = append(tokens, t)
		i++
	}

	tokens = append(tokens, newToken(TokEOF, ""))
	return tokens, nil
}

func scanString(source string, i int) (string, error) {
	start := i
	i++
	for i < len(source) {
		if source[i] == '"' {
			return source[start+1 : i], nil
		}
		i++
	}

	return "", errors.New("unterminated string")
}

func scanIdent(source string, i int) (string, error) {
	start := i
	i++
	for i < len(source) {
		if !isAlphaNumeric(source[i]) {
			return source[start:i], nil
		}
		i++
	}

	return "", errors.New("invalid identifier")
}

func scanInteger(source string, i int) (string, error) {
	start := i
	i++
	for i < len(source) {
		if !isDigit(source[i]) {
			return source[start:i], nil
		}
		i++
	}

	return "", errors.New("invalid integer")
}
