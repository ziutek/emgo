package gotoc

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"reflect"
)

func (cdd *CDD) Type(w *bytes.Buffer, typ types.Type) {
	switch t := typ.(type) {
	case *types.Basic:
		if t.Kind() == types.UnsafePointer {
			w.WriteString("unsafe_Pointer")
		} else {
			types.WriteType(w, nil, t)
		}

	case *types.Named:
		cdd.Name(w, t.Obj())

	case *types.Pointer:
		cdd.Type(w, t.Elem())
		w.WriteByte('*')

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
			cdd.Type(w, f.Type())
			if !f.Anonymous() {
				w.WriteByte(' ')
				cdd.Name(w, f)
			}
			w.WriteString(";\n")
		}
		cdd.il--
		w.WriteByte('}')

	default:
		fmt.Fprintf(w, "<%T>", t)
	}
}

func (cdd *CDD) TypeStr(typ types.Type) string {
	buf := new(bytes.Buffer)
	cdd.Type(buf, typ)
	return buf.String()
}

func (cdd *CDD) Tuple(w *bytes.Buffer, t *types.Tuple, sep string) {
	for i, n := 0, t.Len(); i < n; i++ {
		if i != 0 {
			w.WriteString(sep)
		}
		v := t.At(i)
		cdd.Type(w, v.Type())
		w.WriteByte(' ')
		cdd.Name(w, v)
	}
}
