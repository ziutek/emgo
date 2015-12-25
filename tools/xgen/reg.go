package main

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

type reg struct {
	Name   string
	Size   uint
	Len    int
	Offset uint64
	Bits   []string
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
			fdie(f, "no register size")
		case name:
			fdie(f, "no register name")
		}
		var size uint
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
		offset, err := strconv.ParseUint(offstr, 0, 64)
		if err != nil {
			fdie(f, "bad offset %s: %v", offstr, err)
		} else if offset&uint64(size-1) != 0 {
			fdie(f, "bad offset %s for %s-bit register", offstr, sizstr)
		}
		for offset > nextoff {
			var siz uint = 4
			for nextoff+uint64(siz) > offset || nextoff&uint64(siz-1) != 0 {
				siz >>= 1
			}
			if last := regs[len(regs)-1]; last.Name == "" && last.Size == siz*8 {
				if last.Len == 0 {
					last.Len += 2
				} else {
					last.Len++
				}
			} else {
				regs = append(regs, &reg{
					Size:   siz * 8,
					Offset: nextoff,
				})
			}
			nextoff += uint64(siz)
		}
		regs = append(regs, &reg{
			Name:   name,
			Size:   size * 8,
			Len:    length,
			Offset: offset,
		})
		if length == 0 {
			length = 1
		}
		nextoff += uint64(size) * uint64(length)
	}
	regmap := make(map[string]*reg)
	for _, r := range regs {
		regmap[r.Name] = r
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
