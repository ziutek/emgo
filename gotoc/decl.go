package gotoc

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"strconv"
)

func (cc *CC) Tuple(w *bytes.Buffer, t *types.Tuple, sep string) {
	for i, n := 0, t.Len(); i < n; i++ {
		if i != 0 {
			w.WriteString(sep)
		}
		v := t.At(i)
		cc.Type(w, v.Type())
		w.WriteByte(' ')
		w.WriteString(v.Name())
	}
}

type retVal struct {
	name, typ string
}

func (cc *CC) FuncDecl(d *ast.FuncDecl) {
	f := cc.ti.Objects[d.Name].(*types.Func)

	funcName := cc.NameStr(f)
	s := f.Type().(*types.Signature)

	var wh *bytes.Buffer

	if f.IsExported() || funcName == "main_main" {
		b := d.Body
		d.Body = nil
		printer.Fprint(&cc.wg, cc.fset, d)
		cc.wg.WriteByte('\n')
		d.Body = b

		wh = &cc.wh
	} else {
		wh = &cc.ws
	}
	wc := &cc.wc

	var (
		retT        string
		retList     []retVal
		hasRetNames bool
	)

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
			cc.indent(wh)
			wh.WriteString("typedef struct {\n")
			cc.il++
			for _, v := range retList {
				cc.indent(wh)
				wh.WriteString(v.typ)
				wh.WriteByte(' ')
				wh.WriteString(v.name)
				wh.WriteString(";\n")
			}
			cc.il--
			cc.indent(wh)
			wh.WriteString("} " + retT + ";\n")
		}
	}

	if wh == &cc.ws {
		wh.WriteString("static ") // not exported function
	}

	buf := new(bytes.Buffer)

	cc.indent(buf)
	switch len(retList) {
	case 0:
		buf.WriteString("void")

	case 1:
		buf.WriteString(retList[0].typ)

	default:
		buf.WriteString(retT)
	}

	buf.WriteString(" " + funcName + "(")
	if r := s.Recv(); r != nil {
		cc.Type(buf, r.Type())
		buf.WriteByte(' ')
		buf.WriteString(r.Name())
		if s.Params() != nil {
			buf.WriteString(", ")
		}
	}
	if p := s.Params(); p != nil {
		cc.Tuple(buf, p, ", ")
	}
	buf.WriteByte(')')

	wh.Write(buf.Bytes())
	wh.WriteString(";\n")

	if d.Body == nil {
		return
	}

	buf.WriteTo(wc)
	wc.WriteByte(' ')

	if hasRetNames {
		cc.indent(wc)
		wc.WriteString("{\n")
		cc.il++
		for _, v := range retList {
			cc.indent(wc)
			wc.WriteString(v.typ + " " + v.name + " = {0};\n")
		}
		cc.indent(wc)

		if retT == "" {
			// Inform ReturnStmt that there is one named result
			retT = "_"
		}
	}

	cc.BlockStmt(wc, d.Body, retT)
	wc.WriteByte('\n')

	if hasRetNames {
		cc.il--
		cc.indent(wc)
		wc.WriteString("__end:\n")
		cc.il++

		cc.indent(wc)
		wc.WriteString("return ")
		if len(retList) == 1 {
			wc.WriteString(retList[0].name)
		} else {
			wc.WriteString("(" + retT + ") {")
			for i, v := range retList {
				if i > 0 {
					wc.WriteString(", ")
				}
				wc.WriteString(v.name)
			}
			wc.WriteByte('}')
		}
		wc.WriteString(";\n")

		cc.il--
		wc.WriteString("}\n")
	}
}

func (cc *CC) GenDecl(d *ast.GenDecl) {
	wc := &cc.wc
	wg := &cc.wg

	switch d.Tok {
	case token.IMPORT:
		// Imports are handled at higher level
		break

	case token.CONST:
		for _, s := range d.Specs {
			vs := s.(*ast.ValueSpec)
			vals := vs.Values

			for i, n := range vs.Names {
				c := cc.ti.Objects[n].(*types.Const)
				eval := c.Val().String()

				exported := n.IsExported() && cc.isGlobal(c)

				// Go header
				if exported {
					wg.WriteString("const ")
					wg.WriteString(n.Name)
					wg.WriteByte(' ')
					if cc.GoType(wg, c.Type()) {
						wg.WriteByte(' ')
					}
					wg.WriteString("= ")
					wg.WriteString(eval)

					if cc.OriginComments {
						// Comment about original value
						if i < len(vals) {
							switch v := vals[i].(type) {
							case *ast.BasicLit:
								// It was written before
							default:
								wg.WriteString(" // = ")
								printer.Fprint(wg, cc.fset, v)
							}
						}
					}

					wg.WriteByte('\n')
				}

				// C header
				if exported {
					wh := &cc.wh
					cc.indent(wh)
					wh.WriteString("#define ")
					cc.Name(wh, c)
					wh.WriteByte(' ')
					wh.WriteString(eval)
					wh.WriteByte('\n')
				}

			}
		}

	case token.VAR:
		buf := new(bytes.Buffer)
		for _, s := range d.Specs {
			vs := s.(*ast.ValueSpec)
			vals := vs.Values

			for i, n := range vs.Names {
				v := cc.ti.Objects[n].(*types.Var)
				typ := v.Type()

				if cc.isGlobal(v) && n.IsExported() {
					wg.WriteString("var ")
					wg.WriteString(n.Name)
					wg.WriteByte(' ')
					cc.GoType(wg, typ)
					if cc.OriginComments {
						// Comment about original initial value
						if i < len(vals) {
							wg.WriteString(" // = ")
							printer.Fprint(wg, cc.fset, vals[i])
						}
					}
					wg.WriteByte('\n')
				}

				cc.Type(buf, typ)
				buf.WriteByte(' ')
				cc.Name(buf, v)

				var wh *bytes.Buffer

				if cc.isGlobal(v) {
					if n.IsExported() {
						wh = &cc.wh

						wh.WriteString("extern ")
					} else {
						wh = &cc.ws

						wh.WriteString("static ")
						wc.WriteString("static ")
					}
					wh.Write(buf.Bytes())
					wh.WriteString(";\n")
				}

				cc.indent(wc)
				buf.WriteTo(wc)
				wc.WriteString(" = ")
				if i < len(vals) {
					cc.Expr(wc, vals[i])
				} else {
					wc.WriteString("{0}")
				}
				wc.WriteString(";\n")
			}
		}

	case token.TYPE:
		for _, s := range d.Specs {
			ts := s.(*ast.TypeSpec)

			var wh *bytes.Buffer
			if ts.Name.IsExported() {
				wg := &cc.wg
				wg.WriteString("type ")
				printer.Fprint(wg, cc.fset, ts)
				wg.WriteByte('\n')

				wh = &cc.wh
			} else {
				wh = &cc.ws
			}

			cc.indent(wh)
			wh.WriteString("typedef ")
			cc.Type(wh, cc.ti.Types[ts.Type])
			wh.WriteByte(' ')
			cc.Name(wh, cc.ti.Objects[ts.Name])
			wh.WriteString(";\n")
		}

	default:
		fmt.Fprintf(&cc.wh, "%%%s%%", d.Tok)
	}

}

func (cc *CC) Decl(decl ast.Decl) {
	switch d := decl.(type) {
	case *ast.FuncDecl:
		cc.FuncDecl(d)

	case *ast.GenDecl:
		cc.GenDecl(d)

	default:
		fmt.Fprintf(&cc.wc, "@%v<%T>@", d, d)
	}
}
