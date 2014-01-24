package main

import (
	"code.google.com/p/go.tools/go/exact"
	"code.google.com/p/go.tools/go/importer"
	"code.google.com/p/go.tools/go/types"
	//	"fmt"
	"github.com/ziutek/emgo/gotoc"
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
	GOPATH:      "/home/michal/P/go/github/emgo/egpath",
	Compiler:    "gc",
	ReleaseTags: []string{"go1.1", "go1.2"},
	CgoEnabled:  false,
}

func compile(dir string) {
	bp, err := buildCtx.ImportDir(dir, 0)
	checkErr(err)

	//fmt.Printf("package \"%s\"\n%+v\n", dir, bp)

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

	wp, err := os.Create("__.EXPORTS")
	checkErr(err)
	defer wp.Close()
	wh, err := os.Create("__.h")
	checkErr(err)
	defer wh.Close()
	wc, err := os.Create("_.c")
	checkErr(err)
	defer wc.Close()

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

	exportData := importer.ExportData(pkg)
	_, err = wp.Write(exportData)
	checkErr(err)

	gtc := gotoc.NewGTC(pkg, ti)
	checkErr(gtc.Translate(wh, wc, flist))

	/*
		bt := determineBuildTools()
		if bt == nil {
			die("can't determine build tools for " +
				buildCtx.GOOS + "_" + buildCtx.GOARCH)
		}

		checkErr(bt.Compile("_.c"))
	*/
}

func main() {
	path, err := os.Getwd()
	checkErr(err)
	compile(path)
}
