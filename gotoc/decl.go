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

func (gtc *GTC) FuncDecl(d *ast.FuncDecl, il int) *CDD {
	f := gtc.ti.Objects[d.Name].(*types.Func)
	s := f.Type().(*types.Signature)

	cdd := gtc.newCDD(f, FuncDecl, il)

	var (
		retT        string
		retList     []retVal
		hasRetNames bool
	)

	w := new(bytes.Buffer)
	funcName := cdd.NameStr(f)
	if ret := s.Results(); ret != nil {
		retList = make([]retVal, ret.Len())
		for i, n := 0, ret.Len(); i < n; i++ {
			v := ret.At(i)
			n := cdd.NameStr(v)
			if n == "" {
				n = "__" + strconv.Itoa(i)
			} else {
				hasRetNames = true
			}
			retList[i] = retVal{
				name: n,
				typ:  cdd.TypeStr(v.Type()),
			}
		}

		if len(retList) > 1 {
			retT = "__" + funcName
			cdd.indent(w)
			w.WriteString("typedef struct {\n")
			cdd.il++
			for _, v := range retList {
				cdd.indent(w)
				w.WriteString(v.typ)
				w.WriteByte(' ')
				w.WriteString(v.name)
				w.WriteString(";\n")
			}
			cdd.il--
			cdd.indent(w)
			w.WriteString("} " + retT + ";\n")
		}
	}

	cdd.indent(w)
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
		cdd.Type(w, r.Type())
		w.WriteByte(' ')
		cdd.Name(w, r)
		if s.Params() != nil {
			w.WriteString(", ")
		}
	}
	if p := s.Params(); p != nil {
		cdd.Tuple(w, p, ", ")
	}
	w.WriteByte(')')

	cdd.copyDecl(w, ";\n")

	if d.Body == nil {
		return cdd
	}

	cdd.body = true

	w.WriteByte(' ')

	if hasRetNames {
		cdd.indent(w)
		w.WriteString("{\n")
		cdd.il++
		for _, v := range retList {
			cdd.indent(w)
			w.WriteString(v.typ + " " + v.name + " = {0};\n")
		}
		cdd.indent(w)

		if retT == "" {
			// Inform ReturnStmt that there is one named result
			retT = "_"
		}
	}

	cdd.BlockStmt(w, d.Body, retT)
	w.WriteByte('\n')

	if hasRetNames {
		cdd.il--
		cdd.indent(w)
		w.WriteString("__end:\n")
		cdd.il++

		cdd.indent(w)
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

		cdd.il--
		w.WriteString("}\n")
	}
	cdd.copyDef(w)
	return cdd
}

func (gtc *GTC) GenDecl(d *ast.GenDecl, il int) (cdds []*CDD) {
	w := new(bytes.Buffer)

	switch d.Tok {
	case token.IMPORT:
		// Imports are handled differently
		break

	case token.CONST:
		for _, s := range d.Specs {
			vs := s.(*ast.ValueSpec)

			for _, n := range vs.Names {
				c := gtc.ti.Objects[n].(*types.Const)
				eval := c.Val().String()
				cdd := gtc.newCDD(c, ConstDecl, il)

				// All constants in expressions are evaluated so
				// only exported constants need be translated to C
				if c.IsExported() {
					cdd.indent(w)
					w.WriteString("#define ")
					cdd.Name(w, c)
					w.WriteByte(' ')
					w.WriteString(eval)

					cdd.copyDecl(w, "\n")
					w.Reset()

					cdds = append(cdds, cdd)
				}
			}
		}

	case token.VAR:
		for _, s := range d.Specs {
			vs := s.(*ast.ValueSpec)
			vals := vs.Values

			for i, n := range vs.Names {
				v := gtc.ti.Objects[n].(*types.Var)
				typ := v.Type()
				cdd := gtc.newCDD(v, VarDecl, il)
				name := cdd.NameStr(v)

				cdd.indent(w)
				cdd.Type(w, typ)
				w.WriteByte(' ')
				w.WriteString(name)
				
				cinit := true // true if C declaration can init value
				
				if cdd.gtc.isGlobal(v){
					cdd.copyDecl(w, ";\n") // Global variables need declaration
					if i < len(vals) {
						_, cinit = cdd.gtc.ti.Values[vals[i]]
					}
				}
				if cinit {
					w.WriteString(" = ")
					if i < len(vals) {
						cdd.Expr(w, vals[i])
					} else {
						w.WriteString("{0}")
					}
				}
				w.WriteString(";\n")
				cdd.copyDef(w)

				if !cinit {
					// Runtime initialisation
					w.Reset()
					w.WriteByte('\t')
					w.WriteString(name)
					w.WriteString(" = ")
					cdd.Expr(w, vals[i])
					w.WriteString(";\n")
					cdd.copyInit(w)
				}

				w.Reset()

				cdds = append(cdds, cdd)
			}
		}

	case token.TYPE:
		// TODO: split struct types to Decl and Def
		for _, s := range d.Specs {
			ts := s.(*ast.TypeSpec)
			t := gtc.ti.Objects[ts.Name]
			cdd := gtc.newCDD(t, TypeDecl, il)

			cdd.indent(w)
			w.WriteString("typedef ")
			cdd.Type(w, gtc.ti.Types[ts.Type])
			w.WriteByte(' ')
			cdd.Name(w, t)

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

func (cc *GTC) Decl(decl ast.Decl, il int) []*CDD {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		return []*CDD{cc.FuncDecl(d, il)}

	case *ast.GenDecl:
		return cc.GenDecl(d, il)
	}

	panic(fmt.Sprint("Unknown declaration: ", decl))
}
