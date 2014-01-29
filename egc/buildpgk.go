package main

import (
	"errors"
	"fmt"
	"go/build"
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
		cc: "-nostdinc -fno-exceptions",
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
	f := make([]string, 0, 6*2)
	f = append(f, strings.Fields(cf.Arch)...)
	f = append(f, strings.Fields(cf.OS)...)
	f = append(f, strings.Fields(cf.Dbg)...)
	f = append(f, strings.Fields(cf.Opt)...)
	f = append(f, strings.Fields(cf.Warn)...)
	f = append(f, strings.Fields(cf.Incl)...)

	return f
}

type BuildTools struct {
	CC     string // path to the C compiler
	CFLAGS []string

	LD      string // path to the linker
	LDFLAGS []string

	AR      string // path to the archiver
	ARFLAGS []string
}

const (
	prefix   = "/usr/local/stm32/bin/arm-none-eabi-"
	cc       = prefix + "gcc"
	ld       = prefix + "ld"
	ar       = prefix + "ar"
	ldscript = "stm32f407.ld"
)

func NewBuildTools(ctx *build.Context) (*BuildTools, error) {
	var (
		ldflags string
		ok      bool
	)

	pkgoa := filepath.Join("pkg", ctx.GOOS+"_"+ctx.GOARCH)

	cflags := CFLAGS{
		Dbg: "-g",
		Incl: "-I" + filepath.Join(ctx.GOROOT, pkgoa) +
			" -I" + filepath.Join(ctx.GOROOT, "/src/builtin"),
		Opt:  "-Os -fno-common",
		Warn: "-Wall -Wno-parentheses -Wno-unused-function",
	}

	if cflags.Arch, ok = archMap[ctx.GOARCH]; !ok {
		return nil, errors.New("unknown EGARCH: " + ctx.GOARCH)
	}
	if fl, ok := osMap[ctx.GOOS]; ok {
		cflags.OS = fl.cc
		ldflags = fl.ld
	} else {
		return nil, errors.New("unknown EGOS: " + ctx.GOOS)
	}

	for _, p := range strings.Split(ctx.GOPATH, ":") {
		cflags.Incl += " -I" + filepath.Join(p, pkgoa)
	}

	bt := &BuildTools{
		CC:      cc,
		CFLAGS:  cflags.SplitAll(),
		LD:      ld,
		LDFLAGS: strings.Fields(ldflags),
		AR:      ar,
		ARFLAGS: []string{"rcs"},
	}
	return bt, nil
}

func (bt *BuildTools) Compile(o, c string) error {
	args := append(bt.CFLAGS, "-o", o, "-c", c)
	cmd := exec.Command(bt.CC, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func (bt *BuildTools) Archive(a string, f ...string) error {
	args := append(bt.ARFLAGS, a)
	args = append(args, f...)
	fmt.Println(bt.AR, args)
	cmd := exec.Command(bt.AR, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
