package gotoc_test

import (
	"bytes"
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"testing"

	_ "code.google.com/p/go.tools/go/gcimporter"
	"code.google.com/p/go.tools/go/types"

	"github.com/ziutek/emgo/gotoc"
)

type ddi struct {
	decl, def, init string
}

type sampleDecl struct {
	goDecl string
	c      []ddi
}

func (s sampleDecl) testDecl(fname string) error {
	src := "package foo\n" + s.goDecl

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fname, src, 0)
	if err != nil {
		return err
	}

	ti := &types.Info{
		Types:   make(map[ast.Expr]types.TypeAndValue),
		Objects: make(map[*ast.Ident]types.Object),
	}
	pkg, err := new(types.Config).Check("foo", fset, []*ast.File{f}, ti)
	if err != nil {
		return err
	}

	gtc := gotoc.NewGTC(pkg, ti)
	cdds := gtc.Decl(f.Decls[0], 0)

	if len(cdds) != len(s.c) {
		return s.notMatch(fname, cdds)
	}
	for i, cdd := range cdds {
		if string(cdd.Decl) != s.c[i].decl ||
			string(cdd.Def) != s.c[i].def ||
			string(cdd.Init) != s.c[i].init {

			return s.notMatch(fname, cdds)
		}
	}
	return nil
}

func (s sampleDecl) notMatch(fname string, cdds []*gotoc.CDD) error {
	buf := new(bytes.Buffer)
	buf.WriteString("Go to C translation error\n")
	buf.WriteString("// Go code (" + fname + "):\n")
	buf.WriteString(s.goDecl)
	buf.WriteString("// Generated C code:\n")
	for _, cdd := range cdds {
		if len(cdd.Decl) > 0 {
			buf.WriteString("// decl\n")
			buf.Write(cdd.Decl)
		}
		if len(cdd.Def) > 0 {
			buf.WriteString("// def\n")
			buf.Write(cdd.Def)
		}
		if len(cdd.Init) > 0 {
			buf.WriteString("// init\n")
			buf.Write(cdd.Init)
		}
	}
	buf.WriteString("// Expected C code:\n")
	for _, c := range s.c {
		if len(c.decl) > 0 {
			buf.WriteString("// decl\n")
			buf.WriteString(c.decl)
		}
		if len(c.def) > 0 {
			buf.WriteString("// def:\n")
			buf.WriteString(c.def)
		}
		if len(c.init) > 0 {
			buf.WriteString("// init\n")
			buf.WriteString(c.init)
		}
	}
	buf.WriteString("// end")
	return errors.New(buf.String())
}

var tabDecl = []sampleDecl{
	{
		"type S struct {a, b int}\n",

		[]ddi{{
			decl: "struct foo_S_struct;\n" +
				"typedef struct foo_S_struct foo_S;\n",

			def: "struct foo_S_struct {\n" +
				"	int a;\n" +
				"	int b;\n" +
				"};\n",
		}},
	}, {
		"var A = 3\n",

		[]ddi{{
			decl: "int foo_A;\n",
			def:  "int foo_A = 3;\n",
		}},
	}, {
		"var A = []int{1, 2, 3}\n",

		[]ddi{{
			decl: "__slice foo_A;\n",
			def:  "__slice foo_A;\n",
			init: "	foo_A = (__slice){(int[]){1, 2, 3}, 3, 3};\n",
		}},
	}, {
		"var A = [][2]int{{1, 2}, {3, 4}}\n",

		[]ddi{{
			decl: "__slice foo_A;\n",
			def:  "__slice foo_A;\n",
			init: "	foo_A = (__slice){(int[][2]){{1, 2}, {3, 4}}, 2, 2};\n",
		}},
	},
}

func TestDecl(t *testing.T) {
	for i, s := range tabDecl {
		if err := s.testDecl(strconv.Itoa(i) + ".go"); err != nil {
			t.Error(err)
		}
	}
}

type simpleDecl struct {
	g, c string
}

var tabSimpleDecl = []simpleDecl{
	{"type P *int", "typedef int *foo_P;"},
	{"type A [4]int", "typedef int foo_A[4];"},
	{"type AP [4]*int", "typedef int *foo_AP[4];"},
	{"type PA *[4]int", "typedef int (*foo_PA)[4];"},
	{"type PAP *[4]*int", "typedef int *(*foo_PAP)[4];"},
	{"type AA [4][3]int", "typedef int foo_AA[4][3];"},
	{"type PAA *[4][3]int", "typedef int (*foo_PAA)[4][3];"},
	{"type PAPA *[4]*[3]int", "typedef int (*(*foo_PAPA)[4])[3];"},
	{"type PAPAP *[4]*[3]*int", "typedef int *(*(*foo_PAPAP)[4])[3];"},

	{
		"type F func(a, b int, c byte) byte",
		"typedef byte (*foo_F)(int a, int b, byte c);",
	},

	{"func F(a int)", "void foo_F(int a);"},
	{"func F(a [4]int) uint", "uint foo_F(int a[4]);"},
}

func TestSimpleDecl(t *testing.T) {
	for i, s := range tabSimpleDecl {
		sd := sampleDecl{goDecl: s.g + "\n", c: []ddi{{decl: s.c + "\n"}}}
		if err := sd.testDecl(strconv.Itoa(i) + ".go"); err != nil {
			t.Error(err)
		}
	}
}
