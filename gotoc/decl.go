package gotoc

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"go/ast"
	"go/token"
)

func (gtc *GTC) FuncDecl(d *ast.FuncDecl, il int) (cdds []*CDD) {
	f := gtc.ti.Objects[d.Name].(*types.Func)

	cdd := gtc.newCDD(f, FuncDecl, il)
	cdds = append(cdds, cdd)
	w := new(bytes.Buffer)
	fname := cdd.NameStr(f, true)

	res := cdd.Signature(w, fname, f.Type().(*types.Signature), true)
	if res.cdd != nil {
		cdds = append(cdds, res.cdd)
	}

	init := (f.Name() == "init")

	if !init {
		cdd.copyDecl(w, ";\n")
	}

	if d.Body == nil {
		return
	}

	cdd.body = true

	w.WriteByte(' ')

	if res.hasNames {
		cdd.indent(w)
		w.WriteString("{\n")
		cdd.il++
		for _, v := range res.fields {
			cdd.indent(w)
			cdd.Type(w, v.Type())
			w.WriteByte(' ')
			cdd.Name(w, v, true)
			w.WriteString(" = {0};\n")
		}
		cdd.indent(w)
	}

	end := cdd.BlockStmt(w, d.Body, res.typ)
	w.WriteByte('\n')

	if res.hasNames {
		if end {
			cdd.il--
			cdd.indent(w)
			w.WriteString("__end:\n")
			cdd.il++

			cdd.indent(w)
			w.WriteString("return ")
			if len(res.fields) == 1 {
				cdd.Name(w, res.fields[0], true)
			} else {
				w.WriteString("(" + res.typ + ") {")
				for i, v := range res.fields {
					if i > 0 {
						w.WriteString(", ")
					}
					cdd.Name(w, v, true)
				}
				w.WriteByte('}')
			}
			w.WriteString(";\n")
		}
		cdd.il--
		w.WriteString("}\n")
	}
	cdd.copyDef(w)

	if init {
		cdd.Init = []byte("\t" + fname + "();\n")
	}
	return
}

func (gtc *GTC) GenDecl(d *ast.GenDecl, il int) (cdds []*CDD) {
	w := new(bytes.Buffer)

	switch d.Tok {
	case token.IMPORT:
		// Only for unrefferenced imports
		for _, s := range d.Specs {
			is := s.(*ast.ImportSpec)
			if is.Name != nil && is.Name.Name == "_" {
				cdd := gtc.newCDD(gtc.ti.Objects[is.Name], ImportDecl, il)
				cdds = append(cdds, cdd)
			}
		}

	case token.CONST:
		for _, s := range d.Specs {
			vs := s.(*ast.ValueSpec)

			for _, n := range vs.Names {
				c := gtc.ti.Objects[n].(*types.Const)

				// All constants in expressions are evaluated so
				// only exported constants need be translated to C
				if !c.IsExported() {
					continue
				}

				cdd := gtc.newCDD(c, ConstDecl, il)

				cdd.indent(w)
				w.WriteString("#define ")
				cdd.Name(w, c, true)
				w.WriteByte(' ')
				cdd.Value(w, c.Val(), c.Type())
				cdd.copyDecl(w, "\n")
				w.Reset()

				cdds = append(cdds, cdd)
			}
		}

	case token.VAR:
		for _, s := range d.Specs {
			vs := s.(*ast.ValueSpec)
			vals := vs.Values

			for i, n := range vs.Names {
				v := gtc.ti.Objects[n].(*types.Var)
				typ := v.Type()
				cdd := gtc.newCDD(v, VarDecl, il)
				name := cdd.NameStr(v, true)

				cdd.indent(w)
				cdd.Type(w, typ)
				w.WriteByte(' ')
				w.WriteString(name)

				cinit := true // true if C declaration can init value

				if cdd.gtc.isGlobal(v) {
					cdd.copyDecl(w, ";\n") // Global variables may need declaration
					if i < len(vals) {
						_, cinit = cdd.gtc.ti.Values[vals[i]]
					}
				}
				if cinit {
					w.WriteString(" = ")
					if i < len(vals) {
						cdd.Expr(w, vals[i])
					} else {
						w.WriteString("{0}")
					}
				}
				w.WriteString(";\n")
				cdd.copyDef(w)

				if !cinit {
					// Runtime initialisation
					w.Reset()
					w.WriteByte('\t')
					w.WriteString(name)
					w.WriteString(" = ")
					cdd.Expr(w, vals[i])
					w.WriteString(";\n")
					cdd.copyInit(w)
				}

				w.Reset()

				cdds = append(cdds, cdd)
			}
		}

	case token.TYPE:
		for _, s := range d.Specs {
			ts := s.(*ast.TypeSpec)
			to := gtc.ti.Objects[ts.Name]
			tt := gtc.ti.Types[ts.Type]
			cdd := gtc.newCDD(to, TypeDecl, il)
			name := cdd.NameStr(to, true)

			cdd.indent(w)

			switch typ := tt.(type) {
			case *types.Struct:
				cdd.structDecl(w, name, typ)

			case *types.Signature:
				w.WriteString("typedef ")
				res := cdd.Signature(w, name, typ, false)
				if res.cdd != nil {
					cdds = append(cdds, res.cdd)
				}
				cdd.copyDecl(w, ";\n")

			default:
				w.WriteString("typedef ")
				cdd.Type(w, typ)
				w.WriteByte(' ')
				w.WriteString(name)
				cdd.copyDecl(w, ";\n")
			}
			w.Reset()

			cdds = append(cdds, cdd)
		}

	default:
		// Return fake CDD for unknown declaration
		cdds = []*CDD{{
			Decl: []byte(fmt.Sprintf("@%v (%T)@\n", d.Tok, d)),
		}}
	}
	return
}

func (cdd *CDD) structDecl(w *bytes.Buffer, name string, typ *types.Struct) {
	n := w.Len()

	w.WriteString("struct ")
	w.WriteString(name)
	w.WriteString("_struct;\n")
	cdd.indent(w)
	w.WriteString("typedef struct ")
	w.WriteString(name)
	w.WriteString("_struct ")
	w.WriteString(name)

	cdd.copyDecl(w, ";\n")
	w.Truncate(n)

	w.WriteString("struct ")
	w.WriteString(name)
	w.WriteByte('_')
	cdd.Type(w, typ)
	w.WriteString(";\n")

	cdd.copyDef(w)
	w.Truncate(n)
}

func (cc *GTC) Decl(decl ast.Decl, il int) []*CDD {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return cc.FuncDecl(d, il)

	case *ast.GenDecl:
		return cc.GenDecl(d, il)
	}

	panic(fmt.Sprint("Unknown declaration: ", decl))
}
