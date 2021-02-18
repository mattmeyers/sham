package sham

import (
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
			name: "Success",
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
				t.Errorf("OrderedMap.MarshalJSON() = %v, want %v", got, tt.want)
			}
		})
	}
}
