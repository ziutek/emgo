package main

import (
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"code.google.com/p/go.tools/go/importer"
	"code.google.com/p/go.tools/go/types"
)

type Importer struct {
	tc types.Config // Legacy: used by importSrc.
}

func NewImporter() *Importer {
	imp := new(Importer)
	imp.tc.IgnoreFuncBodies = true
	imp.tc.Import = imp.Import
	return imp
}

func (imp *Importer) Import(imports map[string]*types.Package, path string) (*types.Package, error) {
	//return imp.importSrc(imports, path)
	//return imp.importSrc1(imports, path)
	return imp.importPkg(imports, path)
}

func (imp *Importer) importSrc(imports map[string]*types.Package, path string) (*types.Package, error) {
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

func (imp *Importer) importSrc1(imports map[string]*types.Package, path string) (*types.Package, error) {
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

	bp, err := buildCtx.Import(path, srcDir, build.FindOnly|build.AllowBinary)
	if err != nil {
		return nil, err
	}

	if pkg := imports[path]; pkg != nil && pkg.Complete() {
		return pkg, nil
	}

	buf, err := arReadFile(bp.PkgObj, "_.go")
	if err != nil {
		return nil, err
	}

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", buf, 0)
	if err != nil {
		return nil, err
	}

	pkg, err := imp.tc.Check(path, fset, []*ast.File{file}, nil)
	if err != nil {
		return nil, err
	}

	imports[path] = pkg
	return pkg, nil
}

func (imp *Importer) importPkg(imports map[string]*types.Package, path string) (*types.Package, error) {
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

	bp, err := buildCtx.Import(path, srcDir, build.AllowBinary)
	if err != nil {
		return nil, err
	}

	if pkg := imports[path]; pkg != nil && pkg.Complete() {
		return pkg, nil
	}

	buf, err := loadExports(bp)
	if err != nil {
		return nil, err
	}
	return importer.ImportData(imports, buf)
}

func loadExports(bp *build.Package) ([]byte, error) {
	if ok, err := checkPkg(bp); err != nil {
		return nil, err
	} else if !ok {
		if err := compile(bp); err != nil {
			return nil, err
		}
	}
	return arReadFile(bp.PkgObj, "__.EXPORTS")
}

var builtinCTime time.Time

// checkPkg returns true if package need to be (re)compiled.
func checkPkg(bp *build.Package) (bool, error) {
	oi, err := os.Stat(bp.PkgObj)
	if err != nil {
		if os.IsNotExist(err) {
			return true, nil
		}
		return false, err
	}
	if bp.ImportPath == "builtin" {
		builtinCTime = oi.ModTime()
	} else if !oi.ModTime().After(builtinCTime) {
		return true, nil
	}
	if len(bp.GoFiles) == 0 {
		return false, nil
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
	return hi.ModTime().After(omt), nil
}
