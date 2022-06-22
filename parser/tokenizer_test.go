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
		{
			fmt:  "Len:32/integer, Type:32, Chunk:Len, CRC:32",
			want: []parser.TokenValue{{Type: "Identifier", Value: "Len"}, {Type: "Colon", Value: ":"}, {Type: "Integer", Value: "32"}, {Type: "Slash", Value: "/"}, {Type: "Identifier", Value: "integer"}, {Type: "Comma", Value: ","}, {Type: "Identifier", Value: "Type"}, {Type: "Colon", Value: ":"}, {Type: "Integer", Value: "32"}, {Type: "Comma", Value: ","}, {Type: "Identifier", Value: "Chunk"}, {Type: "Colon", Value: ":"}, {Type: "Identifier", Value: "Len"}, {Type: "Comma", Value: ","}, {Type: "Identifier", Value: "CRC"}, {Type: "Colon", Value: ":"}, {Type: "Integer", Value: "32"}},
		},
		{
			fmt:  "Chunk[Foo:64]:...",
			want: []parser.TokenValue{{Type: "Identifier", Value: "Chunk"}, {Type: "BracketLeft", Value: "["}, {Type: "Identifier", Value: "Foo"}, {Type: "Colon", Value: ":"}, {Type: "Integer", Value: "64"}, {Type: "BracketRight", Value: "]"}, {Type: "Colon", Value: ":"}, {Type: "Dot", Value: "."}, {Type: "Dot", Value: "."}, {Type: "Dot", Value: "."}},
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
