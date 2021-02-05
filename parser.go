package sham

import (
	"errors"
	"fmt"
	"strconv"
)

type Parser struct {
	source string
	tokens []Token
	i      int
}

var errEOF = errors.New("EOF")

func NewParser(d string) *Parser {
	return &Parser{
		source: d,
		tokens: make([]Token, 0),
		i:      0,
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

func (p *Parser) Parse() (Schema, error) {
	tokens, err := Tokenize(p.source)
	if err != nil {
		return Schema{}, err
	}
	p.tokens = tokens

	root := Schema{}

	root.Root, err = p.parseValue()
	if err != nil {
		return Schema{}, err
	}

	return root, nil
}

func (p *Parser) parseValue() (Node, error) {
	var n Node
	var err error

	t := p.current()

	switch t.Type {
	case TokLBrace:
		n, err = p.parseObject()
		if err != nil {
			return Schema{}, err
		}
	case TokLBracket:
		n, err = p.parseArray()
		if err != nil {
			return Schema{}, err
		}
	case TokLParen:
		n, err = p.parseRange()
	case TokIdent:
		g, ok := terminalGenerators[t.Value]
		if !ok {
			return Schema{}, errors.New("unknown terminal generator")
		}
		n = g
	case TokInteger:
		i, _ := strconv.Atoi(t.Value)
		n = Literal{Value: i}
	case TokString:
		n = Literal{Value: t.Value}
	case TokEOF:
		return nil, errors.New("empty input")
	}

	return n, nil
}

func (p *Parser) parseObject() (Object, error) {
	obj := Object{Values: make(map[string]Node)}
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

		obj.Values[k] = n

		t = p.advance()
		if t.Type != TokRBrace && t.Type != TokComma {
			return Object{}, fmt.Errorf(`expected "," or "}", got %v`, t)
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
			return Array{}, nil
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

func (p *Parser) parseRange() (Range, error) {
	r := Range{}
	t := p.current()

	if t = p.advance(); t.Type != TokInteger {
		return Range{}, fmt.Errorf(`expected integer for range min, got %v`, t)
	}

	i, _ := strconv.Atoi(t.Value)
	r.Min = i

	if t = p.advance(); t.Type != TokComma {
		return Range{}, fmt.Errorf(`expected ",", got %v`, t)
	}

	if t = p.advance(); t.Type != TokInteger {
		return Range{}, fmt.Errorf(`expected integer for range max, got %v`, t)
	}

	i, _ = strconv.Atoi(t.Value)
	r.Max = i

	if t = p.advance(); t.Type != TokRParen {
		return Range{}, fmt.Errorf(`expected ")", got %v`, t)
	}

	return r, nil
}
