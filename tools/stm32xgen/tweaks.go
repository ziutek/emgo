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
		case "BKP":
			bkp(p)
		case "I2C":
			i2c(p)
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
		case r.Name == "PRLH", r.Name == "PRLL",
			r.Name == "DIVH", r.Name == "DIVL",
			r.Name == "CNTH", r.Name == "CNTL",
			r.Name == "ALRH", r.Name == "ALRL":
			r.Bits = nil
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
		case strings.HasPrefix(r.Name, "WRPR"):
			r.Bits = nil
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

func bkp(p *Periph) {
	for _, r := range p.Regs {
		if !strings.HasPrefix(r.Name, "DR") {
			continue
		}
		if c := r.Name[2]; c < '0' || c > '9' {
			continue
		}
		r.Bits = nil
	}
}

func i2c(p *Periph) {
	for _, r := range p.Regs {
		switch r.Name {
		case "OAR2":
			for _, b := range r.Bits {
				if b.Name == "ADD2" {
					b.Name = "SECADD1_7"
					break
				}
			}
		case "SR2":
			for _, b := range r.Bits {
				if b.Name == "PEC" {
					b.Name = "PECVAL"
					break
				}
			}
		case "CCR":
			for _, b := range r.Bits {
				if b.Name == "CCR" {
					b.Name = "CCRVAL"
					break
				}
			}
		}
	}
}
