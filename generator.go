package sham

import "github.com/mattmeyers/sham/gen"

type Generator interface {
	Generate() interface{}
}

type GeneratorFunc func() interface{}

func (f GeneratorFunc) Generate() interface{} { return f() }

func stringAdaptor(f func() string) func() interface{} { return func() interface{} { return f() } }
func intAdaptor(f func() int) func() interface{}       { return func() interface{} { return f() } }

var TerminalGenerators = map[string]Generator{
	"name":        GeneratorFunc(stringAdaptor(gen.Name)),
	"firstName":   GeneratorFunc(stringAdaptor(gen.FirstName)),
	"lastName":    GeneratorFunc(stringAdaptor(gen.LastName)),
	"phoneNumber": GeneratorFunc(stringAdaptor(gen.PhoneNumber)),
}
