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
		i:      -1,
	}
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

	switch t := p.advance(); t.Type {
	case TokLBrace:
		root.Root, err = p.parseObject(t)
		if err != nil {
			return Schema{}, err
		}
	case TokIdent:
		g, ok := terminalGenerators[t.Value]
		if !ok {
			return Schema{}, errors.New("unknown terminal generator")
		}
		root.Root = g
	case TokInteger:
		i, _ := strconv.Atoi(t.Value)
		root.Root = Literal{Value: i}
	case TokString:
		root.Root = Literal{Value: t.Value}
	case TokEOF:
		return Schema{}, errors.New("empty input")
	}

	return root, nil
}

func (p *Parser) parseObject(t Token) (Object, error) {
	obj := Object{Values: make(map[string]Node)}

	for t.Type != TokRBrace {
		t = p.advance()
		if t.Type != TokString {
			return Object{}, fmt.Errorf("expected string, got %v", t)
		}

		k, n, err := p.parsePair(t)
		if err != nil {
			return Object{}, err
		}

		obj.Values[k] = n

		t = p.advance()
		if t.Type != TokRBrace && t.Type != TokComma {
			return Object{}, fmt.Errorf("expected \",\", got %v", t)
		}
	}

	return obj, nil
}

func (p *Parser) parsePair(t Token) (string, Node, error) {
	key := t.Value

	if t := p.advance(); t.Type != TokColon {
		return "", nil, fmt.Errorf("expected \":\", got %v", t)
	}

	var n Node
	switch t := p.advance(); t.Type {
	case TokInteger:
		i, _ := strconv.Atoi(t.Value)
		n = Literal{Value: i}
	case TokString:
		n = Literal{Value: t.Value}
	case TokIdent:
		g, ok := terminalGenerators[t.Value]
		if !ok {
			return "", nil, errors.New("unknown terminal generator")
		}
		n = g
	default:
		return "", nil, fmt.Errorf("unknown node type, got %v", t)
	}

	return key, n, nil
}

func (p *Parser) parseArray(t Token) (Array, error) {
	arr := Array{}

	return arr, nil
}
