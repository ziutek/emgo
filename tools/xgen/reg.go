package main

import "strconv"

type decl struct {
	Name string
	Typ  string
	Und  string
}

type reg struct {
	Reg   string
	Bits  string
	Field string
	N     string
	Decls []decl
}

func regs(f string, lines []string) ([]reg, int) {
	var (
		regs []reg
		max  int
	)
	for _, line := range lines {
		num, name := nameval(line, ':')
		if num == "" || name == "" {
			continue
		}
		n, err := strconv.ParseUint(num, 0, 0)
		if err != nil {
			fdie(f, "bad index %s: %v", num, err)
		}
		regs = append(regs, reg{
			Reg:   name,
			Bits:  name + "_Bits",
			Field: name + "_Field",
			N:     num,
		})
		if max < int(n) {
			max = int(n)
		}
	}
	return regs, max + 1
}
