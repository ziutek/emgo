package gotoc

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"code.google.com/p/go.tools/go/types"
)

func (cdd *CDD) ReturnStmt(w *bytes.Buffer, s *ast.ReturnStmt, resultT string) (end bool) {
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

func (cdd *CDD) Stmt(w *bytes.Buffer, stmt ast.Stmt, label, resultT string) (end bool, acds []*CDD) {
	updateEA := func(e bool, a []*CDD) {
		if e {
			end = true
		}
		acds = append(acds, a...)
	}

	cdd.Complexity++
	
	// TODO: really need to rewrite this convoluted code
	switch s := stmt.(type) {
	case *ast.AssignStmt:
		tmpnames := make([]string, len(s.Lhs))
		lhs := make([]string, len(s.Lhs))
		
		if len(s.Lhs) > 1 && len(s.Rhs) == 1 {
			// Tuple on RHS
			tup := cdd.gtc.ti.Types[s.Rhs[0]].Type.(*types.Tuple)
			w.WriteString(cdd.gtc.tn.Name(tup))
			tmpname := "__tmp" + cdd.gtc.uniqueId()
			w.WriteString(" " + tmpname + " = ")
			cdd.Expr(w, s.Rhs[0])
			w.WriteString(";\n")
			for i, n := 0, tup.Len(); i < n; i++ {
				if  s.Tok == token.DEFINE {
					// Type is need for definition
					typ, dim, a := cdd.TypeStr(tup.At(i).Type())
					acds = append(acds, a...)
					ident := s.Lhs[i].(*ast.Ident)
					name := cdd.NameStr(cdd.gtc.ti.Objects[ident], true)
					lhs[i] = typ + " " + dimFuncPtr(name, dim)
				} else {
					lhs[i] = cdd.ExprStr(s.Lhs[i])
				}
				tmpnames[i] = tmpname + "._" + strconv.Itoa(i)
			}
		} else {
			for i := 0; i < len(s.Lhs); i++ {
				var name string
				if s.Tok == token.DEFINE {
					name = cdd.NameStr(cdd.gtc.ti.Objects[s.Lhs[i].(*ast.Ident)], true)
				} else if len(s.Lhs) > 1 {
					name = "__tmp" + cdd.gtc.uniqueId()
				}

				if i != 0 {
					cdd.indent(w)
				}

				if s.Tok == token.DEFINE || len(s.Lhs) > 1 {
					// Type is need for definition or temporary variable
					dim, a := cdd.Type(w, cdd.gtc.ti.Types[s.Rhs[i]].Type)
					acds = append(acds, a...)
					w.WriteByte(' ')
					w.WriteString(dimFuncPtr(name, dim))
					if s.Tok != token.DEFINE {
						lhs[i] = cdd.ExprStr(s.Lhs[i])
					}
				} else {
					cdd.Expr(w, s.Lhs[i])
				}

				switch s.Tok {
				case token.DEFINE:
					w.WriteString(" = ")

				case token.AND_NOT_ASSIGN:
					w.WriteString(" &= ~(")

				default:
					w.WriteString(" " + s.Tok.String() + " ")
				}

				cdd.Expr(w, s.Rhs[i])

				if s.Tok == token.AND_NOT_ASSIGN {
					w.WriteByte(')')
				}
				w.WriteString(";\n")

				tmpnames[i] = name
			}
			if len(s.Lhs) == 1 || s.Tok == token.DEFINE {
				break
			}
		}

		for i, tmpname := range tmpnames {
			cdd.indent(w)
			w.WriteString(lhs[i])
			w.WriteString(" = ")
			w.WriteString(tmpname)
			w.WriteString(";\n")
		}

	case *ast.ExprStmt:
		cdd.Expr(w, s.X)
		w.WriteString(";\n")

	case *ast.IfStmt:
		if s.Init != nil {
			w.WriteString("{\n")
			cdd.il++
			cdd.indent(w)
			updateEA(cdd.Stmt(w, s.Init, "", resultT))
			cdd.indent(w)
		}

		w.WriteString("if (")
		cdd.Expr(w, s.Cond)
		w.WriteString(") ")
		updateEA(cdd.BlockStmt(w, s.Body, resultT))
		if s.Else == nil {
			w.WriteByte('\n')
		} else {
			w.WriteString(" else ")
			updateEA(cdd.Stmt(w, s.Else, "", resultT))
		}

		if s.Init != nil {
			cdd.il--
			cdd.indent(w)
			w.WriteString("}\n")
		}

	case *ast.IncDecStmt:
		w.WriteString(s.Tok.String())
		w.WriteByte('(')
		cdd.Expr(w, s.X)
		w.WriteString(");\n")

	case *ast.BlockStmt:
		updateEA(cdd.BlockStmt(w, s, ""))
		w.WriteByte('\n')

	case *ast.ForStmt:
		if s.Init != nil {
			w.WriteString("{\n")
			cdd.il++
			cdd.indent(w)
			updateEA(cdd.Stmt(w, s.Init, "", resultT))
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
		updateEA(cdd.BlockStmt(w, s.Body, ""))
		w.WriteByte('\n')

		if s.Post != nil {
			cdd.indent(w)
			updateEA(cdd.Stmt(w, s.Post, "", resultT))
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

	case *ast.ReturnStmt:
		if cdd.ReturnStmt(w, s, resultT) {
			end = true
		}

	case *ast.SwitchStmt:
		w.WriteString("switch(0) {\n")
		cdd.indent(w)
		w.WriteString("case 0:;\n")
		cdd.il++

		if s.Init != nil {
			cdd.indent(w)
			updateEA(cdd.Stmt(w, s.Init, "", resultT))
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
		updateEA(cdd.BlockStmt(w, s.Body, resultT))
		w.WriteByte('\n')

		cdd.il--
		cdd.indent(w)
		w.WriteString("}\n")

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
					ftLabel = "__fallthrough" + cdd.gtc.uniqueId()
				}
				w.WriteString("goto " + ftLabel + ";\n")
			} else {
				updateEA(cdd.Stmt(w, stmt, "", resultT))
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
	return
}

func (cdd *CDD) BlockStmt(w *bytes.Buffer, bs *ast.BlockStmt, resultT string) (end bool, acds []*CDD) {
	updateEA := func(e bool, a []*CDD) {
		if e {
			end = true
		}
		acds = append(acds, a...)
	}

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

		case *ast.LabeledStmt:
			cdd.label(w, s.Label.Name, "")
			cdd.indent(w)
			updateEA(cdd.Stmt(w, s.Stmt, s.Label.Name, resultT))

		default:
			cdd.indent(w)
			updateEA(cdd.Stmt(w, s, "", resultT))
		}
	}
	cdd.il--
	cdd.indent(w)
	w.WriteString("}")
	return
}
