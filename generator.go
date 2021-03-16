package sham

import (
	"time"

	"github.com/mattmeyers/sham/gen"
)

// Generator represents the core functionality behind Sham's data generation.
// Any type that implements this interface can be used to generate data. Implementors
// can be either data structures containing more data, or simpler functions that
// directly generate a sinlge piece of data. These latter objects are referred to
// as terminal generators since they are generally found as leaves in the AST.
type Generator interface {
	Generate() interface{}
}

// GeneratorFunc is a simple function type that implements the Generator interface.
// This type can be used to provide single functions as Generators.
type GeneratorFunc func() interface{}

func (f GeneratorFunc) Generate() interface{} { return f() }

func stringAdaptor(f func() string) func() interface{}  { return func() interface{} { return f() } }
func intAdaptor(f func() int) func() interface{}        { return func() interface{} { return f() } }
func timeAdaptor(f func() time.Time) func() interface{} { return func() interface{} { return f() } }

// TerminalGenerators is the standard collection of terminal generators provided by Sham.
var TerminalGenerators = map[string]Generator{
	"name":        GeneratorFunc(stringAdaptor(gen.Name)),
	"firstName":   GeneratorFunc(stringAdaptor(gen.FirstName)),
	"lastName":    GeneratorFunc(stringAdaptor(gen.LastName)),
	"phoneNumber": GeneratorFunc(stringAdaptor(gen.PhoneNumber)),
	"timestamp":   GeneratorFunc(timeAdaptor(gen.Timestamp)),
}
