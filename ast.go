package sham

import (
	"fmt"
	"math/rand"
)

type Node interface {
	Generator
}

type Schema struct {
	Root Node
}

func (s Schema) Generate() interface{} {
	if s.Root == nil {
		return nil
	}

	return s.Root.Generate()
}

type Object struct {
	Values map[string]Node
}

func (o Object) Generate() interface{} {
	out := make(map[string]interface{})
	for k, v := range o.Values {
		out[k] = v.Generate()
	}
	return out
}

type Array struct {
	Range *Range
	Inner Node
}

func (a Array) Generate() interface{} {
	if a.Inner == nil {
		return []interface{}{}
	}

	n := 1
	if a.Range != nil {
		n = a.Range.GetValue()
	}

	out := make([]interface{}, n)
	for i := 0; i < n; i++ {
		out[i] = a.Inner.Generate()
	}
	return out
}

type Range struct {
	Min int
	Max int
}

func (r Range) GetValue() int {
	if r.Min == r.Max {
		return r.Min
	}

	return rand.Intn((r.Max+1)-r.Min) + r.Min
}

func (r Range) Generate() interface{} { return r.GetValue() }

type FormattedString struct {
	Raw    string
	Format string
	Params []Generator
}

func (f FormattedString) Generate() interface{} {
	if len(f.Params) == 0 {
		return f.Raw
	}

	params := make([]interface{}, len(f.Params))
	for i, p := range f.Params {
		params[i] = p.Generate()
	}
	return fmt.Sprintf(f.Format, params...)
}

type Literal struct {
	Value interface{}
}

func (l Literal) Generate() interface{} { return l.Value }
