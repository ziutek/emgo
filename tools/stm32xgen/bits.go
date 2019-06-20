package main

import (
	"regexp"
	"strconv"
	"strings"
)

func addtoreg(pkgs []*Package, bits *Bits) bool {
	name := bits.Name
	var (
		pack   *Package
		periph *Periph
	)
	m := 0
	for _, pkg := range pkgs {
		for _, p := range pkg.Periphs {
			if strings.HasPrefix(name, p.Name+"_") && len(p.Name) > m {
				pack = pkg
				periph = p
				m = len(p.Name)
				name = strings.TrimPrefix(name, periph.Name+"_")
			}
		}
	}
	if periph == nil {
		for _, pkg := range pkgs {
			for _, p := range pkg.Periphs {
				i := strings.LastIndexByte(p.Name, '_')
				if i == -1 {
					continue
				}
				if strings.HasPrefix(name, p.Name[:i+1]) && i > m {
					pack = pkg
					periph = p
					m = i
					name = name[i+1:]
				}
			}
		}
		if periph == nil {
			return false
		}
	}
	var reg *Register
	m = 0
	for _, r := range periph.Regs {
		if strings.HasPrefix(name, r.Name+"_") && len(r.Name) > m {
			reg = r
			m = len(r.Name)
		}
	}
	if reg == nil {
		// Try other peripherals from this package (handles DMA channels/streams).
		switch {
		case strings.HasPrefix(name, "Sx"): // DMA streams.
			name = name[2:]
		}
		m = 0
		for _, p := range pack.Periphs {
			for _, r := range p.Regs {
				if len(r.Name) > m {
					if strings.HasPrefix(name, r.Name+"_") {
						reg = r
						m = len(r.Name)
					} else if strings.HasPrefix(name, r.Name+"1_") {
						reg = r
						m = len(r.Name) + 1
					}
				}
			}
		}
	}
	if reg == nil {
		return false
	}
	bits.Name = ident(name[m+1:])
	for _, b := range reg.Bits {
		if b.Name == bits.Name && b.Mask == bits.Mask && b.LSL == bits.LSL {
			return true // Duplicate.
		}
	}
	reg.Bits = append(reg.Bits, bits)
	return true
}

func bits(r *scanner, pkgs []*Package) {
	maskPos := make(map[string]struct {
		mask uint32
		pos  uint
	})
	bitRegexp := regexp.MustCompile(`_\d$`)
	for r.Scan() {
		line := strings.TrimSpace(r.Text())
		if def := doxy(line, "#define"); def != "" {
			name, mask := split(def)
			if strings.HasPrefix(name, "ADC12") {
				name = "ADC" + name[5:]
			}
			var descr string
			n := strings.Index(mask, "/*")
			if n > 0 {
				descr = strings.TrimSpace(mask[n+2:])
				descr = strings.TrimPrefix(descr, "!<")
				descr = strings.TrimSuffix(descr, "*/")
				descr = strings.TrimSpace(descr)
				mask = mask[:n]
			}
			mask = strings.TrimSpace(mask)
			mask = strings.Trim(mask, "()")
			mask = strings.TrimSpace(mask)
			if n := strings.IndexByte(mask, ')'); n >= 0 {
				mask = strings.TrimSpace(mask[n+1:])
			}
			var bits *Bits
			if mp, ok := maskPos[name]; ok {
				bits = &Bits{Name: name, Mask: mp.mask, LSL: mp.pos, Descr: descr}
			} else {
				msk := strings.HasSuffix(name, "_Msk")
				if msk {
					if n := strings.Index(mask, "<<"); n > 0 {
						mask = strings.TrimSpace(mask[:n])
					}
				} else {
					if bitRegexp.MatchString(name) {
						warn(name, " looks like a bit alias, and we don't know how to handle those, skipping")
						continue
					}
				}
				mask = strings.TrimSuffix(mask, "U")
				m, err := strconv.ParseUint(mask, 0, 32)
				if err != nil {
					warn("Bad bitmask", mask, ":", err)
					continue
				}
				switch {
				case msk:
					key := name[:len(name)-4]
					mp := maskPos[key]
					mp.mask = uint32(m)
					maskPos[key] = mp
					continue
				case strings.HasSuffix(name, "_Pos"):
					key := name[:len(name)-4]
					mp := maskPos[key]
					mp.pos = uint(m)
					maskPos[key] = mp
					continue
				}
				m32 := uint32(m)
				tz := trailingZeros32(m32)
				bits = &Bits{Name: name, Mask: m32 >> tz, LSL: tz, Descr: descr}
			}
			if !addtoreg(pkgs, bits) {
				warn("Can not assign", name, "to any register.")
			}
		}
		if doxy(line, "@addtogroup") != "" {
			break
		}
	}
}
