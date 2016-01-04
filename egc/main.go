package main

import (
	"flag"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
)

var buildCtx = build.Context{
	Compiler:   "gc",
	CgoEnabled: false,
}

var EGCC, EGLD, EGAR string

func getEnv() {
	EGCC = os.Getenv("EGCC")
	if EGCC == "" {
		die("EGCC environment variable not set")
	}
	EGLD = os.Getenv("EGLD")
	if EGLD == "" {
		die("EGLD environment variable not set")
	}
	EGAR = os.Getenv("EGAR")
	if EGAR == "" {
		die("EGAR environment variable not set")
	}
	buildCtx.GOARCH = os.Getenv("EGARCH")
	if buildCtx.GOARCH == "" {
		die("EGARCH environment variable not set")
	}
	if _, ok := archMap[buildCtx.GOARCH]; !ok {
		die("Unknown EGARCH: " + buildCtx.GOARCH)
	}
	buildCtx.GOOS = os.Getenv("EGOS")
	if buildCtx.GOOS == "" {
		die("EGOS environment variable not set")
	}
	if _, ok := osMap[buildCtx.GOOS]; !ok {
		die("Unknown EGOS: " + buildCtx.GOOS)
	}
	buildCtx.GOROOT = os.Getenv("EGROOT")
	if buildCtx.GOROOT == "" {
		die("EGROOT environment variable not set")
	}
	buildCtx.GOPATH = os.Getenv("EGPATH")
	if buildCtx.GOPATH == "" {
		die("EGPATH environment variable not set")
	}
	if egtarget := os.Getenv("EGTARGET"); egtarget != "" {
		buildCtx.BuildTags = []string{egtarget}
		buildCtx.InstallSuffix = egtarget
	}
}

var (
	tmpDir    string
	verbosity int
	optLevel  string
	disableBC bool
)

func usage() {
	fmt.Println("Usage:\n  egc [flags] PKGPATH")
	flag.PrintDefaults()
}

func main() {
	flag.IntVar(&verbosity, "v", 0, "Verbosity level [0...2]")
	flag.StringVar(&optLevel, "O", "s", "GCC optimization level")
	flag.BoolVar(&disableBC, "B", false, "Disable bounds checking")
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
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	if err = egc(path); err != nil {
		logErr(err)
		os.Exit(1)
	}
}
