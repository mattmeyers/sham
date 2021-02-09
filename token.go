package sham

import "fmt"

type TokenType int

const (
	TokInvalid TokenType = iota
	TokEOF
	TokWS
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
	TokFString
	TokInteger
	TokFloat
	TokIdent

	TokNull
	TokTrue
	TokFalse
)

var tokenStrings = map[TokenType]string{
	TokInvalid:  "<INVALID>",
	TokEOF:      "<EOF>",
	TokLBrace:   "{",
	TokRBrace:   "}",
	TokLBracket: "[",
	TokRBracket: "]",
	TokLParen:   "(",
	TokRParen:   ")",
	TokColon:    ":",
	TokComma:    ",",
	TokString:   "<STRING>",
	TokFString:  "<F STRING>",
	TokInteger:  "<INTEGER",
	TokFloat:    "<FLOAT>",
	TokIdent:    "<IDENT>",
	TokNull:     "null",
	TokTrue:     "true",
	TokFalse:    "false",
}

func (t TokenType) String() string {
	s, ok := tokenStrings[t]
	if !ok {
		s = fmt.Sprintf("%%!(UNKNOWN=%d)", t)
	}
	return s
}

type QuoteType byte

const (
	QuoteSingle   QuoteType = '\''
	QuoteDouble   QuoteType = '"'
	QuoteBacktick QuoteType = '`'
)

type Token struct {
	Type  TokenType
	Value string
}

func newToken(t TokenType, v string) Token {
	return Token{Type: t, Value: v}
}

var keywordMap = map[string]TokenType{
	"null":  TokNull,
	"true":  TokTrue,
	"false": TokFalse,
}
