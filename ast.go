package sham

import (
	"math/rand"
)

type Node interface {
	Generator
}

type Schema struct {
	Root Node
}

func (s Schema) Generate() interface{} { return s.Root.Generate() }

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

type Literal struct {
	Value interface{}
}

func (l Literal) Generate() interface{} { return l.Value }
