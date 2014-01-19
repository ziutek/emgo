package main

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"io"
	"os"
	"path/filepath"
	"strconv"
)

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
	//return imp.importSrc(imports, path)
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

	buf, err := readGoExports(bp.PkgObj)
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

var (
	arHeader = []byte("!<arch>\n")
	goexName = []byte("_.go")
)

func readGoExports(aname string) ([]byte, error) {
	a, err := os.Open(aname)
	if err != nil {
		return nil, err
	}
	defer a.Close()

	blen := 16 + 12 + 6 + 6 + 8 + 10 + 2
	buf := make([]byte, blen)

	n, err := io.ReadFull(a, buf[:len(arHeader)])

	if err != nil && err != io.ErrUnexpectedEOF {
		return nil, err
	}
	if err == io.ErrUnexpectedEOF || !bytes.Equal(buf[:n], arHeader) {
		err = fmt.Errorf(
			"%s is too short or doesn't begin from ar header",
			aname,
		)
		return nil, err
	}

	for {
		n, err = io.ReadFull(a, buf)
		if err != nil {
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				err = fmt.Errorf(
					"archive %s doesn't contain %s file",
					aname, goexName,
				)
			}
			return nil, err
		}
		if buf[blen-2] != '`' || buf[blen-1] != '\n' {
			err = fmt.Errorf("bad file header magic in %s", aname)
			return nil, err
		}
		fname := bytes.TrimRight(buf[:16], " ")
		if last := len(fname) - 1; fname[last] == '/' {
			// GNU ar
			fname = fname[:last]
		}
		flen, err := strconv.ParseUint(
			string(bytes.TrimRight(buf[48:58], " ")),
			10, 64,
		)
		if err != nil {
			err = fmt.Errorf(
				"bad file size for %s in %s: %s",
				fname, aname, err,
			)
			return nil, err
		}

		if bytes.Equal(fname, goexName) {
			buf = make([]byte, flen)
			if _, err = io.ReadFull(a, buf); err != nil {
				err = fmt.Errorf(
					"can't read %s file from %s: %s",
					fname, aname, err,
				)
				return nil, err
			}
			return buf, nil
		}

		if flen&1 != 0 {
			flen++
		}
		if _, err = a.Seek(int64(flen), 1); err != nil {
			return nil, err
		}
	}
}
