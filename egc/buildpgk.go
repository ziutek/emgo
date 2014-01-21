package main

import (
//	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
)

var archMap = map[string]string{
	"cortexm3": "-mcpu=cortex-m3 -mthumb -mfloat-abi=soft",
	"cortexm4": "-mcpu=cortex-m4 -mthumb -mfloat-abi=soft",
	"cortexm4f": "-mcpu=cortex-m4 -mthumb -mfloat-abi=hard " +
		"-mfpu=fpv4-sp-d16 -fsingle-precision-constant",
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

func determineBuildTools() *BuildTools {
	var (
		ldflags string
		ok      bool
	)

	cflags := CFLAGS{
		Dbg: "-g",
		Incl: "-I" + path.Join(buildCtx.GOROOT, "/pkg/include") +
			" -I" + path.Join(buildCtx.GOROOT, "/src/builtin") +
			" -I" + path.Join(buildCtx.GOROOT, "/src/pkg"),
		Opt:  "-Os -fno-common",
		Warn: "-Wall -Wno-parentheses -Wno-unused-function",
	}

	if cflags.Arch, ok = archMap[buildCtx.GOARCH]; !ok {
		return nil
	}
	if fl, ok := osMap[buildCtx.GOOS]; ok {
		cflags.OS = fl.cc
		ldflags = fl.ld
	} else {
		return nil
	}

	for _, p := range strings.Split(buildCtx.GOPATH, ":") {
		cflags.Incl += " -I" + path.Join(p, "pkg/include")
	}

	bt := &BuildTools{
		CC:      cc,
		CFLAGS:  cflags.SplitAll(),
		LD:      ld,
		LDFLAGS: strings.Fields(ldflags),
		AR:      ar,
		ARFLAGS: []string{"rcs"},
	}

	//fmt.Printf("build tools:\n %+v\n", bt)

	return bt
}

func (bt *BuildTools) Compile(c string) error {
	args := append(bt.CFLAGS, "-c", c)
	cmd := exec.Command(bt.CC, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
