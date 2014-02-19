package main

import (
	"code.google.com/p/go.tools/go/importer"
	"code.google.com/p/go.tools/go/types"
	//"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

type Importer struct {
	tc types.Config // used by importSrc
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

	//fmt.Printf("\nimport \"%s\"\n%+v\n", path, bp)
	
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

	bp, err := buildCtx.Import(path, srcDir, build.FindOnly|build.AllowBinary)
	if err != nil {
		return nil, err
	}

	//fmt.Printf("\nimport \"%s\"\n%+v\n", path, bp)

	if pkg := imports[path]; pkg != nil && pkg.Complete() {
		return pkg, nil
	}

	buf, err := arReadFile(bp.PkgObj, "__.EXPORTS")
	if err != nil {
		return nil, err
	}
	return importer.ImportData(imports, buf)
}
