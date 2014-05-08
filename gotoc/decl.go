package gotoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"code.google.com/p/go.tools/go/types"
)

func (gtc *GTC) FuncDecl(d *ast.FuncDecl, il int) (cdds []*CDD) {
	f := gtc.object(d.Name).(*types.Func)

	cdd := gtc.newCDD(f, FuncDecl, il)
	w := new(bytes.Buffer)
	fname := cdd.NameStr(f, true)

	sig := f.Type().(*types.Signature)
	res, params := cdd.signature(sig, true)

	w.WriteString(res.typ)
	w.WriteByte(' ')
	w.WriteString(dimFuncPtr(fname+params, res.dim))

	cdds = append(cdds, res.acds...)
	cdds = append(cdds, cdd)

	cdd.init = (f.Name() == "init" && sig.Recv() == nil && !cdd.gtc.isLocal(f))

	if !cdd.init {
		cdd.copyDecl(w, ";\n")
	}

	if d.Body == nil {
		return
	}

	cdd.body = true

	w.WriteByte(' ')

	all := true
	if res.hasNames {
		cdd.indent(w)
		w.WriteString("{\n")
		cdd.il++
		for i, v := range res.fields {
			name := res.names[i]
			if name == "_" && len(res.fields) > 1 {
				all = false
				continue
			}
			cdd.indent(w)
			dim, acds := cdd.Type(w, v.Type())
			cdds = append(cdds, acds...)
			w.WriteByte(' ')
			w.WriteString(dimFuncPtr(name, dim))
			w.WriteString(" = {0};\n")
		}
		cdd.indent(w)
	}

	end, acds := cdd.BlockStmt(w, d.Body, res.typ, sig.Results())
	cdds = append(cdds, acds...)
	w.WriteByte('\n')

	if res.hasNames {
		if end {
			cdd.il--
			cdd.indent(w)
			w.WriteString("end:\n")
			cdd.il++

			cdd.indent(w)
			w.WriteString("return ")
			if len(res.fields) == 1 {
				w.WriteString(res.names[0])
			} else {
				w.WriteString("(" + res.typ + "){")
				comma := false
				for i, name := range res.names {
					if name == "_" {
						continue
					}
					if comma {
						w.WriteString(", ")
					} else {
						comma = true
					}
					if !all {
						w.WriteString("._" + strconv.Itoa(i) + "=")
					}
					w.WriteString(name)
				}
				w.WriteByte('}')
			}
			w.WriteString(";\n")
		}
		cdd.il--
		w.WriteString("}\n")
	}
	cdd.copyDef(w)

	if cdd.init {
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
				cdd := gtc.newCDD(gtc.object(is.Name), ImportDecl, il)
				cdds = append(cdds, cdd)
			}
		}

	case token.CONST:
		for _, s := range d.Specs {
			vs := s.(*ast.ValueSpec)

			for _, n := range vs.Names {
				c := gtc.object(n).(*types.Const)

				// All constants in expressions are evaluated so
				// only exported constants need be translated to C
				if !c.Exported() {
					continue
				}

				cdd := gtc.newCDD(c, ConstDecl, il)

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
				v := gtc.object(n).(*types.Var)
				cdd := gtc.newCDD(v, VarDecl, il)
				name := cdd.NameStr(v, true)

				var val ast.Expr
				if i < len(vals) {
					val = vals[i]
				}
				if i > 0 {
					cdd.indent(w)
				}
				acds := cdd.varDecl(w, v.Type(), cdd.gtc.isGlobal(v), name, val)

				w.Reset()

				cdds = append(cdds, cdd)
				cdds = append(cdds, acds...)
			}
		}

	case token.TYPE:
		for i, s := range d.Specs {
			ts := s.(*ast.TypeSpec)
			to := gtc.object(ts.Name)
			tt := gtc.exprType(ts.Type)
			cdd := gtc.newCDD(to, TypeDecl, il)
			name := cdd.NameStr(to, true)

			if i > 0 {
				cdd.indent(w)
			}

			switch typ := tt.(type) {
			case *types.Struct:
				cdd.structDecl(w, name, typ)

			default:
				w.WriteString("typedef ")
				dim, acds := cdd.Type(w, typ)
				cdds = append(cdds, acds...)
				w.WriteByte(' ')
				w.WriteString(dimFuncPtr(name, dim))
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

func (cdd *CDD) varDecl(w *bytes.Buffer, typ types.Type, global bool, name string, val ast.Expr) (acds []*CDD) {

	dim, acds := cdd.Type(w, typ)
	w.WriteByte(' ')
	w.WriteString(dimFuncPtr(name, dim))

	typ = underlying(typ)

	constInit := true // true if C declaration can init value

	if global {
		cdd.copyDecl(w, ";\n") // Global variables may need declaration
		if val != nil {
			constInit = cdd.exprValue(val) != nil
		}
	}

	if constInit {
		if val != nil {
			w.WriteString(" = ")
			cdd.Expr(w, val, typ)
		} else {
			if t, ok := typ.(*types.Struct); ok && t.NumFields() == 0 {
			} else if t, ok := typ.(*types.Array); ok && t.Len() == 0 {
			} else {
				w.WriteString(" = {0}")
			}
		}
	}
	w.WriteString(";\n")
	cdd.copyDef(w)

	if !constInit {
		// Runtime initialisation
		w.Reset()

		assign := false

		switch t := typ.(type) {
		case *types.Slice:
			switch vt := val.(type) {
			case *ast.CompositeLit:
				aname := "array" + cdd.gtc.uniqueId()
				at := types.NewArray(t.Elem(), int64(len(vt.Elts)))
				o := types.NewVar(vt.Lbrace, cdd.gtc.pkg, aname, at)
				cdd.gtc.pkg.Scope().Insert(o)
				acd := cdd.gtc.newCDD(o, VarDecl, cdd.il)
				av := *vt
				cdd.gtc.ti.Types[&av] = types.TypeAndValue{Type: at} // BUG: thread-unsafe
				n := w.Len()
				acd.varDecl(w, o.Type(), cdd.gtc.isGlobal(o), aname, &av)
				w.Truncate(n)
				acds = append(acds, acd)

				w.WriteByte('\t')
				w.WriteString(name)
				w.WriteString(" = ASLICE(")
				w.WriteString(aname)
				w.WriteString(");\n")

			default:
				assign = true
			}

		case *types.Array:
			w.WriteByte('\t')
			w.WriteString("ACPY(")
			w.WriteString(name)
			w.WriteString(", ")

			switch val.(type) {
			case *ast.CompositeLit:
				w.WriteString("((")
				dim, _ := cdd.Type(w, t.Elem())
				dim = append([]string{"[]"}, dim...)
				w.WriteString("(" + dimFuncPtr("", dim) + "))")
				cdd.Expr(w, val, typ)

			default:
				cdd.Expr(w, val, typ)
			}

			w.WriteString("));\n")

		case *types.Pointer:
			u, ok := val.(*ast.UnaryExpr)
			if !ok {
				assign = true
				break
			}
			c, ok := u.X.(*ast.CompositeLit)
			if !ok {
				assign = true
				break
			}
			cname := "cl" + cdd.gtc.uniqueId()
			ct := cdd.exprType(c)
			o := types.NewVar(c.Lbrace, cdd.gtc.pkg, cname, ct)
			cdd.gtc.pkg.Scope().Insert(o)
			acd := cdd.gtc.newCDD(o, VarDecl, cdd.il)
			n := w.Len()
			acd.varDecl(w, o.Type(), cdd.gtc.isGlobal(o), cname, c)
			w.Truncate(n)
			acds = append(acds, acd)

			w.WriteByte('\t')
			w.WriteString(name)
			w.WriteString(" = &")
			w.WriteString(cname)
			w.WriteString(";\n")

		default:
			assign = true
		}

		if assign {
			// Ordinary assignment gos to the init() function
			cdd.init = true
			w.WriteByte('\t')
			w.WriteString(name)
			w.WriteString(" = ")
			cdd.Expr(w, val, typ)
			w.WriteString(";\n")
		}
		cdd.copyInit(w)
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

	tuple := strings.Contains(name, "$$")

	if tuple {
		cdd.indent(w)
		w.WriteString("#ifndef " + name + "$\n")
		cdd.indent(w)
		w.WriteString("#define " + name + "$\n")
	}
	cdd.indent(w)
	w.WriteString("struct ")
	w.WriteString(name)
	w.WriteByte('_')
	cdd.Type(w, typ)
	w.WriteString(";\n")
	if tuple {
		cdd.indent(w)
		w.WriteString("#endif\n")
	}

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
