package xmlparser

import (
	"testing"

	"github.com/danwhitford/xmlparser/tokeniser"
	"github.com/google/go-cmp/cmp"
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
	}

	for _, tst := range table {
		ter := NewParser(tst.input)
		got, err := ter.Parse()
		if err != nil {
			t.Fatalf("failed on input '%v'. %v.", tst.input, err)
		}
		if diff := cmp.Diff(tst.want, got); diff != "" {
			t.Fatalf("failed on input '%v' with diff '%v'", tst.input, diff)
		}
	}
}
