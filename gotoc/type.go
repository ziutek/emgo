package gotoc

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"

	"code.google.com/p/go.tools/go/types"
)

func dimFuncPtr(name string, dim []string) string {
	if len(dim) == 0 {
		return name
	}
	for i, d := range dim {
		switch d[0] {
		case '*':
			if i == len(dim)-1 || dim[i+1][0] == '*' {
				name = "*" + name
			} else {
				name = "(*" + name + ")"
			}

		case '[':
			name += d

		default:
			name = "(*" + name + ")" + d
		}
	}
	return name
}

func (cdd *CDD) Type(w *bytes.Buffer, typ types.Type) (dim []string, acds []*CDD) {
	direct := true

writeType:
	switch t := typ.(type) {
	case *types.Basic:
		if t.Kind() == types.UnsafePointer {
			w.WriteString("unsafe_Pointer")
		} else {
			types.WriteType(w, nil, t)
		}

	case *types.Named:
		cdd.Name(w, t.Obj(), direct)

	case *types.Pointer:
		typ = t.Elem()
		direct = false
		dim = append(dim, "*")
		goto writeType

	case *types.Struct:
		w.WriteString("struct {\n")
		cdd.il++
		for i, n := 0, t.NumFields(); i < n; i++ {
			f := t.Field(i)
			cdd.indent(w)
			if tag := t.Tag(i); tag != "" {
				w.WriteString(reflect.StructTag(tag).Get("C"))
				w.WriteByte(' ')
			}
			d, a := cdd.Type(w, f.Type())
			acds = append(acds, a...)
			if !f.Anonymous() {
				w.WriteByte(' ')
				name := cdd.NameStr(f, true)
				w.WriteString(dimFuncPtr(name, d))
			}
			w.WriteString(";\n")
		}
		cdd.il--
		cdd.indent(w)
		w.WriteByte('}')

	case *types.Array:
		dim = append(dim, "["+strconv.FormatInt(t.Len(), 10)+"]")
		d, a := cdd.Type(w, t.Elem())
		dim = append(dim, d...)
		acds = append(acds, a...)

	case *types.Slice:
		w.WriteString("__slice")

	case *types.Signature:
		res, params := cdd.signature(t)
		w.WriteString(res.typ)
		dim = append(dim, params)
		dim = append(dim, res.dim...)
		acds = append(acds, res.acds...)

	default:
		fmt.Fprintf(w, "<%T>", t)
	}
	return
}

func (cdd *CDD) TypeStr(typ types.Type) (string, []string, []*CDD) {
	buf := new(bytes.Buffer)
	dim, acds := cdd.Type(buf, typ)
	return buf.String(), dim, acds
}

type retVar struct {
	name, typ string
}

type results struct {
	fields   []*types.Var
	names    []string
	typ      string
	dim      []string
	hasNames bool
	acds     []*CDD
}

func (cdd *CDD) results(tup *types.Tuple) (res results) {
	if tup == nil {
		res.typ = "void"
		return
	}

	n := tup.Len()
	res.fields = make([]*types.Var, n)
	res.names = make([]string, n)

	for i := 0; i < n; i++ {
		v := tup.At(i)
		n := strconv.Itoa(i)
		res.fields[i] = types.NewField(v.Pos(), v.Pkg(), "_"+n, v.Type(), false)

		name := v.Name()
		if name == "" {
			name = "__" + n
		} else {
			res.hasNames = true
		}
		res.names[i] = name
	}

	if n == 1 {
		res.typ, res.dim, res.acds = cdd.TypeStr(res.fields[0].Type())
		return
	}

	var declared bool
	res.typ, declared = cdd.gtc.tn.DeclName(tup)

	if !declared {
		s := types.NewStruct(res.fields, nil)
		o := types.NewTypeName(tup.At(0).Pos(), cdd.gtc.pkg, res.typ, s)

		acd := cdd.gtc.newCDD(o, TypeDecl, cdd.il)
		acd.structDecl(new(bytes.Buffer), res.typ, s)
		res.acds = append(res.acds, acd)

		cdd.DeclUses[o] = true
		cdd.BodyUses[o] = true
	}
	return
}

func (cdd *CDD) signature(sig *types.Signature) (res results, params string) {
	params = "("
	res = cdd.results(sig.Results())
	if r := sig.Recv(); r != nil {
		typ, dim, acds := cdd.TypeStr(r.Type())
		res.acds = append(res.acds, acds...)
		pname := cdd.NameStr(r, true)
		params += typ + " " + dimFuncPtr(pname, dim)
		if sig.Params() != nil {
			params += ", "
		}
	}
	if p := sig.Params(); p != nil {
		for i, n := 0, p.Len(); i < n; i++ {
			if i != 0 {
				params += ", "
			}
			v := p.At(i)
			typ, dim, acds := cdd.TypeStr(v.Type())
			res.acds = append(res.acds, acds...)
			name := cdd.NameStr(v, true)
			if name != "" {
				params += typ + " " + dimFuncPtr(name, dim)
			} else {
				params += typ + dimFuncPtr("", dim)
			}
		}
	}
	params += ")"
	return
}
