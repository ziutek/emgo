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
	N     int
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
		if err != nil || n&3 != 0 {
			fdie(f, "bad index %s: %v", num, err)
		}
		n >>= 2
		regs = append(regs, reg{
			Reg:   name,
			Bits:  name + "_Bits",
			Field: name + "_Field",
			N:     int(n),
		})
		if max < int(n) {
			max = int(n)
		}
	}
	return regs, max + 1
}
