package gotoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"regexp"
	"strconv"
	"strings"
)

func (gtc *GTC) FuncDecl(d *ast.FuncDecl, il int) (cdds []*CDD) {
	f := gtc.object(d.Name).(*types.Func)
	sig := f.Type().(*types.Signature)
	cdd := gtc.newCDD(f, FuncDecl, il)
	cdds = append(cdds, cdd)
	fname := cdd.NameStr(f, true)
	w := new(bytes.Buffer)

	pragmas, cattrs := gtc.pragmas(d)
	for _, p := range pragmas {
		switch p {
		case "inline":
			cdd.Complexity -= cdd.gtc.noinlineThres
		case "noinline":
			cdd.Complexity += cdd.gtc.noinlineThres
		case "export":
			cdd.forceExport = true
		}
	}
	for _, cattr := range cattrs {
		w.WriteString(cattr)
		w.WriteByte('\n')
		cdd.indent(w)
	}

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

	cdd.where = inFuncBody

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
			pragmas, cattrs := gtc.pragmas(d, vs)
			var pconst, pexport bool
			for _, p := range pragmas {
				switch p {
				case "export":
					pexport = true
				case "const":
					pconst = true
				}
			}
			cattr := strings.Join(cattrs, " ")
			for i, n := range vs.Names {
				v := gtc.object(n).(*types.Var)
				cdd := gtc.newCDD(v, VarDecl, il)
				cdd.forceExport = pexport
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
				cdd.varDecl(w, v.Type(), name, val, cattr, pconst, true)
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
			cattrs, pexport := gtc.cattrs(d, s)
			cdd.forceExport = pexport
			switch typ := tt.(type) {
			case *types.Struct:
				cdd.structDecl(w, name, typ, cattrs)
			case *types.Interface:
				cdd.structDecl(w, name, typ, cattrs)
			default:
				w.WriteString("typedef ")
				if cattrs != "" {
					w.WriteString(cattrs)
					w.WriteByte(' ')
				}
				dim := cdd.Type(w, typ)
				w.WriteByte(' ')
				w.WriteString(dimFuncPtr(name, dim))
				cdd.copyDecl(w, ";\n")
			}
			w.Reset()
			typ := to.Type().(*types.Named)
			for i, n := 0, typ.NumMethods(); i < n; i++ {
				// Type uses its methods too.
				cdd.addObject(typ.Method(i), false)
			}
			if ts.Assign == 0 {
				// Not type alias.
				cdd.tinfo(w, typ)
				w.Reset()
				if _, ok := typ.Underlying().(*types.Pointer); !ok {
					cdd.tinfo(w, types.NewPointer(typ))
					w.Reset()
				}
			}
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

func (cdd *CDD) varDecl(w *bytes.Buffer, typ types.Type, name string, val ast.Expr, cattrs string, pconst, permitaa bool) {
	global := (cdd.Typ == VarDecl) && cdd.gtc.isGlobal(cdd.Origin)
	if cattrs != "" {
		w.WriteString(cattrs)
		w.WriteByte(' ')
	}
	dim := cdd.Type(w, typ)
	w.WriteByte(' ')
	if pconst {
		w.WriteString(dimFuncPtr("const "+name, dim))
	} else {
		w.WriteString(dimFuncPtr(name, dim))
	}
	if global {
		cdd.copyDecl(w, ";\n") // Global variables may need declaration
		w.Reset()
		cdd.indent(w)
		w.WriteString("__typeof__(" + name + ") " + name)
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
	if cdd.where == inDecl {
		cdd.where = inVarInit
	}
	w.WriteString(" = ")
	if val != nil {
		cdd.interfaceExpr(w, val, typ, permitaa)
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
	if val == nil {
		return true
	}
	if _, ok := typ.Underlying().(*types.Interface); ok {
		if !types.Identical(typ, cdd.exprType(val)) {
			return false
		}
	}
	if cdd.exprValue(val) != nil {
		return true
	}
	switch v := val.(type) {
	case *ast.Ident:
		if _, ok := typ.(*types.Signature); ok {
			return true
		}
		return v.Name == "nil"
	case *ast.UnaryExpr:
		if v.Op == token.AND {
			if ident, ok := v.X.(*ast.Ident); ok {
				o := cdd.object(ident)
				if o.Parent() == cdd.gtc.pkg.Scope() || cdd.gtc.isImported(o) {
					return true
				}
			}
		}
		return cdd.isConstExpr(v.X, cdd.exprType(v.X))
	case *ast.CompositeLit:
		switch t := typ.Underlying().(type) {
		case *types.Slice:
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
	case *ast.CallExpr:
		t := cdd.exprType(v.Fun)
		if _, ok := t.(*types.Signature); ok {
			return false
		}
		arg := v.Args[0]
		at := cdd.exprType(arg)
		switch typ := t.Underlying().(type) {
		case *types.Slice:
			if _, ok := at.Underlying().(*types.Basic); ok {
				return false // string
			}
		case *types.Basic:
			if typ.Kind() == types.String {
				if _, ok := at.(*types.Basic); !ok {
					return false
				}
			}

			/* Not need because -fno-strict-aliasing
			case *types.Pointer:
				if _, ok := at.Underlying().(*types.Pointer); !ok {
					// Casting unsafe.Pointer
					return false
				}
			*/
		}
		return cdd.isConstExpr(arg, at)
	case *ast.SelectorExpr:
		if t := cdd.exprType(v.X); t != nil {
			// Not package.
			return cdd.gtc.ti.Selections[v].Kind() == types.MethodExpr
		}
		_, ok := typ.(*types.Signature)
		return ok
	}
	return false
}

var tuparrRe = regexp.MustCompile(`(\$\$)|(^\$[0-9]+_\$.)`)

func (cdd *CDD) structDecl(w *bytes.Buffer, name string, typ types.Type, cattrs string) {
	n := w.Len()

	w.WriteString("struct ")
	w.WriteString(name)
	w.WriteString("_struct;\n")
	cdd.indent(w)
	w.WriteString("typedef ")
	if cattrs != "" {
		w.WriteString(cattrs)
		w.WriteByte(' ')
	}
	w.WriteString("struct ")
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
