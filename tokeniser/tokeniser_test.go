package tokeniser

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestTokenise(t *testing.T) {
	table := []struct {
		input string
		want  []Token
	}{
		{
			`foo`,
			[]Token{{T: Keyword, Val: "foo"}},
		},
		{
			`<`,
			[]Token{{T: LB, Val: "<"}},
		},
		{
			`<foo>`,
			[]Token{
				{T: LB, Val: "<"},
				{T: Keyword, Val: "foo"},
				{T: RB, Val: ">"},
			},
		},
		{
			`<foo version="1.0">`,
			[]Token{
				{T: LB, Val: "<"},
				{T: Keyword, Val: "foo"},
				{T: Whitespace, Val: " "},
				{T: Keyword, Val: "version"},
				{T: EQ, Val: "="},
				{T: String, Val: "1.0"},
				{T: RB, Val: ">"},
			},
		},
		{
			`<body>Don't forget me this weekend!</body>`,
			[]Token{
				{T: LB, Val: "<"},
				{T: Keyword, Val: "body"},
				{T: RB, Val: ">"},
				{T: Keyword, Val: "Don't"},
				{T: Whitespace, Val: " "},
				{T: Keyword, Val: "forget"},
				{T: Whitespace, Val: " "},
				{T: Keyword, Val: "me"},
				{T: Whitespace, Val: " "},
				{T: Keyword, Val: "this"},
				{T: Whitespace, Val: " "},
				{T: Keyword, Val: "weekend!"},

				{T: CloB, Val: "</"},
				{T: Keyword, Val: "body"},
				{T: RB, Val: ">"},
			},
		},
		{
			`<?xml version="1.0" encoding="UTF-8"?>`,
			[]Token{
				{T: ProcLB, Val: "<?"},
				{T: Keyword, Val: "xml"},
				{T: Whitespace, Val: " "},

				{T: Keyword, Val: "version"},
				{T: EQ, Val: "="},
				{T: String, Val: "1.0"},
				{T: Whitespace, Val: " "},

				{T: Keyword, Val: "encoding"},
				{T: EQ, Val: "="},
				{T: String, Val: "UTF-8"},
				{T: ProcRB, Val: "?>"},
			},
		},
		{
			`<enclosure length="7500000" type="audio/mpeg"/>`,
			[]Token{ // <enclosure length="7500000" type="audio/mpeg"/>
				{T: LB, Val: "<"},

				{T: Keyword, Val: "enclosure"},
				{T: Whitespace, Val: " "},

				{T: Keyword, Val: "length"},
				{T: EQ, Val: "="},
				{T: String, Val: "7500000"},

				{T: Whitespace, Val: " "},

				{T: Keyword, Val: "type"},
				{T: EQ, Val: "="},
				{T: String, Val: "audio/mpeg"},

				{T: SelfRB, Val: "/>"},
			},
		},
		{
			`<?xml version="1.0" encoding="UTF-8"?>`,
			[]Token{ //<?xml version="1.0" encoding="UTF-8"?>
				{T: ProcLB, Val: "<?"},
				{T: Keyword, Val: "xml"},
				{T: Whitespace, Val: " "},

				{T: Keyword, Val: "version"},
				{T: EQ, Val: "="},
				{T: String, Val: "1.0"},
				{T: Whitespace, Val: " "},

				{T: Keyword, Val: "encoding"},
				{T: EQ, Val: "="},
				{T: String, Val: "UTF-8"},
				{T: ProcRB, Val: "?>"},
			},
		},
		{
			`<foo>"problem"? no</foo>`,
			[]Token{
				{T: LB, Val: "<"},
				{T: Keyword, Val: "foo"},
				{T: RB, Val: ">"},

				{T: String, Val: "problem"},
				{T: Keyword, Val: "?"},
				{T: Whitespace, Val: " "},
				{T: Keyword, Val: "no"},

				{T: CloB, Val: "</"},
				{T: Keyword, Val: "foo"},
				{T: RB, Val: ">"},
			},
		},
		{
			`remember / the`,
			[]Token{
				{T: Keyword, Val: "remember"},
				{T: Whitespace, Val: " "},
				{T: Keyword, Val: "/"},
				{T: Whitespace, Val: " "},
				{T: Keyword, Val: "the"},
			},
		},
		{
			"<itunes:image href=\"https://megaphone.imgix.net/podcasts/2a28152c-2426-11ee-88f1-bb16dd76d190/image/show-cover.jpg?ixlib=rails-4.3.1&amp;max-w=3000&amp;max-h=3000&amp;fit=crop&amp;auto=format,compress\"/>",
			[]Token{
				{T: LB, Val: "<"},
				{T: Keyword, Val: "itunes:image"},

				{T: Whitespace, Val: " "},

				{T: Keyword, Val: "href"},
				{T: EQ, Val: "="},
				{T: String, Val: "https://megaphone.imgix.net/podcasts/2a28152c-2426-11ee-88f1-bb16dd76d190/image/show-cover.jpg?ixlib=rails-4.3.1&amp;max-w=3000&amp;max-h=3000&amp;fit=crop&amp;auto=format,compress"},
				{T: SelfRB, Val: "/>"},
			},
		},
		{
			"<url>https://megaphone.imgix.net/podcasts/00c0a118-2426-11ee-b258-73d331d0123b/image/show-cover.jpg?ixlib=rails-4.3.1</url>",
			[]Token{
				{T: LB, Val: "<"},
				{T: Keyword, Val: "url"},
				{T: RB, Val: ">"},

				{T: Keyword, Val: "https://megaphone.imgix.net/podcasts/00c0a118-2426-11ee-b258-73d331d0123b/image/show-cover.jpg?ixlib"},
				{T: EQ, Val: "="},
				{T: Keyword, Val: "rails-4.3.1"},

				{T: CloB, Val: "</"},
				{T: Keyword, Val: "url"},
				{T: RB, Val: ">"},
			},
		},
	}

	for i, tst := range table {
		t.Run(fmt.Sprintf("Test %d/%d", i, len(table)), func(t *testing.T) {
			ter := NewTokeniser(tst.input)
			got, err := ter.Tokenise()
			if err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tst.want, got); diff != "" {
				t.Fatalf("failed on input '%v' with diff '%v' (-want +got)", tst.input, diff)
			}
		})
	}
}
