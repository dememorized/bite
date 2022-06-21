package parser_test

import (
	"bite/parser"
	"reflect"
	"strings"
	"testing"
)

func FuzzTokenize(f *testing.F) {
	f.Add("0x89, \"PNG\\r\\n\", 0x1A, \"\\n\"")
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = parser.Tokenize(strings.NewReader(s))
	})
}

func TestTokenize(t *testing.T) {
	tests := []struct {
		fmt     string
		want    []parser.TokenValue
		wantErr bool
	}{
		{
			fmt:  "",
			want: nil,
		},
		{
			fmt:  "0x89, \"PNG\\r\\n\", 0x1A, \"\\n\"",
			want: []parser.TokenValue{{Type: "Integer", Value: "0x89"}, {Type: "Comma", Value: ","}, {Type: "String", Value: "\"PNG\\r\\n\""}, {Type: "Comma", Value: ","}, {Type: "Integer", Value: "0x1A"}, {Type: "Comma", Value: ","}, {Type: parser.String, Value: "\"\\n\""}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.fmt, func(t *testing.T) {
			got, err := parser.Tokenize(strings.NewReader(tt.fmt))
			if (err != nil) != tt.wantErr {
				t.Errorf("Tokenize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Tokenize() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
