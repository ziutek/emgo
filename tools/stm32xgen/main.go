// stm32xgen generates STM32 peripheral filex in xgen format.
//
// stm32xgen is usually used this wahy:
//  unifdef -k -f undef.h -D STM32TARGET stm32f4xx.h |stm32xgen
package main

import (
	"fmt"
	"io"
	"os"
	"strings"
	"text/tabwriter"
)

func main() {
	r := newScanner(os.Stdin, "stdin")
	for r.Scan() {
	noscan:
		switch doxy(r.Text(), "@addtogroup") {
		case "Peripheral_registers_structures":
			regs(r)
		default:
			continue
		}
		goto noscan
	}
	checkErr(r.Err())
}

type Register struct {
	Offset int
	Size   string
	Name   string
	Descr  string
}

type Instance struct {
	Name  string
	Addr  string
	Descr string
}

type Periph struct {
	Name  string
	Descr string
	Insts []*Instance
	Regs  []*Register
}

func (p *Periph) WriteTo(w io.Writer) {
	fmt.Fprintf(w, "//\n// Peripheral %s: %s\n", p.Name, p.Descr)
	fmt.Fprintln(w, "// Instances:")
	tw := new(tabwriter.Writer)
	tw.Init(w, 0, 0, 1, ' ', 0)
	for _, inst := range p.Insts {
		fmt.Fprintf(tw, "//  %s\t%s\t%s\n", inst.Name, inst.Addr, inst.Descr)
	}
	tw.Flush()
	fmt.Fprintln(w, "// Registers:")
	for _, r := range p.Regs {
		fmt.Fprintf(
			tw, "//  0x%02X\t%s\t %s\t %s\n",
			r.Offset, r.Size, r.Name, r.Descr,
		)
	}
	tw.Flush()
}

type Package struct {
	Name    string
	Descr   string
	Periphs []*Periph
}

func (pkg *Package) WriteTo(w io.Writer) {
	fmt.Fprintf(
	 	w, "\n// Package %s provides interface to %s.\n",
	 	pkg.Name, pkg.Descr,
	 )
	for _, periph := range pkg.Periphs {
		periph.WriteTo(w)
	}
}

func regs(r *scanner) {
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
			if io != 0 {
				offset += size
				continue
			}
			descr := doxy(line, "/*!<")
			if descr != "" {
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
			regs = append(regs, &Register{offset, typ, reg, descr})
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
				if n := strings.IndexAny(d, " \t"); n > 0 {
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
			continue
		}
		if doxy(line, "@addtogroup") != "" {
			break
		}
	}
	checkErr(r.Err())
	//fmt.Printf("%#v\n", pkgs)
	for _, pkg := range pkgs {
		pkg.WriteTo(checkedWriter{os.Stdout})
	}
}
