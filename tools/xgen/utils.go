package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"unicode"
)

func die(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func fdie(f, format string, args ...interface{}) {
	die(f+": "+format, args...)
}

func checkErr(err error) {
	if err != nil {
		die("Error: %v", err)
	}
}

func split(s string) (first, rest string) {
	s = strings.TrimSpace(s)
	i := strings.IndexFunc(s, unicode.IsSpace)
	if i < 0 {
		return s, ""
	}
	return s[:i], strings.TrimSpace(s[i+1:])
}

func save(fpath string, tpl *template.Template, ctx interface{}) {
	buf := new(bytes.Buffer)
	checkErr(tpl.Execute(buf, ctx))
	src, err := format.Source(buf.Bytes())
	checkErr(err)
	dir := filepath.Dir(fpath)
	base := filepath.Base(fpath)
	f, err := os.Create(filepath.Join(dir, "xgen_"+base))
	checkErr(err)
	defer f.Close()
	_, err = f.Write(src)
	checkErr(err)
}
