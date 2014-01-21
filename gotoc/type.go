package gotoc

import (
	"bytes"
	"code.google.com/p/go.tools/go/types"
	"fmt"
	"reflect"
	"strconv"
)

func (cc *GTC) Type(w *bytes.Buffer, typ types.Type) {
	switch t := typ.(type) {
	case *types.Basic:
		if t.Kind() == types.UnsafePointer {
			w.WriteString("unsafe_Pointer")
		} else {
			types.WriteType(w, nil, t)
		}

	case *types.Named:
		o := t.Obj()
		if p := o.Pkg(); p != nil {
			w.WriteString(upath(p.Path()))
			w.WriteByte('_')
		}
		w.WriteString(o.Name())

	case *types.Pointer:
		cc.Type(w, t.Elem())
		w.WriteByte('*')

	case *types.Struct:
		w.WriteString("struct {\n")
		cc.il++
		for i, n := 0, t.NumFields(); i < n; i++ {
			f := t.Field(i)
			cc.indent(w)
			if tag := t.Tag(i); tag != "" {
				w.WriteString(reflect.StructTag(tag).Get("C"))
				w.WriteByte(' ')
			}
			cc.Type(w, f.Type())
			if !f.Anonymous() {
				w.WriteByte(' ')
				if f.Name() != "_" {
					w.WriteString(f.Name())
				} else {
					w.WriteString("__reserved")
					w.WriteString(strconv.Itoa(i))
				}
			}
			w.WriteString(";\n")
		}
		cc.il--
		w.WriteByte('}')

	default:
		fmt.Fprintf(w, "<%T>", t)
	}
}

func (cc *GTC) TypeStr(typ types.Type) string {
	buf := new(bytes.Buffer)
	cc.Type(buf, typ)
	return buf.String()
}

func (cc *GTC) Tuple(w *bytes.Buffer, t *types.Tuple, sep string) {
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
