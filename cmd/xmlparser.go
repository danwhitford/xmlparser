package main

import (
	"fmt"
	"os"

	"github.com/danwhitford/xmlparser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("need argument for filename")
		return
	}
	path := os.Args[1]
	raw, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("error reading file at '%s'. %v\n", path, err)
	}
	root, err := xmlparser.Parse(string(raw))
	if err != nil {
		fmt.Printf("error parsing xml. %v\n", err)
	}
	root.PrettyPrint(os.Stdout)
}
