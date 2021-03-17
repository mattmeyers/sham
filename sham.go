// Package sham generates pseudorandom data from a supplied schema written in
// the Sham language.
package sham

// Generate parses a Sham schema, and on success, performs a single generation
// of data using the default terminal generators. This function is intended to
// be a simple wrapper for the Sham data generation process. If multiple
// generations are required, or custom terminal generators are needed, then
// a parser should be instantiated with NewParser. After successfully parsing,
// the resulting Schema object can be used to generate data multiple times
// without parsing the schema.
func Generate(schema []byte) (interface{}, error) {
	s, err := NewDefaultParser(schema).Parse()
	if err != nil {
		return nil, err
	}

	return s.Generate(), nil
}
