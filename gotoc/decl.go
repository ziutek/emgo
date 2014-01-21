package gotoc

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"go/ast"
	"go/token"
	"strconv"
)

type retVal struct {
	name, typ string
}

func (cc *GTC) FuncDecl(d *ast.FuncDecl) *CDD {
	f := cc.ti.Objects[d.Name].(*types.Func)
	s := f.Type().(*types.Signature)

	cdd := newCDD(f, FuncDecl)

	var (
		retT        string
		retList     []retVal
		hasRetNames bool
	)

	w := new(bytes.Buffer)
	funcName := cc.NameStr(f)
	if ret := s.Results(); ret != nil {
		retList = make([]retVal, ret.Len())
		for i, n := 0, ret.Len(); i < n; i++ {
			v := ret.At(i)
			n := v.Name()
			if n == "" {
				n = "__" + strconv.Itoa(i)
			} else {
				hasRetNames = true
			}
			retList[i] = retVal{
				name: n,
				typ:  cc.TypeStr(v.Type()),
			}
		}

		if len(retList) > 1 {
			retT = "__" + funcName
			cc.indent(w)
			w.WriteString("typedef struct {\n")
			cc.il++
			for _, v := range retList {
				cc.indent(w)
				w.WriteString(v.typ)
				w.WriteByte(' ')
				w.WriteString(v.name)
				w.WriteString(";\n")
			}
			cc.il--
			cc.indent(w)
			w.WriteString("} " + retT + ";\n")
		}
	}

	cc.indent(w)
	switch len(retList) {
	case 0:
		w.WriteString("void")

	case 1:
		w.WriteString(retList[0].typ)

	default:
		w.WriteString(retT)
	}

	w.WriteString(" " + funcName + "(")
	if r := s.Recv(); r != nil {
		cc.Type(w, r.Type())
		w.WriteByte(' ')
		w.WriteString(r.Name())
		if s.Params() != nil {
			w.WriteString(", ")
		}
	}
	if p := s.Params(); p != nil {
		cc.Tuple(w, p, ", ")
	}
	w.WriteByte(')')

	cdd.copyDecl(w, ";\n")

	if d.Body == nil {
		return cdd
	}

	w.WriteByte(' ')

	if hasRetNames {
		cc.indent(w)
		w.WriteString("{\n")
		cc.il++
		for _, v := range retList {
			cc.indent(w)
			w.WriteString(v.typ + " " + v.name + " = {0};\n")
		}
		cc.indent(w)

		if retT == "" {
			// Inform ReturnStmt that there is one named result
			retT = "_"
		}
	}

	cc.BlockStmt(w, d.Body, retT)
	w.WriteByte('\n')

	if hasRetNames {
		cc.il--
		cc.indent(w)
		w.WriteString("__end:\n")
		cc.il++

		cc.indent(w)
		w.WriteString("return ")
		if len(retList) == 1 {
			w.WriteString(retList[0].name)
		} else {
			w.WriteString("(" + retT + ") {")
			for i, v := range retList {
				if i > 0 {
					w.WriteString(", ")
				}
				w.WriteString(v.name)
			}
			w.WriteByte('}')
		}
		w.WriteString(";\n")

		cc.il--
		w.WriteString("}\n")
	}
	cdd.copyDef(w)
	return cdd
}

func (cc *GTC) GenDecl(d *ast.GenDecl) (cdds []*CDD) {
	w := new(bytes.Buffer)

	switch d.Tok {
	case token.IMPORT:
		// Imports are handled at higher level
		break

	case token.CONST:
		for _, s := range d.Specs {
			vs := s.(*ast.ValueSpec)

			for _, n := range vs.Names {
				c := cc.ti.Objects[n].(*types.Const)
				eval := c.Val().String()
				cdd := newCDD(c, ConstDecl)

				// All constants in expressions are evaluated so
				// only exported constants need be translated to C
				if c.IsExported() {
					cc.indent(w)
					w.WriteString("#define ")
					cc.Name(w, c)
					w.WriteByte(' ')
					w.WriteString(eval)

					cdd.copyDecl(w, "\n")
					w.Reset()
				}
				cdds = append(cdds, cdd)
			}
		}

	case token.VAR:
		for _, s := range d.Specs {
			vs := s.(*ast.ValueSpec)
			vals := vs.Values

			for i, n := range vs.Names {
				v := cc.ti.Objects[n].(*types.Var)
				typ := v.Type()
				cdd := newCDD(v, VarDecl)

				cc.indent(w)
				cc.Type(w, typ)
				w.WriteByte(' ')
				cc.Name(w, v)
				cdd.copyDecl(w, ";\n")

				w.WriteString(" = ")
				if i < len(vals) {
					cc.Expr(w, vals[i])
				} else {
					w.WriteString("{0}")
				}
				w.WriteString(";\n")
				
				cdd.copyDef(w)
				w.Reset()

				cdds = append(cdds, cdd)
			}
		}

	case token.TYPE:
		for _, s := range d.Specs {
			ts := s.(*ast.TypeSpec)
			t := cc.ti.Objects[ts.Name]
			cdd := newCDD(t, TypeDecl)

			cc.indent(w)
			w.WriteString("typedef ")
			cc.Type(w, cc.ti.Types[ts.Type])
			w.WriteByte(' ')
			cc.Name(w, t)
			
			cdd.copyDecl(w, ";\n")
			w.Reset()
			
			cdds = append(cdds, cdd)
		}

	default:
		// Return fake CDD for unknown declaration
		cdds = []*CDD{{
			Decl: []byte(fmt.Sprintf("@%v (%T)@\n", d.Tok, d)),
		}}
	}
	return
}

func (cc *GTC) Decl(decl ast.Decl) []*CDD {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return []*CDD{cc.FuncDecl(d)}

	case *ast.GenDecl:
		return cc.GenDecl(d)
	}

	panic(fmt.Sprint("Unknown declaration: ", decl))
}
