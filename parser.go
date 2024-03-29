package sham

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// Parser maintains the internal state of the language parser. This struct takes
// a slice of tokens as input and produces an AST if and only if the token stream
// represents a valid Sham schema. Part of the internal state is the terminalGenerator
// map. This map of terminal generators must be set prior to the initiating the
// parsing method. If a terminal generator is referenced in the schema, but not
// defined in the terminal generator map, then the parsing process will be halted
// with an error.
//
// To ensure the parser begins with the proper state, one of the constructor functions
// should be used.
type Parser struct {
	TerminalGenerators map[string]Generator
	source             []byte
	tokens             []Token
	i                  int
}

var errEOF = errors.New("EOF")

// NewParser creates a new Parser instance with an empty terminal generator map.
func NewParser(d []byte) *Parser {
	return &Parser{
		TerminalGenerators: make(map[string]Generator),
		source:             d,
		tokens:             make([]Token, 0),
		i:                  0,
	}
}

// NewDefaultParser creates a new Parser instance using the default terminal
// generators map.
func NewDefaultParser(d []byte) *Parser {
	return &Parser{
		TerminalGenerators: TerminalGenerators,
		source:             d,
		tokens:             make([]Token, 0),
		i:                  0,
	}
}

// RegisterGenerators merges a terminal generator map into the parser's internal
// terminal generator map. If a generator is already registered, then the existing
// generator will be overwritten with the new. To avoid parsing errors, all terminal
// generators should be registered prior to parsing.
func (p *Parser) RegisterGenerators(gs map[string]Generator) {
	for k, v := range gs {
		p.TerminalGenerators[k] = v
	}
}

func (p *Parser) current() Token {
	return p.tokens[p.i]
}

func (p *Parser) peek() Token {
	return p.tokens[p.i+1]
}

func (p *Parser) advance() Token {
	t := p.peek()
	p.i++
	return t
}

// Parse generates a new AST from the schema provided to the parser. The parser
// combines multiple steps to generate this structure. The schema is first
// tokenized. If an invalid token is presented (either an unknown character or
// an unterminated sequence), then a scanning error will be returned. Upon success,
// the slice of tokens will be parsed. If the tokens are representative of a valid
// Sham schema, then an AST will be returned. Otherwise and error will be returned.
func (p *Parser) Parse() (Schema, error) {
	tokens, err := Tokenize(p.source)
	if err != nil {
		return Schema{}, err
	}
	p.tokens = tokens

	root, err := p.parseValue()
	if err != nil {
		return Schema{}, err
	}

	return Schema{Root: root}, nil
}

func (p *Parser) parseValue() (Node, error) {
	var n Node
	var err error

	t := p.current()

	switch t.Type {
	case TokLBrace:
		n, err = p.parseObject()
	case TokLBracket:
		n, err = p.parseArray()
	case TokLParen:
		n, err = p.parseRange()
	case TokIdent:
		n, err = p.parseIdent()
	case TokInteger:
		n, err = p.parseInteger()
	case TokFloat:
		n, err = p.parseFloat()
	case TokString:
		n = Literal{Value: t.Value}
	case TokFString:
		n, err = p.parseFString()
	case TokRegex:
		n, err = p.parseRegex()
	case TokNull:
		n = Literal{Value: nil}
	case TokTrue:
		n = Literal{Value: true}
	case TokFalse:
		n = Literal{Value: false}
	case TokEOF:
		err = errors.New("empty input")
	}

	if err != nil {
		return nil, err
	}

	return n, nil
}

func (p *Parser) parseObject() (Object, error) {
	obj := Object{}
	t := p.current()

	if p.peek().Type == TokRBrace {
		return obj, nil
	}

	for {
		t = p.advance()
		if t.Type != TokString {
			return Object{}, fmt.Errorf("expected string, got %v", t)
		}

		k, n, err := p.parsePair()
		if err != nil {
			return Object{}, err
		}

		obj.AppendPair(k, n)

		t = p.advance()
		if t.Type != TokRBrace && t.Type != TokComma {
			return Object{}, fmt.Errorf(`expected "," or "}", got %q`, t)
		} else if t.Type == TokRBrace {
			break
		}
	}

	return obj, nil
}

func (p *Parser) parsePair() (string, Node, error) {
	t := p.current()
	key := t.Value

	if t := p.advance(); t.Type != TokColon {
		return "", nil, fmt.Errorf("expected \":\", got %v", t)
	}

	p.advance()

	n, err := p.parseValue()
	if err != nil {
		return "", nil, err
	}

	return key, n, nil
}

func (p *Parser) parseArray() (Array, error) {
	var err error
	t := p.current()
	arr := Array{}

	if p.peek().Type == TokRBracket {
		p.advance()
		return arr, nil
	}

	if p.peek().Type == TokLParen {
		p.advance()
		r, err := p.parseRange()
		if err != nil {
			return Array{}, err
		}
		arr.Range = &r

		t = p.advance()
		if t.Type != TokComma {
			return Array{}, fmt.Errorf(`expected ",", got %v`, t)
		}
	}

	t = p.advance()
	arr.Inner, err = p.parseValue()
	if err != nil {
		return Array{}, err
	}

	t = p.advance()
	if t.Type != TokRBracket {
		return Array{}, fmt.Errorf(`expected "]", got %v`, t)
	}

	return arr, nil
}

func (p *Parser) parseRange() (r Range, err error) {
	t := p.current()

	if t = p.advance(); t.Type != TokInteger {
		return Range{}, fmt.Errorf(`expected integer for range min, got %v`, t)
	}

	i, _ := strconv.Atoi(t.Value)
	r.Min = i

	t = p.advance()

	if t.Type == TokRParen {
		r.Max = r.Min
		return r, nil
	}

	if t.Type != TokComma {
		return Range{}, fmt.Errorf(`expected ",", got %v`, t)
	}

	if t = p.advance(); t.Type != TokInteger {
		return Range{}, fmt.Errorf(`expected integer for range max, got %v`, t)
	}

	i, err = strconv.Atoi(t.Value)
	if err != nil {
		return Range{}, err
	} else if i < r.Min {
		return Range{}, errors.New("range maximum cannot be less than the minimum")
	}
	r.Max = i

	if t = p.advance(); t.Type != TokRParen {
		return Range{}, fmt.Errorf(`expected ")", got %v`, t)
	}

	return r, nil
}

var fStringRegex = regexp.MustCompile(`{([^{}]*)}`)

func (p *Parser) parseFString() (FormattedString, error) {
	t := p.current()

	matches := fStringRegex.FindAllString(t.Value, -1)
	if len(matches) == 0 {
		return FormattedString{Raw: t.Value}, nil
	}

	params := make([]Generator, len(matches))

	for i, m := range matches {
		g, ok := TerminalGenerators[m[1:len(m)-1]]
		if !ok {
			return FormattedString{}, fmt.Errorf("unknown terminal generator %s in formatted string", m)
		}
		params[i] = g
	}

	format := fStringRegex.ReplaceAllString(t.Value, "%v")

	return FormattedString{
		Raw:    t.Value,
		Format: format,
		Params: params,
	}, nil
}

func (p *Parser) parseRegex() (Regex, error) {
	t := p.current()

	r, err := NewRegex(t.Value)
	if err != nil {
		return Regex{}, err
	}
	return r, nil
}

func (p *Parser) parseInteger() (Literal, error) {
	t := p.current()
	i, err := strconv.Atoi(t.Value)
	if err != nil {
		return Literal{}, err
	}
	return Literal{Value: i}, nil
}

func (p *Parser) parseFloat() (Literal, error) {
	t := p.current()
	f, err := strconv.ParseFloat(t.Value, 10)
	if err != nil {
		return Literal{}, err
	}
	return Literal{Value: f}, nil
}

func (p *Parser) parseIdent() (TerminalGenerator, error) {
	n := p.current().Value
	fn, ok := p.TerminalGenerators[n]
	if !ok {
		return TerminalGenerator{}, fmt.Errorf("unknown terminal generator %q", n)
	}
	return TerminalGenerator{Name: n, fn: fn}, nil
}
