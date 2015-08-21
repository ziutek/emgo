package gotoc_test

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	//_ "golang.org/x/tools/go/gcimporter"

	"github.com/ziutek/emgo/gotoc"
)

type ddi struct {
	decl, def, init string
}

type sampleDecl struct {
	filePos string
	goDecl  string
	c       []*ddi
}

func (s sampleDecl) testDecl() error {
	src := "package foo\n" + s.goDecl

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, s.filePos, src, 0)
	if err != nil {
		return err
	}

	ti := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}
	pkg, err := new(types.Config).Check("foo", fset, []*ast.File{f}, ti)
	if err != nil {
		return err
	}

	gtc := gotoc.NewGTC(fset, pkg, ti, &gotoc.StdSizes{4, 8})
	var cdds []*gotoc.CDD
	for _, d := range f.Decls {
		for _, cdd := range gtc.Decl(d, 0) {
			cdds = append(cdds, cdd.AllCDDS()...)
		}
	}

	if len(cdds) < len(s.c) {
		return s.noCDDError(len(cdds))
	}
	if len(cdds) > len(s.c) {
		return s.noExpError(cdds[len(s.c)])
	}
	for i, cdd := range cdds {
		cddDecl := string(cdd.Decl)
		if cddDecl != s.c[i].decl {
			return s.notMatch("decl", cddDecl, s.c[i].decl)
		}
		cddDef := string(cdd.Def)
		if cddDef != s.c[i].def {
			return s.notMatch("def", cddDef, s.c[i].def)
		}
		cddInit := string(cdd.Init)
		if cddInit != s.c[i].init {
			return s.notMatch("init", cddInit, s.c[i].init)
		}
	}
	return nil
}

func (s sampleDecl) noCDDError(n int) error {
	buf := new(bytes.Buffer)
	buf.WriteString(s.filePos + ": there is no generated code for:\n")
	c := s.c[n]
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
	buf.WriteString("// end\n")
	return errors.New(buf.String())
}

func (s sampleDecl) noExpError(cdd *gotoc.CDD) error {
	buf := new(bytes.Buffer)
	buf.WriteString(s.filePos + ": too much generated code:\n")
	buf.WriteString("// Go code:\n")
	buf.WriteString(s.goDecl)
	buf.WriteString("// Too much generated C code:\n")
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
	buf.WriteString("// end\n")
	return errors.New(buf.String())
}

func (s sampleDecl) notMatch(section, cdd, c string) error {
	buf := new(bytes.Buffer)
	buf.WriteString(s.filePos + ": code not match:\n")
	buf.WriteString("// Go code:\n")
	buf.WriteString(s.goDecl)
	buf.WriteString("// Generated " + section + ":\n")
	buf.WriteString(cdd)
	buf.WriteString("// Expected " + section + ":\n")
	buf.WriteString(c)
	buf.WriteString("// end")
	return errors.New(buf.String())
}

func TestDeclFiles(t *testing.T) {
	dname := "tests"
	dir, err := os.Open(dname)
	if err != nil {
		t.Fatal(err)
	}
	fnames, err := dir.Readdirnames(-1)
	if err != nil {
		t.Fatal(err)
	}
	for _, fname := range fnames {
		if !strings.HasSuffix(fname, ".test") {
			continue
		}
		if err := testDeclFile(filepath.Join(dname, fname)); err != nil {
			t.Error(err)
		}
	}
}

func testDeclFile(fname string) error {
	f, err := os.Open(fname)
	if err != nil {
		return err
	}
	defer f.Close()
	s := bufio.NewScanner(f)
	lnum := 0
	for s.Scan() {
		lnum++
		if s.Text() != "// Go code:" {
			continue
		}
		goDecl := ""
		lineN := lnum
		for s.Scan() {
			lnum++
			line := s.Text()
			if line == "// C code:" {
				break
			}
			goDecl += line + "\n"
		}
		var (
			c []*ddi
			d *ddi
		)
		lastid := 3
		for s.Scan() {
			lnum++
			line := s.Text()
			if line == "// end" {
				sd := sampleDecl{
					filePos: fname + ":" + strconv.Itoa(lineN),
					goDecl:  goDecl,
					c:       c,
				}
				if err := sd.testDecl(); err != nil {
					return err
				}
				break
			}
			var id int
			switch line {
			case "// decl":
				id = 0

			case "// def":
				id = 1

			case "// init":
				id = 2

			default:
				switch lastid {
				case 0:
					d.decl += line + "\n"

				case 1:
					d.def += line + "\n"

				case 2:
					d.init += line + "\n"

				default:
					return fmt.Errorf("%s:%d test syntax error", fname, lnum)
				}
				continue
			}

			if lastid >= id {
				d = new(ddi)
				c = append(c, d)
			}
			lastid = id
		}
	}
	if err := s.Err(); err != nil {
		return err
	}
	return nil
}
