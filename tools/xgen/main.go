package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		die("xgen FILE1.go FILE2.go ...")
	}

	fset := token.NewFileSet()
	mode := parser.PackageClauseOnly | parser.ParseComments
	for _, f := range os.Args[1:] {
		if !strings.HasSuffix(f, ".go") {
			fmt.Fprintln(os.Stderr, "ignoring:", f)
			continue
		}
		a, err := parser.ParseFile(fset, f, nil, mode)
		checkErr(err)
		pkg := a.Name.Name
		for _, cg := range a.Comments {
			for len(cg.List) > 0 {
				c := strings.TrimLeft(cg.List[0].Text, "/* \t")
				switch {
				case strings.HasPrefix(c, "BaseAddr:"):
					mmio(pkg, f, cg.Text())
				default:
					cg.List = cg.List[1:]
					continue
				}
				break
			}
		}
	}
}
