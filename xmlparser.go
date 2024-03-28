package xmlparser

import (
	"fmt"
	"strings"

	"github.com/danwhitford/xmlparser/tokeniser"
)

type parser struct {
	Input []tokeniser.Token
	curr  int
	l     int
}

func Parse(input string) (XmlNode, error) {
	t := tokeniser.NewTokeniser(input)
	tokens, err := t.Tokenise()
	if err != nil {
		return XmlNode{}, fmt.Errorf("error tokenising. %s", err)
	}
	p := newParser(tokens)
	out, err := p.runParser()
	if err != nil {
		return XmlNode{}, fmt.Errorf("error running parser. %s", err)
	}
	return out, nil
}

func newParser(input []tokeniser.Token) parser {
	return parser{
		input,
		0,
		len(input),
	}
}

func (p *parser) runParser() (XmlNode, error) {
	root := XmlNode{}
	root.Attributes = make(map[string]string)
	root.Instructions = make(map[string]map[string]string)

	if p.Peek().T == tokeniser.ProcLB {
		err := p.readProcessingInstruction(&root)
		if err != nil {
			return root, fmt.Errorf("error reading a processing instruction. %s", err)
		}
	}

	for p.curr < p.l {
		if p.Peek().T == tokeniser.Whitespace {
			_, err := p.readNext(tokeniser.Whitespace)
			if err != nil {
				return root, fmt.Errorf("error skipping through whitespace. %v", err)
			}
		} else {
			break
		}
	}

	if p.curr < p.l {
		err := p.readOpeningTag(&root)
		if err != nil {
			return root, fmt.Errorf("error reading opening tag. %v. %#v", err, p.Input[:p.curr+1])
		}
	}

	for p.curr < p.l {
		switch p.Peek().T {
		case tokeniser.Keyword:
			contents, err := p.readContents()
			if err != nil {
				return root, err
			}
			root.Contents = contents
		case tokeniser.LB:
			child, err := p.runParser()
			if err != nil {
				return root, err
			}
			root.Children = append(root.Children, child)
		case tokeniser.CloB:
			err := p.chompClosingTag(root.Name)
			if err != nil {
				return root, err
			}
			return root, nil
		case tokeniser.Whitespace:
			_, err := p.readNext(tokeniser.Whitespace)
			if err != nil {
				return root, err
			}
		default:
			return root, fmt.Errorf("dunno what to do with '%v'", p.Peek())
		}
	}

	return root, nil
}

func (p *parser) readNext(expected tokeniser.TokenType) (tokeniser.Token, error) {
	if p.curr >= p.l {
		return tokeniser.Token{}, fmt.Errorf("at end of input but expecting '%v'", expected)
	}
	t := p.Input[p.curr]
	if t.T != expected {
		return t, fmt.Errorf("token incorrect type. want '%v' got '%v'", expected, t.T)
	}
	p.curr++
	return t, nil
}

func (p *parser) Peek() tokeniser.Token {
	return p.Input[p.curr]
}

func (p *parser) readOpeningTag(root *XmlNode) error {
	_, err := p.readNext(tokeniser.LB)
	if err != nil {
		return fmt.Errorf("failed to read name tag. %v", err)
	}

	nameToken, err := p.readNext(tokeniser.Keyword)
	if err != nil {
		return err
	}

	root.Name = nameToken.Val

	for {
		switch p.Peek().T {
		case tokeniser.RB:
			_, err = p.readNext(tokeniser.RB)
			if err != nil {
				return err
			}
			return nil
		case tokeniser.Whitespace:
			_, err = p.readNext(tokeniser.Whitespace)
			if err != nil {
				return err
			}
		case tokeniser.Keyword:
			key, val, err := p.readAttr()
			if err != nil {
				return fmt.Errorf("error reading attr. %s", err)
			}
			root.Attributes[key] = val
		default:
			return fmt.Errorf("did not expect '%v' while reading opening tag", p.Peek())
		}
	}
}

func (p *parser) readProcessingInstruction(root *XmlNode) error {
	_, err := p.readNext(tokeniser.ProcLB)
	if err != nil {
		return fmt.Errorf("failed to read name tag. %v", err)
	}

	nameToken, err := p.readNext(tokeniser.Keyword)
	if err != nil {
		return err
	}

	root.Instructions[nameToken.Val] = make(map[string]string)

	for {
		switch p.Peek().T {
		case tokeniser.ProcRB:
			_, err = p.readNext(tokeniser.ProcRB)
			if err != nil {
				return err
			}
			return nil
		case tokeniser.Whitespace:
			_, err = p.readNext(tokeniser.Whitespace)
			if err != nil {
				return err
			}
		case tokeniser.Keyword:
			key, val, err := p.readAttr()
			if err != nil {
				return fmt.Errorf("error reading attr. %s", err)
			}
			root.Instructions[nameToken.Val][key] = val
		default:
			return fmt.Errorf("did not expect '%v' while reading processing instruction", p.Peek())
		}
	}
}

func (p *parser) readAttr() (string, string, error) {
	key, err := p.readNext(tokeniser.Keyword)
	if err != nil {
		return "", "", err
	}
	_, err = p.readNext(tokeniser.EQ)
	if err != nil {
		return "", "", err
	}
	val, err := p.readNext(tokeniser.String)
	if err != nil {
		return "", "", err
	}
	return key.Val, val.Val, nil
}

func (p *parser) readContents() (string, error) {
	var sb strings.Builder

	for p.curr < p.l {
		switch p.Peek().T {
		case tokeniser.Keyword:
			t, err := p.readNext(tokeniser.Keyword)
			if err != nil {
				return "", err
			}
			sb.WriteString(t.Val)

		case tokeniser.Whitespace:
			t, err := p.readNext(tokeniser.Whitespace)
			if err != nil {
				return "", err
			}
			sb.WriteString(t.Val)

		default:
			return sb.String(), nil
		}
	}

	return sb.String(), nil
}

func (p *parser) chompClosingTag(rootName string) error {
	_, err := p.readNext(tokeniser.CloB)
	if err != nil {
		return fmt.Errorf("error while chomping at position %d. %v", p.curr, err)
	}
	nameToken, err := p.readNext(tokeniser.Keyword)
	if err != nil {
		return fmt.Errorf("error while chomping. %v", err)
	}
	name := nameToken.Val
	if name != rootName {
		return fmt.Errorf("'%v' did not match '%v'", name, rootName)
	}
	_, err = p.readNext(tokeniser.RB)
	if err != nil {
		return fmt.Errorf("error while chomping. %v", err)
	}
	return nil
}
