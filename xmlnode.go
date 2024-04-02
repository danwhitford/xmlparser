package xmlparser

import (
	"fmt"
	"io"
)

type Attribute struct {
	Key, Value string
}

type Instruction struct {
	Name       string
	Attributes []Attribute
}

type XmlNode struct {
	Name         string
	Children     []XmlNode
	Contents     string
	Attributes   []Attribute
	Instructions []Instruction
}

func (node XmlNode) prettyPrintIndented(sb io.Writer, indent int) {
	var indentString string
	for i := 0; i < indent; i++ {
		indentString += "\t"
	}

	if len(node.Instructions) > 0 {
		for _, instruction := range node.Instructions {
			fmt.Fprintf(sb, "<?%s", instruction.Name)
			for _, attr := range instruction.Attributes {
				fmt.Fprintf(sb, " %s=\"%s\"", attr.Key, attr.Value)
			}
			fmt.Fprintln(sb, "?>")
		}
	}

	fmt.Fprintf(sb, "%s<%s", indentString, node.Name)

	if len(node.Attributes) > 0 {
		for _, attr := range node.Attributes {
			fmt.Fprintf(sb, ` %s="%s"`, attr.Key, attr.Value)
		}
	}

	if node.Contents == "" && len(node.Children) == 0 {
		fmt.Fprintf(sb, "/>\n")
		return
	} else {
		fmt.Fprint(sb, ">")
	}

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
