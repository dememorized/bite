package parser

import (
	"errors"
	"fmt"
	"io"
)

func Parse(format io.Reader) (Root, error) {
	tokens, err := Tokenize(format)
	if err != nil {
		return Root{}, err
	}

	p := &parser{
		tokens: tokens,
		pos:    0,
	}
	return p.parseRoot()
}

type parser struct {
	tokens []TokenValue
	pos    int
}

func (p *parser) parseRoot() (Root, error) {
	nodes := []Node{}

	for ; len(p.tokens) > p.pos; p.pos++ {
		t := p.tokens[p.pos]
		switch t.Type {
		case String, Integer:
			nodes = append(nodes, NodeLiteral{
				Value: t.Value,
				Type:  t.Type.String(),
			})
			p.pos++
		case Identifier:
			nextTok, err := p.lookahead(1)
			if err != nil && !errors.Is(err, io.EOF) {
				return Root{}, err
			}
			var n Node
			switch nextTok.Type {
			case BracketLeft:
				n, err = p.parseList()
			default:
				n, err = p.parseVariable()
			}
			if err != nil {
				return Root{}, err
			}

			nodes = append(nodes, n)
		default:
			return Root{}, fmt.Errorf("unknown token: %s // %s", t.Value, t.Type)
		}

		if p.pos >= len(p.tokens) {
			break
		}
		if p.tokens[p.pos].Type != Comma {
			return Root{}, fmt.Errorf("expected a comma, got %s // %s", t.Value, t.Type)
		}
	}

	return Root{
		Nodes: nodes,
	}, nil
}

func (p *parser) parseList() (NodeList, error) {
	t, err := p.nextToken()
	if err != nil {
		return NodeList{}, err
	}

	name := t.Value

	p.eat(BracketLeft)

	for {
		t2, err := p.nextToken()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return NodeList{}, fmt.Errorf("expected ], got EOF")
			}
			return NodeList{}, err
		}

		if t2.Type == BracketRight {
			break
		}
	}

	length, err := p.parseLen()
	if err != nil {
		return NodeList{}, err
	}

	return NodeList{
		Name:   name,
		Nodes:  nil,
		Length: length,
	}, nil
}

func (p *parser) parseVariable() (NodeVariable, error) {
	t, err := p.nextToken()
	if err != nil {
		return NodeVariable{}, err
	}

	name := t.Value

	length, err := p.parseLen()
	if err != nil {
		return NodeVariable{}, err
	}

	tpe, err := p.parseType()
	if err != nil {
		return NodeVariable{}, err
	}

	return NodeVariable{
		Name:   name,
		Length: length,
		Type:   tpe,
	}, nil
}

func (p *parser) parseLen() (string, error) {
	if p.eat(Colon) {
		t, err := p.nextToken()
		if err != nil {
			return "", err
		}

		switch t.Type {
		case Integer:
			return t.Value, nil
		case Identifier:
			return t.Value, nil
		case Dot:
			if p.eat(Dot, Dot) {
				return "...", nil
			}
			return "", fmt.Errorf("expected three dots in an ellipsis")
		default:
			return "", fmt.Errorf("expected integer or identifier for length, got: %s", t.Type)
		}
	}
	return "", nil
}

var typeKeywords = map[string]struct{}{
	"bytes":   {},
	"integer": {},
	"float":   {},
}

func (p *parser) parseType() (string, error) {
	if p.eat(Slash) {
		t, err := p.nextToken()
		if err != nil {
			return "", err
		}

		switch t.Type {
		case Identifier:
			if _, exists := typeKeywords[t.Value]; !exists {
				return "", fmt.Errorf("expected a correct type")
			}
			return t.Value, nil
		default:
			return "", fmt.Errorf("expected identifier for type, got: %s", t.Type)
		}
	}
	return "", nil
}

func (p *parser) nextToken() (TokenValue, error) {
	t, err := p.lookahead(0)
	if err != nil {
		return TokenValue{}, err
	}
	p.pos++
	return t, nil
}

func (p *parser) lookahead(offset int) (TokenValue, error) {
	if len(p.tokens) <= p.pos+offset {
		return TokenValue{}, io.EOF
	}
	return p.tokens[p.pos+offset], nil
}

func (p *parser) eat(toks ...Token) bool {
	for _, t := range toks {
		if len(p.tokens) <= p.pos {
			return false
		}

		if p.tokens[p.pos].Type != t {
			return false
		}
		p.pos++
	}
	return true
}
