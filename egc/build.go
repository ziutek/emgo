package main

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strings"

	"code.google.com/p/go.tools/go/importer"
	"code.google.com/p/go.tools/go/types"

	"github.com/ziutek/emgo/gotoc"
)

func egc(ppath string) error {
	if err := checkBuiltin(); err != nil {
		return err
	}

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
	return compile(bp)
}

func compile(bp *build.Package) error {
	if verbose > 0 {
		defer fmt.Println(bp.ImportPath)
	}

	// Parse

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

	ppath := bp.ImportPath
	elf := ""
	if bp.Name == "main" {
		elf = filepath.Join(bp.Dir, "main.elf")
		ppath = "main"

		f, err := parser.ParseFile(
			fset, "_importbr.go",
			"package main\n"+
				"import (\n"+
				"	_ \"runtime\"\n"+
				"	_ \"builtin\"\n"+
				")\n",
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
		Sizes:  &types.StdSizes{4, 4}, // BUG: set sizes based on EGARCH
	}
	ti := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
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

	var (
		hpath string
		objs  []string
	)

	if ppath == "main" {
		hpath = filepath.Join(bp.Dir, "_.h")
	} else {
		hpath = filepath.Join(bp.PkgRoot, buildCtx.GOOS+"_"+buildCtx.GOARCH, ppath+".h")
		expath := filepath.Join(work, "__.EXPORTS")
		impath := filepath.Join(work, "__.IMPORTS")
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

	wc, err := os.Create(filepath.Join(bp.Dir, "_.c"))
	if err != nil {
		return err
	}
	defer wc.Close()

	up := strings.Replace(ppath, "/", "$", -1)
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
		if !strings.HasSuffix(h, "+.h") {
			continue
		}
		f, err := os.Open(filepath.Join(bp.Dir, h))
		if err != nil {
			return err
		}
		if _, err = io.WriteString(wh, "\n// included "+h+"\n"); err != nil {
			return err
		}
		if _, err = bufio.NewReader(f).WriteTo(wh); err != nil {
			return err
		}
	}

	if _, err = io.WriteString(wh, "\n#endif\n"); err != nil {
		return err
	}

	var csfiles = []string{"_.c"}

	for _, c := range bp.CFiles {
		if !strings.HasSuffix(c, "+.c") {
			csfiles = append(csfiles, c)
			continue
		}
		f, err := os.Open(filepath.Join(bp.Dir, c))
		if err != nil {
			return err
		}
		if _, err = io.WriteString(wc, "\n// included "+c+"\n"); err != nil {
			return err
		}
		if _, err = bufio.NewReader(f).WriteTo(wc); err != nil {
			return err
		}
	}
	csfiles = append(csfiles, bp.SFiles...)

	// Build (package or binary)

	bt, err := NewBuildTools(&buildCtx)
	if err != nil {
		return err
	}

	if verbose > 1 {
		bt.Log = os.Stdout
	}

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

func checkBuiltin() error {
	bp, err := buildCtx.Import("builtin", "", build.AllowBinary)
	if err != nil {
		return err
	}
	if nc, err := needCompile(bp); err != nil {
		return err
	} else if nc {
		return compile(bp)
	}
	return nil
}
