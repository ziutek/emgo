package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func usage() {
	fmt.Fprintln(
		flag.CommandLine.Output(),
		"Usage: datago [OPTIONS] FILES...",
	)
	flag.PrintDefaults()
}

func main() {
	var (
		elsiz  int
		prefix string
		big    bool
		typ    string
	)
	flag.IntVar(&elsiz, "elsiz", 1, "element size: 1, 2, 4, 8 bytes")
	flag.StringVar(&prefix, "prefix", "", "prefix for variable name")
	flag.BoolVar(&big, "big", false, "big endian")
	flag.Usage = usage
	flag.Parse()
	if flag.NArg() == 0 {
		usage()
		os.Exit(1)
	}
	switch elsiz {
	case 1:
		typ = "byte"
	case 2:
		typ = "uint16"
	case 4:
		typ = "uint32"
	case 8:
		typ = "uint64"
	default:
		fmt.Fprintln(
			flag.CommandLine.Output(), "bad element size",
		)
		usage()
		os.Exit(1)
	}

	w := bufio.NewWriter(os.Stdout)
	defer flush(w)
	hw := newHexWriter(w, elsiz, big)

	for _, fname := range flag.Args() {
		f, err := os.Open(fname)
		checkErr(err)
		defer f.Close()

		name := prefix + strings.Map(escape, filepath.Base(fname))
		printf(w, "//emgo:const\n")
		printf(w, "var %s = [...]%s{\n", name, typ)

		n, err := io.Copy(hw, f)
		checkErr(err)
		checkErr(hw.Flush())

		printf(w, "}\n")
		printf(w, "const %sSize = %d\n", name, n)
	}
}
