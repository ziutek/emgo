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

type CC struct {
	fset *token.FileSet
	pkg  *types.Package
	ti   *types.Info

	imports map[string]*IPkg // imports for whole package
	il      int

	// Result of translation

	wg, // Go exported declarations
	wc, // C implementation
	ws, // C local declarations
	wh bytes.Buffer // C exported declarations
}

func NewCC(fset *token.FileSet, pkg *types.Package, ti *types.Info,
	imports map[string]*IPkg) *CC {

	return &CC{fset: fset, pkg: pkg, ti: ti, imports: imports}
}

// Result returns current result of translation.
// Some kind of result not returned are marks in imports map.
func (cc *CC) Result() (g, c, s, h []byte) {
	return cc.wg.Bytes(), cc.wc.Bytes(), cc.ws.Bytes(), cc.wh.Bytes()
}

// Resets
func (cc *CC) Reset() {
	// Reset buffers
	cc.wg.Reset()
	cc.wc.Reset()
	cc.ws.Reset()
	cc.wh.Reset()
	for _, p := range cc.imports {
		p.Exported = false
	}
	cc.il = 0
}

func (cc *CC) File(f *ast.File) {
	for _, d := range f.Decls {
		cc.Decl(d)
	}
}

// Complie translates files to complete set of C/Go source. It resets cc
// before translation. It writes results of translation to:
//	og - Go "header", contains exported declarations
//	oh - C header, contains exported declarations translated to C
//	oc - C source
func (cc *CC) Compile(og, oh, oc io.Writer, files []*ast.File) error {
	cc.Reset()

	cc.ws.WriteString("\n// Local declarations\n\n")
	cc.wc.WriteString("\n// Implementation\n\n")

	// Translate files - result in wg, wc, ws, wh + marks in imports
	for _, f := range files {
		cc.File(f)
	}

	buf := new(bytes.Buffer)

	// Write Go header file - it contains all exported declarations

	buf.WriteString("// ")
	buf.WriteString("Tu ma być dokumentacja pakietu!")
	buf.WriteByte('\n')

	buf.WriteString("package " + cc.pkg.Name() + "\n")

	buf.WriteString("import (\n")
	for path, ipkg := range cc.imports {
		if !ipkg.Exported {
			continue
		}
		if ipkg.Name != "" {
			buf.WriteString(ipkg.Name)
			buf.WriteByte(' ')
		}
		buf.WriteByte('"')
		buf.WriteString(path)
		buf.WriteString("\"\n")
	}
	buf.WriteString(")\n")

	if _, err := buf.WriteTo(og); err != nil {
		return err
	}
	if _, err := cc.wg.WriteTo(og); err != nil {
		return err
	}

	// Write C .h and .c

	buf.WriteString("#include \"types.h\"\n")
	buf.WriteString("#include \"_.h\"\n\n")

	if _, err := buf.WriteTo(oc); err != nil {
		return err
	}

	up := upath(cc.pkg.Path())
	buf.WriteString("#ifndef " + up + "\n")
	buf.WriteString("#define " + up + "\n\n")

	if _, err := buf.WriteTo(oh); err != nil {
		return err
	}

	for path, ipkg := range cc.imports {
		if path == "unsafe" {
			continue
		}

		buf.WriteString("#include \"")
		buf.WriteString(path)
		buf.WriteString("/_.h\"\n")

		w := oc
		if ipkg.Exported {
			w = oh
		}

		if _, err := buf.WriteTo(w); err != nil {
			return err
		}
	}

	cc.wh.WriteString("\n#endif\n")

	if _, err := cc.wh.WriteTo(oh); err != nil {
		return err
	}

	if _, err := cc.ws.WriteTo(oc); err != nil {
		return err
	}
	if _, err := cc.wc.WriteTo(oc); err != nil {
		return err
	}

	return nil
}

func (cc *CC) indent(w *bytes.Buffer) {
	for i := 0; i < cc.il; i++ {
		w.WriteByte('\t')
	}
}

func (cc *CC) isImported(o types.Object) bool {
	if p := o.Pkg(); p != nil {
		return p != cc.pkg
	}
	return false
}

func (cc *CC) isLocal(o types.Object) bool {
	if cc.isImported(o) {
		return false
	}
	return o.Parent() != cc.pkg.Scope()
}

func (cc *CC) isGlobal(o types.Object) bool {
	return o.Parent() == cc.pkg.Scope()
}
