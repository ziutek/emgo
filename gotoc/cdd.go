package gotoc

import (
	"bytes"
	"go/ast"
	"go/token"
	"io"

	"golang.org/x/tools/go/exact"
	"golang.org/x/tools/go/types"
)

type DeclType int

const (
	FuncDecl DeclType = iota
	VarDecl
	ConstDecl
	TypeDecl
	ImportDecl
)

// CDD stores Go declaration translated to C declaration and definition.
type CDD struct {
	Origin       types.Object // object for this declaration
	DeclUses     map[types.Object]bool
	FuncBodyUses map[types.Object]bool
	Complexity   int

	Typ    DeclType
	Export bool
	Weak   bool
	Inline bool // set by DetermineInline()

	Decl     []byte
	Def      []byte
	Init     []byte
	InitNext *CDD

	gtc *GTC
	il  int // indentation level

	init  bool // true if generated code will be placed in init() function
	fbody bool // true if translation process in function body
	dfsm  int8

	acds []*CDD // additional CDDs
}

func (gtc *GTC) newCDD(o types.Object, t DeclType, il int) *CDD {
	cdd := &CDD{
		Origin:       o,
		Typ:          t,
		DeclUses:     make(map[types.Object]bool),
		FuncBodyUses: make(map[types.Object]bool),
		gtc:          gtc,
		il:           il,
	}
	return cdd
}

func (cdd *CDD) indent(w *bytes.Buffer) {
	for i := 0; i < cdd.il; i++ {
		w.WriteByte('\t')
	}
}

func (cdd *CDD) copyDecl(b *bytes.Buffer, suffix string) {
	n := b.Len()
	b.WriteString(suffix)
	cdd.Decl = append([]byte(nil), b.Bytes()...)
	b.Truncate(n)
}

func (cdd *CDD) copyDef(b *bytes.Buffer) {
	cdd.Def = append([]byte(nil), b.Bytes()...)
}

func (cdd *CDD) prependDef(b *bytes.Buffer) {
	newDef := make([]byte, b.Len()+len(cdd.Def))
	copy(newDef, b.Bytes())
	copy(newDef[b.Len():], cdd.Def)
	cdd.Def = newDef
}

func (cdd *CDD) copyInit(b *bytes.Buffer) {
	cdd.Init = append([]byte(nil), b.Bytes()...)
}

func (cdd *CDD) WriteDecl(wh, wc io.Writer) error {
	if len(cdd.Decl) == 0 {
		return nil
	}

	var prefix string

	switch cdd.Typ {
	case FuncDecl:
		if cdd.Inline {
			prefix = "static inline\n"
		} else if !cdd.Export {
			prefix = "static\n"
		}

	case VarDecl:
		if cdd.Export {
			prefix += "extern "
		} else {
			return nil
		}

	case ConstDecl:
		if !cdd.Export {
			return nil
		}
	}

	w := wc
	if cdd.Export {
		w = wh
	}

	_, err := io.WriteString(w, prefix)
	if err != nil {
		return err
	}
	_, err = w.Write(cdd.Decl)
	return err
}

func (cdd *CDD) WriteDef(wh, wc io.Writer) error {
	if len(cdd.Def) == 0 {
		return nil
	}

	prefix := ""
	w := wc

	switch cdd.Typ {
	case FuncDecl:
		if cdd.Export {
			if cdd.Inline {
				prefix = "static inline\n"
				w = wh
			}
		} else {
			prefix = "static\n"
		}

	case VarDecl:
		if cdd.Weak {
			prefix = "__attribute__((weak)) "
		} else if !cdd.Export {
			prefix = "static "
		}

	case ConstDecl:
		return nil

	case TypeDecl:
		if cdd.Export {
			w = wh
		}
	}

	_, err := io.WriteString(w, prefix)
	if err != nil {
		return err
	}
	_, err = w.Write(cdd.Def)
	return err
}

func (cdd *CDD) DetermineInline() {
	if len(cdd.Def) == 0 || cdd.Complexity < 0 {
		// Declaration only or function marked not for inlining.
		return
	}
	// TODO: Use more information (from il, BodyUses).
	// TODO: Complexity can be better calculated.
	if cdd.Complexity <= cdd.gtc.inlineThres {
		cdd.Inline = true
	}
}

func (cdd *CDD) addObject(o types.Object, direct bool) {
	if o == cdd.Origin || o == nil {
		return
	}
	if o.Pkg() == nil {
		// Don't save references for builtin objects (eg: error type)
		return
	}
	if cdd.init && !cdd.gtc.isImported(o) {
		// Don't save references to package objects if used in init() function.
		// This is mainly for global variables initialization in init().
		return
	}
	if cdd.fbody {
		cdd.FuncBodyUses[o] = direct
	} else {
		cdd.DeclUses[o] = direct
	}
}

func (cdd *CDD) dfs(all map[types.Object]*CDD, out []*CDD) []*CDD {
	if cdd.dfsm > 0 {
		panic("direct cycle in type declaration")
	}
	if cdd.dfsm < 0 {
		return out
	}
	cdd.dfsm = 1
	for o, direct := range cdd.DeclUses {
		if !direct {
			continue
		}
		u, ok := all[o]
		if !ok {
			continue
		}
		out = u.dfs(all, out)
	}
	cdd.dfsm = -1
	return append(out, cdd)
}

func dfs(all map[types.Object]*CDD) []*CDD {
	out := make([]*CDD, 0, len(all))
	for _, cdd := range all {
		out = cdd.dfs(all, out)
	}
	return out
}

func (cdd *CDD) object(ident *ast.Ident) types.Object {
	return cdd.gtc.object(ident)
}

func (cdd *CDD) exprType(e ast.Expr) types.Type {
	return cdd.gtc.exprType(e)
}

func (cdd *CDD) exprValue(e ast.Expr) exact.Value {
	return cdd.gtc.exprValue(e)
}

func (cdd *CDD) exit(pos token.Pos, f string, a ...interface{}) {
	cdd.gtc.exit(pos, f, a...)
}

func (cdd *CDD) notImplemented(n ast.Node, tl ...types.Type) {
	cdd.gtc.notImplemented(n, tl...)
}
func (cdd *CDD) AllCDDS() (cdds []*CDD) {
	for _, a := range cdd.acds {
		cdds = append(cdds, a.AllCDDS()...)
	}
	cdds = append(cdds, cdd)
	return
}
