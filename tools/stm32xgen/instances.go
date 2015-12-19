package main

import (
	"strings"
)

func instances(r *scanner, pkgs []*Package) {
	for r.Scan() {
		line := strings.TrimSpace(r.Text())
		if def := doxy(line, "#define"); def != "" {
			inst, base := split(def)
			base = strings.TrimSpace(strings.Trim(base, "()"))
			n := strings.Index(base, "_TypeDef")
			if n < 0 {
				continue
			}
			periph := base[:n]
			base = strings.TrimSpace(base[n+len("_TypeDef"):])
			base = strings.TrimSpace(strings.TrimLeft(base, "*)"))
		loop:
			for _, pkg := range pkgs {
				for _, p := range pkg.Periphs {
					if p.Name == periph {
						p.Insts = append(
							p.Insts, &Instance{Name: inst, Base: "mmap." + base},
						)
						break loop
					}
				}
			}
		}
		if doxy(line, "@addtogroup") != "" {
			break
		}
	}
}
