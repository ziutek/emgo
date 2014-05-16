package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
)

var tmpDir string

var buildCtx = build.Context{
	GOARCH:      "cortexm4f",
	GOOS:        "noos",
	GOROOT:      "/home/michal/P/go/src/github.com/ziutek/emgo/egroot",
	GOPATH:      "/home/michal/P/go/src/github.com/ziutek/emgo/egpath",
	Compiler:    "gc",
	ReleaseTags: []string{"go1.1", "go1.2"},
	CgoEnabled:  false,
}

var verbose int

func usage() {
	fmt.Println("Usage:\n  egc [flags] PKGPATH")
	flag.PrintDefaults()
}

func main() {
	flag.IntVar(&verbose, "v", 0, "Verbose level [0...2]")
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) > 1 {
		usage()
		os.Exit(1)
	}

	path := "."
	if len(args) == 1 {
		path = args[0]
	}

	var err error

	tmpDir, err = ioutil.TempDir("", "eg-build")
	if err != nil {
		logErr(err)
		return
	}
	defer os.RemoveAll(tmpDir)

	if err = egc(path); err != nil {
		logErr(err)
		return
	}
}
