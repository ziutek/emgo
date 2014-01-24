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
}

func NewGTC(pkg *types.Package, ti *types.Info) *GTC {
	cc := &GTC{
		pkg: pkg,
		ti:  ti,
	}
	return cc
}

func (cc *GTC) File(f *ast.File) (cdds []*CDD) {
	for _, d := range f.Decls {
		// TODO: concurrently?
		cdds = append(cdds, cc.Decl(d, 0)...)
	}
	return
}

func (gtc *GTC) exportDecl(cddm map[types.Object]*CDD, o types.Object) {
	cdd := cddm[o]
	if cdd.Export {
		return
	}
	cdd.Export = true
	for o := range cdd.DeclUses {
		if gtc.isImported(o) {
			continue
		}
		gtc.exportDecl(cddm, o)
	}
}

type imports map[*types.Package]bool

func (i imports) add(pkg *types.Package, export bool) {
	if e, ok := i[pkg]; ok {
		if !e && export {
			i[pkg] = true
		}
	} else {
		i[pkg] = export
	}
}

// Translate translates files to complete set of C/Go source. It resets cc
// before translation. It writes results of translation to:
//	wh - C header, contains exported declarations translated to C
//	wc - C source
func (gtc *GTC) Translate(wh, wc io.Writer, files []*ast.File) error {
	var cdds []*CDD

	for _, f := range files {
		// TODO: do this concurrently
		cdds = append(cdds, gtc.File(f)...)
	}

	cddm := make(map[types.Object]*CDD)
	for _, cdd := range cdds {
		cddm[cdd.Origin] = cdd
		if cdd.Typ == FuncDecl {
			cdd.DetermineInline()
		}
	}

	// Find unexported decls refferenced by inlined
	// code and mark them for export
	for _, cdd := range cdds {
		if cdd.Inline {
			for o := range cdd.BodyUses {
				if gtc.isImported(o) {
					continue
				}
				gtc.exportDecl(cddm, o)
			}
		}
	}

	// Find all external packages refferenced by exported code
	imp := make(imports)
	for _, cdd := range cdds {
		for o := range cdd.DeclUses {
			if gtc.isImported(o) {
				imp.add(o.Pkg(), cdd.Export)
			}
		}
		for o := range cdd.BodyUses {
			if gtc.isImported(o) {
				imp.add(o.Pkg(), cdd.Export && cdd.Inline)
			}
		}
	}
	
	buf := new(bytes.Buffer)

	buf.WriteString("#include \"types.h\"\n")
	buf.WriteString("#include \"__.h\"\n\n")

	if _, err := buf.WriteTo(wc); err != nil {
		return err
	}

	up := upath(gtc.pkg.Path())

	buf.WriteString("#ifndef " + up + "\n")
	buf.WriteString("#define " + up + "\n\n")

	if _, err := buf.WriteTo(wh); err != nil {
		return err
	}

	for pkg, export := range imp {
		path := pkg.Path()
		if path == "unsafe" {
			continue
		}

		buf.WriteString("#include \"")
		buf.WriteString(path)
		buf.WriteString("/__.h\"\n")

		w := wc
		if export {
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

	buf.WriteString("\nvoid " + up + "_init() {\n")
	for _, cdd := range cdds {
		buf.Write(cdd.Init)
	}
	buf.WriteString("}\n")
	if _, err := buf.WriteTo(wc); err != nil {
		return err
	}
	
	return nil
}

func (gtc *GTC) isImported(o types.Object) bool {
	return o.Pkg() != gtc.pkg
}

func (gtc *GTC) isLocal(o types.Object) bool {
	return !gtc.isImported(o) && o.Parent() != gtc.pkg.Scope()
}


func (gtc *GTC) isGlobal(o types.Object) bool {
	return !gtc.isImported(o) && o.Parent() == gtc.pkg.Scope()
}
