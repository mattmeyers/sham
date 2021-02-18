package sham

import (
	"encoding/json"
	"strings"
)

type OrderedMap struct {
	Values map[string]interface{}
	Keys   []string
}

func NewOrderedMap() *OrderedMap {
	return &OrderedMap{
		Values: make(map[string]interface{}),
		Keys:   make([]string, 0),
	}
}

func (m *OrderedMap) Set(k string, v interface{}) {
	if _, ok := m.Values[k]; !ok {
		m.Keys = append(m.Keys, k)
	}

	m.Values[k] = v
}

func (m *OrderedMap) MarshalJSON() ([]byte, error) {
	var sb strings.Builder
	sb.WriteByte('{')

	for i, k := range m.Keys {
		sb.WriteByte('"')
		sb.WriteString(k)
		sb.WriteByte('"')
		sb.WriteByte(':')

		d, err := json.Marshal(m.Values[k])
		if err != nil {
			return nil, err
		}

		sb.Write(d)

		if i < len(m.Keys)-1 {
			sb.WriteByte(',')
		}
	}

	sb.WriteByte('}')

	return []byte(sb.String()), nil
}
