package xmlparser

import (
	"fmt"
	"io"
)

type XmlNode struct {
	Name         string
	Children     []XmlNode
	Contents     string
	Attributes   map[string]string
	Instructions map[string]map[string]string
}

func (node XmlNode) prettyPrintIndented(sb io.Writer, indent int) {
	var indentString string
	for i := 0; i < indent; i++ {
		indentString += "\t"
	}

	fmt.Fprintf(sb, "%s<%s", indentString, node.Name)

	if len(node.Attributes) > 0 {
		fmt.Fprint(sb, " ")
		for k, v := range node.Attributes {
			fmt.Fprintf(sb, `%s="%s"`, k, v)
		}
	}
	fmt.Fprint(sb, ">")

	if node.Contents != "" {
		fmt.Fprintf(sb, "%s</%s>\n", node.Contents, node.Name)
		return
	}

	if len(node.Children) > 0 {
		fmt.Fprintln(sb)
		for _, child := range node.Children {
			child.prettyPrintIndented(sb, indent+1)
		}
	}
	fmt.Fprintf(sb, "%s</%s>\n", indentString, node.Name)
}

func (node XmlNode) PrettyPrint(sb io.Writer) {
	node.prettyPrintIndented(sb, 0)
}
