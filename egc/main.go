package main

import (
	"code.google.com/p/go.tools/go/exact"
	"code.google.com/p/go.tools/go/importer"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"github.com/ziutek/emgo/gotoc"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
)

var tmpDir string

var buildCtx = build.Context{
	GOARCH:      "cortexm4f",
	GOOS:        "none",
	GOROOT:      "/home/michal/P/go/github/emgo/egroot",
	GOPATH:      "/home/michal/P/go/github/emgo/egpath",
	Compiler:    "gc",
	ReleaseTags: []string{"go1.1", "go1.2"},
	CgoEnabled:  false,
}

func compile(ppath string) error {
	var (
		srcDir string
		err    error
	)

	if build.IsLocalImport(ppath) {
		if srcDir, err = os.Getwd(); err != nil {
			return err
		}
	}

	bp, err := buildCtx.Import(ppath, srcDir, 0)
	if err != nil {
		return err
	}

	fmt.Printf("package \"%s\"\n%+v\n\n", ppath, bp)

	flist := make([]*ast.File, len(bp.GoFiles))
	fset := token.NewFileSet()

	for i, fname := range bp.GoFiles {
		fname = filepath.Join(bp.Dir, fname)
		f, err := parser.ParseFile(
			fset, fname, nil, parser.ParseComments,
		)
		if err != nil {
			return err
		}
		flist[i] = f
	}

	ppath = bp.ImportPath
	if bp.Name == "main" {
		ppath = "main"
	}

	tc := &types.Config{Import: NewImporter().Import}
	ti := &types.Info{
		Types:   make(map[ast.Expr]types.Type),
		Values:  make(map[ast.Expr]exact.Value),
		Objects: make(map[*ast.Ident]types.Object),
	}

	pkg, err := tc.Check(ppath, fset, flist, ti)
	if err != nil {
		return err
	}

	work := filepath.Join(tmpDir, ppath)
	if err = os.MkdirAll(work, 0700); err != nil {
		return err
	}

	expath := filepath.Join(work, "__.EXPORTS")
	wp, err := os.Create(expath)
	if err != nil {
		return err
	}
	defer wp.Close()

	hpath := filepath.Join(bp.PkgRoot, buildCtx.GOOS+"_"+buildCtx.GOARCH, ppath+".h")

	if err = os.MkdirAll(filepath.Dir(hpath), 0755); err != nil && !os.IsExist(err) {
		return err
	}

	wh, err := os.Create(hpath)
	if err != nil {
		return err
	}
	defer wh.Close()

	cpath := filepath.Join(bp.Dir, "_.c")
	wc, err := os.Create(cpath)
	if err != nil {
		return err
	}
	defer wc.Close()

	exportData := importer.ExportData(pkg)
	_, err = wp.Write(exportData)
	if err != nil {
		return err
	}
	gtc := gotoc.NewGTC(pkg, ti)
	if err = gtc.Translate(wh, wc, flist); err != nil {
		return err
	}

	bt, err := NewBuildTools(&buildCtx)
	if err != nil {
		return err
	}

	eoh := make([]string, 0, len(bp.CFiles)+3)
	eoh = append(eoh, expath)

	for _, c := range append(bp.CFiles, "_.c") {
		// TODO: avoid recompile up to date obj
		o := filepath.Join(work, c[:len(c)-1]+"o")
		c = filepath.Join(bp.Dir, c)
		if err = bt.Compile(o, c); err != nil {
			return err
		}
		eoh = append(eoh, o)
	}

	eoh = append(eoh, hpath)

	if err = bt.Archive(hpath[:len(hpath)-1]+"a", eoh...); err != nil {
		return err
	}

	return nil
}

func main() {
	path := "."
	if len(os.Args) >= 2 {
		path = os.Args[1]
	}

	var err error

	tmpDir, err = ioutil.TempDir("", "eg-build")
	if err != nil {
		logErr(err)
		return
	}
	//defer os.RemoveAll(tmpDir)

	if err = compile(path); err != nil {
		logErr(err)
		return
	}
}
