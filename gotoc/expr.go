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
	case "_":
		w.WriteString("unused" + cdd.gtc.uniqueId())

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

	case types.PackageObj:
		cdd.Name(w, sel.Obj(), true)

	default: // types.MethodExpr
		notImplemented(e)
	}
	return
}

func (cdd *CDD) builtin(b *types.Builtin, args []ast.Expr) (fun, recv string) {
	name := b.Name()

	switch name {
	case "len":
		switch t := cdd.exprType(args[0]).(type) {
		case *types.Slice, *types.Map, *types.Basic: // Basic == String
			return "len", ""

		case *types.Array:
			return "sizeof", ""

		default:
			panic(t)
		}

	case "copy":
		switch t := cdd.exprType(args[1]).(type) {
		case *types.Basic: // string
			return "STRCPY", ""

		case *types.Slice:
			typ, dim, _ := cdd.TypeStr(t.Elem())
			return "SLICPY", typ + dimFuncPtr("", dim)

		default:
			panic(t)
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
	switch t := cdd.exprType(e.Fun).(type) {
	case *types.Signature:
		fun, _, recv, _ := cdd.funStr(e.Fun, e.Args)
		w.WriteString(fun)
		w.WriteByte('(')
		if recv != "" {
			w.WriteString(recv)
			if len(e.Args) > 0 {
				w.WriteString(", ")
			}
		}
		tup := t.Params()
		for i, a := range e.Args {
			if i != 0 {
				w.WriteString(", ")
			}
			cdd.Expr(w, a, tup.At(i).Type())
		}
		w.WriteByte(')')

	default:
		arg := e.Args[0]
		switch typ := cdd.exprType(e.Fun).(type) {
		case *types.Slice:
			switch cdd.exprType(arg).(type) {
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
		typ := cdd.exprType(e.X)

		pt, isPtr := typ.(*types.Pointer)
		if isPtr {
			w.WriteString("(*")
			typ = pt.Elem()
		}

		var indT types.Type

		switch t := typ.(type) {
		case *types.Basic: // string
			cdd.Expr(w, e.X, nil)
			w.WriteString(".str")

		case *types.Slice:
			w.WriteString("((")
			dim, _ := cdd.Type(w, t.Elem())
			dim = append([]string{"*"}, dim...)
			w.WriteString(dimFuncPtr("", dim))
			w.WriteByte(')')
			cdd.Expr(w, e.X, nil)
			w.WriteString(".arr)")

		case *types.Array:
			cdd.Expr(w, e.X, nil)

		case *types.Map:
			indT = t.Key()
			notImplemented(e)

		default:
			panic(t)
		}

		if isPtr {
			w.WriteByte(')')
		}

		w.WriteByte('[')
		cdd.Expr(w, e.Index, indT)
		w.WriteByte(']')

	case *ast.KeyValueExpr:
		w.WriteByte('.')
		kt := cdd.exprType(e.Key)
		if i, ok := e.Key.(*ast.Ident); ok && kt == nil {
			// e.Key is field name
			w.WriteString(i.Name)
		} else {
			cdd.Expr(w, e.Key, kt)
		}
		w.WriteString(" = ")
		cdd.Expr(w, e.Value, nilT)

	case *ast.ParenExpr:
		w.WriteByte('(')
		cdd.Expr(w, e.X, nilT)
		w.WriteByte(')')

	case *ast.SelectorExpr:
		cdd.SelectorExpr(w, e)

	case *ast.SliceExpr:
		cdd.SliceExpr(w, e)

	case *ast.StarExpr:
		w.WriteByte('*')
		cdd.Expr(w, e.X, nil)

	case *ast.TypeAssertExpr:
		notImplemented(e)

	case *ast.CompositeLit:
		typ := cdd.exprType(e)

		switch t := typ.(type) {
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

		default:
			w.WriteByte('(')
			cdd.Type(w, t)
			w.WriteString("){")
			nilT = nil
		}

		for i, el := range e.Elts {
			if i > 0 {
				w.WriteString(", ")
			}
			cdd.Expr(w, el, nilT)
		}

		switch typ.(type) {
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
			for u, typPtr := range c.BodyUses {
				cdd.BodyUses[u] = typPtr
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
		if e.Low != nil {
			switch {
			case e.High == nil && e.Max == nil:
				w.WriteString("ASLICEL(")

			case e.High != nil && e.Max == nil:
				w.WriteString("ASLICELH(")

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
				w.WriteString("ASLICE(")

			case e.High != nil && e.Max == nil:
				w.WriteString("ASLICEH(")

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
	switch t.(type) {
	case *types.Slice:
		w.WriteString("NILSLICE")

	case *types.Map:
		w.WriteString("NILMAP")

	case *types.Pointer:
		w.WriteString("nil")

	default:
		w.WriteString("{0}")
	}
}

func eq(w *bytes.Buffer, lhs, op, rhs string, ltyp, rtyp types.Type) {
	typ := ltyp
	if typ == types.Typ[types.UntypedNil] {
		typ = rtyp
	}

	switch typ.(type) {
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
