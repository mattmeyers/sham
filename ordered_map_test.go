package sham

import (
	"encoding/xml"
	"reflect"
	"testing"
)

func TestOrderedMap_MarshalJSON(t *testing.T) {
	type fields struct {
		Values map[string]interface{}
		Keys   []string
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name: "Simple map",
			fields: fields{
				Values: map[string]interface{}{
					"a": 1,
					"b": 2,
					"c": 3,
				},
				Keys: []string{"b", "c", "a"},
			},
			want:    []byte(`{"b":2,"c":3,"a":1}`),
			wantErr: false,
		},
		{
			name: "Nested maps",
			fields: fields{
				Values: map[string]interface{}{
					"a": 1,
					"b": 2,
					"c": &OrderedMap{
						Values: map[string]interface{}{
							"d": 5,
							"e": 6,
							"f": 4,
						},
						Keys: []string{"f", "d", "e"},
					},
				},
				Keys: []string{"b", "c", "a"},
			},
			want:    []byte(`{"b":2,"c":{"f":4,"d":5,"e":6},"a":1}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &OrderedMap{
				Values: tt.fields.Values,
				Keys:   tt.fields.Keys,
			}
			got, err := m.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderedMap.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderedMap.MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestOrderedMap_MarshalXML(t *testing.T) {
	tests := []struct {
		name    string
		val     *OrderedMap
		want    []byte
		wantErr bool
	}{
		{
			name: "Simple map",
			val: &OrderedMap{
				Values: map[string]interface{}{
					"a": 1,
					"b": 2,
					"c": 3,
				},
				Keys: []string{"b", "c", "a"},
			},
			want:    []byte(`<b>2</b><c>3</c><a>1</a>`),
			wantErr: false,
		},
		{
			name: "Nested maps",
			val: &OrderedMap{
				Values: map[string]interface{}{
					"a": 1,
					"b": 2,
					"c": &OrderedMap{
						Values: map[string]interface{}{
							"d": 5,
							"e": 6,
							"f": 4,
						},
						Keys: []string{"f", "d", "e"},
					},
				},
				Keys: []string{"b", "c", "a"},
			},
			want:    []byte(`<b>2</b><c><f>4</f><d>5</d><e>6</e></c><a>1</a>`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := xml.Marshal(tt.val)
			if (err != nil) != tt.wantErr {
				t.Errorf("OrderedMap.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("OrderedMap.MarshalJSON() = %s, want %s", got, tt.want)
			}
		})
	}
}
