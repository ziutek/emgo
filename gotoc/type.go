package gotoc

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/tools/go/types"
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

func (cdd *CDD) Type(w *bytes.Buffer, typ types.Type) (dim []string) {
	direct := true

writeType:
	switch t := typ.(type) {
	case *types.Basic:
		if t.Kind() == types.UnsafePointer {
			w.WriteString("unsafe$Pointer")
		} else {
			types.WriteType(w, nil, t)
		}

	case *types.Named:
		if _, ok := t.Underlying().(*types.Interface); ok {
			w.WriteString("interface")
		} else {
			cdd.Name(w, t.Obj(), direct)
		}

	case *types.Pointer:
		typ = t.Elem()
		direct = false
		dim = append(dim, "*")
		goto writeType

	case *types.Struct:
		if t.NumFields() == 0 {
			w.WriteString("empty")
			break
		}
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
			w.WriteByte(' ')
			name := dimFuncPtr(f.Name(), d)
			if name == "_" {
				name += strconv.Itoa(i) + "$"
			}
			w.WriteString(name + ";\n")
		}

		cdd.il--
		cdd.indent(w)
		w.WriteByte('}')

	case *types.Array:
		w.WriteString(cdd.arrayName(t))

	case *types.Slice:
		w.WriteString("slice")

	case *types.Map:
		w.WriteString("map")

	case *types.Chan:
		w.WriteString("chan")

	case *types.Signature:
		res, params := cdd.signature(t, true, noNames)
		w.WriteString(res.typ)
		dim = append(dim, params.String())
		dim = append(dim, res.dim...)

	case *types.Interface:
		w.WriteString("interface")

	default:
		fmt.Fprintf(w, "<%T>", t)
	}
	return
}

func (cdd *CDD) iface(w *bytes.Buffer, it *types.Interface) {
	w.WriteString("struct {\n")
	cdd.il++
	cdd.indent(w)
	w.WriteString("ithead h$;\n")
	for i := 0; i < it.NumMethods(); i++ {
		cdd.indent(w)
		f := it.Method(i)
		d := cdd.Type(w, f.Type())
		w.WriteString(" " + dimFuncPtr(f.Name(), d) + ";\n")
	}
	cdd.il--
	cdd.indent(w)
	w.WriteByte('}')
}

func (cdd *CDD) TypeStr(typ types.Type) (string, []string) {
	buf := new(bytes.Buffer)
	dim := cdd.Type(buf, typ)
	return buf.String(), dim
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
}

func (cdd *CDD) results(tup *types.Tuple) (res results) {
	if tup == nil {
		res.typ = "void"
		return
	}

	n := tup.Len()

	res.names = make([]string, n)
	for i := 0; i < n; i++ {
		name := tup.At(i).Name()
		switch name {
		case "":
			name = "_" + strconv.Itoa(n)

		case "_":
			res.hasNames = true

		default:
			name += "$"
			res.hasNames = true
		}
		res.names[i] = name
	}
	if n > 1 {
		res.typ, res.fields = cdd.tupleName(tup)
		return
	}
	v := tup.At(0)
	field0 := types.NewField(v.Pos(), v.Pkg(), "_0", v.Type(), false)
	res.fields = []*types.Var{field0}
	res.typ, res.dim = cdd.TypeStr(v.Type())
	return
}

const (
	noNames = iota
	numNames
	orgNames
	orgNamesI
)

type param struct {
	typ  string
	name string
}

type params []param

func (prs params) String() string {
	s := make([]string, len(prs))
	for i, p := range prs {
		s[i] = fmt.Sprintf(p.typ, p.name)
	}
	return "(" + strings.Join(s, ", ") + ")"
}

func (cdd *CDD) signature(sig *types.Signature, recv bool, pnames int) (res results, prms params) {
	res = cdd.results(sig.Results())
	if r := sig.Recv(); r != nil && recv {
		var (
			typ string
			dim []string
		)
		if _, ok := r.Type().Underlying().(*types.Interface); ok || pnames == orgNamesI {
			typ = "ival*"
		} else {
			typ, dim = cdd.TypeStr(r.Type())
		}
		var pname string
		switch pnames {
		case numNames:
			pname = "_0"
		case orgNames, orgNamesI:
			pname = cdd.NameStr(r, true)
		}
		if pname == "" {
			prms = append(
				prms,
				param{typ: typ + dimFuncPtr("%s", dim)},
			)
		} else {
			prms = append(
				prms,
				param{typ: typ + " " + dimFuncPtr("%s", dim), name: pname},
			)
		}
	}
	if p := sig.Params(); p != nil {
		for i, n := 0, p.Len(); i < n; i++ {
			v := p.At(i)
			typ, dim := cdd.TypeStr(v.Type())
			var pname string
			switch pnames {
			case numNames:
				pname = "_" + strconv.Itoa(i+1)
			case orgNames, orgNamesI:
				pname = cdd.NameStr(v, true)
				if pname == "_$" {
					pname = "unused" + cdd.gtc.uniqueId()
				}
			}
			if pname == "" {
				prms = append(
					prms,
					param{typ: typ + dimFuncPtr("%s", dim)},
				)
			} else {
				prms = append(
					prms,
					param{typ: typ + " " + dimFuncPtr("%s", dim), name: pname},
				)
			}
		}
	}
	return
}

func escape(s string) (ret string) {
	for {
		i := strings.IndexAny(s, "*()")
		if i == -1 {
			break
		}
		ret += s[:i]
		switch s[i] {
		case '*':
			ret += "$8$"
		case '(':
			ret += "$9$"
		default:
			ret += "$0$"
		}
		s = s[i+1:]
	}
	ret += s
	return
}

func (cdd *CDD) tupleName(tup *types.Tuple) (tupName string, fields []*types.Var) {
	n := tup.Len()
	for i := 0; i < n; i++ {
		if i != 0 {
			tupName += "$$"
		}
		name, dim := cdd.TypeStr(tup.At(i).Type())
		tupName += dimFuncPtr(name, dim)
	}
	tupName = escape(tupName)

	fields = make([]*types.Var, n)
	for i := 0; i < n; i++ {
		v := tup.At(i)
		fields[i] = types.NewField(
			v.Pos(), v.Pkg(), "_"+strconv.Itoa(i), v.Type(), false,
		)
	}

	if o, ok := cdd.gtc.tuples[tupName]; ok {
		cdd.DeclUses[o] = true
		return
	}

	s := types.NewStruct(fields, nil)
	o := types.NewTypeName(tup.At(0).Pos(), cdd.gtc.pkg, tupName, s)
	cdd.gtc.tuples[tupName] = o
	acd := cdd.gtc.newCDD(o, TypeDecl, 0)
	acd.structDecl(new(bytes.Buffer), tupName, s)
	cdd.acds = append(cdd.acds, acd)

	cdd.DeclUses[o] = true

	return
}

func (cdd *CDD) arrayName(a *types.Array) string {
	l := strconv.FormatInt(a.Len(), 10)
	name, dim := cdd.TypeStr(a.Elem())
	name = dimFuncPtr(name, dim)
	name = "$" + l + "_$" + escape(name)

	if o, ok := cdd.gtc.arrays[name]; ok {
		cdd.addObject(o, true)
		return name
	}
	f := types.NewField(0, cdd.gtc.pkg, "arr["+l+"]", a.Elem(), false)
	s := types.NewStruct([]*types.Var{f}, nil)
	o := types.NewTypeName(0, cdd.gtc.pkg, name, s)
	cdd.gtc.arrays[name] = o
	cdd.addObject(o, true)
	acd := cdd.gtc.newCDD(o, TypeDecl, 0)
	cdd.acds = append(cdd.acds, acd)
	acd.structDecl(new(bytes.Buffer), name, s)
	return name
}
