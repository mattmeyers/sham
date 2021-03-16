package sham

import (
	"encoding/json"
	"encoding/xml"
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

func (m *OrderedMap) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if start.Name.Local != "OrderedMap" {
		if err := e.EncodeToken(start); err != nil {
			return err
		}
	}

	for _, key := range m.Keys {
		err := e.EncodeElement(m.Values[key], xml.StartElement{Name: xml.Name{Space: "", Local: key}})
		if err != nil {
			return err
		}
	}

	if start.Name.Local != "OrderedMap" {
		if err := e.EncodeToken(xml.EndElement{Name: start.Name}); err != nil {
			return err
		}
	}

	return e.Flush()
}
