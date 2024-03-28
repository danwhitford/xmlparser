package xmlparser

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestPrintXml(t *testing.T) {
	table := []struct {
		input XmlNode
		want  string
	}{
		{
			XmlNode{
				Name:     "foo",
				Contents: "bar",
			},
			"<foo>bar</foo>\n",
		},
		{
			XmlNode{
				Name:     "foo",
				Contents: "bar",
				Attributes: []Attribute{
					{"version", "1.0"},
				},
			},
			"<foo version=\"1.0\">bar</foo>\n",
		},
		{
			XmlNode{
				Name: "foo",
				Children: []XmlNode{
					{Name: "item", Contents: "Item Contents"},
					{Name: "item", Contents: "Item Contents"},
					{Name: "item", Contents: "Item Contents"},
				},
			},
			"<foo>\n" +
				"\t<item>Item Contents</item>\n" +
				"\t<item>Item Contents</item>\n" +
				"\t<item>Item Contents</item>\n" +
				"</foo>\n",
		},
		{
			XmlNode{
				Name: "parent",
				Children: []XmlNode{
					{Name: "item", Children: []XmlNode{
						{Name: "granditem", Contents: "foo"},
						{Name: "granditem", Contents: "bar"},
					}},
				},
			},
			"<parent>\n" +
				"\t<item>\n" +
				"\t\t<granditem>foo</granditem>\n" +
				"\t\t<granditem>bar</granditem>\n" +
				"\t</item>\n" +
				"</parent>\n",
		},
		{
			XmlNode{
				Name:     "foo",
				Contents: "bar",
				Attributes: []Attribute{
					{"version", "1.0"},
					{"type", "test"},
				},
			},
			"<foo version=\"1.0\" type=\"test\">bar</foo>\n",
		},
	}

	for _, tst := range table {
		var sb strings.Builder
		tst.input.PrettyPrint(&sb)
		got := sb.String()

		if diff := cmp.Diff(tst.want, got); diff != "" {
			t.Fatalf("failed on input '%v' with diff '%v'", tst.input, diff)
		}
	}
}
