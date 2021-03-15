package gen

import (
	"math/rand"
	"strings"
)

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
