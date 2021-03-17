package sham

import (
	"math/rand"
	"regexp/syntax"
)

const maxRepeats int = 10

// NewRegex parses a regular expression. Regular expressions are of the Go flavor
// and use Perl flags.
func NewRegex(pattern string) (Regex, error) {
	s, err := syntax.Parse(pattern, syntax.Perl)
	if err != nil {
		return Regex{}, err
	}

	return Regex{Pattern: pattern, regex: s.Simplify()}, nil
}

// Regex holds a compiled regex. While any valid regular expression can be
// provided, only a subset will actually generate data. Every node in a
// parsed regex leads to a possible choice. During data generation, a
// random path through the parsed expression is taken. Therefore, a complciated
// expression has the potential to lead to wildly different performance on
// repeated generations.
//
// TODO: fully document nodes that can generate data
type Regex struct {
	Pattern string
	regex   *syntax.Regexp
}

// Generate traverses a parsed regular expression and generates data where
// applicable.
func (r Regex) Generate() interface{} {
	return string(r.gen(r.regex))
}

func (r Regex) gen(re *syntax.Regexp) []rune {
	rs := make([]rune, 0)
	switch re.Op {
	case syntax.OpLiteral:
		return re.Rune
	case syntax.OpStar:
		n := rand.Intn(maxRepeats)
		for i := 0; i < n; i++ {
			rs = append(rs, r.gen(re.Sub0[0])...)
		}
	case syntax.OpPlus:
		n := rand.Intn(maxRepeats-1) + 1
		for i := 0; i < n; i++ {
			rs = append(rs, r.gen(re.Sub0[0])...)
		}
	case syntax.OpConcat:
		for _, s := range re.Sub {
			rs = append(rs, r.gen(s)...)
		}
	case syntax.OpAlternate:
		return r.gen(re.Sub[rand.Intn(len(re.Sub))])
	case syntax.OpCapture:
		return r.gen(re.Sub0[0])
	case syntax.OpEmptyMatch:
		return nil
	case syntax.OpCharClass:
		r := fromCharClass(re.Rune)
		rs = append(rs, r)
	case syntax.OpQuest:
		if rand.Float64() < 0.75 {
			return r.gen(re.Sub0[0])
		}
		return nil
	}

	return rs
}

func fromCharClass(class []rune) rune {
	if len(class) == 0 {
		return 0
	}
	min := class[0]
	max := class[len(class)-1]
	return rune(rand.Int31n(max-min) + min)
}
