package main

import (
	"bufio"
	"errors"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ziutek/emgo/egc/importer"
	"github.com/ziutek/emgo/gotoc"
)

func egc(ppath string) error {
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

var uptodate = make(map[string]struct{})

var (
	cortexmSizes = &gotoc.StdSizes{4, 8}

	sizesMap = map[string]types.Sizes{
		"cortexm0":  cortexmSizes,
		"cortexm3":  cortexmSizes,
		"cortexm4":  cortexmSizes,
		"cortexm4f": cortexmSizes,
		"cortexm7f": cortexmSizes,
		"cortexm7d": cortexmSizes,

		"amd64": &gotoc.StdSizes{8, 8},
	}
)

func compile(bp *build.Package) error {
	if ok, err := checkPkg(bp); err != nil {
		return err
	} else if ok {
		return nil
	}
	if verbosity > 0 {
		defer fmt.Println(bp.ImportPath)
	}

	// Parse

	flist := make([]*ast.File, 0, len(bp.GoFiles)+1)
	fset := token.NewFileSet()

	for _, fname := range bp.GoFiles {
		fname = filepath.Join(bp.Dir, fname)
		f, err := parser.ParseFile(fset, fname, nil, parser.ParseComments)
		if err != nil {
			return err
		}
		flist = append(flist, f)
	}

	var iimp string

	ppath := bp.ImportPath
	if bp.Name == "main" {
		ppath = "main"
		iimp = `_ "runtime";_ "internal"`
	} else if bp.ImportPath != "internal" {
		iimp = `_ "internal"`
	}

	f, err := parser.ParseFile(
		fset,
		"_iimports.go",
		"package "+bp.Name+";import("+iimp+")",
		0,
	)
	if err != nil {
		return err
	}
	flist = append(flist, f)

	// Type check

	var tcerrors []string

	tc := &types.Config{
		Importer: NewImporter(),
		Sizes:    sizesMap[buildCtx.GOARCH],
		Error:    func(err error) { tcerrors = append(tcerrors, err.Error()) },
	}

	ti := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}

	pkg, err := tc.Check(ppath, fset, flist, ti)
	if err != nil {
		return errors.New(strings.Join(tcerrors, "\n"))
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

	oat := buildCtx.GOOS + "_" + buildCtx.GOARCH
	if buildCtx.InstallSuffix != "" {
		oat += "_" + buildCtx.InstallSuffix
	}

	if ppath == "main" {
		hpath = filepath.Join(bp.Dir, "__"+oat+".h")
	} else {
		hpath = filepath.Join(bp.PkgRoot, oat, ppath+".h")
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

	cout := "__" + oat + ".c"

	wc, err := os.Create(filepath.Join(bp.Dir, cout))
	if err != nil {
		return err
	}
	defer wc.Close()

	if ppath == "main" {
		h := filepath.Base(hpath)
		_, err = io.WriteString(wc, "#include \""+h+"\"\n\n")
	} else {
		up := gotoc.Upath(ppath)
		_, err = io.WriteString(wh, "#ifndef $"+up+"$\n#define $"+up+"$\n\n")
	}
	if err != nil {
		return err
	}

	gtc := gotoc.NewGTC(fset, pkg, ti, tc.Sizes)
	gtc.SetNoinlineThres(7)
	gtc.SetBoundsCheck(!disableBC)
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

	if ppath != "main" {
		if _, err = io.WriteString(wh, "\n#endif\n"); err != nil {
			return err
		}
	}

	var csfiles = []string{cout}

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

	if verbosity > 1 {
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
		if err := bt.Archive(hpath[:len(hpath)-1]+"a", objs...); err != nil {
			return err
		}
		now := time.Now()
		return os.Chtimes(hpath, now, now)
	}

	imports := make([]string, len(pkg.Imports()))
	for i, p := range pkg.Imports() {
		imports[i] = p.Path()
	}
	return bt.Link(
		filepath.Join(bp.Dir, buildCtx.GOARCH+".elf"),
		imports, objs...,
	)
}

// checkPkg returns true if the package and its dependences are up to date
// (doesn't need to be (re)compiled). It always returns false for main package.
func checkPkg(bp *build.Package) (bool, error) {
	if bp.Name == "main" {
		return false, nil
	}
	if _, ok := uptodate[bp.ImportPath]; ok {
		return true, nil
	}
	pkgobj := bp.PkgObj
	oi, err := os.Stat(pkgobj)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	if len(bp.GoFiles) == 0 {
		uptodate[bp.ImportPath] = struct{}{}
		return true, nil
	}
	src := append(bp.GoFiles, bp.CFiles...)
	src = append(src, bp.HFiles...)
	src = append(src, bp.SFiles...)
	dir := filepath.Join(bp.SrcRoot, bp.ImportPath)
	for _, s := range src {
		si, err := os.Stat(filepath.Join(dir, s))
		if err != nil {
			return false, err
		}
		if !oi.ModTime().After(si.ModTime()) {
			return false, nil
		}
	}
	h := bp.PkgObj[:len(bp.PkgObj)-1] + "h"
	ok, err := checkH(h, oi.ModTime())
	if err != nil {
		return false, err
	}
	if !ok {
		data, err := arReadFile(bp.PkgObj, filepath.Base(h))
		if err != nil {
			return false, err
		}
		if err = ioutil.WriteFile(h, data, 0644); err != nil {
			return false, err
		}
	}
	if bp.ImportPath == "internal" {
		if len(bp.Imports) > 1 || len(bp.Imports) == 1 && bp.Imports[0] != "unsafe" {
			return false, errors.New("internal can't import other packages")
		}
	} else {
		imports := addPkg(bp.Imports, "internal")
		for _, imp := range imports {
			if imp == "unsafe" {
				continue
			}
			ibp, err := buildCtx.Import(imp, dir, build.AllowBinary)
			if err != nil {
				return false, err
			}
			if ok, err := checkPkg(ibp); err != nil {
				return false, err
			} else if !ok {
				return false, nil
			} else {
				pi, err := os.Stat(ibp.PkgObj)
				if err != nil {
					return false, err
				}
				if !oi.ModTime().After(pi.ModTime()) {
					return false, nil
				}
			}
			uptodate[imp] = struct{}{}
		}
	}
	uptodate[bp.ImportPath] = struct{}{}
	return true, nil
}

func checkH(h string, omt time.Time) (bool, error) {
	hi, err := os.Stat(h)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return !omt.Before(hi.ModTime()), nil
}

func addPkg(imports []string, pkg string) []string {
	for _, s := range imports {
		if s == pkg {
			return imports
		}
	}
	return append(imports, pkg)
}
