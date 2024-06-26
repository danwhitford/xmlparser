package xmlparser

import (
	"fmt"
	"testing"

	"github.com/danwhitford/xmlparser/tokeniser"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestParseXml(t *testing.T) {
	table := []struct {
		input []tokeniser.Token
		want  XmlNode
	}{
		{
			[]tokeniser.Token{
				{T: tokeniser.LB, Val: "<"},
				{T: tokeniser.Keyword, Val: "foo"},
				{T: tokeniser.RB, Val: ">"},

				{T: tokeniser.Keyword, Val: "bar"},

				{T: tokeniser.CloB, Val: "</"},
				{T: tokeniser.Keyword, Val: "foo"},
				{T: tokeniser.RB, Val: ">"},
			},
			XmlNode{
				Name:     "foo",
				Contents: "bar",
			},
		},
		{
			[]tokeniser.Token{
				{T: tokeniser.LB, Val: "<"},
				{T: tokeniser.Keyword, Val: "foo"},
				{T: tokeniser.RB, Val: ">"},

				{T: tokeniser.Keyword, Val: "bar"},
				{T: tokeniser.Whitespace, Val: " "},
				{T: tokeniser.Keyword, Val: "baz"},

				{T: tokeniser.CloB, Val: "</"},
				{T: tokeniser.Keyword, Val: "foo"},
				{T: tokeniser.RB, Val: ">"},
			},
			XmlNode{
				Name:     "foo",
				Contents: "bar baz",
			},
		},
		{
			[]tokeniser.Token{ // <foo><bar>baz</bar></foo>
				{T: tokeniser.LB, Val: "<"},
				{T: tokeniser.Keyword, Val: "foo"},
				{T: tokeniser.RB, Val: ">"},

				{T: tokeniser.LB, Val: "<"},
				{T: tokeniser.Keyword, Val: "bar"},
				{T: tokeniser.RB, Val: ">"},

				{T: tokeniser.Keyword, Val: "baz"},

				{T: tokeniser.CloB, Val: "</"},
				{T: tokeniser.Keyword, Val: "bar"},
				{T: tokeniser.RB, Val: ">"},

				{T: tokeniser.CloB, Val: "</"},
				{T: tokeniser.Keyword, Val: "foo"},
				{T: tokeniser.RB, Val: ">"},
			},
			XmlNode{
				Name:     "foo",
				Contents: "",
				Children: []XmlNode{
					{
						Name:     "bar",
						Contents: "baz",
					},
				},
			},
		},
		{
			[]tokeniser.Token{ // <list><item>apples</item><item>pears></item></list>
				{T: tokeniser.LB, Val: "<"},
				{T: tokeniser.Keyword, Val: "list"},
				{T: tokeniser.RB, Val: ">"},

				{T: tokeniser.LB, Val: "<"},
				{T: tokeniser.Keyword, Val: "item"},
				{T: tokeniser.RB, Val: ">"},

				{T: tokeniser.Keyword, Val: "apples"},

				{T: tokeniser.CloB, Val: "</"},
				{T: tokeniser.Keyword, Val: "item"},
				{T: tokeniser.RB, Val: ">"},

				{T: tokeniser.LB, Val: "<"},
				{T: tokeniser.Keyword, Val: "item"},
				{T: tokeniser.RB, Val: ">"},

				{T: tokeniser.Keyword, Val: "pears"},

				{T: tokeniser.CloB, Val: "</"},
				{T: tokeniser.Keyword, Val: "item"},
				{T: tokeniser.RB, Val: ">"},

				{T: tokeniser.CloB, Val: "</"},
				{T: tokeniser.Keyword, Val: "list"},
				{T: tokeniser.RB, Val: ">"},
			},
			XmlNode{
				Name:     "list",
				Contents: "",
				Children: []XmlNode{
					{
						Name:     "item",
						Contents: "apples",
					},
					{
						Name:     "item",
						Contents: "pears",
					},
				},
			},
		},
		{
			[]tokeniser.Token{ // <foo version="1.0">
				{T: tokeniser.LB, Val: "<"},
				{T: tokeniser.Keyword, Val: "foo"},
				{T: tokeniser.Whitespace, Val: " "},
				{T: tokeniser.Keyword, Val: "version"},
				{T: tokeniser.EQ, Val: "="},
				{T: tokeniser.String, Val: "1.0"},
				{T: tokeniser.RB, Val: ">"},
			},
			XmlNode{
				Name:       "foo",
				Contents:   "",
				Children:   nil,
				Attributes: []Attribute{{"version", "1.0"}},
			},
		},
		{
			[]tokeniser.Token{ // <foo version="1.0" type="nonsense">
				{T: tokeniser.LB, Val: "<"},
				{T: tokeniser.Keyword, Val: "foo"},
				{T: tokeniser.Whitespace, Val: " "},
				{T: tokeniser.Keyword, Val: "version"},
				{T: tokeniser.EQ, Val: "="},
				{T: tokeniser.String, Val: "1.0"},
				{T: tokeniser.Keyword, Val: "type"},
				{T: tokeniser.EQ, Val: "="},
				{T: tokeniser.String, Val: "nonsense"},
				{T: tokeniser.RB, Val: ">"},
			},
			XmlNode{
				Name:     "foo",
				Contents: "",
				Children: nil,
				Attributes: []Attribute{
					{"version", "1.0"},
					{"type", "nonsense"},
				},
			},
		},
		{
			[]tokeniser.Token{ //<?xml version="1.0" encoding="UTF-8"?>
				{T: tokeniser.ProcLB, Val: "<?"},
				{T: tokeniser.Keyword, Val: "xml"},
				{T: tokeniser.Whitespace, Val: " "},

				{T: tokeniser.Keyword, Val: "version"},
				{T: tokeniser.EQ, Val: "="},
				{T: tokeniser.String, Val: "1.0"},
				{T: tokeniser.Whitespace, Val: " "},

				{T: tokeniser.Keyword, Val: "encoding"},
				{T: tokeniser.EQ, Val: "="},
				{T: tokeniser.String, Val: "UTF-8"},
				{T: tokeniser.ProcRB, Val: "?>"},
			},
			XmlNode{
				Instructions: []Instruction{
					{
						"xml",
						[]Attribute{
							{"version", "1.0"},
							{"encoding", "UTF-8"},
						},
					},
				},
			},
		},
		{
			[]tokeniser.Token{ //<?xml version="1.0" encoding="UTF-8"?>\n<foo>bar</foo>
				{T: tokeniser.ProcLB, Val: "<?"},
				{T: tokeniser.Keyword, Val: "xml"},
				{T: tokeniser.Whitespace, Val: " "},

				{T: tokeniser.Keyword, Val: "version"},
				{T: tokeniser.EQ, Val: "="},
				{T: tokeniser.String, Val: "1.0"},
				{T: tokeniser.Whitespace, Val: " "},

				{T: tokeniser.Keyword, Val: "encoding"},
				{T: tokeniser.EQ, Val: "="},
				{T: tokeniser.String, Val: "UTF-8"},
				{T: tokeniser.ProcRB, Val: "?>"},

				{T: tokeniser.Whitespace, Val: "\n"},

				{T: tokeniser.LB, Val: "<"},
				{T: tokeniser.Keyword, Val: "foo"},
				{T: tokeniser.RB, Val: ">"},

				{T: tokeniser.Keyword, Val: "bar"},

				{T: tokeniser.CloB, Val: "</"},
				{T: tokeniser.Keyword, Val: "foo"},
				{T: tokeniser.RB, Val: ">"},
			},
			XmlNode{
				Instructions: []Instruction{
					{
						"xml",
						[]Attribute{
							{"version", "1.0"},
							{"encoding", "UTF-8"},
						},
					},
				},
				Name:     "foo",
				Contents: "bar",
			},
		},
		{
			[]tokeniser.Token{ // <enclosure length="7500000" type="audio/mpeg"/>
				{T: tokeniser.LB, Val: "<"},

				{T: tokeniser.Keyword, Val: "enclosure"},
				{T: tokeniser.Whitespace, Val: " "},

				{T: tokeniser.Keyword, Val: "length"},
				{T: tokeniser.EQ, Val: "="},
				{T: tokeniser.String, Val: "7500000"},

				{T: tokeniser.Whitespace, Val: " "},

				{T: tokeniser.Keyword, Val: "type"},
				{T: tokeniser.EQ, Val: "="},
				{T: tokeniser.String, Val: "audio/mpeg"},

				{T: tokeniser.SelfRB, Val: "/>"},
			},
			XmlNode{
				Name: "enclosure",
				Attributes: []Attribute{
					{"length", "7500000"},
					{"type", "audio/mpeg"},
				},
			},
		},
		{
			[]tokeniser.Token{
				{T: tokeniser.LB, Val: "<"},
				{T: tokeniser.Keyword, Val: "url"},
				{T: tokeniser.RB, Val: ">"},

				{T: tokeniser.Keyword, Val: "https://megaphone.imgix.net/podcasts/00c0a118-2426-11ee-b258-73d331d0123b/image/show-cover.jpg?ixlib"},
				{T: tokeniser.EQ, Val: "="},
				{T: tokeniser.Keyword, Val: "rails-4.3.1"},

				{T: tokeniser.CloB, Val: "</"},
				{T: tokeniser.Keyword, Val: "url"},
				{T: tokeniser.RB, Val: ">"},
			},
			XmlNode{
				Name:     "url",
				Contents: "https://megaphone.imgix.net/podcasts/00c0a118-2426-11ee-b258-73d331d0123b/image/show-cover.jpg?ixlib=rails-4.3.1",
			},
		},
	}

	for i, tst := range table {
		t.Run(fmt.Sprintf("Test %d of %d", i+1, len(table)), func(t *testing.T) {
			ter := newParser(tst.input)
			got, err := ter.runParser()
			if err != nil {
				t.Fatalf("failed on input '%v'. %v.", tst.input, err)
			}
			if diff := cmp.Diff(tst.want, got, cmpopts.EquateEmpty()); diff != "" {
				t.Fatalf("failed on input '%v' with diff '%v'", tst.input, diff)
			}
		})
	}
}
