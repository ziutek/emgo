package gotoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
)

func (cdd *CDD) ReturnStmt(w *bytes.Buffer, s *ast.ReturnStmt, resultT string) {
	cdd.indent(w)
	switch len(s.Results) {
	case 0:
		if resultT == "" {
			w.WriteString("return;\n")
		} else {
			w.WriteString("goto __end;\n")
		}

	case 1:
		w.WriteString("return ")
		cdd.Expr(w, s.Results[0])
		w.WriteString(";\n")

	default:
		w.WriteString("return (" + resultT + ") {")
		for i, e := range s.Results {
			if i > 0 {
				w.WriteString(", ")
			}
			cdd.Expr(w, e)
		}
		w.WriteString("};\n")
	}
}

func (cdd *CDD) Stmt(w *bytes.Buffer, stmt ast.Stmt) {
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			panic("unsuported multiple assignment")
		}

		if s.Tok == token.DEFINE {
			cdd.Type(w, cdd.gtc.ti.Types[s.Rhs[0]])
			w.WriteByte(' ')
		}

		cdd.Expr(w, s.Lhs[0])

		switch s.Tok {
		case token.DEFINE:
			w.WriteString(" = ")

		case token.AND_NOT_ASSIGN:
			w.WriteString(" &= ~(")

		default:
			w.WriteString(" " + s.Tok.String() + " ")
		}

		cdd.Expr(w, s.Rhs[0])

		if s.Tok == token.AND_NOT_ASSIGN {
			w.WriteByte(')')
		}
		w.WriteString(";\n")

	case *ast.ExprStmt:

		cdd.Expr(w, s.X)
		w.WriteString(";\n")

	case *ast.IfStmt:
		if s.Init != nil {
			w.WriteString("{\n")
			cdd.il++
			cdd.Stmt(w, s.Init)
		}

		w.WriteString("if (")
		cdd.Expr(w, s.Cond)
		w.WriteString(") ")
		cdd.BlockStmt(w, s.Body, "")
		if s.Else == nil {
			w.WriteByte('\n')
		} else {
			w.WriteString(" else ")
			cdd.Stmt(w, s.Else)
		}

		if s.Init != nil {
			cdd.il--
			cdd.indent(w)
			w.WriteString("}\n")
		}

	case *ast.IncDecStmt:
		w.WriteString(s.Tok.String())
		cdd.Expr(w, s.X)
		w.WriteString(";\n")

	case *ast.BlockStmt:
		cdd.BlockStmt(w, s, "")
		w.WriteByte('\n')

	case *ast.ForStmt:
		braces := s.Init != nil || s.Post != nil

		if braces {
			w.WriteString("{\n")
			cdd.il++
		}
		if s.Init != nil {
			cdd.indent(w)
			cdd.Stmt(w, s.Init)
		}
		if braces {
			cdd.indent(w)
		}

		w.WriteString("while (")
		if s.Cond != nil {
			cdd.Expr(w, s.Cond)
		} else {
			w.WriteString("true")
		}
		w.WriteString(") ")

		cdd.BlockStmt(w, s.Body, "")
		w.WriteByte('\n')

		if s.Post != nil {
			cdd.indent(w)
			cdd.Stmt(w, s.Post)
		}

		if braces {
			cdd.il--
			cdd.indent(w)
			w.WriteString("}\n")
		}

	default:
		fmt.Fprintf(w, "#<%T>#", stmt)
	}

}

func (cdd *CDD) BlockStmt(w *bytes.Buffer, bs *ast.BlockStmt, resultT string) {
	w.WriteString("{\n")
	cdd.il++
	for _, stmt := range bs.List {

		switch s := stmt.(type) {
		case *ast.DeclStmt:
			cdds := cdd.gtc.Decl(s.Decl, cdd.il)
			for _, cdd := range cdds{
				w.Write(cdd.Decl)
			}
			for _, cdd := range cdds{
				w.Write(cdd.Def)
			}
			

		case *ast.ReturnStmt:
			cdd.ReturnStmt(w, s, resultT)

		default:
			cdd.indent(w)
			cdd.Stmt(w, s)
		}
	}
	cdd.il--
	cdd.indent(w)
	w.WriteString("}")
}
