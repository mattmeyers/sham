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
			name: "Tokenize array of objects",
			source: []byte(`[{"a": 123,    "b": true, "c": "def",	"g": -1.56e-6, "h":name}]  `),
			want: []Token{
				{Type: TokLBracket, Value: "["},
				{Type: TokLBrace, Value: "{"},
				{Type: TokString, Value: "a"},
				{Type: TokColon, Value: ":"},
				{Type: TokInteger, Value: "123"},
				{Type: TokComma, Value: ","},
				{Type: TokString, Value: "b"},
				{Type: TokColon, Value: ":"},
				{Type: TokTrue, Value: "true"},
				{Type: TokComma, Value: ","},
				{Type: TokString, Value: "c"},
				{Type: TokColon, Value: ":"},
				{Type: TokString, Value: "def"},
				{Type: TokComma, Value: ","},
				{Type: TokString, Value: "g"},
				{Type: TokColon, Value: ":"},
				{Type: TokFloat, Value: "-1.56e-6"},
				{Type: TokComma, Value: ","},
				{Type: TokString, Value: "h"},
				{Type: TokColon, Value: ":"},
				{Type: TokIdent, Value: "name"},
				{Type: TokRBrace, Value: "}"},
				{Type: TokRBracket, Value: "]"},
			},
			wantErr: false,
		},
		{
			name:   "Tokenize range",
			source: []byte(`(1, 2)`),
			want: []Token{
				{Type: TokLParen, Value: "("},
				{Type: TokInteger, Value: "1"},
				{Type: TokComma, Value: ","},
				{Type: TokInteger, Value: "2"},
				{Type: TokRParen, Value: ")"},
			},
			wantErr: false,
		},
		{
			name:   "Tokenize fstring",
			source: []byte("`foo ${bar}`"),
			want: []Token{
				{Type: TokFString, Value: "foo ${bar}"},
			},
			wantErr: false,
		},
		{
			name:   "Tokenize regex",
			source: []byte(`/^abc123.*\/(|)$/`),
			want: []Token{
				{Type: TokRegex, Value: `^abc123.*\/(|)$`},
			},
			wantErr: false,
		},
		{
			name:    "Unterminated string",
			source:  []byte(`"abc`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Unterminated regex",
			source:  []byte(`/ab`),
			want:    nil,
			wantErr: true,
		},
		{
			name:    "Unknown token",
			source:  []byte(`{}>`),
			want:    nil,
			wantErr: true,
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
