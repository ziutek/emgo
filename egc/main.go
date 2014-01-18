package main

import (
	"code.google.com/p/go.tools/go/exact"
	"code.google.com/p/go.tools/go/types"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

const FmtGoExports = true // use false for production

var buildCtx = build.Context{
	GOARCH:      "amd64",
	GOOS:        "linux",
	GOROOT:      "/home/michal/P/go/github/emgo/egroot",
	GOPATH:      "/home/michal/P/go/github/emgo",
	Compiler:    "gc",
	ReleaseTags: []string{"go1.1", "go1.2"},
	CgoEnabled:  false,
}

type Importer struct {
	tc types.Config
}

func NewImporter() *Importer {
	imp := new(Importer)
	imp.tc.IgnoreFuncBodies = true
	imp.tc.Import = imp.Import
	return imp
}

func (imp *Importer) Import(imports map[string]*types.Package, path string) (*types.Package, error) {
	if path == "unsafe" {
		return types.Unsafe, nil
	}

	var (
		srcDir string
		err    error
	)

	if build.IsLocalImport(path) {
		srcDir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}

	bp, err := buildCtx.Import(path, srcDir, 0)
	if err != nil {
		return nil, err
	}

	if pkg := imports[path]; pkg != nil && pkg.Complete() {
		return pkg, nil
	}

	fset := token.NewFileSet()
	files := make([]*ast.File, len(bp.GoFiles))

	for i, fname := range bp.GoFiles {
		fname = filepath.Join(bp.Dir, fname)
		files[i], err = parser.ParseFile(fset, fname, nil, 0)
		if err != nil {
			return nil, err
		}
	}

	pkg, err := imp.tc.Check(path, fset, files, nil)
	if err != nil {
		return nil, err
	}

	imports[path] = pkg
	return pkg, nil
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

	checkErr(Compile(og, oh, oc, pkg, fset, flist, ti))

	og.Close()
	oh.Close()
	oc.Close()
}

func main() {
	path, err := os.Getwd()
	checkErr(err)
	compile(path)
}
