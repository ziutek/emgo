package main

import (
	"strings"
)

func tweaks(pkg *Package) {
	for _, p := range pkg.Periphs {
		switch p.Name {
		case "RTC":
			rtc(p)
		case "FLASH":
			flash(p)
		case "EXTI":
			exti(p)
		}
	}
}

func rtc(p *Periph) {
	regs := make([]*Register, 0, 20)
	var bkpr *Register
	for _, r := range p.Regs {
		switch {
		case r.Name == "ALRMAR":
			r.Name = "ALRMR"
			r.Len = 2
			r.Descr = "Alarm A, B registers"
			for _, b := range r.Bits {
				b.Name = "A" + b.Name
			}
		case r.Name == "ALRMASSR":
			r.Name = "ALRMSSR"
			r.Len = 2
			r.Descr = "Alarm A, B subsecond registers"
			for _, b := range r.Bits {
				b.Name = "A" + b.Name
			}
		case strings.HasPrefix(r.Name, "ALRMB"):
			continue
		case r.Name == "TSTR" || r.Name == "TSDR" || r.Name == "TSSSR":
			for _, b := range r.Bits {
				b.Name = "T" + b.Name
			}
		case r.Name == "BKP0R":
			bkpr = r
			bkpr.Name = "BKPR"
			bkpr.Len = 1
			bkpr.Descr = "Backup registers"
		case strings.HasPrefix(r.Name, "BKP"):
			bkpr.Len++
			continue
		}
		regs = append(regs, r)
	}
	p.Regs = regs
}

func flash(p *Periph) {
	regs := make([]*Register, 0, 20)
	var optcr *Register
	for _, r := range p.Regs {
		switch {
		case r.Name == "OPTCR":
			optcr = r
			optcr.Len = 1
			optcr.Descr = "Option control registers"
		case strings.HasPrefix(r.Name, "OPTCR"):
			optcr.Len++
			continue
		}
		regs = append(regs, r)
	}
	p.Regs = regs
}

func exti(p *Periph) {
	for _, r := range p.Regs {
		switch r.Name {
		case "IMR":
			for _, b := range r.Bits {
				b.Name = strings.Replace(b.Name, "MR", "IL", 1)
			}
		case "EMR":
			for _, b := range r.Bits {
				b.Name = strings.Replace(b.Name, "MR", "EL", 1)
			}
		case "FTSR":
			for _, b := range r.Bits {
				b.Name = strings.Replace(b.Name, "TR", "TF", 1)
			}

		}
	}
}
