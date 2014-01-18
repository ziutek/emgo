package main

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"go/ast"
	"go/token"
	"io"
)

type ipkg struct {
	name     string
	exported bool
}

type CC struct {
	wg, // Go exported declarations
	wc, // C implementation
	ws, // C local declarations
	wh bytes.Buffer // C exported declarations

	pkg     *types.Package
	imports map[string]*ipkg // imports for whole package

	ti   *types.Info
	fset *token.FileSet

	il int
}

func Compile(og, oh, oc io.Writer, pkg *types.Package, fset *token.FileSet, files []*ast.File, ti *types.Info) error {
	cc := &CC{
		pkg:     pkg,
		imports: make(map[string]*ipkg),
		ti:      ti,
		fset:    fset,
	}

	// Package imports

	for _, f := range files {
		for _, i := range f.Imports {
			path := i.Path.Value
			path = path[1 : len(path)-1]
			p := cc.imports[path]
			if p == nil {
				p = new(ipkg)
				cc.imports[path] = p
			}
			if i.Name != nil {
				// Local name is allways unambiguous
				p.name = i.Name.Name
			}
		}
	}

	// Build package
	
	cc.ws.WriteString("\n// Local declarations\n\n")
	cc.wc.WriteString("\n// Implementation\n\n")

	for _, f := range files {
		for _, d := range f.Decls {
			cc.Decl(d)
		}
	}

	buf := new(bytes.Buffer)

	// Write Go header file that contains all exported symbols

	buf.WriteString("// ")
	buf.WriteString("Tu ma być dokumentacja pakietu!")
	buf.WriteByte('\n')

	buf.WriteString("package " + pkg.Name() + "\n")

	buf.WriteString("import (\n")
	for path, ipkg := range cc.imports {
		if !ipkg.exported {
			continue
		}
		if ipkg.name != "" {
			buf.WriteString(ipkg.name)
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

	up := upath(pkg.Path())
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
		if ipkg.exported {
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
