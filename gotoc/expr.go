package gotoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"strconv"
	"strings"
)

func writeInt(w *bytes.Buffer, ev constant.Value, k types.BasicKind, sizInt int64) {
	if k == types.Uintptr {
		u, _ := constant.Uint64Val(ev)
		w.WriteString("0x")
		w.WriteString(strconv.FormatUint(u, 16))
		return
	}
	s := ev.String()
	if s[0] == '-' {
		w.WriteByte('(')
	}
	switch {
	case k == types.Int32 || k == types.Int && sizInt == 4:
		if s == "-2147483648" {
			w.WriteString("-2147483647L-1L")
		} else {
			w.WriteString(s + "L")
		}
	case k == types.Uint32:
		w.WriteString(s + "UL")
	case k == types.Int64 || k == types.Int && sizInt == 8:
		if s == "-9223372036854775808" {
			w.WriteString("-9223372036854775807LL-1LL")
		} else {
			w.WriteString(s + "LL")
		}
	case k == types.Uint64:
		w.WriteString(s + "ULL")
	default:
		w.WriteString(s)
	}
	if s[0] == '-' {
		w.WriteByte(')')
	}
}

func writeFloat(w *bytes.Buffer, ev constant.Value, k types.BasicKind) {
	if k == types.Float32 {
		f, _ := constant.Float32Val(ev)
		w.WriteString(strconv.FormatFloat(float64(f), 'e', -1, 32))
		w.WriteByte('F')
	} else {
		f, _ := constant.Float64Val(ev)
		w.WriteString(strconv.FormatFloat(f, 'e', -1, 64))
	}
}

func (cdd *CDD) Value(w *bytes.Buffer, ev constant.Value, t types.Type) {
	if o, ok := t.(*types.Named); ok {
		cdd.addObject(o.Obj(), false)
	}
	k := t.Underlying().(*types.Basic).Kind()
	switch {
	case k <= types.Bool || k == types.UntypedBool:
		w.WriteString(ev.String())
	case k <= types.Uintptr || k == types.UntypedInt || k == types.UntypedRune:
		writeInt(w, ev, k, cdd.gtc.sizInt)
	case k <= types.Float64 || k == types.UntypedFloat:
		writeFloat(w, ev, k)
	case k <= types.Complex128 || k == types.UntypedComplex:
		w.WriteByte('(')
		writeFloat(w, constant.Real(ev), k)
		im := constant.Imag(ev)
		if constant.Sign(im) != -1 {
			w.WriteByte('+')
		}
		writeFloat(w, im, k)
		w.WriteString("i)")
	case k == types.String || k == types.UntypedString:
		if cdd.constInit {
			w.WriteString("EGSTR(")
		} else {
			w.WriteString("EGSTL(")
		}
		str := ev.ExactString()
		// Handle problem with infinite hexadecimal escapes in C.
		for {
			i := strings.Index(str, `\x`)
			if i == -1 || len(str) <= i+4 {
				break
			}
			if i > 0 && str[i-1] == '\\' {
				w.WriteString(str[:i+1])
				str = str[i+1:]
				continue
			}
			d := str[i+4]
			if d >= '0' && d <= '9' || d >= 'a' && d <= 'f' ||
				d >= 'A' && d <= 'F' {

				w.WriteString(str[:i+4])
				w.WriteString(`""`) // Separate escape from subsequent char.
				str = str[i+4:]
				continue
			}
			w.WriteString(str[:i+5])
			str = str[i+5:]
		}
		w.WriteString(str)
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
		w.WriteString(Upath(o.Pkg().Path()))
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
	var hasPrefix bool
	if p := obj.Pkg(); p != nil && !cdd.gtc.isLocal(obj) {
		cdd.addObject(obj, direct)
		w.WriteString(Upath(obj.Pkg().Path()))
		w.WriteByte('$')
		hasPrefix = true
	}
	name := obj.Name()
	if name == "init" {
		if _, ok := obj.(*types.Func); ok && hasPrefix {
			w.WriteString(cdd.gtc.uniqueId() + name)
			return
		}
	}
	w.WriteString(name)
	if !hasPrefix && name != "error" {
		w.WriteByte('$')
	}
}

func (cdd *CDD) NameStr(o types.Object, direct bool) string {
	buf := new(bytes.Buffer)
	cdd.Name(buf, o, direct)
	return buf.String()
}

func (cdd *CDD) SelectorExpr(w *bytes.Buffer, e *ast.SelectorExpr, permitaa bool) (fun, recvt types.Type, recvs string) {
	sel := cdd.gtc.ti.Selections[e]
	if sel == nil {
		cdd.Name(w, cdd.object(e.Sel), true)
		return
	}
	s := cdd.ExprStr(e.X, nil, permitaa)
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

func (cdd *CDD) SelectorExprStr(e *ast.SelectorExpr, permitaa bool) (s string, fun, recvt types.Type, recvs string) {
	buf := new(bytes.Buffer)
	fun, recvt, recvs = cdd.SelectorExpr(buf, e, permitaa)
	s = buf.String()
	return
}

func (cdd *CDD) builtin(b *types.Builtin, args []ast.Expr) (fun, recv string) {
	name := b.Name()

	switch name {
	case "len":
		switch t := cdd.exprType(args[0]).Underlying().(type) {
		case *types.Slice, *types.Basic: // Basic == String
			cdd.Complexity--
			return "len", ""

		case *types.Array:
			cdd.Complexity--
			panic("builtin len(array) isn't handled as constant expression")
			// return "", strconv.FormatInt(t.Len(), 10)

		case *types.Chan:
			return "clen", ""

		default:
			cdd.notImplemented(ast.NewIdent("len"), t)
		}

	case "cap":
		cdd.Complexity--
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
		ft, rt, rs = cdd.SelectorExpr(buf, f, false)
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
	fs = cdd.ExprStr(fe, nil, false)
	ft = cdd.exprType(fe)
	return
}

func (cdd *CDD) CallExpr(w *bytes.Buffer, e *ast.CallExpr, permitaa bool) {
	switch t := cdd.exprType(e.Fun).(type) {
	case *types.Signature:
		c := cdd.call(e, t, false)
		if c.rcv.r != "" || c.arr.r != "" || c.tup.t != nil {
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
		if c.tup.t != nil {
			cdd.Type(w, c.tup.t)
			w.WriteString(" " + c.tup.l + " = " + indent(1, c.tup.r) + ";\n")
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
		if c.rcv.r != "" || c.arr.r != "" || c.tup.t != nil {
			w.WriteString(";\n")
			cdd.il--
			cdd.indent(w)
			w.WriteString("})")
		}

	default:
		cdd.Complexity--
		arg := e.Args[0]
		at := cdd.exprType(arg)
		switch typ := t.Underlying().(type) {
		case *types.Slice:
			switch at.Underlying().(type) {
			case *types.Basic: // string
				w.WriteString("BYTES(")
				cdd.Expr(w, arg, typ, true)
				w.WriteByte(')')
				cdd.Complexity += 1e3 // BUG: Workaround because alloca in BYTES.
			default: // slice
				w.WriteByte('(')
				cdd.Expr(w, arg, typ, permitaa)
				w.WriteByte(')')
			}

		case *types.Interface:
			cdd.interfaceExpr(w, arg, t, permitaa)

		default:
			if b, ok := typ.(*types.Basic); ok && b.Kind() == types.String {
				if _, ok := at.(*types.Slice); ok {
					// string(bytes)
					w.WriteString("NEWSTR(")
					cdd.Expr(w, arg, typ, true)
					w.WriteByte(')')
					break
				}
			}
			/*
				// Not need because -fno-strict-aliasing
				if _, ok := typ.(*types.Pointer); ok {
					if _, ok := at.Underlying().(*types.Pointer); !ok {
						// Casting unsafe.Pointer
						// TODO: Don't use CAST if can check that type under
						// unsafe.Pointer is pointer that points to type of the same
						// size. This will allow clasify such cast as constant.
						w.WriteString("CAST(")
						dim := cdd.Type(w, t)
						w.WriteString(dimFuncPtr("", dim))
						w.WriteString(", ")
						cdd.Expr(w, arg, t)
						w.WriteByte(')')
						break
					}
				}
			*/
			w.WriteString("((")
			dim := cdd.Type(w, t)
			w.WriteString(dimFuncPtr("", dim))
			w.WriteString(")(")
			cdd.Expr(w, arg, t, permitaa)
			w.WriteString("))")
		}
	}
}

func (cdd *CDD) Expr(w *bytes.Buffer, expr ast.Expr, nilT types.Type, permitaa bool) {
	if t := cdd.gtc.ti.Types[expr]; t.Value != nil {
		// Constant expression
		cdd.Value(w, t.Value, t.Type)
		return
	}

	cdd.Complexity++

	switch e := expr.(type) {
	case *ast.BinaryExpr:
		op := e.Op.String()
		ltyp := cdd.exprType(e.X)
		rtyp := cdd.exprType(e.Y)

		lhs := cdd.ExprStr(e.X, ltyp, permitaa)
		rhs := cdd.ExprStr(e.Y, rtyp, permitaa)

		if t, ok := ltyp.Underlying().(*types.Basic); ok && t.Kind() == types.String {
			w.WriteString("(cmpstr(" + lhs + ", " + rhs + ") " + op + " 0)")
			break
		}
		if e.Op == token.EQL || e.Op == token.NEQ {
			cdd.eq(w, lhs, op, rhs, ltyp, rtyp)
			break
		}
		switch e.Op {
		case token.AND_NOT:
			op = "&~"
			fallthrough
		case token.AND, token.OR, token.XOR, token.SHL, token.SHR:
			w.WriteByte('(')
			cdd.Type(w, cdd.exprType(e.X))
			w.WriteByte(')')
		}
		w.WriteString("(" + lhs + op + rhs + ")")

	case *ast.UnaryExpr:
		if e.Op == token.ARROW {
			t := cdd.exprType(e.X).(*types.Chan).Elem()
			if tup, ok := cdd.exprType(e).(*types.Tuple); ok {
				tn, _ := cdd.tupleName(tup)
				w.WriteString("RECVOK(" + tn + ", ")
				cdd.Expr(w, e.X, nil, true)
				w.WriteByte(')')
			} else {
				w.WriteString("RECV(")
				dim := cdd.Type(w, t)
				w.WriteString(dimFuncPtr("", dim))
				w.WriteString(", ")
				cdd.Expr(w, e.X, nil, true)
				w.WriteString(", ")
				zeroVal(w, t)
				w.WriteByte(')')
			}
			break
		}
		if e.Op == token.AND {
			cdd.Complexity--
			cdd.ptrExpr(w, e.X, permitaa)
			break
		}
		if e.Op == token.XOR {
			w.WriteByte('(')
			cdd.Type(w, cdd.exprType(e.X))
			w.WriteString(")(~")
			cdd.Expr(w, e.X, nil, permitaa)
			w.WriteByte(')')
			break
		}
		w.WriteString(e.Op.String())
		cdd.Expr(w, e.X, nil, permitaa)

	case *ast.CallExpr:
		cdd.CallExpr(w, e, permitaa)

	case *ast.Ident:
		cdd.Complexity--
		if e.Name == "nil" {
			cdd.Nil(w, nilT)
			break
		}
		if o := cdd.object(e); o != nil {
			if _, ok := o.(*types.Func); ok {
				w.WriteByte('&')
			}
			cdd.Name(w, o, true)
			break
		}
		w.WriteString(e.Name)
		if e.Name != "_" {
			w.WriteByte('$')
		}

	case *ast.IndexExpr:
		cdd.indexExpr(w, cdd.exprType(e.X), cdd.ExprStr(e.X, nil, false), e.Index, "")

	case *ast.KeyValueExpr:
		kt := cdd.exprType(e.Key)
		if i, ok := e.Key.(*ast.Ident); ok && kt == nil {
			// e.Key is field name
			w.WriteByte('.')
			w.WriteString(i.Name)
		} else {
			w.WriteByte('[')
			cdd.Expr(w, e.Key, kt, true)
			w.WriteByte(']')
		}
		w.WriteString(" = ")
		cdd.interfaceExpr(w, e.Value, nilT, permitaa)

	case *ast.ParenExpr:
		cdd.Complexity--
		w.WriteByte('(')
		cdd.Expr(w, e.X, nilT, permitaa)
		w.WriteByte(')')

	case *ast.SelectorExpr:
		/*
				Added by commit bc72e116df48bb89f97284a6d7203850287465fa but don't
				remember why!?

				if o := cdd.object(e.Sel); o != nil {
				if _, ok := o.(*types.Func); ok {
					w.WriteByte('&')
				}
			}
		*/
		s, fun, recvt, recvs := cdd.SelectorExprStr(e, permitaa)
		if recvt == nil {
			w.WriteString(s)
			cdd.Complexity--
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
		cdd.Expr(w, e.X, nil, permitaa)

	case *ast.TypeAssertExpr:
		w.WriteString("({\n")
		cdd.il++
		ityp := cdd.exprType(e.X)
		cdd.indent(w)
		cdd.varDecl(w, ityp, "_i", e.X, "", false, permitaa)
		w.WriteByte('\n')
		cdd.indent(w)
		iempty := (cdd.gtc.methodSet(ityp).Len() == 0)
		typ := cdd.exprType(e.Type)
		etup, _ := cdd.exprType(e).(*types.Tuple)
		if _, ok := typ.Underlying().(*types.Interface); ok {
			if etup != nil {
				tn, _ := cdd.tupleName(etup)
				w.WriteString(tn + " _ret = {};\n")
				cdd.indent(w)
				w.WriteString("_ret._1 = ")
			} else {
				w.WriteString("if (!")
			}
			w.WriteString("implements(")
			if iempty {
				w.WriteString("_i.itab, &")
			} else {
				w.WriteString("TINFO(_i), &")
			}
			w.WriteString(cdd.tinameDU(typ))
			w.WriteByte(')')
			if etup != nil {
				w.WriteString(";\n")
				cdd.indent(w)
				w.WriteString("if (_ret._1) _ret._0 = ")
				cdd.interfaceES(w, nil, "_i", e.Pos(), ityp, typ, permitaa)
				w.WriteString(";\n")
				cdd.indent(w)
				w.WriteString("_ret;\n")
			} else {
				w.WriteString(") panicIC();\n")
				cdd.indent(w)
				cdd.interfaceES(w, nil, "_i", e.Pos(), ityp, typ, permitaa)
				w.WriteString(";\n")
			}
		} else {
			if etup != nil {
				w.WriteString("bool _ok = ")
			} else {
				w.WriteString("if (!")
			}
			if iempty {
				w.WriteString("(_i.itab == &")
			} else {
				w.WriteString("(TINFO(_i) == &")
			}
			w.WriteString(cdd.tinameDU(typ))
			w.WriteByte(')')
			if etup != nil {
				w.WriteString(";\n")
				tn, _ := cdd.tupleName(etup)
				cdd.indent(w)
				w.WriteString("(" + tn + "){IVAL(_i, ")
				dim := cdd.Type(w, typ)
				w.WriteString(dimFuncPtr("", dim))
				w.WriteString("), _ok};\n")
			} else {
				w.WriteString(") panicIC();\n")
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
		cdd.compositeLit(w, e, nilT, permitaa)

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
			for u, typPtr := range c.DefUses {
				cdd.DefUses[u] = typPtr
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

func (cdd *CDD) newVar(name string, typ types.Type, global bool, val ast.Expr, permitaa bool) {
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
	acd.varDecl(new(bytes.Buffer), typ, name, val, "", false, permitaa)
}

func (cdd *CDD) ptrExpr(w *bytes.Buffer, e ast.Expr, permitaa bool) {
	_, iscl := e.(*ast.CompositeLit)
	if iscl && !permitaa {
		cdd.exit(
			e.Pos(), "can not use pointer to composite literal in this context",
		)
	}
	w.WriteByte('&')
	if !iscl || cdd.Typ != VarDecl || !cdd.gtc.isGlobal(cdd.Origin) {
		cdd.Expr(w, e, nil, permitaa)
		return
	}
	name := "_cl" + cdd.gtc.uniqueId()
	cdd.newVar(name, cdd.exprType(e), true, e, true)
	w.WriteString(name)
}

func (cdd *CDD) compositeLit(w *bytes.Buffer, e *ast.CompositeLit, nilT types.Type, permitaa bool) {
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
		cdd.elts(w, e.Elts, t.Elem(), nil, permitaa)
		w.WriteString("}}")
		if !cdd.constInit {
			w.WriteByte(')')
		}

	case *types.Slice:
		if !permitaa {
			cdd.exit(e.Pos(), "can not use slice literal in this context")
		}
		alen := cdd.clArrayLen(e.Elts)
		slen := strconv.FormatInt(alen, 10)
		if cdd.constInit || (cdd.Typ != VarDecl) || !cdd.gtc.isGlobal(cdd.Origin) {
			w.WriteString("CSLICE(" + slen + ", ((")
			dim := cdd.Type(w, t.Elem())
			dim = append([]string{"[]"}, dim...)
			w.WriteString(dimFuncPtr("", dim))
			w.WriteString("){")
			cdd.elts(w, e.Elts, t.Elem(), nil, true)
			w.WriteString("}))")
			break
		}
		aname := "_cl" + cdd.gtc.uniqueId()
		w.WriteString("ASLICE(&" + aname + ")")
		typ := types.NewArray(t.Elem(), alen)
		tv := cdd.gtc.ti.Types[e]
		tv.Type = typ
		cdd.gtc.ti.Types[e] = tv
		cdd.newVar(aname, typ, true, e, true)

	case *types.Struct:
		if !cdd.constInit {
			w.WriteString("((")
			cdd.Type(w, typ)
			w.WriteByte(')')
		}
		w.WriteByte('{')
		cdd.elts(w, e.Elts, nil, t, permitaa)
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
			k, _ = constant.Int64Val(cdd.gtc.exprValue(kv.Key))

		}
		if k++; k > n {
			n = k
		}
	}
	return n
}

func (cdd *CDD) elts(w *bytes.Buffer, elts []ast.Expr, nilT types.Type, st *types.Struct, permitaa bool) {
	if len(elts) == 0 {
		return
	}
	if nilT != nil {
		for i, el := range elts {
			if i > 0 {
				w.WriteString(", ")
			}
			cdd.interfaceExpr(w, el, nilT, permitaa)
		}
		return
	}
	// Struct
	if _, ok := elts[0].(*ast.KeyValueExpr); !ok {
		for i, el := range elts {
			if i > 0 {
				w.WriteString(", ")
			}
			cdd.interfaceExpr(w, el, st.Field(i).Type(), permitaa)
		}
		return
	}
	for i, el := range elts {
		if i > 0 {
			w.WriteString(", ")
		}
		key := el.(*ast.KeyValueExpr).Key.(*ast.Ident).Name
		for k := 0; k < st.NumFields(); k++ {
			f := st.Field(k)
			if st.Field(k).Name() == key {
				cdd.interfaceExpr(w, el, f.Type(), permitaa)
				break
			}
		}
	}
}

// indexExpr assumes that if idx == nil then ids is checked index.
func (cdd *CDD) indexExpr(w *bytes.Buffer, typ types.Type, xs string, idx ast.Expr, ids string) {
	pt, isPtr := typ.(*types.Pointer)
	if isPtr {
		typ = pt.Elem()
	}
	var indT types.Type

	switch t := typ.Underlying().(type) {
	case *types.Basic: // string
		if cdd.gtc.boundsCheck && idx != nil {
			w.WriteString("STRIDXC(")
		} else {
			w.WriteString("STRIDX(")
		}
	case *types.Slice:
		if cdd.gtc.boundsCheck && idx != nil {
			w.WriteString("SLIDXC(")
		} else {
			w.WriteString("SLIDX(")
		}
		dim := cdd.Type(w, t.Elem())
		dim = append([]string{"*"}, dim...)
		w.WriteString(dimFuncPtr("", dim))
		w.WriteString(", ")
	case *types.Array:
		if cdd.gtc.boundsCheck && idx != nil &&
			!cdd.isConstExpr(idx, types.Typ[types.UntypedInt]) {
			w.WriteString("AIDXC(")
		} else {
			w.WriteString("AIDX(")
		}
		if !isPtr {
			w.WriteByte('&')
		}
	case *types.Map:
		indT = t.Key()
		cdd.notImplemented(&ast.IndexExpr{}, t)
	default:
		panic(t)
	}
	w.WriteString(xs)
	w.WriteString(", ")
	if idx != nil {
		cdd.Expr(w, idx, indT, true)
	} else {
		w.WriteString(ids)
	}
	w.WriteByte(')')
}

// indexExprStr assumes that if idx == nil then ids is checked index.
func (cdd *CDD) indexExprStr(typ types.Type, xs string, idx ast.Expr, ids string) string {
	buf := new(bytes.Buffer)
	cdd.indexExpr(buf, typ, xs, idx, ids)
	return buf.String()
}

func (cdd *CDD) SliceExpr(w *bytes.Buffer, e *ast.SliceExpr) {
	sx := cdd.ExprStr(e.X, nil, false)

	typ := cdd.exprType(e.X)
	pt, isPtr := typ.(*types.Pointer)
	if isPtr {
		typ = pt.Elem()
	}
	var slice *types.Slice
	empty := e.Low == nil && e.High == nil && e.Max == nil
	switch t := typ.Underlying().(type) {
	case *types.Array:
		if !isPtr {
			sx = "&" + sx
		}
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
		cdd.Expr(w, e.Low, nil, true)
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
		cdd.Expr(w, e.High, nil, true)
	}
	if e.Max != nil {
		w.WriteString(", ")
		cdd.Expr(w, e.Max, nil, true)
	}
	w.WriteByte(')')
}

func (cdd *CDD) ExprStr(expr ast.Expr, nilT types.Type, permitaa bool) string {
	buf := new(bytes.Buffer)
	cdd.Expr(buf, expr, nilT, permitaa)
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

func (cdd *CDD) eq(w *bytes.Buffer, lhs, op, rhs string, ltyp, rtyp types.Type) {
	typ := ltyp
	if typ == unil {
		typ = rtyp
	}
	switch t := typ.Underlying().(type) {
	case *types.Interface:
		if op == "!=" {
			w.WriteByte('!')
		}
		if rtyp == unil {
			w.WriteString("ISNILI(" + lhs + ")")
		} else if ltyp == unil {
			w.WriteString("ISNILI(" + rhs + ")")
		} else {
			_, li := ltyp.Underlying().(*types.Interface)
			_, ri := rtyp.Underlying().(*types.Interface)
			if li && ri {
				w.WriteString("EQUALI(" + lhs + ", " + rhs + ")")
			} else if li {
				panic("TODO")
			} else {
				panic("TODO")
			}
		}
		return
	case *types.Basic:
		if t.Kind() != types.String {
			break
		}
		w.WriteString("(cmpstr(" + lhs + ", " + rhs + ") " + op + " 0)")
		return
	case *types.Struct:
		lo := " &&\n"
		if op == "!=" {
			lo = " ||\n"
		}
		w.WriteString("({\n")
		cdd.il++
		cdd.indent(w)
		cdd.Type(w, ltyp)
		id := cdd.gtc.uniqueId()
		lv := "_l" + id
		rv := "_r" + id
		w.WriteString(" " + lv + " = " + lhs + "; ")
		cdd.Type(w, rtyp)
		w.WriteString(" " + rv + " = " + rhs + ";\n")
		n := t.NumFields()
		for i := 0; i < n; {
			f := t.Field(i)
			ft := f.Type()
			fn := f.Name()
			cdd.indent(w)
			cdd.eq(w, lv+"."+fn, op, rv+"."+fn, ft, ft)
			i++
			if i != n {
				w.WriteString(lo)
			} else {
				w.WriteString(";\n")
			}
		}
		cdd.il--
		cdd.indent(w)
		w.WriteString("})")
		return
	case *types.Array:
		if op == "!=" {
			w.WriteByte('!')
		}
		w.WriteString("EQUALA(" + lhs + ", " + rhs + ")")
		return
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
	}
	w.WriteString("(" + lhs + " " + op + " " + rhs + ")")
}

func (cdd *CDD) interfaceES(w *bytes.Buffer, ex ast.Expr, es string, epos token.Pos, etyp, ityp types.Type, permitaa bool) {
	simple := ityp == nil || etyp == nil || types.Identical(ityp, etyp)
	if !simple {
		_, ok := ityp.Underlying().(*types.Interface)
		simple = !ok
	}
	if !simple {
		b, ok := (etyp).(*types.Basic)
		simple = ok && b.Kind() == types.UntypedNil
	}
	if simple {
		if ex != nil {
			cdd.Expr(w, ex, ityp, permitaa)
		} else {
			w.WriteString(es)
		}
		return
	}
	iempty := (cdd.gtc.methodSet(ityp).Len() == 0)
	if _, ok := ityp.(*types.Interface); ok && !iempty {
		cdd.exit(
			epos, "not supported assignment to non-empty unnamed interface",
		)
	}
	if ex != nil {
		es = cdd.ExprStr(ex, ityp, false)
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

func (cdd *CDD) interfaceESstr(ex ast.Expr, es string, epos token.Pos, etyp, ityp types.Type, permitaa bool) string {
	buf := new(bytes.Buffer)
	cdd.interfaceES(buf, ex, es, epos, etyp, ityp, permitaa)
	return buf.String()
}

func (cdd *CDD) interfaceExpr(w *bytes.Buffer, expr ast.Expr, ityp types.Type, permitaa bool) {
	etyp := cdd.exprType(expr)
	epos := expr.Pos()
	cdd.interfaceES(w, expr, "", epos, etyp, ityp, permitaa)
}

func (cdd *CDD) interfaceExprStr(expr ast.Expr, ityp types.Type, permitaa bool) string {
	buf := new(bytes.Buffer)
	cdd.interfaceExpr(buf, expr, ityp, permitaa)
	return buf.String()
}
