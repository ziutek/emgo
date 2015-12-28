package main

import (
	"fmt"
	"io"
	"strings"
)

type MemGroup struct {
	Descr string
	Bases []*MemBase
}

func (g *MemGroup) WriteTo(w io.Writer) {
	if g.Descr != "" {
		fmt.Fprintln(w, "//", g.Descr)
	}
	fmt.Fprintln(w, "const (")
	for _, b := range g.Bases {
		fmt.Fprintf(w, "\t%s uintptr = %s", b.Name, b.Addr)
		if b.Descr != "" {
			fmt.Fprintln(w, " //", b.Descr)
		} else {
			fmt.Fprintln(w)
		}
	}
	fmt.Fprintln(w, ")")
}

type MemBase struct {
	Name  string
	Addr  string
	Descr string
}

func memmap(r *scanner) []*MemGroup {
	var groups []*MemGroup
	group := new(MemGroup)
	for r.Scan() {
		line := strings.TrimSpace(r.Text())
		if strings.HasPrefix(line, "/*!<") {
			descr := strings.TrimPrefix(line, "/*!<")
			descr = strings.TrimSpace(strings.TrimSuffix(descr, "*/"))
			groups = append(groups, group)
			group = &MemGroup{Descr: descr}
			continue
		}
		if def := doxy(line, "#define"); def != "" {
			name, addr := split(def)
			var descr string
			if n := strings.Index(addr, "/*!<"); n > 0 {
				descr = addr[n+4:]
				descr = strings.TrimSpace(strings.TrimSuffix(descr, "*/"))
				addr = strings.TrimSpace(addr[:n])
			}
			addr = strings.Trim(addr, "()")
			if n := strings.Index(addr, ")"); n >= 0 {
				addr = addr[n+1:]
			}
			addr = strings.TrimSpace(addr)
			group.Bases = append(group.Bases, &MemBase{name, addr, descr})
			continue
		}
		if doxy(line, "@addtogroup") != "" {
			break
		}
	}
	return groups
}

func saveMmap(mmap []*MemGroup) {
	mkdir("mmap")
	chdir("mmap")
	defer chdir("..")
	w := create("mmap.go")
	defer w.Close()
	fmt.Fprintln(
		w, "// Package mmap provides base memory adresses for all peripherals.",
	)
	fmt.Fprintln(w, "package mmap")
	w.donotedit()
	for _, g := range mmap {
		g.WriteTo(w)
	}
}
