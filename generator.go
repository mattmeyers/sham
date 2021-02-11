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

var TerminalGenerators = map[string]Generator{
	"name":        GeneratorFunc(Name),
	"firstName":   GeneratorFunc(FirstName),
	"lastName":    GeneratorFunc(LastName),
	"phoneNumber": GeneratorFunc(PhoneNumber),
}

func getRandomString(vals []string) string { return vals[rand.Intn(len(vals))] }

func Name() interface{} {
	return getRandomString(firstNames) + " " + getRandomString(lastNames)
}

func FirstName() interface{} {
	return getRandomString(firstNames)
}

func LastName() interface{} {
	return getRandomString(lastNames)
}

func PhoneNumber() interface{} {
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
