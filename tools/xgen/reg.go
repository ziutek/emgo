package main

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"
)

type reg struct {
	Name   string
	Size   uint64
	Offset uint64
	Bits   []string
}

func regs(f string, lines []string, decls []ast.Decl) []*reg {
	var (
		regs    []*reg
		nextoff uint64
	)
	for _, line := range lines {
		offstr, name := split(line)
		name, _ = split(name)
		if offstr == "" || name == "" {
			continue
		}
		size := "32"
		var siz uint64
		switch size {
		case "32":
			siz = 4
		case "16":
			siz = 2
		case "8":
			siz = 1
		default:
			fdie(f, "bad size %s: not 8, 16, 32", size)
		}
		offset, err := strconv.ParseUint(offstr, 0, 0)
		if err != nil {
			fdie(f, "bad offset %s: %v", offstr, err)
		} else if offset&(siz-1) != 0 {
			fdie(f, "bad offset %s for %s-bit register", offstr, size)
		}
		for offset > nextoff {
			var siz uint64 = 4
			for nextoff+siz > offset || nextoff&(siz-1) != 0 {
				siz >>= 1
			}
			regs = append(regs, &reg{
				Size:   siz * 8,
				Offset: nextoff,
			})
			nextoff += siz
		}
		regs = append(regs, &reg{
			Name:   name,
			Size:   siz * 8,
			Offset: offset,
		})
		nextoff += siz
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
			for _, id := range v.Names {
				r.Bits = append(r.Bits, id.Name)
			}
		}
	}
	return regs
}
