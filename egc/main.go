package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
)

var buildCtx = build.Context{
	Compiler:    "gc",
	ReleaseTags: []string{"go1.1", "go1.2", "go1.3", "go1.4"},
	CgoEnabled:  false,
}

var EGCC, EGLD, EGAR string

func getEnv() {
	if EGCC = os.Getenv("EGCC"); EGCC == "" {
		die("EGCC environment variable not set")
	}
	if EGLD = os.Getenv("EGLD"); EGLD == "" {
		die("EGLD environment variable not set")
	}
	if EGAR = os.Getenv("EGAR"); EGAR == "" {
		die("EGAR environment variable not set")
	}

	if buildCtx.GOARCH = os.Getenv("EGARCH"); buildCtx.GOARCH == "" {
		die("EGARCH environment variable not set")
	}
	if buildCtx.GOOS = os.Getenv("EGOS"); buildCtx.GOOS == "" {
		die("EGOS environment variable not set")
	}
	if buildCtx.GOROOT = os.Getenv("EGROOT"); buildCtx.GOROOT == "" {
		die("EGROOT environment variable not set")
	}
	if buildCtx.GOPATH = os.Getenv("EGPATH"); buildCtx.GOPATH == "" {
		die("EGPATH environment variable not set")
	}

}

var (
	tmpDir    string
	verbosity int
	optLevel  string
)

func usage() {
	fmt.Println("Usage:\n  egc [flags] PKGPATH")
	flag.PrintDefaults()
}

func main() {
	flag.IntVar(&verbosity, "v", 0, "Verbosity level [0...2]")
	flag.StringVar(&optLevel, "O", "s", "GCC optimization level")
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) > 1 {
		usage()
		os.Exit(1)
	}

	getEnv()

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
