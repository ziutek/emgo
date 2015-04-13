package gotoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"golang.org/x/tools/go/exact"
	"golang.org/x/tools/go/types"
)

func writeInt(w *bytes.Buffer, ev exact.Value, k types.BasicKind) {
	if k == types.Uintptr {
		u, _ := exact.Uint64Val(ev)
		w.WriteString("0x")
		w.WriteString(strconv.FormatUint(u, 16))
		return
	}
	s := ev.String()
	if s[0] == '-' {
		w.WriteByte('(')
	}
	switch k {
	case types.Int32:
		if s == "-2147483648" {
			w.WriteString("-2147483647L-1L")
		} else {
			w.WriteString(s + "L")
		}
	case types.Uint32:
		w.WriteString(s + "UL")
	case types.Int64:
		if s == "-9223372036854775808" {
			w.WriteString("-9223372036854775807LL-1LL")
		} else {
			w.WriteString(s + "LL")
		}
	case types.Uint64:
		w.WriteString(s + "ULL")
	default:
		w.WriteString(s)
	}
	if s[0] == '-' {
		w.WriteByte(')')
	}
}

func writeFloat(w *bytes.Buffer, ev exact.Value, k types.BasicKind) {
	w.WriteByte('(')
	if k == types.Float32 {
		f, _ := exact.Float32Val(ev)
		w.WriteString(strconv.FormatFloat(float64(f), 'e', -1, 32))
		w.WriteByte('F')
	} else {
		f, _ := exact.Float64Val(ev)
		w.WriteString(strconv.FormatFloat(f, 'e', -1, 64))
	}
	w.WriteByte(')')
}

func (cdd *CDD) Value(w *bytes.Buffer, ev exact.Value, t types.Type) {
	k := t.Underlying().(*types.Basic).Kind()
	switch {
	case k <= types.Bool || k == types.UntypedBool:
		w.WriteString(ev.String())
	case k <= types.Uintptr || k == types.UntypedInt || k == types.UntypedRune:
		writeInt(w, ev, k)
	case k <= types.Float64 || k == types.UntypedFloat:
		writeFloat(w, ev, k)
	case k <= types.Complex128 || k == types.UntypedComplex:
		writeFloat(w, exact.Real(ev), k)
		im := exact.Imag(ev)
		if exact.Sign(im) != -1 {
			w.WriteByte('+')
		}
		writeFloat(w, im, k)
		w.WriteByte('i')
	case k == types.String || k == types.UntypedString:
		w.WriteString("EGSTR(")
		w.WriteString(ev.String())
		w.WriteByte(')')
	default:
		fmt.Println("Kind", k)
		w.WriteString(ev.String())
	}
}

func (cdd *CDD) Name(w *bytes.Buffer, obj types.Object, direct bool) {
	/*if obj == nil {
		w.WriteByte('_')
		return
	}*/
	switch o := obj.(type) {
	case *types.PkgName:
		// Imported package name in SelectorExpr: pkgname.Name
		w.WriteString(upath(o.Pkg().Path()))
		return

	case *types.Func:
		s := o.Type().(*types.Signature)
		if r := s.Recv(); r != nil {
			t := r.Type()
			if p, ok := t.(*types.Pointer); ok {
				t = p.Elem()
				direct = false
			}
			cdd.Type(w, t)
			w.WriteByte('$')
			w.WriteString(o.Name())
			if !cdd.gtc.isLocal(t.(*types.Named).Obj()) {
				cdd.addObject(o, direct)
			}
			return
		}
	}

	if p := obj.Pkg(); p != nil && !cdd.gtc.isLocal(obj) {
		cdd.addObject(obj, direct)
		w.WriteString(upath(obj.Pkg().Path()))
		w.WriteByte('$')
	}
	name := obj.Name()
	switch name {
	case "init":
		w.WriteString(cdd.gtc.uniqueId() + name)

	default:
		w.WriteString(name)
		if cdd.gtc.isLocal(obj) {
			w.WriteByte('$')
		}
	}
}

func (cdd *CDD) NameStr(o types.Object, direct bool) string {
	buf := new(bytes.Buffer)
	cdd.Name(buf, o, direct)
	return buf.String()
}

func (cdd *CDD) SelectorExpr(w *bytes.Buffer, e *ast.SelectorExpr) (fun, recvt types.Type, recvs string) {
	sel := cdd.gtc.ti.Selections[e]
	if sel == nil {
		cdd.Name(w, cdd.object(e.Sel), true)
		return
	}
	s := cdd.ExprStr(e.X, nil)
	index := sel.Index()
	rt := sel.Recv()
	for _, id := range index[:len(index)-1] {
		if p, ok := rt.(*types.Pointer); ok {
			rt = p.Elem()
			s += "->"
		} else {
			s += "."
		}
		f := rt.Underlying().(*types.Struct).Field(id)
		s += f.Name()
		rt = f.Type()
	}
	rpt, isPtr := rt.(*types.Pointer)
	switch sel.Kind() {
	case types.FieldVal:
		w.WriteString(s)
		if isPtr {
			w.WriteString("->")
		} else {
			w.WriteByte('.')
		}
		w.WriteString(e.Sel.Name)

	case types.MethodVal:
		fun = sel.Obj().Type()
		rtyp := fun.(*types.Signature).Recv().Type()

		switch rtyp.Underlying().(type) {
		case *types.Interface:
			// Method with interface receiver.
			w.WriteString(e.Sel.Name)
			recvs = s
			recvt = rt

		case *types.Pointer:
			// Method with pointer receiver.
			cdd.Name(w, sel.Obj(), true)
			if isPtr {
				recvs = s
				recvt = rt
			} else {
				recvs = "&" + s
				recvt = types.NewPointer(rt)
			}
		default:
			// Method with non-pointer receiver.
			cdd.Name(w, sel.Obj(), true)
			if isPtr {
				recvs = "*" + s
				recvt = rpt.Elem()
			} else {
				recvs = s
				recvt = rt
			}
		}

	case types.MethodExpr:
		cdd.Name(w, sel.Obj(), true)

	default:
		cdd.notImplemented(e)
	}
	return
}

func (cdd *CDD) SelectorExprStr(e *ast.SelectorExpr) (s string, fun, recvt types.Type, recvs string) {
	buf := new(bytes.Buffer)
	fun, recvt, recvs = cdd.SelectorExpr(buf, e)
	s = buf.String()
	return
}

func (cdd *CDD) builtin(b *types.Builtin, args []ast.Expr) (fun, recv string) {
	name := b.Name()

	switch name {
	case "len":
		switch t := cdd.exprType(args[0]).Underlying().(type) {
		case *types.Slice, *types.Basic: // Basic == String
			return "len", ""

		case *types.Array:
			panic("builtin len(array) isn't handled as constant expression")
			// return "", strconv.FormatInt(t.Len(), 10)

		case *types.Chan:
			return "clen", ""

		default:
			cdd.notImplemented(ast.NewIdent("len"), t)
		}

	case "cap":
		switch t := cdd.exprType(args[0]).Underlying().(type) {
		case *types.Slice:
			return "cap", ""

		case *types.Chan:
			return "ccap", ""

		default:
			cdd.notImplemented(ast.NewIdent("cap"), t)
		}

	case "copy":
		switch t := cdd.exprType(args[1]).Underlying().(type) {
		case *types.Basic: // string
			return "STRCPY", ""

		case *types.Slice:
			typ, dim := cdd.TypeStr(t.Elem())
			return "SLICPY", typ + dimFuncPtr("", dim)

		default:
			panic(t)
		}

	case "new":
		typ, dim := cdd.TypeStr(cdd.exprType(args[0]))
		args[0] = nil
		return "NEW", typ + dimFuncPtr("", dim)

	case "make":
		a0t := cdd.exprType(args[0])
		args[0] = nil

		switch t := a0t.Underlying().(type) {
		case *types.Slice:
			typ, dim := cdd.TypeStr(t.Elem())
			name := "MAKESLI"
			if len(args) == 3 {
				name = "MAKESLIC"
			}
			return name, typ + dimFuncPtr("", dim)

		case *types.Chan:
			typ, dim := cdd.TypeStr(t.Elem())
			typ += dimFuncPtr("", dim)
			if len(args) == 1 {
				typ += ", 0"
			}
			return "MAKECHAN", typ

		case *types.Map:
			typ, dim := cdd.TypeStr(t.Key())
			k := typ + dimFuncPtr("", dim)
			typ, dim = cdd.TypeStr(t.Elem())
			e := typ + dimFuncPtr("", dim)
			name := "MAKEMAP"
			if len(args) == 2 {
				name = "MAKEMAPC"
			}
			return name, k + ", " + e

		default:
			cdd.notImplemented(ast.NewIdent(name))
		}

	}

	return name, ""
}

func (cdd *CDD) funStr(fe ast.Expr, args []ast.Expr) (fs string, ft types.Type, rs string, rt types.Type) {
	switch f := fe.(type) {
	case *ast.SelectorExpr:
		buf := new(bytes.Buffer)
		ft, rt, rs = cdd.SelectorExpr(buf, f)
		fs = buf.String()
		return

	case *ast.Ident:
		switch o := cdd.object(f).(type) {
		case *types.Builtin:
			fs, rs = cdd.builtin(o, args)

		default:
			fs = cdd.NameStr(o, true)
			ft = o.Type()
		}
		return
	}
	fs = cdd.ExprStr(fe, nil)
	ft = cdd.exprType(fe)
	return
}

func (cdd *CDD) CallExpr(w *bytes.Buffer, e *ast.CallExpr) {
	switch t := cdd.exprType(e.Fun).(type) {
	case *types.Signature:
		c := cdd.call(e, t, false)
		if c.rcv.r != "" || c.arr.r != "" {
			w.WriteString("({\n")
			cdd.il++
			cdd.indent(w)
		}
		if c.rcv.r != "" {
			dim := cdd.Type(w, c.rcv.t)
			w.WriteString(" " + dimFuncPtr(c.rcv.l, dim) + " = ")
			w.WriteString(indent(1, c.rcv.r) + ";\n")
			cdd.indent(w)
		}
		if c.arr.r != "" {
			argv := c.args
			if c.fun.r != "" {
				argv = append([]arg{c.fun}, c.args...)
			}
			for i, arg := range argv {
				if i == len(argv)-1 {
					// Variadic function.
					dim := cdd.Type(w, c.arr.t)
					w.WriteString(" " + dimFuncPtr(c.arr.l, dim) + " = ")
					w.WriteString(indent(1, c.arr.r) + ";\n")
					cdd.indent(w)
				}
				if arg.r == "" {
					continue // Don't evaluate
				}
				dim := cdd.Type(w, arg.t)
				w.WriteString(" " + dimFuncPtr(arg.l, dim) + " = ")
				w.WriteString(indent(1, arg.r) + ";\n")
				cdd.indent(w)
			}
		}
		w.WriteString(c.fun.l + "(")
		for i, arg := range c.args {
			if i > 0 {
				w.WriteString(", ")
			}
			w.WriteString(arg.l)
		}
		w.WriteString(")")
		if c.rcv.r != "" || c.arr.r != "" {
			w.WriteString(";\n")
			cdd.il--
			cdd.indent(w)
			w.WriteString("})")
		}

	default:
		arg := e.Args[0]
		switch typ := t.Underlying().(type) {
		case *types.Slice:
			switch cdd.exprType(arg).Underlying().(type) {
			case *types.Basic: // string
				w.WriteString("BYTES(")
				cdd.Expr(w, arg, typ)
				w.WriteByte(')')

			default: // slice
				w.WriteByte('(')
				cdd.Expr(w, arg, typ)
				w.WriteByte(')')
			}

		case *types.Interface:
			cdd.interfaceExpr(w, arg, t)

		default:
			w.WriteString("CAST(")
			dim := cdd.Type(w, t)
			w.WriteString(dimFuncPtr("", dim))
			w.WriteString(", ")
			cdd.Expr(w, arg, t)
			w.WriteString(")")
		}
	}
}

func (cdd *CDD) Expr(w *bytes.Buffer, expr ast.Expr, nilT types.Type) {
	cdd.Complexity++

	if t := cdd.gtc.ti.Types[expr]; t.Value != nil {
		// Constant expression
		cdd.Value(w, t.Value, t.Type)
		return
	}

	switch e := expr.(type) {
	case *ast.BinaryExpr:
		op := e.Op.String()
		ltyp := cdd.exprType(e.X)
		rtyp := cdd.exprType(e.Y)

		lhs := cdd.ExprStr(e.X, ltyp)
		rhs := cdd.ExprStr(e.Y, rtyp)

		if op == "==" || op == "!=" {
			eq(w, lhs, op, rhs, ltyp, rtyp)
			break
		}
		// BUG: strings
		if op == "&^" {
			op = "&~"
		}
		w.WriteString("(" + lhs + op + rhs + ")")

	case *ast.UnaryExpr:
		if e.Op == token.ARROW {
			t := cdd.exprType(e.X).(*types.Chan).Elem()
			if tup, ok := cdd.exprType(e).(*types.Tuple); ok {
				tn, _ := cdd.tupleName(tup)
				w.WriteString("RECVOK(" + tn + ", ")
				cdd.Expr(w, e.X, nil)
				w.WriteByte(')')
			} else {
				w.WriteString("RECV(")
				dim := cdd.Type(w, t)
				w.WriteString(dimFuncPtr("", dim))
				w.WriteString(", ")
				cdd.Expr(w, e.X, nil)
				w.WriteString(", ")
				zeroVal(w, t)
				w.WriteByte(')')
			}
			break
		}
		if e.Op == token.AND {
			cdd.ptrExpr(w, e.X)
			break
		}
		op := e.Op.String()
		if e.Op == token.XOR {
			op = "~"
		}
		w.WriteString(op)
		cdd.Expr(w, e.X, nil)

	case *ast.CallExpr:
		cdd.CallExpr(w, e)

	case *ast.Ident:
		if e.Name == "nil" {
			cdd.Nil(w, nilT)
			break
		}
		if o := cdd.object(e); o != nil {
			cdd.Name(w, o, true)
			break
		}
		w.WriteString(e.Name)
		if e.Name != "_" {
			w.WriteByte('$')
		}

	case *ast.IndexExpr:
		cdd.indexExpr(w, cdd.exprType(e.X), cdd.ExprStr(e.X, nil), e.Index)

	case *ast.KeyValueExpr:
		kt := cdd.exprType(e.Key)
		if i, ok := e.Key.(*ast.Ident); ok && kt == nil {
			// e.Key is field name
			w.WriteByte('.')
			w.WriteString(i.Name)
		} else {
			w.WriteByte('[')
			cdd.Expr(w, e.Key, kt)
			w.WriteByte(']')
		}
		w.WriteString(" = ")
		cdd.interfaceExpr(w, e.Value, nilT)

	case *ast.ParenExpr:
		w.WriteByte('(')
		cdd.Expr(w, e.X, nilT)
		w.WriteByte(')')

	case *ast.SelectorExpr:
		s, fun, recvt, recvs := cdd.SelectorExprStr(e)
		if recvt == nil {
			w.WriteString(s)
			break
		}
		sig := fun.(*types.Signature)
		w.WriteString("({")
		res, params := cdd.signature(sig, false, numNames)
		w.WriteString(res.typ)
		w.WriteByte(' ')
		w.WriteString(dimFuncPtr("func"+params.String(), res.dim))
		w.WriteString(" { return " + s + "(" + recvs)
		if p := sig.Params(); p != nil {
			for i := 1; i <= p.Len(); i++ {
				w.WriteString(", _" + strconv.Itoa(i))
			}
		}
		w.WriteString("); } func;})")

	case *ast.SliceExpr:
		cdd.SliceExpr(w, e)

	case *ast.StarExpr:
		w.WriteByte('*')
		cdd.Expr(w, e.X, nil)

	case *ast.TypeAssertExpr:
		w.WriteString("({\n")
		cdd.il++
		ityp := cdd.exprType(e.X)
		cdd.indent(w)
		cdd.varDecl(w, ityp, "_i", e.X)
		w.WriteByte('\n')
		cdd.indent(w)
		w.WriteString("bool _ok = ")
		etyp := cdd.exprType(e)
		iempty := (cdd.gtc.methodSet(ityp).Len() == 0)
		typ := cdd.exprType(e.Type)
		if _, ok := typ.Underlying().(*types.Interface); ok {
			if cdd.gtc.methodSet(typ).Len() == 0 {
				w.WriteString("true;\n")
			} else {
				w.WriteString("implements(")
				if iempty {
					w.WriteString("_i.itab$, &")
				} else {
					w.WriteString("TINFO(_i), &")
				}
				w.WriteString(cdd.tinameDU(typ))
				w.WriteString(");\n")
			}
			cdd.indent(w)
			if t, ok := etyp.(*types.Tuple); ok {
				tn, _ := cdd.tupleName(t)
				w.WriteString("(" + tn + "){")
				cdd.interfaceES(w, "_i", e.Pos(), ityp, typ)
				w.WriteString(", _ok};\n")
			} else {
				w.WriteString("if (!_ok) panicIC();\n")
				cdd.indent(w)
				cdd.interfaceES(w, "_i", e.Pos(), ityp, typ)
				w.WriteString(";\n")
			}
		} else {
			if iempty {
				w.WriteString("(_i.itab$ == &")
			} else {
				w.WriteString("(TINFO(_i) == &")
			}
			w.WriteString(cdd.tinameDU(typ))
			w.WriteString(");\n")
			cdd.indent(w)
			if t, ok := etyp.(*types.Tuple); ok {
				tn, _ := cdd.tupleName(t)
				w.WriteString("(" + tn + "){IVAL(_i, ")
				dim := cdd.Type(w, typ)
				w.WriteString(dimFuncPtr("", dim))
				w.WriteString("), _ok};\n")
			} else {
				w.WriteString("if (!_ok) panicIC();\n")
				cdd.indent(w)
				w.WriteString("IVAL(_i, ")
				dim := cdd.Type(w, typ)
				w.WriteString(dimFuncPtr("", dim))
				w.WriteString(");\n")
			}
		}
		cdd.il--
		cdd.indent(w)
		w.WriteString("})")

	case *ast.CompositeLit:
		cdd.compositeLit(w, e, nilT)

	case *ast.FuncLit:
		fname := "func"

		fd := &ast.FuncDecl{
			Name: &ast.Ident{NamePos: e.Type.Func, Name: fname},
			Type: e.Type,
			Body: e.Body,
		}
		sig := cdd.exprType(e).(*types.Signature)
		cdd.gtc.ti.Defs[fd.Name] = types.NewFunc(e.Type.Func, cdd.gtc.pkg, fname, sig)

		w.WriteString("({\n")
		cdd.il++

		cdds := cdd.gtc.FuncDecl(fd, cdd.il)
		for _, c := range cdds {
			for u, typPtr := range c.FuncBodyUses {
				cdd.FuncBodyUses[u] = typPtr
			}
			cdd.indent(w)
			w.Write(c.Def)
		}

		cdd.indent(w)
		w.WriteString(fname + "$;\n")

		cdd.il--
		cdd.indent(w)
		w.WriteString("})")

	default:
		fmt.Fprintf(w, "!%v<%T>!", e, e)
	}
	return
}

func (cdd *CDD) newVar(name string, typ types.Type, global bool, val ast.Expr) {
	var pos token.Pos
	if val != nil {
		pos = val.Pos()
	}
	o := types.NewVar(pos, cdd.gtc.pkg, name, typ)
	if global {
		cdd.gtc.pkg.Scope().Insert(o)
	}
	acd := cdd.gtc.newCDD(o, VarDecl, 0)
	cdd.acds = append(cdd.acds, acd)
	acd.varDecl(new(bytes.Buffer), typ, name, val)
}

func (cdd *CDD) ptrExpr(w *bytes.Buffer, e ast.Expr) {
	w.WriteByte('&')
	cl, ok := e.(*ast.CompositeLit)
	if !ok || cdd.Typ != VarDecl || !cdd.gtc.isGlobal(cdd.Origin) {
		cdd.Expr(w, e, nil)
		return
	}
	name := "_cl" + cdd.gtc.uniqueId()
	w.WriteString(name)
	cdd.newVar(name, cdd.exprType(cl), true, cl)
}

func (cdd *CDD) compositeLit(w *bytes.Buffer, e *ast.CompositeLit, nilT types.Type) {
	typ := cdd.exprType(e)

	switch t := typ.Underlying().(type) {
	case *types.Array:
		if !cdd.constInit {
			w.WriteString("((")
			dim := cdd.Type(w, typ)
			w.WriteString(dimFuncPtr("", dim))
			w.WriteByte(')')
		}
		w.WriteString("{{")
		cdd.elts(w, e.Elts, t.Elem(), nil)
		w.WriteString("}}")
		if !cdd.constInit {
			w.WriteByte(')')
		}

	case *types.Slice:
		alen := cdd.clArrayLen(e.Elts)
		slen := strconv.FormatInt(alen, 10)
		if !cdd.constInit {
			w.WriteString("CSLICE(" + slen + ", ((")
			dim := cdd.Type(w, t.Elem())
			dim = append([]string{"[]"}, dim...)
			w.WriteString(dimFuncPtr("", dim))
			w.WriteString("){")
			cdd.elts(w, e.Elts, t.Elem(), nil)
			w.WriteString("}))")
			break
		}
		aname := "_cl" + cdd.gtc.uniqueId()
		w.WriteString("{&" + aname + ", " + slen + ", " + slen + "}")

		typ := types.NewArray(t.Elem(), alen)
		tv := cdd.gtc.ti.Types[e]
		tv.Type = typ
		cdd.gtc.ti.Types[e] = tv
		cdd.newVar(aname, typ, true, e)

	case *types.Struct:
		if !cdd.constInit {
			w.WriteString("((")
			cdd.Type(w, typ)
			w.WriteByte(')')
		}
		w.WriteByte('{')
		cdd.elts(w, e.Elts, nil, t)
		w.WriteByte('}')
		if !cdd.constInit {
			w.WriteByte(')')
		}
	default:
		cdd.notImplemented(e, t)
	}
}

func (cdd *CDD) clArrayLen(elems []ast.Expr) int64 {
	if len(elems) == 0 {
		return 0
	}
	var n, k int64
	for _, e := range elems {
		if kv, ok := e.(*ast.KeyValueExpr); ok {
			k, _ = exact.Int64Val(cdd.gtc.exprValue(kv.Key))

		}
		if k++; k > n {
			n = k
		}
	}
	return n
}

func (cdd *CDD) elts(w *bytes.Buffer, elts []ast.Expr, nilT types.Type, st *types.Struct) {
	for i, el := range elts {
		if i > 0 {
			w.WriteString(", ")
		}
		if nilT != nil {
			cdd.interfaceExpr(w, el, nilT)
		} else {
			cdd.interfaceExpr(w, el, st.Field(i).Type())
		}
	}
}

func (cdd *CDD) indexExpr(w *bytes.Buffer, typ types.Type, xs string, idx ast.Expr) {
	pt, isPtr := typ.(*types.Pointer)
	if isPtr {
		typ = pt.Elem()
	}
	var indT types.Type

	switch t := typ.Underlying().(type) {
	case *types.Basic: // string
		if cdd.gtc.boundsCheck {
			w.WriteString("STRIDXC(")
		} else {
			w.WriteString("STRIDX(")
		}
	case *types.Slice:
		if cdd.gtc.boundsCheck {
			w.WriteString("SLIDXC(")
		} else {
			w.WriteString("SLIDX(")
		}
		dim := cdd.Type(w, t.Elem())
		dim = append([]string{"*"}, dim...)
		w.WriteString(dimFuncPtr("", dim))
		w.WriteString(", ")
	case *types.Array:
		if cdd.gtc.boundsCheck {
			w.WriteString("AIDXC(")
		} else {
			w.WriteString("AIDX(")
		}
	case *types.Map:
		indT = t.Key()
		cdd.notImplemented(&ast.IndexExpr{}, t)
	default:
		panic(t)
	}
	if isPtr {
		w.WriteByte('*')
	}
	w.WriteString(xs)
	w.WriteString(", ")
	cdd.Expr(w, idx, indT)
	w.WriteByte(')')
}

func (cdd *CDD) indexExprStr(typ types.Type, xs string, idx ast.Expr) string {
	buf := new(bytes.Buffer)
	cdd.indexExpr(buf, typ, xs, idx)
	return buf.String()
}

func (cdd *CDD) SliceExpr(w *bytes.Buffer, e *ast.SliceExpr) {
	sx := cdd.ExprStr(e.X, nil)

	typ := cdd.exprType(e.X)
	pt, isPtr := typ.(*types.Pointer)
	if isPtr {
		typ = pt.Elem()
		sx = "(*" + sx + ")"
	}
	var slice *types.Slice
	empty := e.Low == nil && e.High == nil && e.Max == nil
	switch t := typ.Underlying().(type) {
	case *types.Array:
		if empty {
			w.WriteString("ASLICE(" + sx + ")")
			return
		}
		w.WriteByte('A')
	case *types.Basic: // string
		if empty {
			w.WriteString(sx)
			return
		}
		w.WriteByte('S')
	case *types.Slice:
		if empty {
			w.WriteString(sx)
			return
		}
		slice = t
	}
	w.WriteString("SLICE")
	if e.Low != nil {
		switch {
		case e.High == nil && e.Max == nil:
			w.WriteByte('L')
		case e.High != nil && e.Max == nil:
			w.WriteString("LH")
		case e.High == nil && e.Max != nil:
			w.WriteByte('M')
		default:
			w.WriteString("LHM")
		}
		if cdd.gtc.boundsCheck {
			w.WriteString("C(")
		} else {
			w.WriteByte('(')
		}
		w.WriteString(sx)
		if slice != nil {
			w.WriteString(", ")
			dim := cdd.Type(w, slice.Elem())
			dim = append([]string{"*"}, dim...)
			w.WriteString(dimFuncPtr("", dim))
		}
		w.WriteString(", ")
		cdd.Expr(w, e.Low, nil)
	} else {
		switch {
		case e.High != nil && e.Max == nil:
			w.WriteByte('H')
		case e.High == nil && e.Max != nil:
			w.WriteByte('M')
		default:
			w.WriteString("HM")
		}
		if cdd.gtc.boundsCheck {
			w.WriteString("C(")
		} else {
			w.WriteByte('(')
		}
		w.WriteString(sx)
	}
	if e.High != nil {
		w.WriteString(", ")
		cdd.Expr(w, e.High, nil)
	}
	if e.Max != nil {
		w.WriteString(", ")
		cdd.Expr(w, e.Max, nil)
	}
	w.WriteByte(')')
}

func (cdd *CDD) ExprStr(expr ast.Expr, nilT types.Type) string {
	buf := new(bytes.Buffer)
	cdd.Expr(buf, expr, nilT)
	return buf.String()
}

func (cdd *CDD) Nil(w *bytes.Buffer, t types.Type) {
	switch t.Underlying().(type) {
	case *types.Slice:
		w.WriteString("NILSLICE")

	case *types.Map:
		w.WriteString("NILMAP")

	case *types.Chan:
		w.WriteString("NILCHAN")

	case *types.Pointer, *types.Basic, *types.Signature:
		// Pointer or unsafe.Pointer
		w.WriteString("nil")

	case *types.Interface:
		w.WriteByte('(')
		cdd.Type(w, t)
		w.WriteString("){}")

	default:
		w.WriteString("'unknown nil")
	}
}

var unil = types.Typ[types.UntypedNil]

func eq(w *bytes.Buffer, lhs, op, rhs string, ltyp, rtyp types.Type) {
	typ := ltyp
	if typ == unil {
		typ = rtyp
	}
	switch t := typ.Underlying().(type) {
	case *types.Slice:
		nilv := "nil"
		sel := ".arr"
		if rtyp == unil {
			lhs += sel
			rhs = nilv
		} else {
			lhs = nilv
			rhs += sel
		}
	case *types.Interface:
		nilv := "NILI"
		if op == "!=" {
			w.WriteByte('!')
		}
		if rtyp == unil {
			rhs = nilv
		} else if ltyp == unil {
			lhs = nilv
		}
		w.WriteString("EQUALI(" + lhs + ", " + rhs + ")")
		return
	case *types.Basic:
		if t.Kind() != types.String {
			break
		}
		if op == "!=" {
			w.WriteByte('!')
		}
		w.WriteString("equals(" + lhs + ", " + rhs + ")")
		return
	}
	w.WriteString("(" + lhs + " " + op + " " + rhs + ")")
}

func (cdd *CDD) interfaceES(w *bytes.Buffer, es string, epos token.Pos, etyp, ityp types.Type) {
	if ityp == nil || etyp == nil || types.Identical(ityp, etyp) {
		w.WriteString(es)
		return
	}
	if _, ok := ityp.Underlying().(*types.Interface); !ok {
		// Result isn't an interface.
		w.WriteString(es)
		return
	}
	if b, ok := (etyp).(*types.Basic); ok && b.Kind() == types.UntypedNil {
		// Expr. result is nil interface.
		w.WriteString(es)
		return
	}
	iempty := (cdd.gtc.methodSet(ityp).Len() == 0)
	if _, ok := ityp.(*types.Interface); ok && !iempty {
		cdd.exit(epos, "not supported assignment to non-empty unnamed interface")
	}
	if _, ok := etyp.Underlying().(*types.Interface); ok {
		eempty := (cdd.gtc.methodSet(etyp).Len() == 0)
		switch {
		case iempty && eempty:
			w.WriteString(es)
		case iempty:
			w.WriteString("ICONVERTIE(" + es + ")")
		case eempty:
			w.WriteString("ICONVERTEI(" + es + ",  " + cdd.tinameDU(ityp) + ")")
		default:
			w.WriteString("ICONVERTII(" + es + ",  " + cdd.tinameDU(ityp) + ")")
		}
		return
	}
	if cdd.gtc.siz.Sizeof(etyp) > cdd.gtc.sizIval {
		cdd.exit(
			epos,
			"value of type %v is too large to assign to interface variable",
			etyp,
		)
	}
	if iempty {
		w.WriteString("INTERFACE(" + es + ", &" + cdd.tinameDU(etyp) + ")")
	} else {
		ets, its := cdd.tinameDU(etyp), cdd.tinameDU(ityp)
		w.WriteString("IASSIGN(" + es + ", " + ets + ", " + its + ")")
	}
}

func (cdd *CDD) interfaceExpr(w *bytes.Buffer, expr ast.Expr, ityp types.Type) {
	es := cdd.ExprStr(expr, ityp)
	etyp := cdd.exprType(expr)
	epos := expr.Pos()
	cdd.interfaceES(w, es, epos, etyp, ityp)
}

func (cdd *CDD) interfaceExprStr(expr ast.Expr, ityp types.Type) string {
	buf := new(bytes.Buffer)
	cdd.interfaceExpr(buf, expr, ityp)
	return buf.String()
}
