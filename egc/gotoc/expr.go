package gotoc

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"go/ast"
	"go/token"
)

func (cc *CC) Name(w *bytes.Buffer, obj types.Object) {
	switch o := obj.(type) {
	case *types.PkgName:
		w.WriteString(upath(o.Pkg().Path()))
		return

	case *types.Func:
		s := o.Type().(*types.Signature)
		if r := s.Recv(); r != nil {
			t := r.Type()
			if p, ok := t.(*types.Pointer); ok {
				t = p.Elem()
			}
			cc.Type(w, t)
			w.WriteByte('_')
			w.WriteString(o.Name())
			return
		}
	}

	if cc.isImported(obj) || cc.isGlobal(obj) {
		w.WriteString(upath(obj.Pkg().Path()))
		w.WriteByte('_')
	}
	w.WriteString(obj.Name())
}

func (cc *CC) NameStr(o types.Object) string {
	buf := new(bytes.Buffer)
	cc.Name(buf, o)
	return buf.String()
}

func (cc *CC) BasicLit(w *bytes.Buffer, l *ast.BasicLit) {
	switch l.Kind {
	case token.STRING:
		w.WriteString("_GOSTR(")
		w.WriteString(l.Value)
		w.WriteByte(')')

	case token.IMAG:
		notImplemented(l)

	default:
		w.WriteString(l.Value)
	}
}

func (cc *CC) SelectorExpr(w *bytes.Buffer, e *ast.SelectorExpr) ast.Expr {
	xt := cc.ti.Types[e.X]
	sel := cc.ti.Objects[e.Sel]

	switch s := sel.Type().(type) {
	case *types.Signature:
		if recv := s.Recv(); recv != nil {
			cc.Name(w, sel)
			if _, ok := recv.Type().(*types.Pointer); !ok {
				return e.X
			}
			if _, ok := xt.(*types.Pointer); ok {
				return e.X
			}
			return &ast.UnaryExpr{Op: token.AND, X: e.X}
		}
		cc.Expr(w, e.X)
		w.WriteByte('_')
		w.WriteString(e.Sel.Name)

	default:
		cc.Expr(w, e.X)
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

func (cc *CC) Expr(w *bytes.Buffer, expr ast.Expr) {
	if v, ok := cc.ti.Values[expr]; ok {
		// Constant expression
		w.WriteString(v.String())
		return
	}

	switch e := expr.(type) {
	case *ast.BasicLit:
		cc.BasicLit(w, e)

	case *ast.BinaryExpr:
		cc.Expr(w, e.X)
		op := e.Op.String()
		if op == "&^" {
			op = "&~"
		}
		w.WriteString(op)
		cc.Expr(w, e.Y)

	case *ast.CallExpr:
		var recv ast.Expr

		switch cc.ti.Types[e.Fun].(type) {
		case *types.Signature:
			switch f := e.Fun.(type) {
			case *ast.SelectorExpr:
				recv = cc.SelectorExpr(w, f)

			default:
				cc.Expr(w, f)
			}

		default:
			w.WriteByte('(')
			cc.Type(w, cc.ti.Types[e.Fun])
			w.WriteByte(')')
		}

		w.WriteByte('(')
		if recv != nil {
			cc.Expr(w, recv)
			if len(e.Args) > 0 {
				w.WriteString(", ")
			}
		}

		for i, a := range e.Args {
			if i != 0 {
				w.WriteString(", ")
			}
			cc.Expr(w, a)
		}
		w.WriteByte(')')

	case *ast.Ident:
		cc.Name(w, cc.ti.Objects[e])

	case *ast.IndexExpr:
		cc.Expr(w, e.X)
		switch cc.ti.Types[e.X].(type) {
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
		cc.Expr(w, e.Index)
		w.WriteByte(']')

	case *ast.KeyValueExpr:
		w.WriteByte('.')
		cc.Expr(w, e.Key)
		w.WriteString(" = ")
		cc.Expr(w, e.Value)

	case *ast.ParenExpr:
		w.WriteByte('(')
		cc.Expr(w, e.X)
		w.WriteByte(')')

	case *ast.SelectorExpr:
		cc.SelectorExpr(w, e)

	case *ast.SliceExpr:
		notImplemented(e)

	case *ast.StarExpr:
		w.WriteByte('*')
		cc.Expr(w, e.X)

	case *ast.TypeAssertExpr:
		notImplemented(e)

	case *ast.UnaryExpr:
		op := e.Op.String()
		if op == "^" {
			op = "~"
		}
		w.WriteString(op)
		cc.Expr(w, e.X)

	case *ast.CompositeLit:
		w.WriteByte('(')
		cc.Type(w, cc.ti.Types[e])
		w.WriteString("){")

		for i, el := range e.Elts {
			if i > 0 {
				w.WriteString(", ")
			}
			cc.Expr(w, el)
		}

		w.WriteByte('}')

	default:
		fmt.Fprintf(w, "!%v<%T>!", e, e)
	}
}
