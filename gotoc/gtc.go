package gotoc

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"go/ast"
	"io"
)

// GTC stores information from type checker need for translation.
type GTC struct {
	pkg *types.Package
	ti  *types.Info
}

func NewGTC(pkg *types.Package, ti *types.Info) *GTC {
	return &GTC{pkg: pkg, ti: ti}
}

func (cc *GTC) File(f *ast.File) (cdds []*CDD) {
	for _, d := range f.Decls {
		// TODO: concurrently?
		cdds = append(cdds, cc.Decl(d, 0)...)
	}
	return
}

func (gtc *GTC) export(cddm map[types.Object]*CDD, cdd *CDD) {
	if cdd.Export {
		return
	}
	cdd.Export = true
	for o := range cdd.DeclUses {
		if gtc.isImported(o) {
			continue
		}
		gtc.export(cddm, cddm[o])
	}
	if !cdd.Inline {
		return
	}
	for o := range cdd.BodyUses {
		if gtc.isImported(o) {
			continue
		}
		gtc.export(cddm, cddm[o])
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

// Translate translates files to C source.
// It writes results of translation to:
//	wh - C header, contains exported and inlined declarations translated to C,
//	wc - remaining C source.
func (gtc *GTC) Translate(wh, wc io.Writer, files []*ast.File) error {
	var cdds []*CDD

	for _, f := range files {
		// TODO: do this concurrently
		cdds = append(cdds, gtc.File(f)...)
	}

	cddm := make(map[types.Object]*CDD)

	// Determine inline for any function except main.main()
	for _, cdd := range cdds {
		fmt.Printf(
			"Origin: %s <%d>\n DeclUses: %+v\n BodyUses: %+v\n",
			cdd.Origin, cdd.Typ, cdd.DeclUses, cdd.BodyUses,
		)
		if cdd.Typ == ImportDecl {
			continue
		}
		o := cdd.Origin
		cddm[o] = cdd
		if cdd.Typ == FuncDecl && (o.Pkg().Name() != "main" || o.Name() != "main") {
			cdd.DetermineInline()
		}
	}

	// Export code need by exported declarations and inlined function bodies
	for _, cdd := range cdds {
		if cdd.Typ == ImportDecl {
			continue
		}
		o := cdd.Origin
		if o.IsExported() || (o.Pkg().Name() == "main" && o.Name() == "main") {
			gtc.export(cddm, cdd)
		}
	}

	// Find all imported packages refferenced by exported code.
	imp := make(imports)
	for _, cdd := range cdds {
		if cdd.Typ == ImportDecl {
			// Package imported as _
			imp.add(cdd.Origin.Pkg(), false)
			continue
		}
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

	pkgName := gtc.pkg.Name()
	buf := new(bytes.Buffer)

	buf.WriteString("#include \"runtime/types.h\"\n")
	buf.WriteString("#include \"")
	if pkgName == "main" {
		buf.WriteByte('_')
	} else {
		buf.WriteString(gtc.pkg.Path())
	}
	buf.WriteString(".h\"\n\n")

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
		buf.WriteString(path + ".h\"\n")

		w := wc
		if export {
			w = wh
		}

		if _, err := buf.WriteTo(w); err != nil {
			return err
		}
	}

	// Write declarations (types first).
	for _, cdd := range cdds {
		if cdd.Typ != TypeDecl {
			continue
		}
		if err := cdd.WriteDecl(wh, wc); err != nil {
			return err
		}
	}
	for _, cdd := range cdds {
		if cdd.Typ == TypeDecl {
			continue
		}
		if err := cdd.WriteDecl(wh, wc); err != nil {
			return err
		}
	}

	// Write definitions (order: types, variables, functions)
	for _, cdd := range cdds {
		if cdd.Typ != TypeDecl {
			continue
		}
		if err := cdd.WriteDef(wh, wc); err != nil {
			return err
		}
	}
	for _, cdd := range cdds {
		if cdd.Typ != VarDecl {
			continue
		}
		if err := cdd.WriteDef(wh, wc); err != nil {
			return err
		}
	}
	for _, cdd := range cdds {
		if cdd.Typ != FuncDecl {
			continue
		}
		if err := cdd.WriteDef(wh, wc); err != nil {
			return err
		}
	}

	// Initialisation function.
	buf.WriteString("void " + up + "_init();\n\n")
	buf.WriteString("#endif\n")
	if _, err := buf.WriteTo(wh); err != nil {
		return err
	}
	buf.WriteString("void " + up + "_init() {\n")
	m := buf.Len()
	if pkgName != "main" {
		buf.WriteString("\tstatic bool called = false;\n")
		buf.WriteString("\tif (called) {\n\t\treturn;\n\t}\n\tcalled = true;\n")
	}
	n := buf.Len()

	for i := range imp {
		buf.WriteString("\t" + upath(i.Path()) + "_init();\n")
	}

	for _, cdd := range cdds {
		if cdd.Typ == VarDecl {
			buf.Write(cdd.Init)
		}
	}
	for _, cdd := range cdds {
		if cdd.Typ == FuncDecl {
			buf.Write(cdd.Init)
		}
	}
	if buf.Len() == n {
		// no imports, no inits
		buf.Truncate(m)
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
