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
			w.WriteString("structE")
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
			if f.Name() == "_" {
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

func (cdd *CDD) tiname(typ types.Type) string {
	switch t := typ.(type) {
	case *types.Pointer:
		return escape("*") + cdd.tiname(t.Elem())
	case *types.Signature:
		return "func$$$" + cdd.signame(t)
	case *types.Slice:
		return "slice$$" + cdd.tiname(t.Elem())
	case *types.Chan:
		return "chan$$" + cdd.tiname(t.Elem())
	case *types.Map:
		return "map$$" + cdd.tiname(t.Key()) + cdd.tiname(t.Elem())
	case *types.Named:
		obj := t.Obj()
		name := cdd.NameStr(obj, false)
		if cdd.gtc.isLocal(obj) {
			name = upath(obj.Pkg().Path()) + "$" + name + "$" +
				strconv.Itoa(int(obj.Pos()))
		} else {
			name += "$$"
		}
		return name
	default:
		name, _ := cdd.TypeStr(typ)
		switch name {
		case "byte":
			name = "uint8"
		case "rune":
			name = "int32"
		}
		return name + "$$"
	}
}

var basicKinds = []string{
	"Invalid",
	"Bool",
	"Int",
	"Int8",
	"Int16",
	"Int32",
	"Int64",
	"Uint",
	"Uint8",
	"Uint16",
	"Uint32",
	"Uint64",
	"Uintptr",
	"Float32",
	"Float64",
	"Complex64",
	"Complex128",
	"String",
	"UnsafePointer",
}

func (cdd *CDD) tinfo(w *bytes.Buffer, typ types.Type) string {
	tname := cdd.tiname(typ)
	basic := false

	switch t := typ.(type) {
	case *types.Basic:
		basic = true
	case *types.Pointer, *types.Slice:
		switch p := t.(interface {
			Elem() types.Type
		}).Elem().(type) {
		case *types.Basic:
			basic = true
		case *types.Named:
			basic = (p.Obj().Pkg() == nil)
		}
	case *types.Named:
		basic = (t.Obj().Pkg() == nil)
	}
	if basic && cdd.gtc.pkg.Path() != "builtin" {
		// Generate tinfo for "basic" types only in builtin package.
		return tname
	}
	if o, ok := cdd.gtc.tinfos[tname]; ok {
		cdd.addObject(o, true)
		return tname
	}
	v := types.NewVar(0, cdd.gtc.pkg, tname, typ)
	cdd.gtc.tinfos[tname] = v
	cdd.addObject(v, true)
	acd := cdd.gtc.newCDD(v, VarDecl, 0)
	cdd.acds = append(cdd.acds, acd)
	cdd = nil

	nt, named := typ.(*types.Named)
	if named {
		acd.addObject(nt.Obj(), true)
	}
	w.WriteString("const\ntinfo ")
	w.WriteString(tname)
	acd.copyDecl(w, ";\n")
	w.WriteString(" = {\n")
	acd.il++
	acd.indent(w)
	w.WriteString("{\n")
	acd.il++
	if named {
		acd.indent(w)
		w.WriteString(".name = EGSTR(\"" + typ.String() + "\"),\n")
	}
	acd.indent(w)
	w.WriteString(".size = " + strconv.FormatInt(acd.gtc.siz.Sizeof(typ), 10) + ",\n")
	var (
		kind  string
		elems []types.Type
	)
	switch t := typ.Underlying().(type) {
	case *types.Basic:
		k := t.Kind()
		if k < types.Invalid || k >= types.UntypedBool {
			panic(k)
		}
		kind = basicKinds[k]
	case *types.Array:
		kind = "Array"
	case *types.Chan:
		kind = "Chan"
	case *types.Signature:
		kind = "Func"
	case *types.Interface:
		kind = "Interface"
	case *types.Map:
		kind = "Map"
	case *types.Pointer:
		kind = "Ptr"
		elems = []types.Type{t.Elem()}
	case *types.Slice:
		kind = "Slice"
		elems = []types.Type{t.Elem()}
	case *types.Struct:
		kind = "Struct"
		elems = make([]types.Type, t.NumFields())
		for i := range elems {
			elems[i] = t.Field(i).Type()
		}
	default:
		panic(t)
	}
	acd.indent(w)
	w.WriteString(".kind = " + kind)
	if len(elems) > 0 {
		w.WriteString(",\n")
		acd.indent(w)
		w.WriteString(".elems = CSLICE(")
		w.WriteString(strconv.Itoa(len(elems)))
		w.WriteString(", ((const tinfo*[]){\n")
		acd.il++
		for i, e := range elems {
			if i != 0 {
				w.WriteString(",\n")
			}
			acd.indent(w)
			w.WriteByte('&')
			w.WriteString(acd.tinameDU(e))
		}
		w.WriteByte('\n')
		acd.il--
		acd.indent(w)
		w.WriteString("}))")
	}
	// BUG: Following code doesn't work in case of structs that contains embeded
	// types with methods. Use types.MethodSet to fix it.
	var (
		suff    string
		methods []*types.Func
		it      *types.Interface
	)
	if t, ok := typ.(*types.Pointer); ok {
		nt, named = t.Elem().(*types.Named)
		suff = "$0"
	} else {
		suff = "$1"
	}
	acd.Export = basic
	acd.Weak = !basic && !named
	if named {
		if it, _ = nt.Underlying().(*types.Interface); it == nil {
			methods = make([]*types.Func, 0, nt.NumMethods())
			for i := 0; i < cap(methods); i++ {
				m := nt.Method(i)
				if types.Identical(m.Type().(*types.Signature).Recv().Type(), typ) {
					methods = append(methods, m)
				}
			}
		} else {
			methods = make([]*types.Func, it.NumMethods())
			for i := range methods {
				methods[i] = it.Method(i)
			}
		}
	}
	if len(methods) > 0 {
		w.WriteString(",\n")
		acd.indent(w)
		w.WriteString(".methods = CSLICE(")
		w.WriteString(strconv.Itoa(len(methods)))
		w.WriteString(", ((const minfo*[]){\n")
		acd.il++
		for i, m := range methods {
			if i != 0 {
				w.WriteString(",\n")
			}
			acd.indent(w)
			w.WriteByte('&')
			w.WriteString(acd.minfo(m))
		}
		w.WriteByte('\n')
		acd.il--
		acd.indent(w)
		w.WriteString("}))")
		if it == nil {
			w.WriteByte('\n')
			acd.il--
			acd.indent(w)
			w.WriteString("}, {\n")
			acd.il++
			for i, m := range methods {
				if i != 0 {
					w.WriteString(",\n")
				}
				acd.indent(w)
				acd.Name(w, m, true)
				w.WriteString(suff)
			}
		}
	}
	w.WriteByte('\n')
	acd.il--
	acd.indent(w)
	w.WriteString("}\n")
	acd.il--
	acd.indent(w)
	w.WriteString("};\n")
	acd.copyDef(w)
	return tname
}

func (cdd *CDD) tinameDU(typ types.Type) string {
	t := typ
	if p, ok := typ.(*types.Pointer); ok {
		t = p.Elem()
	}
	if n, ok := t.(*types.Named); !ok || n.Obj().Pkg() == nil {
		return cdd.tinfo(new(bytes.Buffer), typ)
	}
	return cdd.tiname(typ)
}

func (cdd *CDD) prname(tup *types.Tuple) string {
	var prname string
	for i := 0; i < tup.Len(); i++ {
		prname += cdd.tiname(tup.At(i).Type())
	}
	return prname
}

func (cdd *CDD) minfo(f *types.Func) string {
	var mname string
	sig := f.Type().(*types.Signature)
	mname = f.Name() + "$$$" + cdd.signame(sig)
	if o, ok := cdd.gtc.minfos[mname]; ok {
		cdd.addObject(o, true)
		return mname
	}
	v := types.NewVar(0, cdd.gtc.pkg, mname, sig)
	cdd.gtc.minfos[mname] = v
	cdd.addObject(v, true)
	acd := cdd.gtc.newCDD(v, VarDecl, 0)
	cdd.acds = append(cdd.acds, acd)
	cdd = nil
	acd.Weak = true
	w := new(bytes.Buffer)
	w.WriteString("__attribute__((section(\".unused\"))) const\n")
	w.WriteString("minfo " + mname)
	acd.copyDecl(w, ";\n")
	acd.Def = acd.Decl
	return mname
}

func (cdd *CDD) signame(sig *types.Signature) string {
	signame := cdd.prname(sig.Params())
	if sig.Variadic() {
		signame += "variadic$$"
	}
	return signame + "$" + cdd.prname(sig.Results())
}
