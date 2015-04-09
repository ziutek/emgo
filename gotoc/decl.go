package gotoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strconv"

	"golang.org/x/tools/go/types"
)

func (gtc *GTC) FuncDecl(d *ast.FuncDecl, il int) (cdds []*CDD) {
	f := gtc.object(d.Name).(*types.Func)
	sig := f.Type().(*types.Signature)
	cdd := gtc.newCDD(f, FuncDecl, il)
	cdds = append(cdds, cdd)
	fname := cdd.NameStr(f, true)
	w := new(bytes.Buffer)

	res, params := cdd.signature(sig, true, orgNames)

	w.WriteString(res.typ)
	w.WriteByte(' ')
	w.WriteString(dimFuncPtr(fname+params.String(), res.dim))

	cdd.initFunc = (f.Name() == "init" && sig.Recv() == nil && !cdd.gtc.isLocal(f))

	if !cdd.initFunc {
		cdd.copyDecl(w, ";\n")
	}

	if d.Body == nil {
		return
	}

	cdd.fbody = true

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
			t := v.Type()
			dim := cdd.Type(w, t)
			w.WriteByte(' ')
			w.WriteString(dimFuncPtr(name, dim))
			w.WriteString(" = ")
			zeroVal(w, t)
			w.WriteString(";\n")
		}
		cdd.indent(w)
	}

	end := cdd.BlockStmt(w, d.Body, res.typ, sig.Results())
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

	if cdd.initFunc {
		cdd.Init = []byte("\t" + fname + "();\n")
	}
	return
}

func (gtc *GTC) GenDecl(d *ast.GenDecl, il int) (cdds []*CDD) {
	w := new(bytes.Buffer)

	switch d.Tok {
	case token.IMPORT:
		/*
			// Only for unrefferenced imports
			for _, s := range d.Specs {
				is := s.(*ast.ImportSpec)
				if is.Name != nil && is.Name.Name == "_" {
					fmt.Println("is:", is.Name, is.Path)
					cdd := gtc.newCDD(gtc.object(is.Name), ImportDecl, il)
					cdds = append(cdds, cdd)
				}
			}
		*/

	case token.CONST:
		for _, s := range d.Specs {
			vs := s.(*ast.ValueSpec)

			for _, n := range vs.Names {
				c := gtc.object(n).(*types.Const)

				// All constants in expressions are evaluated so
				// only exported constants need be translated to C
				if gtc.isLocal(c) || !c.Exported() {
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
		indent := false
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
					if t, ok := cdd.exprType(val).(*types.Tuple); ok {
						gtc.notImplemented(s, t)
					}
				}
				if indent {
					cdd.indent(w)
				} else {
					indent = true
				}
				cdd.varDecl(w, v.Type(), name, val)
				w.Reset()
				cdds = append(cdds, cdd)
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
			case *types.Struct, *types.Interface:
				cdd.structDecl(w, name, typ)
			default:
				w.WriteString("typedef ")
				dim := cdd.Type(w, typ)
				w.WriteByte(' ')
				w.WriteString(dimFuncPtr(name, dim))
				cdd.copyDecl(w, ";\n")
			}
			typ := to.Type()
			w.Reset()
			cdd.tinfo(w, typ)
			if _, ok := typ.Underlying().(*types.Pointer); !ok {
				w.Reset()
				cdd.tinfo(w, types.NewPointer(typ))
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

func zeroVal(w *bytes.Buffer, typ types.Type) {
	switch t := typ.Underlying().(type) {
	case *types.Struct, *types.Array, *types.Slice, *types.Chan, *types.Interface:
		w.WriteString("{}")

	case *types.Pointer, *types.Signature:
		w.WriteString("nil")

	case *types.Basic:
		switch t.Kind() {
		case types.String:
			w.WriteString("{}")
		case types.Bool:
			w.WriteString("false")
		default:
			w.WriteByte('0')
		}

	default:
		w.WriteByte('0')
	}
}

func (cdd *CDD) varDecl(w *bytes.Buffer, typ types.Type, name string, val ast.Expr) {
	global := (cdd.Typ == VarDecl) && cdd.gtc.isGlobal(cdd.Origin)
	dim := cdd.Type(w, typ)
	w.WriteByte(' ')
	w.WriteString(dimFuncPtr(name, dim))
	if global {
		cdd.copyDecl(w, ";\n") // Global variables may need declaration
		cdd.constInit = (val == nil || cdd.isConstExpr(val, typ))
		if !cdd.constInit {
			w.WriteString(";\n")
			cdd.copyDef(w)
			w.Reset()
			cdd.il++
			cdd.indent(w)
			w.WriteString(name)
		}
	}

	w.WriteString(" = ")
	if val != nil {
		cdd.interfaceExpr(w, val, typ)
	} else {
		zeroVal(w, typ)
	}
	w.WriteByte(';')
	if cdd.Typ != VarDecl {
		return
	}
	w.WriteByte('\n')
	if global && !cdd.constInit {
		cdd.il--
		cdd.copyInit(w)
	} else {
		cdd.copyDef(w)
	}
}

// isConstExpr returns true if val can be represented as C constant expr.
// typ is destination type.
func (cdd *CDD) isConstExpr(val ast.Expr, typ types.Type) bool {
	if cdd.exprValue(val) != nil {
		return true
	}
	if _, ok := typ.Underlying().(*types.Interface); ok {
		if !types.Identical(typ, cdd.exprType(val)) {
			return false
		}
	}
	switch v := val.(type) {
	case *ast.Ident:
		return (v.Name == "nil")
	case *ast.UnaryExpr:
		return v.Op == token.AND
	case *ast.CompositeLit:
		switch t := typ.Underlying().(type) {
		case *types.Slice:
			return true
		case *types.Array:
			if len(v.Elts) == 0 {
				return true
			}
			elemt := t.Elem()
			for _, e := range v.Elts {
				if kv, ok := e.(*ast.KeyValueExpr); ok {
					if !cdd.isConstExpr(kv.Value, elemt) {
						return false
					}
				} else {
					if !cdd.isConstExpr(e, elemt) {
						return false
					}
				}
			}
			return true
		case *types.Struct:
			if len(v.Elts) == 0 {
				return true
			}
			if _, ok := v.Elts[0].(*ast.KeyValueExpr); ok {
				for _, e := range v.Elts {
					key := e.(*ast.KeyValueExpr).Key.(*ast.Ident).Name
					var elemt types.Type
					for i := 0; i < t.NumFields(); i++ {
						if t.Field(i).Name() == key {
							elemt = t.Field(i).Type()
							break
						}
					}
					if !cdd.isConstExpr(e.(*ast.KeyValueExpr).Value, elemt) {
						return false
					}
				}
			} else {
				for i, e := range v.Elts {
					if !cdd.isConstExpr(e, t.Field(i).Type()) {
						return false
					}
				}
			}
			return true
		default:
			cdd.gtc.notImplemented(val, typ)
		}
	}
	return false
}

var tuparrRe = regexp.MustCompile(`(\$\$)|(^\$[0-9]+_\$.)`)

func (cdd *CDD) structDecl(w *bytes.Buffer, name string, typ types.Type) {
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

	tuparr := tuparrRe.MatchString(name)
	if tuparr {
		cdd.indent(w)
		w.WriteString("#ifndef " + name + "$\n")
		cdd.indent(w)
		w.WriteString("#define " + name + "$\n")
	}
	cdd.indent(w)
	w.WriteString("struct ")
	w.WriteString(name)
	w.WriteByte('_')
	if it, ok := typ.(*types.Interface); ok {
		cdd.iface(w, it)
	} else {
		cdd.Type(w, typ)
	}
	w.WriteString(";\n")
	if tuparr {
		cdd.indent(w)
		w.WriteString("#endif\n")
	}

	cdd.copyDef(w)
	w.Truncate(n)
	return
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
