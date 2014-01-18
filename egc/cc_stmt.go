package main

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"go/ast"
	"go/token"
	"reflect"
	"strconv"
)

func (cc *CC) Type(w *bytes.Buffer, typ types.Type) {
	switch t := typ.(type) {
	case *types.Basic:
		if t.Kind() == types.UnsafePointer {
			w.WriteString("unsafe_Pointer")
		} else {
			types.WriteType(w, nil, t)
		}

	case *types.Named:
		o := t.Obj()
		if p := o.Pkg(); p != nil {
			w.WriteString(upath(p.Path()))
			w.WriteByte('_')
		}
		w.WriteString(o.Name())

	case *types.Pointer:
		cc.Type(w, t.Elem())
		w.WriteByte('*')

	case *types.Struct:
		w.WriteString("struct {\n")
		cc.il++
		for i, n := 0, t.NumFields(); i < n; i++ {
			f := t.Field(i)
			cc.indent(w)
			if tag := t.Tag(i); tag != "" {
				w.WriteString(reflect.StructTag(tag).Get("C"))
				w.WriteByte(' ')
			}
			cc.Type(w, f.Type())
			if !f.Anonymous() {
				w.WriteByte(' ')
				if f.Name() != "_" {
					w.WriteString(f.Name())
				} else {
					w.WriteString("__reserved")
					w.WriteString(strconv.Itoa(i))
				}
			}
			w.WriteString(";\n")
		}
		cc.il--
		w.WriteByte('}')

	default:
		fmt.Fprintf(w, "<%T>", t)
	}
}

func (cc *CC) TypeStr(typ types.Type) string {
	buf := new(bytes.Buffer)
	cc.Type(buf, typ)
	return buf.String()
}

func (cc *CC) ReturnStmt(w *bytes.Buffer, s *ast.ReturnStmt, resultT string) {
	cc.indent(w)
	switch len(s.Results) {
	case 0:
		if resultT == "" {
			w.WriteString("return;\n")
		} else {
			w.WriteString("goto __end;\n")
		}

	case 1:
		w.WriteString("return ")
		cc.Expr(w, s.Results[0])
		w.WriteString(";\n")

	default:
		w.WriteString("return (" + resultT + ") {")
		for i, e := range s.Results {
			if i > 0 {
				w.WriteString(", ")
			}
			cc.Expr(w, e)
		}
		w.WriteString("};\n")
	}
}

func (cc *CC) Stmt(w *bytes.Buffer, stmt ast.Stmt) {

	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			panic("unsuported multiple assignment")
		}

		if s.Tok == token.DEFINE {
			cc.Type(w, cc.ti.Types[s.Rhs[0]])
			w.WriteByte(' ')
		}

		cc.Expr(w, s.Lhs[0])

		switch s.Tok {
		case token.DEFINE:
			w.WriteString(" = ")

		case token.AND_NOT_ASSIGN:
			w.WriteString(" &= ~(")

		default:
			w.WriteString(" " + s.Tok.String() + " ")
		}

		cc.Expr(w, s.Rhs[0])

		if s.Tok == token.AND_NOT_ASSIGN {
			w.WriteByte(')')
		}
		w.WriteString(";\n")

	case *ast.ExprStmt:

		cc.Expr(w, s.X)
		w.WriteString(";\n")

	case *ast.IfStmt:
		if s.Init != nil {
			w.WriteString("{\n")
			cc.il++
			cc.Stmt(w, s.Init)
		}

		w.WriteString("if (")
		cc.Expr(w, s.Cond)
		w.WriteString(") ")
		cc.BlockStmt(w, s.Body, "")
		if s.Else == nil {
			w.WriteByte('\n')
		} else {
			w.WriteString(" else ")
			cc.Stmt(w, s.Else)
		}

		if s.Init != nil {
			cc.il--
			cc.indent(w)
			w.WriteString("}\n")
		}

	case *ast.IncDecStmt:
		w.WriteString(s.Tok.String())
		cc.Expr(w, s.X)
		w.WriteString(";\n")

	case *ast.BlockStmt:
		cc.BlockStmt(w, s, "")
		w.WriteByte('\n')

	case *ast.ForStmt:
		braces := s.Init != nil || s.Post != nil

		if braces {
			w.WriteString("{\n")
			cc.il++
		}
		if s.Init != nil {
			cc.indent(w)
			cc.Stmt(w, s.Init)
		}
		if braces {
			cc.indent(w)
		}

		w.WriteString("while (")
		if s.Cond != nil {
			cc.Expr(w, s.Cond)
		} else {
			w.WriteString("true")
		}
		w.WriteString(") ")

		cc.BlockStmt(w, s.Body, "")
		w.WriteByte('\n')

		if s.Post != nil {
			cc.indent(w)
			cc.Stmt(w, s.Post)
		}

		if braces {
			cc.il--
			cc.indent(w)
			w.WriteString("}\n")
		}

	default:
		fmt.Fprintf(w, "#<%T>#", stmt)
	}

}

func (cc *CC) BlockStmt(w *bytes.Buffer, bs *ast.BlockStmt, resultT string) {
	w.WriteString("{\n")
	cc.il++
	for _, stmt := range bs.List {

		switch s := stmt.(type) {
		case *ast.DeclStmt:
			cc.Decl(s.Decl)

		case *ast.ReturnStmt:
			cc.ReturnStmt(w, s, resultT)

		default:
			cc.indent(w)
			cc.Stmt(w, s)
		}
	}
	cc.il--
	cc.indent(w)
	w.WriteString("}")
}
