// Package sham generates pseudorandom data from a supplied schema written in
// the Sham language.
package sham

func Generate(schema []byte) (interface{}, error) {
	s, err := NewParser(schema).Parse()
	if err != nil {
		return nil, err
	}

	return s.Generate(), nil
}
