package main

import (
	"code.google.com/p/go.tools/go/importer"
	"code.google.com/p/go.tools/go/types"
	"errors"
	"go/build"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var archMap = map[string]string{
	"cortexm3":  "-mcpu=cortex-m3 -mthumb -mfloat-abi=soft",
	"cortexm4":  "-mcpu=cortex-m4 -mthumb -mfloat-abi=soft",
	"cortexm4f": "-mcpu=cortex-m4 -mthumb -mfloat-abi=hard -mfpu=fpv4-sp-d16",
}

var osMap = map[string]struct{ cc, ld string }{
	"none": {
		cc: "-ffreestanding -nostdinc -fno-exceptions -nostartfiles",
		ld: "-nostdlib  -nodefaultlibs  -nostartfiles",
	},
}

type CFLAGS struct {
	Arch string // architecture flags
	OS   string // OS related flags
	Dbg  string // debuging flags
	Opt  string // optimization flags
	Warn string // warning flags
	Incl string // include flags
}

func (cf *CFLAGS) SplitAll() []string {
	var f []string
	f = append(f, strings.Fields(cf.Arch)...)
	f = append(f, strings.Fields(cf.OS)...)
	f = append(f, strings.Fields(cf.Dbg)...)
	f = append(f, strings.Fields(cf.Opt)...)
	f = append(f, strings.Fields(cf.Warn)...)
	f = append(f, strings.Fields(cf.Incl)...)
	return f
}

type LDFLAGS struct {
	OS     string // OS related flags
	Script string // path to linker script or empty if default script
}

func (lf *LDFLAGS) SplitAll() []string {
	var f []string
	f = append(f, strings.Fields(lf.OS)...)
	if lf.Script != "" {
		f = append(f, "-T", lf.Script)
	}
	return f
}

type BuildTools struct {
	CC     string // path to the C compiler
	CFLAGS []string

	LD      string // path to the linker
	LDFLAGS []string

	AR      string // path to the archiver
	ARFLAGS []string

	Log io.Writer

	importPaths []string
}

const (
	prefix     = "/usr/local/stm32/bin/arm-none-eabi-"
	EGCC       = prefix + "gcc"
	EGLD       = prefix + "ld"
	EGAR       = prefix + "ar"
	EGLDSCRIPT = "stm32f407.ld"
)

func NewBuildTools(ctx *build.Context) (*BuildTools, error) {
	pkgoa := filepath.Join("pkg", ctx.GOOS+"_"+ctx.GOARCH)

	importPaths := append([]string{ctx.GOROOT}, strings.Split(ctx.GOPATH, ":")...)
	for i, p := range importPaths {
		importPaths[i] = filepath.Join(p, pkgoa)
	}

	cflags := CFLAGS{
		Dbg:  "-g",
		Opt:  "-Os -fno-common",
		Warn: "-Wall -Wno-parentheses -Wno-unused-function -Wno-unused-variable -Wno-missing-braces -Wno-unused-label",
		Incl: "-I" + filepath.Join(ctx.GOROOT, "src"),
	}
	ldflags := LDFLAGS{
		Script: EGLDSCRIPT,
	}

	if fl, ok := archMap[ctx.GOARCH]; ok {
		cflags.Arch = fl
	} else {
		return nil, errors.New("unknown EGARCH: " + ctx.GOARCH)
	}
	if fl, ok := osMap[ctx.GOOS]; ok {
		cflags.OS = fl.cc
		ldflags.OS = fl.ld
	} else {
		return nil, errors.New("unknown EGOS: " + ctx.GOOS)
	}

	for _, p := range importPaths {
		cflags.Incl += " -I" + p
	}

	ldflags.Script = EGLDSCRIPT

	bt := &BuildTools{
		CC:      EGCC,
		CFLAGS:  cflags.SplitAll(),
		LD:      EGLD,
		LDFLAGS: ldflags.SplitAll(),
		AR:      EGAR,
		ARFLAGS: []string{"rcs"},

		importPaths: importPaths,
	}
	return bt, nil
}

func (bt *BuildTools) logCmd(cmd string, args []string) {
	if bt.Log == nil {
		return
	}
	io.WriteString(bt.Log, cmd)
	for _, a := range args {
		io.WriteString(bt.Log, " "+a)
	}
	bt.Log.Write([]byte{'\n'})
}

func (bt *BuildTools) Compile(o, c string) error {
	args := append(bt.CFLAGS)
	if odir := filepath.Dir(o); filepath.Base(odir) == "main" {
		args = append(args, "-I", odir)
	}
	args = append(args, "-o", o, "-c", c)

	bt.logCmd(bt.CC, args)

	cmd := exec.Command(bt.CC, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (bt *BuildTools) Archive(a string, f ...string) error {
	os.Remove(a)

	args := append(bt.ARFLAGS, a)
	args = append(args, f...)

	bt.logCmd(bt.AR, args)

	cmd := exec.Command(bt.AR, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (bt *BuildTools) getImports(known map[string]*types.Package, add []*types.Package) ([]string, error) {
	var a []string

	for _, ipkg := range add {
		ppath := ipkg.Path()
		if known[ppath] != nil {
			continue
		}
		var (
			apath string
			data  []byte
			err   error
		)
		for _, ipath := range bt.importPaths {
			apath = filepath.Join(ipath, ppath+".a")
			data, err = arReadFile(apath, "__.EXPORTS")
			if err == nil {
				break
			}
			if !os.IsNotExist(err) {
				return nil, err
			}
			apath = ""
		}
		if apath == "" {
			return nil, errors.New("can't find compiled package for " + ppath)
		}
		a = append(a, apath)
		ipkg, err := importer.ImportData(known, data)
		if err != nil {
			return nil, err
		}
		ia, err := bt.getImports(known, ipkg.Imports())
		if err != nil {
			return nil, err
		}
		a = append(a, ia...)
	}
	return a, nil
}

func (bt *BuildTools) Link(e string, imports []*types.Package, o ...string) error {
	args := append(bt.LDFLAGS, "-o", e)
	args = append(args, o...)

	// Find all imported packages with all nested imports
	a, err := bt.getImports(make(map[string]*types.Package), imports)
	if err != nil {
		return err
	}
	args = append(args, a...)

	bt.logCmd(bt.LD, args)

	cmd := exec.Command(bt.LD, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
