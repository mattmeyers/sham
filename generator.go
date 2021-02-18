package sham

import (
	"math/rand"
	"strings"
)

type Generator interface {
	Generate() interface{}
}

type GeneratorFunc func() interface{}

func (f GeneratorFunc) Generate() interface{} { return f() }

func stringAdaptor(f func() string) func() interface{} { return func() interface{} { return f() } }
func intAdaptor(f func() int) func() interface{}       { return func() interface{} { return f() } }

var TerminalGenerators = map[string]Generator{
	"name":        GeneratorFunc(stringAdaptor(Name)),
	"firstName":   GeneratorFunc(stringAdaptor(FirstName)),
	"lastName":    GeneratorFunc(stringAdaptor(LastName)),
	"phoneNumber": GeneratorFunc(stringAdaptor(PhoneNumber)),
}

func getRandomString(vals []string) string { return vals[rand.Intn(len(vals))] }

func Name() string {
	return getRandomString(firstNames) + " " + getRandomString(lastNames)
}

func FirstName() string {
	return getRandomString(firstNames)
}

func LastName() string {
	return getRandomString(lastNames)
}

func PhoneNumber() string {
	const digits string = "1234567890"
	var sb strings.Builder

	for i := 0; i < 12; i++ {
		if i == 3 || i == 7 {
			sb.WriteRune('-')
		} else {
			sb.WriteByte(digits[rand.Intn(len(digits))])
		}
	}

	return sb.String()
}
