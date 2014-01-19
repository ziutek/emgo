package main

import (
	"./gotoc"
	"code.google.com/p/go.tools/go/exact"
	"code.google.com/p/go.tools/go/types"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

var buildCtx = build.Context{
	GOARCH:      "cortexm4f",
	GOOS:        "none",
	GOROOT:      "/home/michal/P/go/github/emgo/egroot",
	GOPATH:      "/home/michal/P/go/github/emgo",
	Compiler:    "gc",
	ReleaseTags: []string{"go1.1", "go1.2"},
	CgoEnabled:  false,
}

func compile(dir string) {
	bp, err := buildCtx.ImportDir(dir, 0)
	checkErr(err)

	flist := make([]*ast.File, len(bp.GoFiles))
	fset := token.NewFileSet()

	for i, fname := range bp.GoFiles {
		fname = filepath.Join(bp.Dir, fname)
		f, err := parser.ParseFile(
			fset, fname, nil, parser.ParseComments,
		)
		checkErr(err)
		flist[i] = f
	}

	og, err := os.Create("_.go")
	checkErr(err)
	oh, err := os.Create("_.h")
	checkErr(err)
	oc, err := os.Create("_.c")
	checkErr(err)

	tc := &types.Config{Import: NewImporter().Import}
	ti := &types.Info{
		Types:   make(map[ast.Expr]types.Type),
		Values:  make(map[ast.Expr]exact.Value),
		Objects: make(map[*ast.Ident]types.Object),
	}

	path := bp.ImportPath
	if bp.Name == "main" {
		path = "main"
	}

	pkg, err := tc.Check(path, fset, flist, ti)
	checkErr(err)

	cc := gotoc.NewCC(fset, pkg, ti, MakeImports(flist))
	checkErr(cc.Compile(og, oh, oc, flist))

	og.Close()
	oh.Close()
	oc.Close()
}

func main() {
	path, err := os.Getwd()
	checkErr(err)
	compile(path)
}
