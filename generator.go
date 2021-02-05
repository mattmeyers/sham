package sham

import "math/rand"

type Generator interface {
	Generate() interface{}
}

var terminalGenerators = map[string]Generator{
	"name": Name{},
}

func getRandomString(vals []string) string { return vals[rand.Intn(len(vals))] }

type Name struct{}

func (n Name) Generate() interface{} {
	return getRandomString(firstNames) + " " + getRandomString(lastNames)
}
