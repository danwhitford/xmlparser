package xmlparser

import (
	"fmt"
	"strings"

	"github.com/danwhitford/xmlparser/tokeniser"
)

type XmlNode struct {
	Name     string
	Children []XmlNode
	Contents string
}

type Parser struct {
	Input []tokeniser.Token
	curr  int
	l     int
}

func NewParser(input []tokeniser.Token) Parser {
	return Parser{
		input,
		0,
		len(input),
	}
}

func (p *Parser) Parse() (XmlNode, error) {
	root := XmlNode{}

	name, err := p.readName()
	if err != nil {
		return root, fmt.Errorf("error at '%d'. %v", p.curr, err)
	}
	root.Name = name

	for p.curr < p.l {
		switch p.Peek().T {
		case tokeniser.Keyword:
			contents, err := p.readContents()
			if err != nil {
				return root, err
			}
			root.Contents = contents
		case tokeniser.LB:
			child, err := p.Parse()
			if err != nil {
				return root, err
			}
			root.Children = append(root.Children, child)
		case tokeniser.CloB:
			err = p.chompClosingTag(root.Name)
			if err != nil {
				return root, err
			}
			return root, nil
		default:
			return root, fmt.Errorf("dunno what to do with '%v'", p.Peek())
		}		
	}
	
	return root, nil
}

func (p *Parser) readNext(expected tokeniser.TokenType) (tokeniser.Token, error) {
	t := p.Input[p.curr]
	if t.T != expected {
		return t, fmt.Errorf("token incorrect type. want '%v' got '%v'", expected, t.T)
	}
	p.curr++
	return t, nil
}

func (p *Parser) Peek() tokeniser.Token {
	return p.Input[p.curr]
}

func (p *Parser) readName() (string, error) {
	_, err := p.readNext(tokeniser.LB)
	if err != nil {
		return "", fmt.Errorf("failed to read name tag. %v", err)
	}

	nameToken, err := p.readNext(tokeniser.Keyword)
	if err != nil {
		return "", err
	}
	name, ok := nameToken.Val.(string)
	if !ok {
		return "", fmt.Errorf("could not convert token to string %v", nameToken)
	}

	_, err = p.readNext(tokeniser.RB)
	if err != nil {
		return "", err
	}

	return name, nil
}

func (p *Parser) readContents() (string, error) {
	var sb strings.Builder

	for p.curr < p.l {
		switch p.Peek().T {
		case tokeniser.Keyword:
			t, err := p.readNext(tokeniser.Keyword)
			if err != nil {
				return "", err
			}
			s, ok := t.Val.(string)
			if !ok {
				return "", err
			}
			sb.WriteString(s)

		case tokeniser.Whitespace:
			t, err := p.readNext(tokeniser.Whitespace)
			if err != nil {
				return "", err
			}
			s, ok := t.Val.(string)
			if !ok {
				return "", fmt.Errorf("could not convert token to string '%v'", t)
			}
			sb.WriteString(s)

		default:
			return sb.String(), nil
		}
	}

	return sb.String(), nil
}

func (p *Parser) chompClosingTag(rootName string) error {
	_, err := p.readNext(tokeniser.CloB)
	if err != nil {
		return  fmt.Errorf("error while chomping at position %d. %v", p.curr, err)
	}
	nameToken, err := p.readNext(tokeniser.Keyword)
	if err != nil {
		return  fmt.Errorf("error while chomping. %v", err)
	}
	name, ok := nameToken.Val.(string)
	if !ok {
		return  fmt.Errorf("could not convert token to string '%v'", nameToken)
	}
	if name != rootName {
		return  fmt.Errorf("'%v' did not match '%v'", name, rootName)
	}
	_, err = p.readNext(tokeniser.RB)
	if err != nil {
		return  fmt.Errorf("error while chomping. %v", err)
	}
	return nil
}