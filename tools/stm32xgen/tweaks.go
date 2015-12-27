package main

import (
	"strings"
)

func tweaks(pkg *Package) {
	for _, p := range pkg.Periphs {
		switch p.Name {
		case "RTC":
			rtc(p)
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
		case r.Name == "ALRMASSR":
			r.Name = "ALRMSSR"
			r.Len = 2
			r.Descr = "Alarm A, B subsecond registers"
		case strings.HasPrefix(r.Name, "ALRMB"):
			continue
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
