package main

import (
	"errors"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"

	"github.com/ziutek/emgo/egc/importer"
)

type Importer struct {
	tc      types.Config // Legacy: used by importSrc.
	imports map[string]*types.Package
}

func NewImporter() *Importer {
	imp := new(Importer)
	imp.imports = make(map[string]*types.Package)
	imp.tc.IgnoreFuncBodies = true
	imp.tc.Importer = imp
	return imp
}

func (imp *Importer) Import(path string) (*types.Package, error) {
	//return imp.importSrc(path)
	//return imp.importSrc1(path)
	return imp.importPkg(path)
}

func (imp *Importer) importPkg(path string) (*types.Package, error) {
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
	if bp.PkgObj == "" {
		return nil, errors.New("Emgo does not support local imports")
	}
	if pkg := imp.imports[path]; pkg != nil && pkg.Complete() {
		return pkg, nil
	}
	buf, err := loadExports(bp)
	if err != nil {
		return nil, err
	}
	_, pkg, err := importer.ImportData(imp.imports, buf)
	return pkg, err
}

func loadExports(bp *build.Package) ([]byte, error) {
	if err := compile(bp); err != nil {
		return nil, err
	}
	uptodate[bp.ImportPath] = struct{}{}
	return arReadFile(bp.PkgObj, "__.EXPORTS")
}

func (imp *Importer) importSrc(path string) (*types.Package, error) {
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
	if bp.PkgObj == "" {

	}
	if pkg := imp.imports[path]; pkg != nil && pkg.Complete() {
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

	imp.imports[path] = pkg
	return pkg, nil
}

func (imp *Importer) importSrc1(path string) (*types.Package, error) {
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

	if pkg := imp.imports[path]; pkg != nil && pkg.Complete() {
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

	imp.imports[path] = pkg
	return pkg, nil
}
