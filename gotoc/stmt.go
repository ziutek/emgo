package gotoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

func (cdd *CDD) ReturnStmt(w *bytes.Buffer, s *ast.ReturnStmt, resultT string) (end bool) {
	cdd.indent(w)
	switch len(s.Results) {
	case 0:
		if resultT == "void" {
			w.WriteString("return;\n")
		} else {
			w.WriteString("goto __end;\n")
			end = true
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
	return
}

func (cdd *CDD) label(w *bytes.Buffer, label, suffix string) {
	cdd.il--
	cdd.indent(w)
	w.WriteString(label)
	w.WriteString(suffix)
	w.WriteString(":;\n")
	cdd.il++
}

func (cdd *CDD) Stmt(w *bytes.Buffer, stmt ast.Stmt, label, resultT string) {
	cdd.Complexity++

	switch s := stmt.(type) {
	case *ast.AssignStmt:
		if len(s.Lhs) != 1 || len(s.Rhs) != 1 {
			panic("unsuported multiple assignment")
		}

		var dim []int64
		if s.Tok == token.DEFINE {
			dim = cdd.Type(w, cdd.gtc.ti.Types[s.Rhs[0]].Type)
			w.WriteByte(' ')
		}

		cdd.Expr(w, s.Lhs[0])
		writeDim(w, dim)

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
			cdd.Stmt(w, s.Init, "", resultT)
		}

		w.WriteString("if (")
		cdd.Expr(w, s.Cond)
		w.WriteString(") ")
		cdd.BlockStmt(w, s.Body, resultT)
		if s.Else == nil {
			w.WriteByte('\n')
		} else {
			w.WriteString(" else ")
			cdd.Stmt(w, s.Else, "", resultT)
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
		if s.Init != nil {
			w.WriteString("{\n")
			cdd.il++
			cdd.indent(w)
			cdd.Stmt(w, s.Init, "", resultT)
		}
		if label != "" {
			cdd.label(w, label, "_continue")
		}
		if s.Init != nil {
			cdd.indent(w)
		}
		w.WriteString("while (")
		if s.Cond != nil {
			cdd.Expr(w, s.Cond)
		} else {
			w.WriteString("true")
		}
		w.WriteString(") ")

		if s.Post != nil {
			w.WriteString("{\n")
			cdd.il++
			cdd.indent(w)
		}

		cdd.BlockStmt(w, s.Body, "")
		w.WriteByte('\n')

		if s.Post != nil {
			cdd.indent(w)
			cdd.Stmt(w, s.Post, "", resultT)
			cdd.il--
			cdd.indent(w)
			w.WriteString("}\n")
		}

		if s.Init != nil {
			cdd.il--
			cdd.indent(w)
			w.WriteString("}\n")
		}

		if label != "" {
			cdd.label(w, label, "_break")
		}

	case *ast.SwitchStmt:
		w.WriteString("do {\n")
		cdd.il++

		if s.Init != nil {
			cdd.indent(w)
			cdd.Stmt(w, s.Init, "", resultT)
		}

		cdd.indent(w)
		if s.Tag != nil {
			cdd.Type(w, cdd.gtc.ti.Types[s.Tag].Type)
			w.WriteString(" __tag = ")
			cdd.Expr(w, s.Tag)
			w.WriteString(";\n")
		} else {
			w.WriteString("bool __tag = true;\n")
		}

		cdd.indent(w)
		cdd.BlockStmt(w, s.Body, resultT)
		w.WriteByte('\n')

		cdd.il--
		cdd.indent(w)
		w.WriteString("} while(false);\n")

		if label != "" {
			cdd.label(w, label, "_break")
		}

	case *ast.CaseClause:
		if s.List != nil {
			w.WriteString("if (")
			for i, e := range s.List {
				if i != 0 {
					w.WriteString(" || ")
				}
				w.WriteString("__tag == ")
				cdd.Expr(w, e)
			}
			w.WriteString(") ")
		}
		w.WriteString("{\n")
		cdd.il++

		var ftLabel string
		for _, stmt := range s.Body {
			cdd.indent(w)
			if bs, ok := stmt.(*ast.BranchStmt); ok && bs.Tok == token.FALLTHROUGH {
				if ftLabel == "" {
					ftLabel = "__fallthrough" + strconv.Itoa(int(s.End()))
				}
				w.WriteString("goto " + ftLabel + ";\n")
			} else {
				cdd.Stmt(w, stmt, "", resultT)
			}
		}

		cdd.indent(w)
		w.WriteString("break;\n")

		cdd.il--
		cdd.indent(w)
		w.WriteString("}\n")
		if ftLabel != "" {
			cdd.il--
			cdd.indent(w)
			w.WriteString(ftLabel + ":;\n")
			cdd.il++
		}

	case *ast.BranchStmt:
		if s.Label == nil {
			w.WriteString(s.Tok.String())
		} else {
			w.WriteString("goto " + s.Label.Name)
			switch s.Tok {
			case token.BREAK:
				w.WriteString("_break")
			case token.CONTINUE:
				w.WriteString("_continue")
			}
		}
		w.WriteString(";\n")

	default:
		fmt.Fprintf(w, "#<%T>#", stmt)
	}

}

func (cdd *CDD) BlockStmt(w *bytes.Buffer, bs *ast.BlockStmt, resultT string) (end bool) {
	w.WriteString("{\n")
	cdd.il++
	for _, stmt := range bs.List {
		switch s := stmt.(type) {
		case *ast.DeclStmt:
			cdds := cdd.gtc.Decl(s.Decl, cdd.il)
			for _, c := range cdds {
				for u, typPtr := range c.BodyUses {
					cdd.BodyUses[u] = typPtr
				}
				w.Write(c.Decl)
			}
			for _, c := range cdds {
				w.Write(c.Def)
			}

		case *ast.ReturnStmt:
			if cdd.ReturnStmt(w, s, resultT) {
				end = true
			}

		case *ast.LabeledStmt:
			cdd.label(w, s.Label.Name, "")
			cdd.indent(w)
			cdd.Stmt(w, s.Stmt, s.Label.Name, resultT)

		default:
			cdd.indent(w)
			cdd.Stmt(w, s, "", resultT)
		}
	}
	cdd.il--
	cdd.indent(w)
	w.WriteString("}")
	return
}
