package sham

import "math/rand"

type Generator interface {
	Generate() interface{}
}

var TerminalGenerators = map[string]Generator{
	"name":      Name{},
	"firstName": FirstName{},
	"lastName":  LastName{},
}

func getRandomString(vals []string) string { return vals[rand.Intn(len(vals))] }

type Name struct{}

func (n Name) Generate() interface{} {
	return getRandomString(firstNames) + " " + getRandomString(lastNames)
}

type FirstName struct{}

func (n FirstName) Generate() interface{} {
	return getRandomString(firstNames)
}

type LastName struct{}

func (n LastName) Generate() interface{} {
	return getRandomString(lastNames)
}
