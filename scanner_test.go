package sham

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		source  []byte
		want    []Token
		wantErr bool
	}{
		{
			name:    "Empty input",
			source:  []byte(``),
			want:    []Token{},
			wantErr: false,
		},
		{
			name:   "Tokenize simple object",
			source: []byte(`{"abc": 123, "def": true, "ghi": "xyz"}`),
			want: []Token{
				{Type: TokLBrace, Value: "{"},
				{Type: TokString, Value: "abc"},
				{Type: TokColon, Value: ":"},
				{Type: TokInteger, Value: "123"},
				{Type: TokComma, Value: ","},
				{Type: TokString, Value: "def"},
				{Type: TokColon, Value: ":"},
				{Type: TokTrue, Value: "true"},
				{Type: TokComma, Value: ","},
				{Type: TokString, Value: "ghi"},
				{Type: TokColon, Value: ":"},
				{Type: TokString, Value: "xyz"},
				{Type: TokRBrace, Value: "}"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Tokenize(tt.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("Tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize() = %v, want %v", got, tt.want)
			}
		})
	}
}
