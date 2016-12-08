package main

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

type reg struct {
	Name    string
	BitSiz  int
	Len     int
	Offset  uint64
	Bits    []string
	SubRegs []*reg
	BitRegs []*reg
}

func registers(f string, lines []string, decls []ast.Decl) ([]*reg, []string) {
	var (
		regs    []*reg
		nextoff uint64
	)
loop:
	for len(lines) > 0 {
		line := strings.TrimSpace(lines[0])
		switch line {
		case "Import:", "Instances:":
			break loop
		}
		lines = lines[1:]
		if line == "" {
			continue
		}
		offstr, line := split(line)
		sizstr, line := split(line)
		name, _ := split(line)
		switch "" {
		case offstr:
			fdie(f, "no register offset")
		case sizstr:
			fdie(f, "no register bit size")
		case name:
			fdie(f, "no register name")
		}
		var size int
		switch sizstr {
		case "32":
			size = 4
		case "16":
			size = 2
		case "8":
			size = 1
		default:
			fdie(f, "bad register size %s: not 8, 16, 32", sizstr)
		}
		length := 0
		if name[len(name)-1] == ']' {
			n := strings.IndexByte(name, '[')
			if n <= 0 {
				fdie(f, "bad register name: %s", name)
			}
			l, err := strconv.ParseUint(name[n+1:len(name)-1], 0, 0)
			if err != nil {
				fdie(f, "bad array length in %s", name)
			}
			name = name[:n]
			length = int(l)
		}
		var subregs []*reg
		if name[len(name)-1] == '}' {
			n := strings.IndexByte(name, '{')
			if n <= 0 {
				fdie(f, "bad register name: %s", name)
			}
			for _, sname := range strings.Split(name[n+1:len(name)-1], ",") {
				subregs = append(
					subregs, &reg{Name: sname, BitSiz: size * 8, Len: 1},
				)
			}
			name = name[:n]
		}
		offset, err := strconv.ParseUint(offstr, 0, 64)
		if err != nil {
			fdie(f, "bad offset %s: %v", offstr, err)
		} else if offset&uint64(size-1) != 0 {
			fdie(f, "bad offset %s for %s-bit register", offstr, sizstr)
		}
		for offset > nextoff {
			siz := 4
			for nextoff+uint64(siz) > offset || nextoff&uint64(siz-1) != 0 {
				siz >>= 1
			}
			var lastres *reg
			if len(regs) > 0 {
				lastres = regs[len(regs)-1]
				if lastres.Name != "" || lastres.BitSiz != siz*8 {
					lastres = nil
				}
			}
			if lastres != nil {
				if lastres.Len == 0 {
					lastres.Len = 2
				} else {
					lastres.Len++
				}
			} else {
				regs = append(regs, &reg{
					BitSiz: siz * 8,
					Offset: nextoff,
				})
			}
			nextoff += uint64(siz)
		}
		r := &reg{
			Name:    name,
			BitSiz:  size * 8,
			Len:     length,
			Offset:  offset,
			SubRegs: subregs,
			BitRegs: subregs,
		}
		if len(subregs) == 0 {
			r.BitRegs = []*reg{r}
		}
		if length == 0 {
			length = 1
		}
		nextoff += uint64(size) * uint64(length) * uint64(len(r.BitRegs))
		regs = append(regs, r)
	}
	regmap := make(map[string]*reg)
	for _, r := range regs {
		for _, br := range r.BitRegs {
			regmap[br.Name] = br
		}
	}
	for _, d := range decls {
		g, ok := d.(*ast.GenDecl)
		if !ok || g.Tok != token.CONST {
			continue
		}
		for _, s := range g.Specs {
			v := s.(*ast.ValueSpec)
			t, ok := v.Type.(*ast.Ident)
			if !ok {
				continue
			}
			i := strings.LastIndexByte(t.Name, '_')
			if i < 0 {
				continue
			}
			if t.Name[i+1:] != "Bits" {
				continue
			}
			r := regmap[t.Name[:i]]
			if r == nil {
				continue
			}
			if v.Comment == nil {
				continue
			}
			var n int
			for cl := v.Comment.List; n < len(v.Comment.List); n++ {
				if c := cl[n]; c != nil && strings.HasPrefix(c.Text, "//+") {
					break
				}
			}
			if n == len(v.Comment.List) {
				continue
			}
			for _, id := range v.Names {
				r.Bits = append(r.Bits, id.Name)
			}
		}
	}
	return regs, lines
}
