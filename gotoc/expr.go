package gotoc

import (
	"bytes"
	"code.google.com/p/go.tools/go/exact"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
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

//k := c.Type().(*types.Basic).Kind()
//ev := c.Val()

func (cdd *CDD) Value(w *bytes.Buffer, ev exact.Value, t types.Type) {
	k := t.Underlying().(*types.Basic).Kind()

	w.WriteByte('(')

	switch ev.Kind() {
	case exact.Int:
		writeInt(w, ev, k)

	case exact.Float:
		writeFloat(w, ev, k)

	case exact.Complex:

	default:
		w.WriteString(ev.String())
	}

	w.WriteByte(')')
}

func (cdd *CDD) addObject(o types.Object) {
	if o == cdd.Origin {
		return
	}
	if cdd.body {
		cdd.BodyUses[o] = struct{}{}
	} else {
		cdd.DeclUses[o] = struct{}{}
	}
}

func (cdd *CDD) Name(w *bytes.Buffer, obj types.Object) {
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
			}
			cdd.Type(w, t)
			w.WriteByte('_')
			w.WriteString(o.Name())
			if !cdd.gtc.isLocal(t.(*types.Named).Obj()) {
				cdd.addObject(o)
			}
			return
		}
	}

	if p := obj.Pkg(); p != nil && !cdd.gtc.isLocal(obj) {
		cdd.addObject(obj)
		w.WriteString(upath(obj.Pkg().Path()))
		w.WriteByte('_')
	}
	name := obj.Name()
	switch name {
	case "_":
		w.WriteString("__unused")
		w.WriteString(uniqueId(obj))
		return

	case "init":
		name = uniqueId(obj) + name
	}
	w.WriteString(name)
}

func (cdd *CDD) NameStr(o types.Object) string {
	buf := new(bytes.Buffer)
	cdd.Name(buf, o)
	return buf.String()
}

func (cdd *CDD) BasicLit(w *bytes.Buffer, l *ast.BasicLit) {
	switch l.Kind {
	case token.STRING:
		w.WriteString("__GOSTR(")
		w.WriteString(l.Value)
		w.WriteByte(')')

	default:
		notImplemented(l)
	}
}

func (cdd *CDD) SelectorExpr(w *bytes.Buffer, e *ast.SelectorExpr) ast.Expr {
	xt := cdd.gtc.ti.Types[e.X]
	sel := cdd.gtc.ti.Objects[e.Sel]

	switch s := sel.Type().(type) {
	case *types.Signature:
		cdd.Name(w, sel)
		if recv := s.Recv(); recv != nil {
			if _, ok := recv.Type().(*types.Pointer); !ok {
				// Method has non pointer receiver so there is guaranteed
				// that e.X isn't a pointer.
				return e.X
			}
			// Method has pointer receiver.
			if _, ok := xt.(*types.Pointer); ok {
				return e.X // e.X is pointer
			}
			// e.X isn't a pointer so create pointer to it.
			return &ast.UnaryExpr{Op: token.AND, X: e.X}
		}

	default:
		cdd.Expr(w, e.X)
		switch xt.(type) {
		case *types.Named:
			w.WriteByte('.')

		case *types.Pointer:
			w.WriteString("->")

		default:
			w.WriteByte('_')

		}
		w.WriteString(e.Sel.Name)
	}
	return nil
}

func (cdd *CDD) Expr(w *bytes.Buffer, expr ast.Expr) {
	cdd.Complexity++

	if ev := cdd.gtc.ti.Values[expr]; ev != nil {
		// Constant expression
		cdd.Value(w, ev, cdd.gtc.ti.Types[expr])
		return
	}

	switch e := expr.(type) {
	case *ast.BasicLit:
		cdd.BasicLit(w, e)

	case *ast.BinaryExpr:
		cdd.Expr(w, e.X)
		op := e.Op.String()
		if op == "&^" {
			op = "&~"
		}
		w.WriteString(op)
		cdd.Expr(w, e.Y)

	case *ast.CallExpr:
		var recv ast.Expr

		switch cdd.gtc.ti.Types[e.Fun].(type) {
		case *types.Signature:
			switch f := e.Fun.(type) {
			case *ast.SelectorExpr:
				recv = cdd.SelectorExpr(w, f)

			default:
				cdd.Expr(w, f)
			}

		default:
			w.WriteByte('(')
			cdd.Type(w, cdd.gtc.ti.Types[e.Fun])
			w.WriteByte(')')
		}

		w.WriteByte('(')
		if recv != nil {
			cdd.Expr(w, recv)
			if len(e.Args) > 0 {
				w.WriteString(", ")
			}
		}

		for i, a := range e.Args {
			if i != 0 {
				w.WriteString(", ")
			}
			cdd.Expr(w, a)
		}
		w.WriteByte(')')

	case *ast.Ident:
		cdd.Name(w, cdd.gtc.ti.Objects[e])

	case *ast.IndexExpr:
		cdd.Expr(w, e.X)
		switch cdd.gtc.ti.Types[e.X].(type) {
		case *types.Basic: // string
			w.WriteString(".str")
		case *types.Slice:
			w.WriteString(".sli")
		case *types.Array:
			// use C arrays
		default:
			notImplemented(e)
		}
		w.WriteByte('[')
		cdd.Expr(w, e.Index)
		w.WriteByte(']')

	case *ast.KeyValueExpr:
		w.WriteByte('.')
		cdd.Expr(w, e.Key)
		w.WriteString(" = ")
		cdd.Expr(w, e.Value)

	case *ast.ParenExpr:
		w.WriteByte('(')
		cdd.Expr(w, e.X)
		w.WriteByte(')')

	case *ast.SelectorExpr:
		cdd.SelectorExpr(w, e)

	case *ast.SliceExpr:
		notImplemented(e)

	case *ast.StarExpr:
		w.WriteByte('*')
		cdd.Expr(w, e.X)

	case *ast.TypeAssertExpr:
		notImplemented(e)

	case *ast.UnaryExpr:
		op := e.Op.String()
		if op == "^" {
			op = "~"
		}
		w.WriteString(op)
		cdd.Expr(w, e.X)

	case *ast.CompositeLit:
		w.WriteByte('(')
		cdd.Type(w, cdd.gtc.ti.Types[e])
		w.WriteString("){")

		for i, el := range e.Elts {
			if i > 0 {
				w.WriteString(", ")
			}
			cdd.Expr(w, el)
		}

		w.WriteByte('}')

	default:
		fmt.Fprintf(w, "!%v<%T>!", e, e)
	}
}
