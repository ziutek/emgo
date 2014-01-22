package gotoc

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"go/ast"
	"go/token"
	"io"
)

type IPkg struct {
	Name     string // imported package name
	Exported bool   // is this name exported
}

func MakeImports(files []*ast.File) map[string]*IPkg {
	imports := make(map[string]*IPkg)
	for _, f := range files {
		for _, i := range f.Imports {
			path := i.Path.Value
			path = path[1 : len(path)-1]
			p := imports[path]
			if p == nil {
				p = new(IPkg)
				imports[path] = p
			}
			if i.Name != nil {
				// Local name is allways unambiguous
				p.Name = i.Name.Name
			}
		}
	}
	return imports
}

// GTC stores state of Go to C translator.
type GTC struct {
	fset *token.FileSet
	pkg  *types.Package
	ti   *types.Info

	imports map[string]*IPkg // imports for whole package
}

func NewGTC(fset *token.FileSet, pkg *types.Package, ti *types.Info, imports map[string]*IPkg) *GTC {
	cc := &GTC{
		fset:    fset,
		pkg:     pkg,
		ti:      ti,
		imports: imports,
	}
	return cc
}

// Resets
func (cc *GTC) Reset() {
	// Reset buffers
	for _, p := range cc.imports {
		p.Exported = false
	}
}

func (cc *GTC) File(f *ast.File) (cdds []*CDD) {
	for _, d := range f.Decls {
		// TODO: concurrently?
		cdds = append(cdds, cc.Decl(d, 0)...)
	}
	return
}

// Translate translates files to complete set of C/Go source. It resets cc
// before translation. It writes results of translation to:
//	wh - C header, contains exported declarations translated to C
//	wc - C source
func (cc *GTC) Translate(wh, wc io.Writer, files []*ast.File) error {
	cc.Reset()

	var cdds []*CDD

	for _, f := range files {
		// TODO: do this concurrently
		cdds = append(cdds, cc.File(f)...)
	}

	export := make(map[types.Object]bool)
	for _, cdd := range cdds {
		export[cdd.Origin] = cdd.Export
	}

	buf := new(bytes.Buffer)

	buf.WriteString("#include \"types.h\"\n")
	buf.WriteString("#include \"__.h\"\n\n")

	if _, err := buf.WriteTo(wc); err != nil {
		return err
	}

	up := upath(cc.pkg.Path())

	buf.WriteString("#ifndef " + up + "\n")
	buf.WriteString("#define " + up + "\n\n")

	if _, err := buf.WriteTo(wh); err != nil {
		return err
	}

	for path, ipkg := range cc.imports {
		if path == "unsafe" {
			continue
		}

		buf.WriteString("#include \"")
		buf.WriteString(path)
		buf.WriteString("/__.h\"\n")

		w := wc
		if ipkg.Exported {
			w = wh
		}

		if _, err := buf.WriteTo(w); err != nil {
			return err
		}
	}
	
	for _, cdd := range cdds {
		if err := cdd.WriteDecl(wh, wc); err != nil {
			return err
		}
	}
	for _, cdd := range cdds {
		if err := cdd.WriteDef(wh, wc); err != nil {
			return err
		}
	}

	if _, err := io.WriteString(wh, "\n#endif\n"); err != nil {
		return err
	}

	return nil
}

func (cc *GTC) isImported(o types.Object) bool {
	return o.Pkg() != cc.pkg
}

func (cc *GTC) isLocal(o types.Object) bool {
	if cc.isImported(o) {
		return false
	}
	return o.Parent() != cc.pkg.Scope()
}

func (cc *GTC) isGlobal(o types.Object) bool {
	return o.Parent() == cc.pkg.Scope()
}
