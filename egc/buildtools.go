package main

import (
	"errors"
	"go/build"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var archMap = map[string]string{
	"cortexm0":  "-mcpu=cortex-m0 -mthumb -mfloat-abi=soft",
	"cortexm3":  "-mcpu=cortex-m3 -mthumb -mfloat-abi=soft",
	"cortexm4":  "-mcpu=cortex-m4 -mthumb -mfloat-abi=soft",
	"cortexm4f": "-mcpu=cortex-m4 -mthumb -mfloat-abi=hard -mfpu=fpv4-sp-d16",
	"cortexm7f": "-mcpu=cortex-m7 -mthumb -mfloat-abi=hard -mfpu=fpv5-sp-d16",
	"cortexm7d": "-mcpu=cortex-m7 -mthumb -mfloat-abi=hard -mfpu=fpv5-dp-d16",
	"amd64":     "",
}

var osMap = map[string]struct{ cc, ld string }{
	"noos": {
		cc: "-ffreestanding -nostdinc -fno-exceptions -nostartfiles -fno-strict-aliasing",
		ld: "-nostdlib -nodefaultlibs  -nostartfiles -lgcc",
	},
	"linux": {
		cc: "-ffreestanding -nostdinc -fno-exceptions -nostartfiles -fno-strict-aliasing",
		ld: "-nostdlib -nodefaultlibs -nostartfiles -lgcc",
	},
}

type CFLAGS struct {
	Arch string // architecture flags
	OS   string // OS related flags
	Dbg  string // debuging flags
	Opt  string // optimization flags
	Warn string // warning flags
	Incl string // -I flags
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
	Incl   string // -L flags
	Script string // Path to linker script or empty if default script
	Opt    string // Optimization flags
}

func (lf *LDFLAGS) SplitAll() []string {
	var f []string
	f = append(f, strings.Fields(lf.OS)...)
	f = append(f, strings.Fields(lf.Opt)...)
	f = append(f, strings.Fields(lf.Incl)...)
	if lf.Script != "" {
		f = append(f, "-T", lf.Script)
	}
	return f
}

type BuildTools struct {
	CC     string // path to the C compiler
	CFLAGS []string

	LD       string // path to the linker
	LDFLAGS  []string
	LDlibgcc string

	AR      string // path to the archiver
	ARFLAGS []string

	Log io.Writer

	importPaths []string
}

func NewBuildTools(ctx *build.Context) (*BuildTools, error) {
	cflags := CFLAGS{
		Dbg:  "-g",
		Opt:  "-O" + optLevel + " -fplan9-extensions -fno-delete-null-pointer-checks -fno-common -freg-struct-return -ffunction-sections -fdata-sections",
		Warn: "-Wall -Wno-parentheses -Wno-unused-function -Wno-unused-variable -Wno-unused-label -Wno-maybe-uninitialized -Wno-unused-local-typedefs",
		Incl: "-I" + filepath.Join(ctx.GOROOT, "egc"),
	}
	if fl, ok := archMap[ctx.GOARCH]; ok {
		cflags.Arch = fl
	} else {
		return nil, errors.New("unknown EGARCH: " + ctx.GOARCH)
	}

	ldflags := LDFLAGS{
		Incl: "-L" + filepath.Join(ctx.GOROOT, "ld"),
		Opt:  "-gc-sections",
	}
	if ctx.GOOS == "noos" {
		ldflags.Script = "script.ld"
	}

	if fl, ok := osMap[ctx.GOOS]; ok {
		cflags.OS = fl.cc
		ldflags.OS = fl.ld
	} else {
		return nil, errors.New("unknown EGOS: " + ctx.GOOS)
	}
	oat := buildCtx.GOOS + "_" + buildCtx.GOARCH
	if buildCtx.InstallSuffix != "" {
		oat += "_" + buildCtx.InstallSuffix
	}
	pkgoat := filepath.Join("pkg", oat)
	importPaths := append([]string{ctx.GOROOT}, strings.Split(ctx.GOPATH, ":")...)
	for i, p := range importPaths {
		ldflags.Incl += " -L" + filepath.Join(p, "ld")
		p = filepath.Join(p, pkgoat)
		cflags.Incl += " -I" + p
		importPaths[i] = p
	}

	bt := &BuildTools{
		CC:      EGCC,
		CFLAGS:  cflags.SplitAll(),
		LD:      EGLD,
		LDFLAGS: ldflags.SplitAll(),
		AR:      EGAR,
		ARFLAGS: []string{"rcs"},

		importPaths: importPaths,
	}

	for i, f := range bt.LDFLAGS {
		if f == "-lgcc" {
			libgcc, err := exec.Command(
				EGCC,
				append(strings.Fields(cflags.Arch), "-print-libgcc-file-name")...,
			).Output()
			if err != nil {
				return nil, err
			}
			bt.LDFLAGS = append(bt.LDFLAGS[:i], bt.LDFLAGS[i+1:]...)
			bt.LDlibgcc = strings.TrimSpace(string(libgcc))
			break
		}
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
	args := bt.CFLAGS
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

func (bt *BuildTools) getImports(known map[string]struct{}, add []string) ([]string, error) {
	var a []string

	for _, ppath := range add {
		if _, ok := known[ppath]; ok || ppath == "unsafe" {
			continue
		}
		var (
			apath string
			data  []byte
			err   error
		)
		for _, ipath := range bt.importPaths {
			apath = filepath.Join(ipath, ppath+".a")
			data, err = arReadFile(apath, "__.IMPORTS")
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
		known[apath] = struct{}{}
		a = append(a, apath)
		ia, err := bt.getImports(known, strings.Fields(string(data)))
		if err != nil {
			return nil, err
		}
		a = append(a, ia...)
	}
	return a, nil
}

func (bt *BuildTools) Link(e string, imports []string, o ...string) error {
	args := append(bt.LDFLAGS, "-o", e)
	args = append(args, o...)

	// Find all imported packages with all nested imports
	a, err := bt.getImports(make(map[string]struct{}), imports)
	if err != nil {
		return err
	}
	args = append(args, a...)

	if bt.LDlibgcc != "" {
		// Ugly trick to insert libgcc.a before last internal.a.
		lasta := len(args) - 1
		args = append(args[:lasta], bt.LDlibgcc, args[lasta])
	}

	bt.logCmd(bt.LD, args)

	cmd := exec.Command(bt.LD, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
