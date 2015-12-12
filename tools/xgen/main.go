package main

import (
	"fmt"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

func xgen(f string) {
	fset := token.NewFileSet()
	a, err := parser.ParseFile(fset, f, nil, parser.ParseComments)
	checkErr(err)
	pkg := a.Name.Name
	for _, cg := range a.Comments {
		for len(cg.List) > 0 {
			c := strings.TrimLeft(cg.List[0].Text, "/*")
			c = strings.TrimSpace(c)
			switch {
			case strings.HasPrefix(c, "BaseAddr:"):
				one(pkg, f, cg.Text())
			case strings.HasPrefix(c, "Peripheral:"):
				multi(pkg, f, cg.Text(), a.Decls)
			default:
				cg.List = cg.List[1:]
				continue
			}
			return
		}
	}
}

// TODO: Checking masks for registers.

func main() {
	if len(os.Args) < 2 {
		die("xgen FILE1.go FILE2.go ...")
	}

	for _, f := range os.Args[1:] {
		if !strings.HasSuffix(f, ".go") {
			fmt.Fprintln(os.Stderr, "ignoring:", f)
			continue
		}
		xgen(f)
	}
}
