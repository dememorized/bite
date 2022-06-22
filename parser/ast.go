package parser

import (
	"fmt"
	"strings"
)

type Node interface {
	astNode()
}

type NodeLiteral struct {
	Value string
	Type  string
}

func (n NodeLiteral) String() string {
	return n.Value
}

func (n NodeLiteral) astNode() {}

type NodeVariable struct {
	Name   string
	Length string
	Type   string
}

func (n NodeVariable) String() string {
	s := n.Name
	if n.Length != "" {
		s += ":" + n.Length
	}
	if n.Type != "" {
		s += "/" + n.Type
	}
	return s
}

func (n NodeVariable) astNode() {}

type NodeList struct {
	Name   string
	Nodes  []Node
	Length string
}

func (n NodeList) String() string {
	s := n.Name + "["

	strs := []string{}
	for _, node := range n.Nodes {
		strs = append(strs, fmt.Sprintf("%s", node))
	}
	s += strings.Join(strs, ", ")
	s += "]"

	if n.Length != "" {
		s += ":" + n.Length
	}
	return s
}

func (n NodeList) astNode() {}

type Root struct {
	Nodes []Node
}

func (r Root) String() string {
	strs := []string{}
	for _, n := range r.Nodes {
		strs = append(strs, fmt.Sprintf("%s", n))
	}
	return strings.Join(strs, ", ")
}
