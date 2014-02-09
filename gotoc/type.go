package gotoc

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"reflect"
	"strconv"
)

func writeDim(w *bytes.Buffer, dim []int64) {
	for _, d := range dim {
		w.WriteByte('[')
		w.WriteString(strconv.FormatInt(d, 10))
		w.WriteByte(']')
	}
}

func writeStars(w *bytes.Buffer, dim []int64) {
	for _ = range dim {
		w.WriteByte('*')
	}
}

func (cdd *CDD) Type(w *bytes.Buffer, typ types.Type) (dim []int64) {
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
		defer w.WriteByte('*')
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
			d := cdd.Type(w, f.Type())
			if !f.Anonymous() {
				w.WriteByte(' ')
				cdd.Name(w, f, true)
				writeDim(w, d)
			}
			w.WriteString(";\n")
		}
		cdd.il--
		cdd.indent(w)
		w.WriteByte('}')

	case *types.Array:
		dim = append(dim, t.Len())
		dim = append(dim, cdd.Type(w, t.Elem())...)

	case *types.Slice:
		w.WriteString("__slice")

	default:
		fmt.Fprintf(w, "<%T>", t)
	}
	return
}

func (cdd *CDD) TypeStr(typ types.Type) string {
	buf := new(bytes.Buffer)
	cdd.Type(buf, typ)
	return buf.String()
}

func (cdd *CDD) Tuple(w *bytes.Buffer, t *types.Tuple) {
	for i, n := 0, t.Len(); i < n; i++ {
		if i != 0 {
			w.WriteString(", ")
		}
		v := t.At(i)
		cdd.Type(w, v.Type())
		w.WriteByte(' ')
		cdd.Name(w, v, true)
	}
}

type retVar struct {
	name, typ string
}

type results struct {
	fields []*types.Var
	//list     []retVar
	typ      string
	hasNames bool
	cdd      *CDD
}

/*func (res *results) writeStruct() {
	w := new(bytes.Buffer)
	cdd := res.cdd
	cdd.indent(w)
	w.WriteString("typedef struct {\n")
	cdd.il++
	for _, v := range res.list {
		cdd.indent(w)
		w.WriteString(v.typ)
		w.WriteByte(' ')
		w.WriteString(v.name)
		w.WriteString(";\n")
	}
	cdd.il--
	cdd.indent(w)
	w.WriteString("} " + res.typ + ";\n")

	cdd.copyDef(w)
}*/

func (cdd *CDD) results(tup *types.Tuple, fname string) (res results) {
	if tup == nil {
		res.typ = "void"
		return
	}

	n := tup.Len()
	//res.list = make([]retVar, n)
	res.fields = make([]*types.Var, n)

	for i := 0; i < n; i++ {
		v := tup.At(i)
		name := v.Name()
		if name == "" {
			name = "__" + strconv.Itoa(i)
		} else {
			res.hasNames = true
		}
		res.fields[i] = types.NewField(v.Pos(), v.Pkg(), name, v.Type(), false)
	}

	if n == 1 {
		res.typ = cdd.TypeStr(res.fields[0].Type())
		return
	}

	res.typ = "__" + fname
	s := types.NewStruct(res.fields, nil)
	o := types.NewTypeName(tup.At(0).Pos(), cdd.gtc.pkg, res.typ, s)
	res.cdd = cdd.gtc.newCDD(o, TypeDecl, cdd.il)

	res.cdd.structDecl(new(bytes.Buffer), res.typ, s)

	cdd.DeclUses[o] = true
	cdd.BodyUses[o] = true
	return
}

func (cdd *CDD) Signature(w *bytes.Buffer, name string, sig *types.Signature, decl bool) (res results) {
	res = cdd.results(sig.Results(), name)

	w.WriteString(res.typ)
	w.WriteByte(' ')
	if decl {
		w.WriteString(name)
	} else {
		w.WriteString("(*" + name + ")")
	}
	w.WriteByte('(')
	if r := sig.Recv(); r != nil {
		cdd.Type(w, r.Type())
		w.WriteByte(' ')
		cdd.Name(w, r, true)
		if sig.Params() != nil {
			w.WriteString(", ")
		}
	}
	if p := sig.Params(); p != nil {
		cdd.Tuple(w, p)
	}
	w.WriteByte(')')
	return
}
