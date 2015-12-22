package main

import (
	"fmt"
	"strconv"
	"strings"
	"text/tabwriter"
	"unicode"
)

type Bits struct {
	Name  string
	Mask  uint32
	LSL   uint
	Descr string
	Val   bool
}

type Register struct {
	Offset int
	Size   string
	Name   string
	Descr  string
	Bits   []*Bits
}

func (r *Register) fixbits() {
	for i, m := range r.Bits {
		if m.Val {
			continue
		}
		mask := m.Mask << m.LSL
		for _, v := range r.Bits[i+1:] {
			if v.Mask == 0 {
				v.LSL = m.LSL
				v.Val = true
			} else if v.Mask<<v.LSL&mask != 0 {
				if v.LSL > m.LSL {
					v.Mask <<= v.LSL - m.LSL
					v.LSL = m.LSL
				}
				v.Val = true
			}
		}
	}
}

type Instance struct {
	Name  string
	Base  string
	Descr string
}

type Periph struct {
	Name  string
	Descr string
	Insts []*Instance
	Regs  []*Register
}

func (p *Periph) Save(base, pkgname string) {
	w := create(strings.ToLower(p.Name + ".go"))
	defer w.Close()
	fmt.Fprintf(w, "// Peripheral: %s_Periph  %s.\n", p.Name, p.Descr)
	fmt.Fprintln(w, "// Instances:")
	tw := new(tabwriter.Writer)
	tw.Init(w, 0, 0, 1, ' ', 0)
	for _, inst := range p.Insts {
		fmt.Fprintf(tw, "//  %s\t %s\t", inst.Name, inst.Base)
		if inst.Descr != "" {
			fmt.Fprintf(tw, "%s.\n", inst.Descr)
		} else {
			fmt.Fprintln(tw)
		}
	}
	tw.Flush()
	fmt.Fprintln(w, "// Registers:")
	for _, r := range p.Regs {
		fmt.Fprintf(tw, "//  0x%02X\t%s\t %s\t", r.Offset, r.Size, r.Name)
		if r.Descr != "" {
			fmt.Fprintf(tw, "%s.\n", r.Descr)
		} else {
			fmt.Fprintln(tw)
		}
	}
	tw.Flush()
	fmt.Fprintln(w, "// Import:")
	fmt.Fprintln(w, "// ", base+"/mmap")

	fmt.Fprintln(w, "package", pkgname)
	fmt.Fprintln(w)
	for _, r := range p.Regs {
		if len(r.Bits) == 0 {
			continue
		}
		r.fixbits()
		fmt.Fprintln(w, "\nconst (")
		for _, b := range r.Bits {
			fmt.Fprintf(w, "\t%s", b.Name)
			if !b.Val {
				fmt.Fprintf(w, " %s_Bits", r.Name)
			}
			fmt.Fprintf(w, " = 0x%02X << %d", b.Mask, b.LSL)
			if b.Descr != "" {
				fmt.Fprintf(w, " // %s.\n", b.Descr)
			} else {
				fmt.Fprintln(w)
			}
		}
		fmt.Fprintln(w, ")")
	}
}

type Package struct {
	Name    string
	Descr   string
	Periphs []*Periph
}

func (pkg *Package) saveDoc() {
	w := create("0_doc.go")
	defer w.Close()
	fmt.Fprintf(
		w, "// Package %s provides interface to %s.\npackage %s\n",
		pkg.Name, pkg.Descr, pkg.Name,
	)
}

func (pkg *Package) Save(base string) {
	mkdir(pkg.Name)
	chdir(pkg.Name)
	defer chdir("..")
	pkg.saveDoc()
	for _, periph := range pkg.Periphs {
		periph.Save(base, pkg.Name)
	}
}

func peripherals(r *scanner) []*Package {
	var (
		pkgs   []*Package
		pkg    *Package
		brief  string
		pbase  string
		regs   []*Register
		offset int
	)
	for r.Scan() {
		line := strings.TrimSpace(r.Text())
		if bri := doxy(line, "@brief"); bri != "" {
			brief = bri
			continue
		}
		if io := strings.Index(line, "__IO"); io == 0 ||
			strings.HasPrefix(line, "uint") {
			if io == 0 {
				line = strings.TrimSpace(line[len("__IO"):])
			}
			n := strings.IndexByte(line, ';')
			if n < 0 {
				r.Die("';' expected after register name")
			}
			tr := strings.Fields(line[:n])
			if len(tr) != 2 {
				r.Die("wrong number of fields before ';'")
			}
			line = line[n+1:]
			typ, reg := tr[0], tr[1]
			var size int
			switch typ {
			case "uint32_t":
				typ = "32"
				size = 4
			case "uint16_t":
				typ = "16"
				size = 2
			case "uint8_t":
				typ = " 8"
				size = 1
			default:
				r.Die("unknown type:", typ)
			}
			if n := len(reg) - 1; n >= 0 && reg[n] == ']' {
				m := strings.Index(reg, "[")
				if m < 0 {
					die("Bad register name:", reg)
				}
				n, err := strconv.ParseUint(reg[m+1:n], 0, 32)
				checkErr(err)
				size *= int(n)
			}
			if io != 0 {
				offset += size
				continue
			}
			var descr string
			if n := strings.Index(line, "/*"); n > 0 {
				descr = strings.TrimPrefix(line[n+2:], "!<")
				if n := strings.LastIndex(descr, "ddress offset:"); n > 0 {
					descr = descr[:n-1]
				} else if n := strings.LastIndex(descr, "*/"); n >= 0 {
					descr = descr[:n]
				}
				descr = strings.TrimSpace(descr)
				if n := len(descr); n > 0 {
					n--
					switch descr[n] {
					case '.', ',', ';':
						descr = descr[:n]
					}
				}
			}
			regs = append(
				regs,
				&Register{Offset: offset, Size: typ, Name: reg, Descr: descr},
			)
			offset += size
			continue
		}
		if strings.HasPrefix(line, "}") {
			line = strings.TrimSpace(line[1:])
			n := strings.Index(line, "_TypeDef;")
			if n < 0 {
				r.Die("name of type (*_TypeDef) expected after '}'")
			}
			periph := line[:n]
			pb := periph
			if n := strings.IndexByte(pb, '_'); n > 0 {
				pb = periph[:n]
			}
			if pbase != pb {
				pbase = pb
				pkg = &Package{Name: strings.ToLower(pbase)}
				pkgs = append(pkgs, pkg)
			}
			if periph == pbase {
				pkg.Descr = brief
			}
			for _, reg := range regs {
				d := reg.Descr
				if n := strings.IndexFunc(d, unicode.IsSpace); n > 0 {
					if d[:n] == periph {
						d = strings.TrimSpace(d[n+1:])
						reg.Descr = upperFirst(d)
					}
				}
			}
			pkg.Periphs = append(
				pkg.Periphs,
				&Periph{Name: periph, Descr: brief, Regs: regs},
			)
			regs = nil
			offset = 0
			continue
		}
		if doxy(line, "@addtogroup") != "" {
			break
		}
	}
	checkErr(r.Err())
	return pkgs
}
