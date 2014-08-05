package gotoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"code.google.com/p/go.tools/go/exact"
	"code.google.com/p/go.tools/go/types"
)

func writeInt(w *bytes.Buffer, ev exact.Value, k types.BasicKind) {
	if k == types.Uintptr {
		u, _ := exact.Uint64Val(ev)
		w.WriteString("0x")
		w.WriteString(strconv.FormatUint(u, 16))
		return
	}

	w.WriteString(ev.String())
	switch k {
	case types.Int32:
		w.WriteByte('L')

	case types.Uint32:
		w.WriteString("UL")

	case types.Int64:
		w.WriteString("LL")

	case types.Uint64:
		w.WriteString("ULL")
	}
}

func writeFloat(w *bytes.Buffer, ev exact.Value, k types.BasicKind) {
	v, _ := exact.Int64Val(exact.Num(ev))
	w.WriteString(strconv.FormatInt(v, 10))
	v, _ = exact.Int64Val(exact.Denom(ev))
	if v != 1 {
		w.WriteByte('/')
		w.WriteString(strconv.FormatInt(v, 10))
	}
	w.WriteByte('.')
	if k == types.Float32 {
		w.WriteByte('F')
	}
}

func (cdd *CDD) Value(w *bytes.Buffer, ev exact.Value, t types.Type) {
	k := t.Underlying().(*types.Basic).Kind()

	// TODO: use t instead ev.Kind() in following switch
	switch ev.Kind() {
	case exact.Int:
		writeInt(w, ev, k)

	case exact.Float:
		writeFloat(w, ev, k)

	case exact.Complex:
		switch k {
		case types.Complex64:
			k = types.Float32
		case types.Complex128:
			k = types.Float64
		default:
			k = types.UntypedFloat
		}
		writeFloat(w, exact.Real(ev), k)
		im := exact.Imag(ev)
		if exact.Sign(im) != -1 {
			w.WriteByte('+')
		}
		writeFloat(w, im, k)
		w.WriteByte('i')

	case exact.String:
		w.WriteString("EGSTR(")
		w.WriteString(ev.String())
		w.WriteByte(')')

	default:
		w.WriteString(ev.String())
	}
}

func (cdd *CDD) Name(w *bytes.Buffer, obj types.Object, direct bool) {
	if obj == nil {
		w.WriteByte('_')
		return
	}
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

func (cdd *CDD) SelectorExpr(w *bytes.Buffer, e *ast.SelectorExpr) (fun types.Type, recv ast.Expr) {
	sel := cdd.gtc.ti.Selections[e]
	if sel == nil {
		cdd.Name(w, cdd.object(e.Sel), true)
		return
	}
	switch sel.Kind() {
	case types.FieldVal:
		cdd.Expr(w, e.X, nil)
		if sel.Indirect() {
			w.WriteString("->")
		} else {
			w.WriteByte('.')
		}
		w.WriteString(e.Sel.Name)

	case types.MethodVal:
		fun = sel.Obj().Type()
		cdd.Name(w, sel.Obj(), true)
		rtyp := fun.(*types.Signature).Recv().Type()
		if _, ok := rtyp.(*types.Pointer); ok {
			// Method with pointer receiver.
			if sel.Indirect() {
				recv = e.X
			} else {
				recv = &ast.UnaryExpr{Op: token.AND, X: e.X}
			}
		} else {
			// Method with non-pointer receiver.
			if sel.Indirect() {
				recv = &ast.UnaryExpr{Op: token.MUL, X: e.X}
			} else {
				recv = e.X
			}
		}

	case types.MethodExpr:
		cdd.Name(w, sel.Obj(), true)

	default:
		notImplemented(e)
	}
	return
}

func (cdd *CDD) SelectorExprStr(e *ast.SelectorExpr) (s string, fun types.Type, recv ast.Expr) {
	buf := new(bytes.Buffer)
	fun, recv = cdd.SelectorExpr(buf, e)
	s = buf.String()
	return
}

func (cdd *CDD) builtin(b *types.Builtin, args []ast.Expr) (fun, recv string) {
	name := b.Name()

	switch name {
	case "len":
		switch t := underlying(cdd.exprType(args[0])).(type) {
		case *types.Slice, *types.Basic: // Basic == String
			return "len", ""

		case *types.Array:
			return "", strconv.FormatInt(t.Len(), 10)

		case *types.Chan:
			return "clen", ""

		default:
			notImplemented(ast.NewIdent("len"), t)
		}

	case "cap":
		switch t := underlying(cdd.exprType(args[0])).(type) {
		case *types.Slice:
			return "cap", ""

		case *types.Chan:
			return "ccap", ""

		default:
			notImplemented(ast.NewIdent("cap"), t)
		}

	case "copy":
		switch t := underlying(cdd.exprType(args[1])).(type) {
		case *types.Basic: // string
			return "STRCPY", ""

		case *types.Slice:
			typ, dim, _ := cdd.TypeStr(t.Elem())
			return "SLICPY", typ + dimFuncPtr("", dim)

		default:
			panic(t)
		}

	case "new":
		typ, dim, _ := cdd.TypeStr(cdd.exprType(args[0]))
		args[0] = nil
		return "NEW", typ + dimFuncPtr("", dim)

	case "make":
		a0t := cdd.exprType(args[0])
		args[0] = nil

		switch t := underlying(a0t).(type) {
		case *types.Slice:
			typ, dim, _ := cdd.TypeStr(t.Elem())
			name := "MAKESLI"
			if len(args) == 3 {
				name = "MAKESLIC"
			}
			return name, typ + dimFuncPtr("", dim)

		case *types.Chan:
			typ, dim, _ := cdd.TypeStr(t.Elem())
			typ += dimFuncPtr("", dim)
			if len(args) == 1 {
				typ += ", 0"
			}
			return "MAKECHAN", typ

		case *types.Map:
			typ, dim, _ := cdd.TypeStr(t.Key())
			k := typ + dimFuncPtr("", dim)
			typ, dim, _ = cdd.TypeStr(t.Elem())
			e := typ + dimFuncPtr("", dim)
			name := "MAKEMAP"
			if len(args) == 2 {
				name = "MAKEMAPC"
			}
			return name, k + ", " + e

		default:
			notImplemented(ast.NewIdent(name))
		}

	}

	return name, ""
}

func (cdd *CDD) funStr(fe ast.Expr, args []ast.Expr) (fs string, ft types.Type, rs string, re ast.Expr) {
	switch f := fe.(type) {
	case *ast.SelectorExpr:
		buf := new(bytes.Buffer)
		ft, re = cdd.SelectorExpr(buf, f)
		fs = buf.String()
		if re != nil {
			rs = cdd.ExprStr(re, nil)
		}
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
	switch t := underlying(cdd.exprType(e.Fun)).(type) {
	case *types.Signature:
		fun, _, recv, _ := cdd.funStr(e.Fun, e.Args)
		if fun == "" {
			w.WriteString(recv)
		}
		w.WriteString(fun)
		w.WriteByte('(')
		comma := false
		if recv != "" {
			w.WriteString(recv)
			comma = true
		}
		tup := t.Params()
		for i, a := range e.Args {
			if a == nil {
				// builtin can set type args to nil
				continue
			}
			if comma {
				w.WriteString(", ")
			} else {
				comma = true
			}
			var at types.Type
			// Builtin functions may not spefify type for all parameters.
			if i < tup.Len() {
				at = tup.At(i).Type()
			}
			cdd.Expr(w, a, at)
			i++
		}
		w.WriteByte(')')

	default:
		arg := e.Args[0]
		switch typ := t.(type) {
		case *types.Slice:
			switch underlying(cdd.exprType(arg)).(type) {
			case *types.Basic: // string
				w.WriteString("NEWSTR(")
				cdd.Expr(w, arg, typ)
				w.WriteByte(')')

			default: // slice
				w.WriteByte('(')
				cdd.Expr(w, arg, typ)
				w.WriteByte(')')
			}

		default:
			w.WriteString("((")
			dim, _ := cdd.Type(w, typ)
			w.WriteString(dimFuncPtr("", dim))
			w.WriteString(")(")
			cdd.Expr(w, arg, typ)
			w.WriteString("))")
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
		op := e.Op.String()

		if op == "<-" {
			t := cdd.exprType(e.X).(*types.Chan).Elem()
			if tup, ok := cdd.exprType(e).(*types.Tuple); ok {
				tn, _, _ := cdd.tupleName(tup)
				w.WriteString("RECVOK(" + tn + ", ")
				cdd.Expr(w, e.X, nil)
				w.WriteByte(')')
			} else {
				w.WriteString("RECV(")
				dim, _ := cdd.Type(w, t)
				w.WriteString(dimFuncPtr("", dim))
				w.WriteString(", ")
				cdd.Expr(w, e.X, nil)
				w.WriteString(", ")
				zeroVal(w, t)
				w.WriteByte(')')
			}
			break
		}

		if op == "^" {
			op = "~"
		}
		w.WriteString(op)
		cdd.Expr(w, e.X, nil)

	case *ast.CallExpr:
		cdd.CallExpr(w, e)

	case *ast.Ident:
		if e.Name == "nil" {
			cdd.Nil(w, nilT)
		} else {
			cdd.Name(w, cdd.object(e), true)
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
		cdd.Expr(w, e.Value, nilT)

	case *ast.ParenExpr:
		w.WriteByte('(')
		cdd.Expr(w, e.X, nilT)
		w.WriteByte(')')

	case *ast.SelectorExpr:
		s, fun, recv := cdd.SelectorExprStr(e)
		if recv == nil {
			w.WriteString(s)
			break
		}
		sig := fun.(*types.Signature)
		w.WriteString("({")
		res, params := cdd.signature(sig, false, numNames)
		w.WriteString(res.typ)
		w.WriteByte(' ')
		w.WriteString(dimFuncPtr("func"+params, res.dim))
		w.WriteString(" { return " + s + "(")
		cdd.Expr(w, recv, nil)
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
		notImplemented(e)

	case *ast.CompositeLit:
		typ := cdd.exprType(e)

		switch t := underlying(typ).(type) {
		case *types.Array:
			w.WriteByte('{')
			nilT = t.Elem()

		case *types.Slice:
			w.WriteString("(slice){(")
			dim, _ := cdd.Type(w, t.Elem())
			dim = append([]string{"[]"}, dim...)
			w.WriteString(dimFuncPtr("", dim))
			w.WriteString("){")
			nilT = t.Elem()

		case *types.Struct:
			w.WriteByte('(')
			cdd.Type(w, typ)
			w.WriteString("){")
			nilT = nil

		default:
			notImplemented(e, t)
		}

		for i, el := range e.Elts {
			if i > 0 {
				w.WriteString(", ")
			}
			if nilT != nil {
				cdd.Expr(w, el, nilT)
			} else {
				cdd.Expr(w, el, underlying(typ).(*types.Struct).Field(i).Type())
			}
		}

		switch underlying(typ).(type) {
		case *types.Slice:
			w.WriteByte('}')
			plen := ", " + strconv.Itoa(len(e.Elts))
			w.WriteString(plen)
			w.WriteString(plen)
			w.WriteByte('}')

		default:
			w.WriteByte('}')
		}

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
}

func (cdd *CDD) indexExpr(w *bytes.Buffer, typ types.Type, xs string, idx ast.Expr) {
	pt, isPtr := typ.(*types.Pointer)
	if isPtr {
		w.WriteString("(*")
		typ = pt.Elem()
	}

	var indT types.Type

	switch t := typ.Underlying().(type) {
	case *types.Basic: // string
		w.WriteString(xs + ".str")

	case *types.Slice:
		w.WriteString("((")
		dim, _ := cdd.Type(w, t.Elem())
		dim = append([]string{"*"}, dim...)
		w.WriteString(dimFuncPtr("", dim))
		w.WriteByte(')')
		w.WriteString(xs + ".arr)")

	case *types.Array:
		w.WriteString(xs)

	case *types.Map:
		indT = t.Key()
		notImplemented(&ast.IndexExpr{}, t)

	default:
		panic(t)
	}

	if isPtr {
		w.WriteByte(')')
	}
	w.WriteByte('[')
	cdd.Expr(w, idx, indT)
	w.WriteByte(']')
}

func (cdd *CDD) SliceExpr(w *bytes.Buffer, e *ast.SliceExpr) {
	sx := cdd.ExprStr(e.X, nil)

	typ := cdd.exprType(e.X)
	pt, isPtr := typ.(*types.Pointer)
	if isPtr {
		typ = pt.Elem()
		sx = "(*" + sx + ")"
	}

	switch t := typ.(type) {
	case *types.Slice:
		if e.Low == nil && e.High == nil && e.Max == nil {
			w.WriteString(sx)
			break
		}

		if e.Low != nil {
			switch {
			case e.High == nil && e.Max == nil:
				w.WriteString("SLICEL(")

			case e.High != nil && e.Max == nil:
				w.WriteString("SLICELH(")

			case e.High == nil && e.Max != nil:
				w.WriteString("SLICEM(")

			default:
				w.WriteString("SLICELHM(")
			}
			w.WriteString(sx)
			w.WriteString(", ")
			dim, _ := cdd.Type(w, t.Elem())
			dim = append([]string{"*"}, dim...)
			w.WriteString(dimFuncPtr("", dim))
			w.WriteString(", ")
			cdd.Expr(w, e.Low, nil)
		} else {
			switch {
			case e.High != nil && e.Max == nil:
				w.WriteString("SLICEH(")

			case e.High == nil && e.Max != nil:
				w.WriteString("SLICEM(")

			default:
				w.WriteString("SLICEHM(")
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

	case *types.Array:
		alen := strconv.FormatInt(t.Len(), 10) + ", "
		if e.Low != nil {
			switch {
			case e.High == nil && e.Max == nil:
				w.WriteString("ASLICEL(" + alen)

			case e.High != nil && e.Max == nil:
				w.WriteString("ASLICELH(" + alen)

			case e.High == nil && e.Max != nil:
				w.WriteString("ASLICEM(")

			default:
				w.WriteString("ASLICELHM(")
			}
			w.WriteString(sx)
			w.WriteString(", ")
			cdd.Expr(w, e.Low, nil)
		} else {
			switch {
			case e.High == nil && e.Max == nil:
				w.WriteString("ASLICE(" + alen)

			case e.High != nil && e.Max == nil:
				w.WriteString("ASLICEH(" + alen)

			case e.High == nil && e.Max != nil:
				w.WriteString("ASLICEM(")

			default:
				w.WriteString("ASLICEHM(")
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

	case *types.Basic: // string
		if e.Low == nil && e.High == nil {
			w.WriteString(sx)
			break
		}
		switch {
		case e.Low == nil:
			w.WriteString("SSLICEH(")

		case e.High == nil:
			w.WriteString("SSLICEL(")

		default:
			w.WriteString("SSLICELH(")
		}

		w.WriteString(sx)

		if e.Low != nil {
			w.WriteString(", ")
			cdd.Expr(w, e.Low, nil)
		}
		if e.High != nil {
			w.WriteString(", ")
			cdd.Expr(w, e.High, nil)
		}

		w.WriteByte(')')

	default:
		panic(e)
	}
}

func (cdd *CDD) ExprStr(expr ast.Expr, nilT types.Type) string {
	buf := new(bytes.Buffer)
	cdd.Expr(buf, expr, nilT)
	return buf.String()
}

func (cdd *CDD) Nil(w *bytes.Buffer, t types.Type) {
	switch underlying(t).(type) {
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
		w.WriteString("NILI")

	default:
		w.WriteString("'unknown nil")
	}
}

func eq(w *bytes.Buffer, lhs, op, rhs string, ltyp, rtyp types.Type) {
	typ := ltyp
	if typ == types.Typ[types.UntypedNil] {
		typ = rtyp
	}

	switch underlying(typ).(type) {
	case *types.Slice:
		if rtyp == types.Typ[types.UntypedNil] {
			lhs += ".arr"
			rhs = "nil"
		} else {
			lhs = "nil"
			rhs += ".arr"
		}
	}
	w.WriteString(lhs + " " + op + " " + rhs)
}
