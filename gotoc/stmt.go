package gotoc

import (
	"bytes"
	"go/ast"
	"go/token"
	"go/types"
	"strconv"
)

func (cdd *CDD) ReturnStmt(w *bytes.Buffer, s *ast.ReturnStmt, resultT string, tup *types.Tuple) (end bool) {
	switch len(s.Results) {
	case 0:
		if resultT == "void" {
			w.WriteString("return;\n")
		} else {
			w.WriteString("goto end;\n")
			end = true
		}

	case 1:
		w.WriteString("return ")
		if tup.Len() != 1 {
			retTyp := tup
			eTyp := cdd.exprType(s.Results[0])
			if !types.Identical(retTyp, eTyp) {
				w.WriteString("({\n")
				cdd.il++
				cdd.indent(w)
				tn, fields := cdd.tupleName(eTyp.(*types.Tuple))
				tmp := "_tmp" + cdd.gtc.uniqueId()
				w.WriteString(tn + " " + tmp + " = ")
				cdd.Expr(w, s.Results[0], eTyp)
				w.WriteString(";\n")
				cdd.indent(w)
				w.WriteString("(" + resultT + "){")
				for i, v := range fields {
					if i != 0 {
						w.WriteString(", ")
					}
					tmpf := tmp + "." + v.Name()
					cdd.interfaceES(
						w, tmpf, s.Pos(), v.Type(), tup.At(i).Type(),
					)
				}
				w.WriteString("};\n")
				cdd.il--
				cdd.indent(w)
				w.WriteString("})")
			} else {
				cdd.Expr(w, s.Results[0], retTyp)
			}
		} else {
			retTyp := tup.At(0).Type()
			cdd.interfaceExpr(w, s.Results[0], retTyp)
		}
		w.WriteString(";\n")

	default:
		w.WriteString("return (" + resultT + "){")
		for i, expr := range s.Results {
			if i > 0 {
				w.WriteString(", ")
			}
			cdd.interfaceExpr(w, expr, tup.At(i).Type())
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

var untypedNil = types.Typ[types.UntypedNil]

func (cdd *CDD) Stmt(w *bytes.Buffer, stmt ast.Stmt, label, resultT string, tup *types.Tuple) (end bool) {
	updateEnd := func(e bool) {
		if e {
			end = true
		}
	}

	cdd.Complexity++

	switch s := stmt.(type) {
	case *ast.DeclStmt:
		cdds := cdd.gtc.Decl(s.Decl, cdd.il)
		for _, c := range cdds {
			for u, typPtr := range c.FuncBodyUses {
				cdd.FuncBodyUses[u] = typPtr
			}
			w.Write(c.Decl)
			cdd.acds = append(cdd.acds, c.acds...)
		}
		for _, c := range cdds {
			w.Write(c.Def)
		}

	case *ast.AssignStmt:
		rhs := make([]string, len(s.Lhs))
		typ := make([]types.Type, len(s.Lhs))

		rhsIsTuple := len(s.Lhs) > 1 && len(s.Rhs) == 1

		if rhsIsTuple {
			tup := cdd.exprType(s.Rhs[0]).(*types.Tuple)
			/*tupName, _ := cdd.tupleName(tup)
			w.WriteString(tupName)
			tupName = "_tmp" + cdd.gtc.uniqueId()
			w.WriteString(" " + tupName + " = ")
			cdd.Expr(w, s.Rhs[0], nil)*/
			tupName := "_tmp" + cdd.gtc.uniqueId()
			cdd.varDecl(w, tup, tupName, s.Rhs[0], "")
			w.WriteByte('\n')
			cdd.indent(w)
			for i, n := 0, tup.Len(); i < n; i++ {
				rhs[i] = tupName + "._" + strconv.Itoa(i)
				if s.Tok == token.DEFINE {
					if o := cdd.gtc.ti.Defs[s.Lhs[i].(*ast.Ident)]; o != nil {
						typ[i] = o.Type()
					} // else Lhs[i] was defined before.
				}
			}
		} else {
			for i, e := range s.Rhs {
				if s.Tok == token.DEFINE {
					var t types.Type
					if o := cdd.gtc.ti.Defs[s.Lhs[i].(*ast.Ident)]; o != nil {
						t = o.Type()
						typ[i] = t
					} else {
						// Lhs[i] was defined before.
						t = cdd.exprType(s.Lhs[i])
					}
					rhs[i] = cdd.ExprStr(e, t)
					continue
				}
				t := cdd.exprType(s.Lhs[i])
				rhs[i] = cdd.interfaceExprStr(e, t)
				typ[i] = t
			}
		}

		lhs := make([]string, len(s.Lhs))

		if s.Tok == token.DEFINE {
			for i, e := range s.Lhs {
				name := cdd.NameStr(cdd.object(e.(*ast.Ident)), true)
				if name == "_$" {
					lhs[i] = "_"
				} else if typ[i] != nil {
					t, dim := cdd.TypeStr(typ[i])
					lhs[i] = t + " " + dimFuncPtr(name, dim)
				} else {
					lhs[i] = name
				}
			}
		} else {
			for i, e := range s.Lhs {
				lhs[i] = cdd.ExprStr(e, nil)
			}
		}

		if len(s.Rhs) == len(s.Lhs) && len(s.Lhs) > 1 && s.Tok != token.DEFINE {
			for i, t := range typ {
				if i > 0 {
					cdd.indent(w)
				}
				if lhs[i] == "_" {
					w.WriteString("(void)(")
					w.WriteString(rhs[i])
					w.WriteString(");\n")
				} else {
					dim := cdd.Type(w, t)
					tmp := "_tmp" + cdd.gtc.uniqueId()
					w.WriteString(" " + dimFuncPtr(tmp, dim))
					w.WriteString(" = " + rhs[i] + ";\n")
					rhs[i] = tmp
				}
			}
			cdd.indent(w)
		}

		var atok string
		switch s.Tok {
		case token.DEFINE:
			atok = " = "

		case token.AND_NOT_ASSIGN:
			atok = " &= "
			rhs[0] = "~(" + rhs[0] + ")"

		default:
			atok = " " + s.Tok.String() + " "
		}
		indent := false
		for i := 0; i < len(lhs); i++ {
			li := lhs[i]
			if li == "_" && rhsIsTuple {
				continue
			}
			if indent {
				cdd.indent(w)
			} else {
				indent = true
			}
			if li == "_" {
				w.WriteString("(void)(")
				w.WriteString(rhs[i])
				w.WriteString(");\n")
			} else {
				w.WriteString(li)
				w.WriteString(atok)
				w.WriteString(rhs[i])
				w.WriteString(";\n")
			}
		}

	case *ast.ExprStmt:
		cdd.Expr(w, s.X, nil)
		w.WriteString(";\n")

	case *ast.IfStmt:
		if s.Init != nil {
			w.WriteString("{\n")
			cdd.il++
			cdd.indent(w)
			updateEnd(cdd.Stmt(w, s.Init, "", resultT, tup))
			cdd.indent(w)
		}

		w.WriteString("if (")
		cdd.Expr(w, s.Cond, nil)
		w.WriteString(") ")
		updateEnd(cdd.BlockStmt(w, s.Body, resultT, tup))
		if s.Else == nil {
			w.WriteByte('\n')
		} else {
			w.WriteString(" else ")
			updateEnd(cdd.Stmt(w, s.Else, "", resultT, tup))
		}

		if s.Init != nil {
			cdd.il--
			cdd.indent(w)
			w.WriteString("}\n")
		}

	case *ast.IncDecStmt:
		w.WriteString(s.Tok.String())
		w.WriteByte('(')
		cdd.Expr(w, s.X, nil)
		w.WriteString(");\n")

	case *ast.BlockStmt:
		updateEnd(cdd.BlockStmt(w, s, resultT, tup))
		w.WriteByte('\n')

	case *ast.ForStmt:
		if s.Init != nil {
			w.WriteString("{\n")
			cdd.il++
			cdd.indent(w)
			updateEnd(cdd.Stmt(w, s.Init, "", resultT, tup))
			cdd.indent(w)
		}
		w.WriteString("for (;")
		if s.Cond != nil {
			cdd.Expr(w, s.Cond, nil)
		}
		w.WriteByte(';')
		if s.Post != nil {
			w.WriteString(" ({\n")
			cdd.il++
			cdd.indent(w)
			updateEnd(cdd.Stmt(w, s.Post, "", resultT, tup))
			cdd.il--
			cdd.indent(w)
			w.WriteString("})")
		}
		w.WriteString(") ")
		if label != "" {
			w.WriteString("{\n")
			cdd.il++
			cdd.indent(w)
		}
		updateEnd(cdd.BlockStmt(w, s.Body, resultT, tup))
		w.WriteByte('\n')
		if label != "" {
			cdd.label(w, label, "_continue")
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

	case *ast.RangeStmt:
		w.WriteString("{\n")
		cdd.il++
		xt := cdd.exprType(s.X)
		xs := "_x"
		xl := ""

		array := false
		switch t := xt.Underlying().(type) {
		case *types.Array:
			array = true
			xl = strconv.FormatInt(t.Len(), 10)

		case *types.Pointer:
			array = true
			xl = strconv.FormatInt(t.Elem().Underlying().(*types.Array).Len(), 10)
		}

		if v, ok := s.Value.(*ast.Ident); ok && v.Name == "_" {
			s.Value = nil
		}

		switch e := s.X.(type) {
		case *ast.Ident:
			xs = cdd.NameStr(cdd.object(e), true)

		default:
			if s.Value != nil || !array {
				cdd.indent(w)
				cdd.varDecl(w, xt, xs, e, "")
				w.WriteByte('\n')
			}
		}

		if !array {
			xl = "len(" + xs + ")"
		}

		var (
			ks, vs string
			group  bool
		)

		t := xt.Underlying()

		cdd.indent(w)
		switch ct := t.(type) {
		case *types.Slice, *types.Array, *types.Pointer:
			if s.Key == nil {
				ks = "_$"
			} else {
				ks = cdd.ExprStr(s.Key, nil)
			}
			if s.Tok == token.DEFINE {
				w.WriteString("int ")
			}
			w.WriteString(ks + " = 0;\n")

			cdd.indent(w)
			w.WriteString("for (;" + ks + " < " + xl + "; ++" + ks + ") ")

			if s.Value != nil {
				vs = cdd.indexExprStr(xt, xs, s.Key)
			}
			group = (s.Value != nil || label != "")
			if group {
				w.WriteString("{\n")
				cdd.il++
			}

		case *types.Chan:
			tup := types.NewTuple(
				types.NewVar(stmt.Pos(), cdd.gtc.pkg, "", ct.Elem()),
				types.NewVar(stmt.Pos(), cdd.gtc.pkg, "", types.Typ[types.Bool]),
			)
			tn, _ := cdd.tupleName(tup)
			w.WriteString("for (;;) {\n")
			cdd.il++
			cdd.indent(w)
			w.WriteString(tn + " _vok = RECVOK(" + tn + ", " + xs + ");\n")
			cdd.indent(w)
			w.WriteString("if (!_vok._1) break;\n")
			s.Value = s.Key
			if s.Value != nil {
				vs = "_vok._0"
			}
			group = true

		default:
			cdd.notImplemented(s, xt)
		}

		if s.Value != nil {
			cdd.indent(w)
			if s.Tok == token.DEFINE {
				if pt, ok := t.(*types.Pointer); ok {
					t = pt.Elem().Underlying()
				}
				vt := t.(interface {
					Elem() types.Type
				}).Elem()
				dim := cdd.Type(w, vt)
				w.WriteByte(' ')
				w.WriteString(dimFuncPtr(cdd.ExprStr(s.Value, nil), dim))
			} else {
				cdd.Expr(w, s.Value, nil)
			}

			w.WriteString(" = " + vs + ";\n")
		}
		if group {
			cdd.indent(w)
		}

		updateEnd(cdd.BlockStmt(w, s.Body, resultT, tup))
		w.WriteByte('\n')

		if label != "" {
			cdd.label(w, label, "_continue")
		}
		if group {
			cdd.il--
			cdd.indent(w)
			w.WriteString("}\n")
		}

		cdd.il--
		cdd.indent(w)
		w.WriteString("}\n")

		if label != "" {
			cdd.label(w, label, "_break")
		}

	case *ast.ReturnStmt:
		updateEnd(cdd.ReturnStmt(w, s, resultT, tup))

	case *ast.SwitchStmt:
		w.WriteString("switch(0){case 0:{\n")
		cdd.il++

		if s.Init != nil {
			cdd.indent(w)
			updateEnd(cdd.Stmt(w, s.Init, "", resultT, tup))
		}

		cdd.indent(w)

		var typ types.Type
		if s.Tag != nil {
			typ = cdd.exprType(s.Tag)
			cdd.varDecl(w, typ, "_tag", s.Tag, "")
			w.WriteByte('\n')
		} else {
			typ = types.Typ[types.Bool]
			w.WriteString("bool _tag = true;\n")
		}

		for _, stmt := range s.Body.List {
			cdd.indent(w)

			cs := stmt.(*ast.CaseClause)
			if len(cs.List) > 0 {
				w.WriteString("if (")
				for i, e := range cs.List {
					if i != 0 {
						w.WriteString(" || ")
					}
					eq(w, "_tag", "==", cdd.ExprStr(e, typ), typ, typ)
				}
				w.WriteString(") ")
			}
			w.WriteString("{\n")
			cdd.il++

			brk := true
			if n := len(cs.Body) - 1; n >= 0 {
				bs, ok := cs.Body[n].(*ast.BranchStmt)
				if ok && bs.Tok == token.FALLTHROUGH {
					brk = false
					cs.Body = cs.Body[:n]
				}
			}
			for _, s := range cs.Body {
				cdd.indent(w)
				updateEnd(cdd.Stmt(w, s, "", resultT, tup))
			}
			if brk {
				cdd.indent(w)
				w.WriteString("break;\n")
			}

			cdd.il--
			cdd.indent(w)
			w.WriteString("}\n")
		}

		cdd.il--
		cdd.indent(w)
		w.WriteString("}}\n")

		if label != "" {
			cdd.label(w, label, "_break")
		}

	case *ast.TypeSwitchStmt:
		w.WriteString("switch(0){case 0:{\n")
		cdd.il++

		if s.Init != nil {
			cdd.indent(w)
			updateEnd(cdd.Stmt(w, s.Init, "", resultT, tup))
		}
		var (
			x   ast.Expr
			lhs string
		)
		switch a := s.Assign.(type) {
		case *ast.ExprStmt:
			x = a.X.(*ast.TypeAssertExpr).X
		case *ast.AssignStmt:
			x = a.Rhs[0].(*ast.TypeAssertExpr).X
			lhs = cdd.ExprStr(a.Lhs[0], nil)
		default:
			panic(a)
		}
		ityp := cdd.exprType(x)
		iempty := (cdd.gtc.methodSet(ityp).Len() == 0)
		cdd.indent(w)
		cdd.varDecl(w, cdd.exprType(x), "_tag", x, "")
		w.WriteByte('\n')
		for _, stmt := range s.Body.List {
			cdd.indent(w)
			cs := stmt.(*ast.CaseClause)
			caseTyp := ityp
			if cs.List != nil {
				w.WriteString("if (")
				for i, e := range cs.List {
					et := cdd.exprType(e)
					if i == 0 {
						if len(cs.List) == 1 && et != untypedNil {
							caseTyp = et
						}
					} else {
						w.WriteString(" || ")
					}
					switch et.Underlying().(type) {
					case *types.Interface:
						if cdd.gtc.methodSet(et).Len() == 0 {
							w.WriteString("true")
							break
						}
						w.WriteString("implements(")
						if iempty {
							w.WriteString("_tag.itab$, &")
						} else {
							w.WriteString("TINFO(_tag), &")
						}
						w.WriteString(cdd.tinameDU(et))
						w.WriteByte(')')
					default:
						if iempty {
							w.WriteString("_tag.itab$ == ")
						} else {
							w.WriteString("TINFO(_tag) == ")
						}
						if et == untypedNil {
							w.WriteString("nil")
						} else {
							w.WriteString("&" + cdd.tinameDU(et))
						}
					}
				}
				w.WriteString(") ")
			}
			w.WriteString("{\n")
			cdd.il++
			if lhs != "" {
				typ, dim := cdd.TypeStr(caseTyp)
				cdd.indent(w)
				w.WriteString(typ + " " + dimFuncPtr(lhs, dim))
				w.WriteString(" = ")
				if _, ok := caseTyp.Underlying().(*types.Interface); ok {
					cdd.interfaceES(w, "_tag", cs.Case, ityp, caseTyp)
				} else {
					w.WriteString("IVAL(_tag, " + typ + dimFuncPtr("", dim) + ")")
				}
				w.WriteString(";\n")
				cdd.indent(w)
				w.WriteString("{\n")
				cdd.il++
			}
			for _, s := range cs.Body {
				cdd.indent(w)
				updateEnd(cdd.Stmt(w, s, "", resultT, tup))
			}
			if lhs != "" {
				cdd.il--
				cdd.indent(w)
				w.WriteString("}\n")
			}
			cdd.indent(w)
			w.WriteString("break;\n")

			cdd.il--
			cdd.indent(w)
			w.WriteString("}\n")
		}

		cdd.il--
		cdd.indent(w)
		w.WriteString("}}\n")

		if label != "" {
			cdd.label(w, label, "_break")
		}

	case *ast.BranchStmt:
		if s.Label == nil {
			w.WriteString(s.Tok.String())
		} else {
			w.WriteString("goto " + s.Label.Name + "$")
			switch s.Tok {
			case token.BREAK:
				w.WriteString("_break")
			case token.CONTINUE:
				w.WriteString("_continue")
			}
		}
		w.WriteString(";\n")

	case *ast.GoStmt:
		cdd.GoStmt(w, s)

	case *ast.SendStmt:
		et := cdd.exprType(s.Chan).(*types.Chan).Elem()
		w.WriteString("SEND(")
		cdd.Expr(w, s.Chan, nil)
		w.WriteString(", ")
		dim := cdd.Type(w, et)
		w.WriteString(dimFuncPtr("", dim))
		w.WriteString(", ")
		cdd.Expr(w, s.Value, et)
		w.WriteString(");\n")

	case *ast.SelectStmt:
		w.WriteString("switch(0){case 0:{\n")
		cdd.il++
		cdd.indent(w)
		w.WriteString("__label__ ")
		dflt := false
		for i, stmt := range s.Body.List {
			if i != 0 {
				w.WriteString(", ")
			}
			if stmt.(*ast.CommClause).Comm == nil {
				dflt = true
				w.WriteString("dflt")
			} else {
				w.WriteString("case" + strconv.Itoa(i))
			}
		}
		w.WriteString(";\n")

		for i, stmt := range s.Body.List {
			switch s := stmt.(*ast.CommClause).Comm.(type) {
			case nil:

			case *ast.SendStmt:
				cdd.indent(w)
				w.WriteString("SENDINIT(" + strconv.Itoa(i) + ", ")
				cdd.Expr(w, s.Chan, nil)
				w.WriteString(", ")
				et := cdd.exprType(s.Chan).(*types.Chan).Elem()
				dim := cdd.Type(w, et)
				dimFuncPtr("", dim)
				w.WriteString(", ")
				cdd.Expr(w, s.Value, et)
				w.WriteString(");\n")

			default:
				cdd.indent(w)
				w.WriteString("RECVINIT(" + strconv.Itoa(i) + ", ")
				var c ast.Expr
				switch r := s.(type) {
				case *ast.AssignStmt:
					c = r.Rhs[0].(*ast.UnaryExpr).X
				case *ast.ExprStmt:
					c = r.X.(*ast.UnaryExpr).X
				default:
					cdd.notImplemented(s)
				}
				cdd.Expr(w, c, nil)
				w.WriteString(", ")
				et := cdd.exprType(c).(*types.Chan).Elem()
				dim := cdd.Type(w, et)
				dimFuncPtr("", dim)
				w.WriteString(");\n")
			}
		}

		cdd.indent(w)
		n := len(s.Body.List)
		if dflt {
			w.WriteString("NBSELECT(\n")
			n--
		} else {
			w.WriteString("SELECT(\n")
		}

		cdd.il++
		for i, stmt := range s.Body.List {
			s := stmt.(*ast.CommClause).Comm
			switch s.(type) {
			case nil:
				continue

			case *ast.SendStmt:
				cdd.indent(w)
				w.WriteString("SENDCOMM(" + strconv.Itoa(i) + ")")

			default:
				cdd.indent(w)
				w.WriteString("RECVCOMM(" + strconv.Itoa(i) + ")")
			}
			if n--; n > 0 {
				w.WriteByte(',')
			}
			w.WriteByte('\n')
		}
		cdd.il--
		cdd.indent(w)
		w.WriteString(");\n")

		for i, stmt := range s.Body.List {
			cc := stmt.(*ast.CommClause)
			s := cc.Comm
			cdd.indent(w)
			switch s.(type) {
			case nil:
				w.WriteString("dflt")
			default:
				w.WriteString("case" + strconv.Itoa(i))
			}
			w.WriteString(":{\n")
			cdd.il++
			switch s := s.(type) {
			case nil:

			case *ast.SendStmt:
				cdd.indent(w)
				w.WriteString("SELSEND(" + strconv.Itoa(i) + ");\n")

			case *ast.AssignStmt:
				cdd.indent(w)
				name := cdd.ExprStr(s.Lhs[0], nil)
				if len(s.Lhs) == 1 {
					if name != "_$" {
						if s.Tok == token.DEFINE {
							dim := cdd.Type(w, cdd.exprType(s.Rhs[0]))
							w.WriteString(" " + dimFuncPtr(name, dim))
						} else {
							w.WriteString(name)
						}
						w.WriteString(" = ")
					}
					w.WriteString("SELRECV(" + strconv.Itoa(i) + ");\n")
				} else {
					ok := cdd.ExprStr(s.Lhs[1], nil)
					tmp := ""
					var tup *types.Tuple
					if name != "_$" || ok != "_$" {
						tup = cdd.exprType(s.Rhs[0]).(*types.Tuple)
						tupName, _ := cdd.tupleName(tup)
						w.WriteString(tupName + " ")
						tmp = "_tmp" + cdd.gtc.uniqueId()
						w.WriteString(tmp + " = ")
					}
					w.WriteString("SELRECVOK(" + strconv.Itoa(i) + ");\n")
					if name != "_$" {
						cdd.indent(w)
						if s.Tok == token.DEFINE {
							dim := cdd.Type(w, tup.At(0).Type())
							w.WriteString(" " + dimFuncPtr(name, dim))
						} else {
							w.WriteString(name)
						}
						w.WriteString(" = " + tmp + "._0;\n")
					}
					if ok != "_$" {
						cdd.indent(w)
						if s.Tok == token.DEFINE {
							w.WriteString("bool ")
						}
						w.WriteString(ok + " = " + tmp + "._1;\n")
					}
				}

			case *ast.ExprStmt:
				cdd.indent(w)
				w.WriteString("SELRECV(" + strconv.Itoa(i) + ");\n")

			default:
				cdd.notImplemented(s)
			}
			for _, s = range cc.Body {
				cdd.indent(w)
				updateEnd(cdd.Stmt(w, s, "", resultT, tup))
			}
			cdd.indent(w)
			w.WriteString("break;\n")
			cdd.il--
			cdd.indent(w)
			w.WriteString("}\n")
		}
		cdd.il--
		cdd.indent(w)
		w.WriteString("}}\n")

	default:
		cdd.notImplemented(s)
	}
	return
}

type arg struct {
	t types.Type
	l string
	r string
}

type call struct {
	rcv  arg
	fun  arg
	arr  arg
	tup  arg
	args []arg
}

func (cdd *CDD) call(e *ast.CallExpr, t *types.Signature, eval bool) *call {
	c := new(call)
	n := len(e.Args) + 1 // +1 for variadic function without any parameter.
	fs, ft, rs, rt := cdd.funStr(e.Fun, e.Args)
	if t != nil {
		ft = t
	}
	if fs == "" {
		panic("fs == \"\", rs == \"" + rs + "\"")
	}
	if rs != "" {
		n++
	}
	ri := false
	if rt != nil {
		_, ri = rt.Underlying().(*types.Interface)
	}
	c.args = make([]arg, n)
	n = 0
	if ri {
		in, ok := rt.(*types.Named)
		if !ok {
			panic("unimplemented: call method of unnamed interface")
		}
		// Interface receiver
		cast := "((" + cdd.NameStr(in.Obj(), false) + "*)"
		if _, ok := e.Fun.(*ast.SelectorExpr).X.(*ast.Ident); ok && !eval {
			c.fun.l = cast + "(" + rs + ".itab$))->" + fs
			c.args[n] = arg{types.Typ[types.Uintptr], "&" + rs + ".val$", ""}
		} else {
			c.rcv = arg{rt, "_r", rs}
			c.fun.l = cast + "_r.itab$)->" + fs
			c.args[n] = arg{types.Typ[types.Uintptr], "&_r" + ".val$", ""}
		}
		n++
	} else if rs == "" {
		if eval {
			// Call of function or function variable.
			c.fun = arg{ft, "_f", fs}
			if fident, ok := e.Fun.(*ast.Ident); ok {
				if _, ok = cdd.object(fident).(*types.Var); !ok {
					// Ordinary function call
					c.fun = arg{nil, fs, ""}
				}
			}
		} else {
			c.fun.l = fs
		}
	} else {
		// Method call.
		if eval {
			c.rcv = arg{rt, "_r", rs}
			c.args[n] = arg{rt, "_r", ""}
		} else {
			c.args[n].l = rs
		}
		n++
		c.fun.l = fs
	}
	sig := ft.(*types.Signature)
	tup := sig.Params()
	alen := tup.Len()
	if len(e.Args) == 1 {
		a0 := e.Args[0]
		if atup, _ := cdd.exprType(a0).(*types.Tuple); atup != nil {
			c.tup.t = atup
			c.tup.l = "_tup"
			c.tup.r = cdd.ExprStr(a0, c.tup.t)
			for i := 0; i < alen; i++ {
				it := tup.At(i).Type()
				et := atup.At(i).Type()
				ai := "_tup._" + strconv.Itoa(i)
				s := cdd.interfaceESstr(ai, a0.Pos(), et, it)
				if eval || c.arr.t != nil {
					c.args[n] = arg{it, "_" + strconv.Itoa(i), s}
				} else {
					c.args[n].l = s
				}
				n++
			}
			c.args = c.args[:n]
			return c
		}
	}
	variadic := sig.Variadic() && !e.Ellipsis.IsValid()
	if variadic {
		c.arr.t = tup.At(alen - 1).Type().(*types.Slice).Elem()
		c.arr.l = "_a[]"
	}
	for i, a := range e.Args {
		if a == nil {
			// builtin can set type args to nil
			continue
		}
		if variadic && i >= alen-1 {
			if c.arr.r != "" {
				c.arr.r += ", "
			}
			c.arr.r += cdd.interfaceExprStr(a, c.arr.t)
			continue
		}
		var at types.Type
		if i < alen {
			at = tup.At(i).Type()
		} else {
			// Builtin functions may not spefify type of all parameters.
			at = cdd.exprType(a)
		}
		s := cdd.interfaceExprStr(a, at)
		if eval || c.arr.t != nil {
			c.args[n] = arg{at, "_" + strconv.Itoa(i), s}
		} else {
			c.args[n].l = s
		}
		n++
	}
	if c.arr.t != nil {
		if c.arr.r == "" {
			c.args[n].l = "NILSLICE"
			c.arr = arg{}
			if !eval {
				argv := c.args[:n]
				if rs != "" {
					argv = argv[1:]
				}
				for i, a := range argv {
					argv[i] = arg{nil, a.r, ""}
				}
			}
		} else {
			c.args[n].l = "CSLICE(" + strconv.Itoa(len(e.Args)-alen+1) + ", _a)"
			c.arr.r = "{" + c.arr.r + "}"
		}
		n++
	}
	c.args = c.args[:n]
	return c
}

func (cdd *CDD) GoStmt(w *bytes.Buffer, s *ast.GoStmt) {
	c := cdd.call(s.Call, nil, true)

	if c.fun.r == "" && len(c.args) == 0 {
		// Fast path: ordinary function without parameters.
		w.WriteString("GO(" + c.fun.l + "(), false);\n")
		return
	}

	argv := c.args
	if c.fun.r != "" {
		argv = append([]arg{c.fun}, c.args...)
	}

	w.WriteString("{\n")
	cdd.il++

	cdd.indent(w)
	w.WriteString("void wrap(")
	comma := false
	for i, arg := range argv {
		if arg.t == nil {
			continue // const. expr.
		}
		if comma {
			w.WriteString(", ")
		} else {
			comma = true
		}
		if i == len(argv)-1 && c.arr.t != nil {
			// Variadic function.
			arg = c.arr
		}
		dim := cdd.Type(w, arg.t)
		w.WriteString(" " + dimFuncPtr(arg.l, dim))
	}
	w.WriteString(") {\n")
	cdd.il++
	cdd.indent(w)
	w.WriteString("goready();\n")
	cdd.indent(w)
	w.WriteString(c.fun.l + "(")
	for i, arg := range c.args {
		if i > 0 {
			w.WriteString(", ")
		}
		w.WriteString(arg.l)
	}
	w.WriteString(");\n")
	cdd.il--
	cdd.indent(w)
	w.WriteString("}\n")
	if c.tup.t != nil {
		cdd.indent(w)
		cdd.Type(w, c.tup.t)
		w.WriteString(" " + c.tup.l + " = " + indent(1, c.tup.r) + ";\n")
	}
	if c.rcv.r != "" {
		argv = append([]arg{c.rcv}, argv...)
	}
	for i, arg := range argv {
		if i == len(argv)-1 && c.arr.t != nil {
			// Variadic function.
			cdd.indent(w)
			dim := cdd.Type(w, c.arr.t)
			w.WriteString(" " + dimFuncPtr(c.arr.l, dim) + " = ")
			w.WriteString(indent(1, c.arr.r) + ";\n")
		}
		if arg.r == "" {
			continue // Don't evaluate
		}
		cdd.indent(w)
		dim := cdd.Type(w, arg.t)
		w.WriteString(" " + dimFuncPtr(arg.l, dim) + " = ")
		w.WriteString(indent(1, arg.r) + ";\n")
	}
	if c.rcv.r != "" {
		argv = argv[1:]
	}
	cdd.indent(w)
	w.WriteString("GO(wrap(")
	comma = false
	for i, arg := range argv {
		if arg.t == nil {
			continue // const. expr.
		}
		if comma {
			w.WriteString(", ")
		} else {
			comma = true
		}
		if i == len(argv)-1 && c.arr.t != nil {
			// Variadic function.
			arg = c.arr
		}
		w.WriteString(arg.l)
	}
	w.WriteString("), true);\n")

	cdd.il--
	cdd.indent(w)
	w.WriteString("}\n")

	return
}

func (cdd *CDD) BlockStmt(w *bytes.Buffer, bs *ast.BlockStmt, resultT string, tup *types.Tuple) (end bool) {
	updateEnd := func(e bool) {
		if e {
			end = true
		}
	}

	w.WriteString("{\n")
	cdd.il++
	for _, stmt := range bs.List {
		switch s := stmt.(type) {
		case *ast.LabeledStmt:
			label := s.Label.Name + "$"
			cdd.label(w, label, "")
			cdd.indent(w)
			updateEnd(cdd.Stmt(w, s.Stmt, label, resultT, tup))

		default:
			m := w.Len()
			cdd.indent(w)
			n := w.Len()
			updateEnd(cdd.Stmt(w, s, "", resultT, tup))
			if w.Len() == n {
				w.Truncate(m)
			}
		}
	}
	cdd.il--
	cdd.indent(w)
	w.WriteString("}")
	return
}
