package gotoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

// GTC stores information from type checker need for translation.
type GTC struct {
	fset          *token.FileSet
	pkg           *types.Package
	ti            *types.Info
	noinlineThres int
	boundsCheck   bool
	nextInt       chan int
	siz           types.Sizes
	sizPtr        int64
	sizIval       int64

	// TODO: safe concurent acces is need
	tuples  map[string]types.Object
	arrays  map[string]types.Object
	itables map[string]types.Object
	tinfos  map[string]types.Object
	minfos  map[string]types.Object
	cmap    ast.CommentMap
	defs    map[types.Object]ast.Node
}

func NewGTC(fset *token.FileSet, pkg *types.Package, ti *types.Info, siz types.Sizes) *GTC {
	c := make(chan int, 1)
	go nextIntGen(c)
	sizIval := siz.Sizeof(&types.Slice{})
	if sizIval < 16 {
		sizIval = 16 // complex128
	}
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
		cmap:    make(ast.CommentMap),
		defs:    make(map[types.Object]ast.Node),
		siz:     siz,
		sizPtr:  siz.Sizeof(types.NewPointer(types.NewStruct(nil, nil))),
		sizIval: sizIval,
	}
}

func (cc *GTC) SetNoinlineThres(thres int) {
	cc.noinlineThres = thres
}

func (cc *GTC) SetBoundsCheck(bc bool) {
	cc.boundsCheck = bc
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
	cdd.Export = true
	for o := range cdd.DeclUses {
		if gtc.isImported(o) {
			continue
		}
		c := cddm[o]
		if c == nil {
			cdd.exit(o.Pos(), "cddm[o] == nil")
		}
		if !c.Export {
			gtc.export(cddm, c)
		}
	}
	if cdd.Inline {
		for o := range cdd.DefUses {
			if gtc.isImported(o) {
				continue
			}
			if c := cddm[o]; !c.Export {
				gtc.export(cddm, cddm[o])
			}
		}
	}
}

func (gtc *GTC) makeDefs(node ast.Node) bool {
	switch n := node.(type) {
	case *ast.FuncDecl:
		if obj := gtc.ti.Defs[n.Name]; obj != nil {
			gtc.defs[obj] = n
		}
	case *ast.TypeSpec:
		if obj := gtc.ti.Defs[n.Name]; obj != nil {
			gtc.defs[obj] = n
		}
		return false
	case *ast.ValueSpec:
		for _, name := range n.Names {
			if obj := gtc.ti.Defs[name]; obj != nil {
				gtc.defs[obj] = n
			}
		}
		return false
	}
	return true
}

type imports map[*types.Package]bool

// Translate translates files to C source.
// It writes results of translation to:
//	wh - C header, contains exported and inlined declarations translated to C,
//	wc - remaining C source.
func (gtc *GTC) Translate(wh, wc io.Writer, files []*ast.File) error {
	var cdds []*CDD

	for _, f := range files {
		for k, v := range ast.NewCommentMap(gtc.fset, f, f.Comments) {
			gtc.cmap[k] = v
		}
		ast.Inspect(f, gtc.makeDefs)
	}
	for _, f := range files {
		// TODO: do this concurrently
		cdds = append(cdds, gtc.File(f)...)
	}

	cddm := make(map[types.Object]*CDD)

	// Determine inline for any function except main.main()
	for _, cdd := range cdds {
		if cdd.Typ == ImportDecl {
			continue
		}
		o := cdd.Origin
		cddm[o] = cdd
		if cdd.Typ == FuncDecl && (o.Pkg().Name() != "main" || o.Name() != "main") {
			cdd.DetermineInline()
		}
	}

	// Export code need by exported declarations and inlined function bodies.
	for _, cdd := range cdds {
		if cdd.Export || cdd.Typ == ImportDecl {
			continue
		}
		o := cdd.Origin
		exported := o.Exported()
		if exported {
			if f, ok := o.(*types.Func); ok {
				if r := f.Type().(*types.Signature).Recv(); r != nil {
					rt := r.Type()
					if p, ok := rt.(*types.Pointer); ok {
						rt = p.Elem()
					}
					exported = rt.(*types.Named).Obj().Exported()
				}
			}
		}
		if exported || (o.Pkg().Name() == "main" && o.Name() == "main") ||
			cdd.forceExport {
			gtc.export(cddm, cdd)
		}
	}

	// Classify all imported packages.
	imp := make(imports)
	for _, p := range gtc.pkg.Imports() {
		imp[p] = false
	}

	for _, cdd := range cdds {
		if !cdd.Export {
			continue
		}
		for o := range cdd.DeclUses {
			if gtc.isImported(o) {
				if o.Pkg() == nil {
					fmt.Printf("nil pkg: %#v\n", o)
				}
				imp[o.Pkg()] = true
			}
		}
		if cdd.Inline {
			for o := range cdd.DefUses {
				if gtc.isImported(o) {
					imp[o.Pkg()] = true
				}
			}
		}
	}

	pkgmain := gtc.pkg.Name() == "main"

	w := wc
	if pkgmain {
		w = wh
	}
	_, err := io.WriteString(
		w, "#include <internal/types.h>\n#include <internal.h>\n",
	)
	if err != nil {
		return err
	}
	for pkg, export := range imp {
		path := pkg.Path()
		if path == "unsafe" || path == "internal" {
			continue
		}
		w := wc
		if export || pkgmain {
			w = wh
		}
		if _, err = io.WriteString(w, "#include <"+path+".h>\n"); err != nil {
			return err
		}
	}

	if !pkgmain && gtc.pkg.Path() != "builtin" {
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
	up := Upath(gtc.pkg.Path())
	if _, err = io.WriteString(wh, "void "+up+"$init();\n"); err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	buf.WriteString("void " + up + "$init() {\n")
	if !pkgmain {
		buf.WriteString("\tstatic bool called = false;\n")
		buf.WriteString("\tif (called) {\n\t\treturn;\n\t}\n\tcalled = true;\n")
	}
	n := buf.Len()

	for i := range imp {
		buf.WriteString("\t" + Upath(i.Path()) + "$init();\n")
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

func (gtc *GTC) exprValue(e ast.Expr) constant.Value {
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
	return types.NewMethodSet(t)
}

type pragmas []string

func (prs pragmas) Contains(s string) bool {
	for _, p := range prs {
		if p == s {
			return true
		}
	}
	return false
}

func (gtc *GTC) pragmas(nodes ...ast.Node) (prs pragmas, cattrs []string) {
	for _, n := range nodes {
		for _, cg := range gtc.cmap[n] {
			if cg == nil {
				continue
			}
			for _, c := range cg.List {
				s := strings.TrimSpace(c.Text)
				if strings.HasPrefix(s, "//") {
					s = strings.TrimLeftFunc(s[2:], unicode.IsSpace)
					switch {
					case strings.HasPrefix(s, "c:"):
						s = strings.TrimLeftFunc(s[2:], unicode.IsSpace)
						if s == "" {
							gtc.exit(n.Pos(), "empty C attribute")
						}
						cattrs = append(cattrs, s)
					case strings.HasPrefix(s, "emgo:"):
						s = strings.TrimLeftFunc(s[5:], unicode.IsSpace)
						if s == "" {
							gtc.exit(n.Pos(), "empty Emgo pragma")
						}
						prs = append(prs, s)
					}
				}
			}
		}
	}
	return
}

func (gtc *GTC) cattrs(nodes ...ast.Node) (string, bool) {
	pragmas, cattrs := gtc.pragmas(nodes...)
	var pexport bool
	for _, p := range pragmas {
		if p == "export" {
			pexport = true
			break
		}
	}
	return strings.Join(cattrs, " "), pexport
}
