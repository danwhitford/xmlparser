package tokeniser

import (
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
	}

	for _, tst := range table {
		ter := NewTokeniser(tst.input)
		got, err := ter.Tokenise()
		if err != nil {
			t.Fatal(err)
		}
		if diff := cmp.Diff(tst.want, got); diff != "" {
			t.Fatalf("failed on input '%v' with diff '%v'", tst.input, diff)
		}
	}
}
