package parser_test

import (
	"bite/parser"
	"reflect"
	"strings"
	"testing"
)

func FuzzParser(f *testing.F) {
	f.Add("0x89, \"PNG\\r\\n\", 0x1A, \"\\n\"")
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = parser.Parse(strings.NewReader(s))
	})
}

func TestParser(t *testing.T) {
	tests := []struct {
		fmt     string
		want    parser.Root
		wantErr bool
	}{
		{
			fmt:  "",
			want: parser.Root{Nodes: []parser.Node{}},
		},
		{
			fmt:  "0x89, \"PNG\\r\\n\", 0x1A, \"\\n\"",
			want: parser.Root{Nodes: []parser.Node{parser.NodeLiteral{Value: "0x89", Type: "Integer"}, parser.NodeLiteral{Value: "\"PNG\\r\\n\"", Type: "String"}, parser.NodeLiteral{Value: "0x1A", Type: "Integer"}, parser.NodeLiteral{Value: "\"\\n\"", Type: "String"}}},
		},
		{
			fmt:  "Len:32/integer, Type:32, Chunk:Len, CRC:32",
			want: parser.Root{Nodes: []parser.Node{parser.NodeVariable{Name: "Len", Length: "32", Type: "integer"}, parser.NodeVariable{Name: "Type", Length: "32", Type: ""}, parser.NodeVariable{Name: "Chunk", Length: "Len", Type: ""}, parser.NodeVariable{Name: "CRC", Length: "32", Type: ""}}},
		},
		{
			fmt:  "Var:...",
			want: parser.Root{Nodes: []parser.Node{parser.NodeVariable{Name: "Var", Length: "..."}}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.fmt, func(t *testing.T) {
			got, err := parser.Parse(strings.NewReader(tt.fmt))
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() got = %#v, want %#v", got, tt.want)
			}
		})
	}
}
