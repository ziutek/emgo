package gotoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"io"
	"os"
	"strconv"

	"golang.org/x/tools/go/exact"
	"golang.org/x/tools/go/types"
)

// GTC stores information from type checker need for translation.
type GTC struct {
	fset        *token.FileSet
	pkg         *types.Package
	ti          *types.Info
	inlineThres int
	nextInt     chan int
	siz         types.Sizes
	sizPtr      int64
	sizIval     int64

	msetc types.MethodSetCache

	// TODO: safe concurent acces is need
	tuples  map[string]types.Object
	arrays  map[string]types.Object
	itables map[string]types.Object
	tinfos  map[string]types.Object
	minfos  map[string]types.Object
}

func NewGTC(fset *token.FileSet, pkg *types.Package, ti *types.Info, siz types.Sizes) *GTC {
	c := make(chan int, 1)
	go nextIntGen(c)
	return &GTC{
		fset:    fset,
		pkg:     pkg,
		ti:      ti,
		nextInt: c,
		tuples:  make(map[string]types.Object),
		arrays:  make(map[string]types.Object),
		itables: make(map[string]types.Object),
		tinfos:  make(map[string]types.Object),
		minfos:  make(map[string]types.Object),
		siz:     siz,
		sizPtr:  siz.Sizeof(types.NewPointer(types.NewStruct(nil, nil))),
		sizIval: siz.Sizeof(types.Typ[types.Complex128]),
	}
}

func (cc *GTC) SetInlineThres(thres int) {
	cc.inlineThres = thres
}

func (gtc *GTC) File(f *ast.File) (cdds []*CDD) {
	for _, d := range f.Decls {
		// TODO: concurrently?
		for _, cdd := range gtc.Decl(d, 0) {
			cdds = append(cdds, cdd.AllCDDS()...)
		}
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
		if cddm[o] == nil {
			cdd.exit(o.Pos(), "cddm[o] == nil")
		}
		gtc.export(cddm, cddm[o])
	}
	if cdd.Typ != FuncDecl || !cdd.Inline {
		return
	}
	for o := range cdd.FuncBodyUses {
		if gtc.isImported(o) {
			continue
		}
		gtc.export(cddm, cddm[o])
	}
}

type imports map[*types.Package]bool

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
		/*fmt.Printf(
			"Origin: %s <%d>\n DeclUses: %+v\n BodyUses: %+v\n",
			cdd.Origin, cdd.Typ, cdd.DeclUses, cdd.BodyUses,
		)*/
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
		if o.Exported() || (o.Pkg().Name() == "main" && o.Name() == "main") {
			gtc.export(cddm, cdd)
		}
	}

	// Classify all imported packages.
	imp := make(imports)
	for _, p := range gtc.pkg.Imports() {
		imp[p] = false
	}

	for _, cdd := range cdds {
		/*if cdd.Typ == ImportDecl {
			// Package imported as _
			//fmt.Println(cdd.Origin.Pkg())
			imp.add(cdd.Origin.Pkg(), false)
			continue
		}*/
		for o := range cdd.DeclUses {
			if gtc.isImported(o) {
				if cdd.Export {
					if o.Pkg() == nil {
						fmt.Printf("nil pkg: %#v\n", o)
					}
					imp[o.Pkg()] = true
				}
			}
		}
		for o := range cdd.FuncBodyUses {
			if gtc.isImported(o) && cdd.Export && cdd.Inline {
				imp[o.Pkg()] = true
			}
		}
	}
	_, err := io.WriteString(
		wc,
		"#include <internal/types.h>\n#include <builtin.h>\n",
	)
	if err != nil {
		return err
	}
	for pkg, export := range imp {
		path := pkg.Path()
		if path == "unsafe" || path == "builtin" {
			continue
		}
		w := wc
		if export {
			w = wh
		}
		if _, err = io.WriteString(w, "#include <"+path+".h>\n"); err != nil {
			return err
		}
	}
	pkgName := gtc.pkg.Name()
	if pkgName == "main" {
		_, err = io.WriteString(wc, "\n#include \"_.h\"\n")
	} else if gtc.pkg.Path() != "builtin" {
		_, err = io.WriteString(wc, "\n#include <"+gtc.pkg.Path()+".h>\n")
	}
	if err != nil {
		return err
	}

	var tcs, vcs, fcs, ccs []*CDD
	cddm = make(map[types.Object]*CDD)

	for _, cdd := range cdds {
		switch cdd.Typ {
		case TypeDecl:
			cddm[cdd.Origin] = cdd

		case VarDecl:
			vcs = append(vcs, cdd)

		case FuncDecl:
			fcs = append(fcs, cdd)

		case ConstDecl:
			ccs = append(ccs, cdd)
		}
	}

	/*for _, v := range cddm {
		fmt.Printf("** origin: %p\n    uses:", v.Origin)
		for p, ok := range v.DeclUses {
			fmt.Printf(" %p:%t", p, ok)
		}
		fmt.Printf("\ndecl:\n%s\n", v.Decl)
		fmt.Printf("def:\n%s\n", v.Def)
		fmt.Printf("init:\n%s\n", v.Init)
	}*/

	tcs = dfs(cddm)

	if err := write("// type decl\n", wh, wc); err != nil {
		return err
	}
	for _, cdd := range tcs {
		if err := cdd.WriteDecl(wh, wc); err != nil {
			return err
		}
	}
	if err := write("// var  decl\n", wh, wc); err != nil {
		return err
	}
	cddm = make(map[types.Object]*CDD)
	for _, cdd := range vcs {
		cddm[cdd.Origin] = cdd
		if err := cdd.WriteDecl(wh, wc); err != nil {
			return err
		}
	}
	if err := write("// func decl\n", wh, wc); err != nil {
		return err
	}
	for _, cdd := range fcs {
		if err := cdd.WriteDecl(wh, wc); err != nil {
			return err
		}
	}
	if err := write("// const decl\n", wh, wc); err != nil {
		return err
	}
	for _, cdd := range ccs {
		if err := cdd.WriteDecl(wh, wc); err != nil {
			return err
		}
	}

	if err := write("// type def\n", wh, wc); err != nil {
		return err
	}
	for _, cdd := range tcs {
		if err := cdd.WriteDef(wh, wc); err != nil {
			return err
		}
	}
	if err := write("// var  def\n", wh, wc); err != nil {
		return err
	}
	for _, cdd := range vcs {
		if err := cdd.WriteDef(wh, wc); err != nil {
			return err
		}
	}
	if err := write("// func def\n", wh, wc); err != nil {
		return err
	}
	for _, cdd := range fcs {
		if err := cdd.WriteDef(wh, wc); err != nil {
			return err
		}
	}

	if err := write("// init\n", wh, wc); err != nil {
		return err
	}
	up := upath(gtc.pkg.Path())
	if _, err = io.WriteString(wh, "void "+up+"$init();\n"); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	buf.WriteString("void " + up + "$init() {\n")
	if pkgName != "main" {
		buf.WriteString("\tstatic bool called = false;\n")
		buf.WriteString("\tif (called) {\n\t\treturn;\n\t}\n\tcalled = true;\n")
	}
	n := buf.Len()

	for i := range imp {
		buf.WriteString("\t" + upath(i.Path()) + "$init();\n")
	}

	for _, i := range gtc.ti.InitOrder {
		for _, l := range i.Lhs {
			if cdd := cddm[l]; cdd != nil {
				cdd.writeInits(buf)
			}
			/*for cdd != nil {
				buf.Write(cdd.Init)
				cdd = cdd.InitNext
			}*/
		}
	}
	for _, cdd := range fcs {
		buf.Write(cdd.Init)
	}
	if buf.Len() == n {
		// no imports, no inits
		_, err := io.WriteString(wh, "#define "+up+"$init()\n")
		return err
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

func nextIntGen(c chan<- int) {
	i := 0
	for {
		c <- i
		i++
	}
}

func (gtc *GTC) uniqueId() string {
	return strconv.FormatInt(int64(<-gtc.nextInt), 10)
}

func (gtc *GTC) object(ident *ast.Ident) types.Object {
	o := gtc.ti.Defs[ident]
	if o == nil {
		o = gtc.ti.Uses[ident]
	}
	return o
}

func (gtc *GTC) exprType(e ast.Expr) types.Type {
	return gtc.ti.Types[e].Type
}

func (gtc *GTC) exprValue(e ast.Expr) exact.Value {
	return gtc.ti.Types[e].Value
}

func (gtc *GTC) exit(pos token.Pos, f string, a ...interface{}) {
	fmt.Fprint(os.Stderr, gtc.fset.Position(pos), " ")
	fmt.Fprintf(os.Stderr, f+"\n", a...)
	os.Exit(1)
}

func (gtc *GTC) notImplemented(n ast.Node, tl ...types.Type) {
	fmt.Fprint(os.Stderr, gtc.fset.Position(n.Pos()))
	fmt.Fprintf(os.Stderr, "not implemented: %T\n", n)
	for _, t := range tl {
		fmt.Fprintf(os.Stderr, "	in case of: %T\n", t)
	}
	os.Exit(1)
}

func (gtc *GTC) methodSet(t types.Type) *types.MethodSet {
	return gtc.msetc.MethodSet(t)
}
