package parser

import (
	"errors"
	"fmt"
	"io"
)

type Token string

const (
	Unknown    Token = ""
	Integer          = "Integer"
	String           = "String"
	Identifier       = "Identifier"
	Comma            = "Comma"
	Colon            = "Colon"
	Slash            = "Slash"
	Dot              = "Dot"
	Whitespace       = "Whitespace"
)

func (t Token) String() string {
	if t == "" {
		return "Unknown"
	}
	return string(t)
}

type TokenValue struct {
	Type  Token
	Value string
}

func Tokenize(format io.Reader) ([]TokenValue, error) {
	t := tokenizer{}
	buffer := make([]byte, 64)

	cont := true
	current := []byte{}

	for cont {
		n, err := format.Read(buffer)
		if err != nil {
			if !errors.Is(err, io.EOF) {
				return nil, err
			}
			cont = false
		}

		var bytes []byte
		if n != 0 {
			bytes = buffer[:n]
		}

		for _, b := range bytes {
			ok, tok, err := swallow(current, b)
			if err != nil {
				return nil, err
			}
			if ok {
				if tok == Whitespace {
					current = []byte{b}
					continue
				}

				t.tokens = append(t.tokens, TokenValue{tok, string(current)})
				current = []byte{b}
				continue
			}

			current = append(current, b)
		}
	}

	if len(current) != 0 {
		ok, tok, err := swallow(current, '\x00')
		if err != nil {
			return nil, err
		}
		if ok {
			if tok != Whitespace {
				t.tokens = append(t.tokens, TokenValue{tok, string(current)})
			}
		} else {
			return nil, fmt.Errorf("unfinished type: %s, [%v]", tok, current)
		}
	}

	return t.tokens, nil
}

var (
	whitespaceMap = map[byte]struct{}{
		' ':    {},
		'\n':   {},
		'\r':   {},
		'\x00': {},
	}
)

func swallow(bytes []byte, lookahead byte) (done bool, tok Token, err error) {
	if len(bytes) == 0 {
		return
	}

	switch b := bytes[0]; b {
	case ',':
		return true, Comma, nil
	case ':':
		return true, Colon, nil
	case '.':
		return true, Dot, nil
	case ' ':
		tok = Whitespace

		if _, exists := whitespaceMap[lookahead]; exists {
			return false, tok, nil
		}
		return true, tok, nil
	case '"':
		if len(bytes) >= 2 && bytes[len(bytes)-1] == '"' {
			if len(bytes) >= 3 && bytes[len(bytes)-2] == '\\' && bytes[len(bytes)-3] != '\\' {
				return false, String, nil
			}
			return true, String, nil
		}
		return false, String, nil
	}

	if bytes[0] >= '0' && bytes[0] <= '9' {
		if len(bytes) == 1 {
			switch lookahead {
			case 'x', 'o', 'b':
				return false, Integer, nil
			case '.':
				return false, Integer, fmt.Errorf("bite does not support floating point values")
			}
			if lookahead >= '0' && lookahead <= '9' {
				return false, Integer, nil
			}
		}

		if _, exists := whitespaceMap[lookahead]; exists {
			return true, Integer, nil
		}
		if lookahead == ',' || lookahead == '/' {
			return true, Integer, nil
		}

		if len(bytes) > 1 && inBase([...]byte{bytes[0], bytes[1]}, lookahead) {
			return false, Integer, nil
		}

		if lookahead >= '0' && lookahead <= '9' {
			return false, Integer, nil
		}

		return true, Integer, fmt.Errorf("unknown next character immediately after integer: %c [%x]", lookahead, lookahead)
	}

	return true, Unknown, fmt.Errorf("unknown character: %c [%x]", lookahead, lookahead)
}

func inBase(base [2]byte, lookahead byte) bool {
	if base[0] != '0' {
		return lookahead >= '0' && lookahead <= '9'
	}

	switch base[1] {
	case 'x':
		return lookahead >= '0' && lookahead <= '9' ||
			lookahead >= 'a' && lookahead <= 'f' ||
			lookahead >= 'A' && lookahead <= 'F'
	case 'o':
		return lookahead >= '0' && lookahead <= '7'
	case 'b':
		return lookahead == '1' || lookahead == '0'
	}

	return lookahead >= '0' && lookahead <= '9'
}

type tokenizer struct {
	row    int
	column int

	tokens []TokenValue
}
