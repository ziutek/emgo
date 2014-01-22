package gotoc

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

func (cdd *CDD) Name(w *bytes.Buffer, obj types.Object) {
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
			cdd.Type(w, t)
			w.WriteByte('_')
			w.WriteString(o.Name())
			return
		}
	}

	if cdd.gtc.isImported(obj) || cdd.gtc.isGlobal(obj) {
		w.WriteString(upath(obj.Pkg().Path()))
		w.WriteByte('_')
	}
	name := obj.Name()
	if name == "_" {
		w.WriteString("__unused")
		w.WriteString(strconv.Itoa(cdd.un))
		cdd.un++
	} else {
		w.WriteString(name)
	}
}

func (cdd *CDD) NameStr(o types.Object) string {
	buf := new(bytes.Buffer)
	cdd.Name(buf, o)
	return buf.String()
}

func (cdd *CDD) BasicLit(w *bytes.Buffer, l *ast.BasicLit) {
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

func (cdd *CDD) SelectorExpr(w *bytes.Buffer, e *ast.SelectorExpr) ast.Expr {
	xt := cdd.gtc.ti.Types[e.X]
	sel := cdd.gtc.ti.Objects[e.Sel]

	switch s := sel.Type().(type) {
	case *types.Signature:
		if recv := s.Recv(); recv != nil {
			cdd.Name(w, sel)
			if _, ok := recv.Type().(*types.Pointer); !ok {
				return e.X
			}
			if _, ok := xt.(*types.Pointer); ok {
				return e.X
			}
			return &ast.UnaryExpr{Op: token.AND, X: e.X}
		}
		cdd.Expr(w, e.X)
		w.WriteByte('_')
		w.WriteString(e.Sel.Name)

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
	if v, ok := cdd.gtc.ti.Values[expr]; ok {
		// Constant expression
		w.WriteString(v.String())
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
