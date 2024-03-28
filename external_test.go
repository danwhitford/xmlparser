package xmlparser_test

import (
	_ "embed"
	"strings"
	"testing"

	"github.com/danwhitford/xmlparser"
	"github.com/google/go-cmp/cmp"
)

//go:embed "examplerss.xml"
var exampleRss string

func TestExamplePodFeed(t *testing.T) {
	root, err := xmlparser.Parse(exampleRss)
	if err != nil {
		t.Fatalf("did not want an error. %s", err)
	}

	var sb strings.Builder
	root.PrettyPrint(&sb)
	got := sb.String()

	if diff := cmp.Diff(exampleRss, got); diff != "" {
		t.Fatalf("wanted indentical parse/deparse but got diff %s", diff)
	}
}
