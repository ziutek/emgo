package gotoc

import (
	"bytes"
	"fmt"
	"go/types"
	"reflect"
	"strconv"
	"strings"
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
		switch t.Kind() {
		case types.UnsafePointer:
			w.WriteString("unsafe$Pointer")
		case types.UntypedString:
			w.WriteString("string")
		case types.Int:
			w.WriteString("int_")
		default:
			types.WriteType(w, t, nil)
		}

	case *types.Named:
		if _, ok := t.Underlying().(*types.Interface); ok {
			w.WriteString("interface")
			cdd.addObject(t.Obj(), direct)
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
				//fmt.Println(cdd.gtc.fset.Position(f.Pos()), tag)
				w.WriteString(reflect.StructTag(tag).Get("c"))
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

	case *types.Tuple:
		tn, _ := cdd.tupleName(t)
		w.WriteString(tn)

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
			if pname = cdd.NameStr(r, true); pname == "" {
				pname = "_0"
			}
		}
		if pname == "" && pnames != orgNamesI {
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
				switch pname {
				case "":
					pname = "_" + strconv.Itoa(i+1)
				case "_$":
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
		i := strings.IndexAny(s, " ,*(){}\n; \t")
		if i == -1 {
			break
		}
		ret += s[:i]
		switch s[i] {
		case '*':
			ret += "$8$"
		case '(':
			ret += "$9$"
		case ')':
			ret += "$0$"
		case ',':
			ret += "$1$"
		case '{':
			ret += "$2$"
		case '}':
			ret += "$3$"
		case '\n', ';':
			ret += "$4$"
		case ' ', '\t':
			// Nothing.
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
		//cdd.DeclUses[o] = true
		cdd.addObject(o, true)
		return
	}

	s := types.NewStruct(fields, nil)
	o := types.NewTypeName(tup.At(0).Pos(), cdd.gtc.pkg, tupName, s)
	cdd.gtc.tuples[tupName] = o
	acd := cdd.gtc.newCDD(o, TypeDecl, 0)
	acd.structDecl(new(bytes.Buffer), tupName, s, "")
	cdd.acds = append(cdd.acds, acd)

	//cdd.DeclUses[o] = true
	cdd.addObject(o, true)

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
	acd.structDecl(new(bytes.Buffer), name, s, "")
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
	case *types.Struct:
		s, _ := cdd.TypeStr(typ)
		return "anon$$" + escape(s)
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
	var builtin, named bool
	switch t := typ.(type) {
	case *types.Basic:
		builtin = true
	case *types.Pointer:
		switch e := t.Elem().(type) {
		case *types.Basic:
			builtin = true
		case *types.Named:
			named = true
			builtin = (e.Obj().Pkg() == nil)
		}
	case *types.Slice:
		switch e := t.Elem().(type) {
		case *types.Basic:
			builtin = true
		case *types.Named:
			builtin = (e.Obj().Pkg() == nil)
		}
	case *types.Named:
		named = true
		builtin = (t.Obj().Pkg() == nil)
	}
	if builtin && cdd.gtc.pkg.Path() != "internal" {
		// Generate tinfo for builtin types only in internal package.
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

	acd.Export = builtin
	acd.Weak = !builtin && !named
	w.WriteString("const tinfo ")
	w.WriteString(tname)
	acd.copyDecl(w, ";\n")
	w.WriteString(" = {\n")
	acd.il++
	acd.indent(w)
	w.WriteString("{\n")
	acd.il++
	if nt, ok := typ.(*types.Named); ok {
		acd.addObject(nt.Obj(), true)
		acd.indent(w)
		w.WriteString(".name = EGSTR(\"" + nt.String() + "\"),\n")
	}
	//acd.indent(w)
	//w.WriteString(".size = " + strconv.FormatInt(acd.gtc.siz.Sizeof(typ), 10) + ",\n")
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
		kind = "Array - " + strconv.FormatInt(t.Len(), 10)
		elems = []types.Type{t.Elem()}
	case *types.Chan:
		kind = "Chan"
		elems = []types.Type{t.Elem()}
	case *types.Signature:
		kind = "Func"
	case *types.Interface:
		kind = "Interface"
	case *types.Map:
		kind = "Map"
		elems = []types.Type{t.Key(), t.Elem()}
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
	mset := acd.gtc.methodSet(typ)
	_, isi := typ.Underlying().(*types.Interface)
	if mset.Len() > 0 && (isi || acd.gtc.siz.Sizeof(typ) <= acd.gtc.sizIval) {
		acd.Weak = false
		w.WriteString(",\n")
		acd.indent(w)
		w.WriteString(".methods = CSLICE(")
		w.WriteString(strconv.Itoa(mset.Len()))
		w.WriteString(", ((const minfo*[]){\n")
		acd.il++
		for i, c := 0, false; i < mset.Len(); i++ {
			f := mset.At(i).Obj().(*types.Func)
			if !f.Exported() {
				pragmas, _ := acd.gtc.pragmas(acd.gtc.defs[f])
				if !pragmas.Contains("minfo") {
					continue
				}
			}
			if c {
				w.WriteString(",\n")
			} else {
				c = true
			}
			acd.indent(w)
			w.WriteByte('&')
			w.WriteString(acd.minfo(f))
		}
		w.WriteByte('\n')
		acd.il--
		acd.indent(w)
		w.WriteString("}))")
		if !isi {
			w.WriteByte('\n')
			acd.il--
			acd.indent(w)
			w.WriteString("}, {\n")
			acd.il++
			for i, c := 0, false; i < mset.Len(); i++ {
				method := mset.At(i)
				f := method.Obj().(*types.Func)
				if !f.Exported() {
					pragmas, _ := acd.gtc.pragmas(acd.gtc.defs[f])
					if !pragmas.Contains("minfo") {
						continue
					}
				}
				if c {
					w.WriteString(",\n")
				} else {
					c = true
				}
				acd.indent(w)
				w.WriteString(acd.imethod(method))
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
	acd.NoAlloc = true
	w := new(bytes.Buffer)
	w.WriteString("const minfo " + mname)
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

func (cdd *CDD) imethod(sel *types.Selection) string {
	rcv := sel.Recv()
	fun := sel.Obj().(*types.Func)
	sig := fun.Type().(*types.Signature)
	var fname string
	if p, ok := rcv.(*types.Pointer); ok {
		fname, _ = cdd.TypeStr(p.Elem())
		fname += "$" + fun.Name() + "$0"
	} else {
		fname, _ = cdd.TypeStr(rcv)
		fname += "$" + fun.Name() + "$1"
	}
	fi := types.NewFunc(fun.Pos(), cdd.Origin.Pkg(), fname, sig)
	//cdd.addObject(fi, false) - commented to avoid export fi.
	acd := cdd.gtc.newCDD(fi, FuncDecl, 0)
	acd.Complexity = cdd.gtc.noinlineThres
	cdd.acds = append(cdd.acds, acd)

	cdd = nil

	w := new(bytes.Buffer)
	res, params := acd.signature(sig, true, orgNamesI)

	w.WriteString(res.typ)
	w.WriteByte(' ')
	w.WriteString(dimFuncPtr(fname+params.String(), res.dim))
	acd.copyDecl(w, ";\n")

	w.WriteString(" {\n")
	acd.il++

	var s string

	if _, ok := rcv.(*types.Pointer); ok {
		ts, dim := acd.TypeStr(rcv)
		s = "((" + ts + dimFuncPtr("", dim) + ")" + params[0].name + "->ptr)"
	} else {
		ts, dim := acd.TypeStr(types.NewPointer(rcv))
		s = "(*(" + ts + dimFuncPtr("", dim) + ")" + params[0].name + ")"
	}
	index := sel.Index()

	for _, id := range index[:len(index)-1] {
		if p, ok := rcv.(*types.Pointer); ok {
			rcv = p.Elem()
			s += "->"
		} else {
			s += "."
		}
		f := rcv.Underlying().(*types.Struct).Field(id)
		s += f.Name()
		rcv = f.Type()
	}

	if _, ok := rcv.Underlying().(*types.Interface); ok {
		acd.indent(w)
		acd.Type(w, rcv)
		w.WriteString(" _r = " + s + ";\n")
		acd.indent(w)
		w.WriteString("return ((")
		acd.Name(w, rcv.(*types.Named).Obj(), false)
		w.WriteString("*)_r.itab$)->")
		w.WriteString(fun.Name())
		w.WriteString("(&_r.val$")
	} else {
		acd.indent(w)
		w.WriteString("return ")
		acd.Name(w, fun, true)
		w.WriteByte('(')
		_, sigpr := sig.Recv().Type().(*types.Pointer)
		if _, ok := rcv.(*types.Pointer); ok {
			if !sigpr {
				w.WriteByte('*')
			}
		} else {
			if sigpr {
				w.WriteByte('&')
			}
		}
		w.WriteString(s)
	}
	for i := 1; i < len(params); i++ {
		w.WriteString(", " + params[i].name)
	}

	w.WriteString(");\n")

	acd.il--
	acd.indent(w)
	w.WriteString("}\n")
	acd.copyDef(w)

	return fname
}
