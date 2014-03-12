package main

import (
	"bufio"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/go.tools/go/importer"
	"code.google.com/p/go.tools/go/types"

	"github.com/ziutek/emgo/gotoc"
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
	// Parse

	srcDir := ""
	if build.IsLocalImport(ppath) {
		var err error
		if srcDir, err = os.Getwd(); err != nil {
			return err
		}
	}

	bp, err := buildCtx.Import(ppath, srcDir, 0)
	if err != nil {
		return err
	}

	flist := make([]*ast.File, 0, len(bp.GoFiles)+1)
	fset := token.NewFileSet()

	for _, fname := range bp.GoFiles {
		fname = filepath.Join(bp.Dir, fname)
		f, err := parser.ParseFile(fset, fname, nil, 0)
		if err != nil {
			return err
		}
		flist = append(flist, f)
	}

	ppath = bp.ImportPath
	elf := ""
	if bp.Name == "main" {
		elf = filepath.Join(bp.Dir, "main.elf")
		ppath = "main"

		f, err := parser.ParseFile(
			fset, "_importruntime.go",
			"package main\nimport _ \"runtime\"\n",
			0,
		)
		if err != nil {
			return err
		}
		flist = append(flist, f)
	}

	// Type check

	tc := &types.Config{
		Import: NewImporter().Import,
		Sizes: &types.StdSizes{4, 4}, // BUG: set sizes based on EGARCH
	}
	ti := &types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}

	pkg, err := tc.Check(ppath, fset, flist, ti)
	if err != nil {
		return err
	}

	// Translate to C

	work := filepath.Join(tmpDir, ppath)
	if err = os.MkdirAll(work, 0700); err != nil {
		return err
	}

	cpath := filepath.Join(bp.Dir, "_.c")
	wc, err := os.Create(cpath)
	if err != nil {
		return err
	}
	defer wc.Close()

	var (
		hpath string
		objs  []string
	)
	csfiles := append(bp.CFiles, bp.SFiles...)
	csfiles = append(csfiles, "_.c")

	if ppath == "main" {
		hpath = filepath.Join(bp.Dir, "_.h")
		objs = make([]string, 0, len(csfiles))
	} else {
		hpath = filepath.Join(bp.PkgRoot, buildCtx.GOOS+"_"+buildCtx.GOARCH, ppath+".h")
		expath := filepath.Join(work, "__.EXPORTS")
		impath := filepath.Join(work, "__.IMPORTS")
		objs = make([]string, 0, len(csfiles)+3)
		objs = append(objs, expath, impath, hpath)

		err = os.MkdirAll(filepath.Dir(hpath), 0755)
		if err != nil && !os.IsExist(err) {
			return err
		}
		wp, err := os.Create(expath)
		if err != nil {
			return err
		}
		edata := importer.ExportData(pkg)
		_, err = wp.Write(edata)
		if err != nil {
			return err
		}
		wp.Close()
		wp, err = os.Create(impath)
		if err != nil {
			return err
		}
		for _, p := range pkg.Imports() {
			if _, err := io.WriteString(wp, p.Path()+"\n"); err != nil {
				return err
			}
		}
		wp.Close()
	}

	wh, err := os.Create(hpath)
	if err != nil {
		return err
	}
	defer wh.Close()

	up := strings.Replace(ppath, "/", "_", -1)
	_, err = io.WriteString(wh, "#ifndef "+up+"\n#define "+up+"\n\n")
	if err != nil {
		return err
	}

	gtc := gotoc.NewGTC(pkg, ti)
	gtc.SetInlineThres(12)
	if err = gtc.Translate(wh, wc, flist); err != nil {
		return err
	}

	for _, h := range bp.HFiles {
		f, err := os.Open(filepath.Join(bp.Dir, h))
		if err != nil {
			return err
		}
		if _, err = bufio.NewReader(f).WriteTo(wh); err != nil {
			return err
		}
		if _, err = wh.Write([]byte{'\n'}); err != nil {
			return err
		}
	}

	if _, err = io.WriteString(wh, "#endif\n"); err != nil {
		return err
	}

	// Build (package or binary)

	bt, err := NewBuildTools(&buildCtx)
	if err != nil {
		return err
	}

	//bt.Log = os.Stdout

	for _, c := range csfiles {
		// TODO: avoid recompile up to date objects
		o := filepath.Join(work, c[:len(c)-1]+"o")
		c = filepath.Join(bp.Dir, c)
		if err = bt.Compile(o, c); err != nil {
			return err
		}
		objs = append(objs, o)
	}

	if ppath != "main" {
		return bt.Archive(hpath[:len(hpath)-1]+"a", objs...)
	}

	imports := make([]string, len(pkg.Imports()))
	for i, p := range pkg.Imports() {
		imports[i] = p.Path()
	}
	return bt.Link(elf, imports, objs...)
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
	defer os.RemoveAll(tmpDir)

	if err = compile(path); err != nil {
		logErr(err)
		return
	}
}
