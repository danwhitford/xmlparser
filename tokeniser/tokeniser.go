package tokeniser

import (
	"strings"
	"unicode"
)

//go:generate stringer -type=TokenType
type TokenType int

const (
	Keyword TokenType = iota
	LB
	RB
	EQ
	String
	CloB
	Whitespace
)

type Token struct {
	T   TokenType
	Val any
}

type Tokeniser struct {
	Input string
	curr  int
	l     int
}

func NewTokeniser(input string) Tokeniser {
	return Tokeniser{
		input,
		0,
		len(input),
	}
}

func (t *Tokeniser) Tokenise() ([]Token, error) {
	var tokens []Token
	for t.curr < t.l {
		switch t.Input[t.curr] {
		case ' ':
			token, err := t.getWhitespace()
			if err != nil {
				return tokens, err
			}
			tokens = append(tokens, token)
		case '<':
			if t.curr+1 < t.l && t.Input[t.curr+1] == '/' {
				tokens = append(tokens, Token{CloB, "</"})
				t.curr += 2
			} else {
				tokens = append(tokens, Token{LB, "<"})
				t.curr++
			}
		case '>':
			tokens = append(tokens, Token{RB, ">"})
			t.curr++
		case '=':
			tokens = append(tokens, Token{EQ, "="})
			t.curr++
		case '"':
			token, err := t.getString()
			if err != nil {
				return tokens, err
			}
			tokens = append(tokens, token)
		default:
			keyword, err := t.getKeyword()
			if err != nil {
				return tokens, err
			}
			tokens = append(tokens, keyword)
		}
	}

	return tokens, nil
}

func (t *Tokeniser) getKeyword() (Token, error) {
	var sb strings.Builder
	for t.curr < t.l {
		peek := t.Input[t.curr]
		if peek == '>' || peek == ' ' || peek == '=' || peek == '<' {
			break
		}
		sb.WriteByte(t.Input[t.curr])
		t.curr++
	}
	return Token{
		T:   Keyword,
		Val: sb.String(),
	}, nil
}

func (t *Tokeniser) getString() (Token, error) {
	var sb strings.Builder
	t.curr++ // eat the opening quotes
	for t.curr < t.l {
		peek := t.Input[t.curr]
		if peek == '"' {
			t.curr++
			break
		}
		sb.WriteByte(t.Input[t.curr])
		t.curr++
	}
	return Token{
		T:   String,
		Val: sb.String(),
	}, nil
}

func (t *Tokeniser) getWhitespace() (Token, error) {
	var sb strings.Builder
	for t.curr < t.l {
		peek := t.Input[t.curr]
		if !unicode.IsSpace(rune(peek)) {
			break
		}
		sb.WriteByte(t.Input[t.curr])
		t.curr++
	}
	return Token{
		T:   Whitespace,
		Val: sb.String(),
	}, nil
}
