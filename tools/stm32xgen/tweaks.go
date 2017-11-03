package main

import (
	"strings"
)

func lastTweaks(pkg *Package) {
	for _, p := range pkg.Periphs {
		for _, r := range p.Regs {
			fixbits(r)
		}
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
		case "TIM":
			tim(p)
		case "FSMC_Bank1", "FMC_Bank1":
			fsmcBank1_(p)
		case "FSMC_Bank1E", "FMC_Bank1E":
			fsmcBank1e(p)
		case "FMC_Bank5_6":
			fsmcBank5_6(p)
		case "ADC":
			adc(p)
		case "ADC_Common":
			adcCommon(p)
		case "PWR":
			pwr(p)
		case "DMA_Request":
			dmaRequest(p)
		}
	}
}

func fixbits(r *Register) {
	if len(r.Bits) == 1 && r.Bits[0].Name == r.Name {
		r.Bits = nil
		return
	}
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

func rtc(p *Periph) {
	regs := make([]*Register, 0, len(p.Regs))
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
	regs := make([]*Register, 0, len(p.Regs))
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

func tim(p *Periph) {
	for _, r := range p.Regs {
		switch r.Name {
		case "CCMR1":
			for _, b := range r.Bits {
				if b == nil {
					continue
				}
				switch b.Name {
				case "IC1PSC", "IC1F", "IC2PSC", "IC2F":
					b.Val = false
				}
			}
		case "CCMR2":
			for _, b := range r.Bits {
				if b == nil {
					continue
				}
				switch b.Name {
				case "IC3PSC", "IC3F", "IC4PSC", "IC4F":
					b.Val = false
				}
			}
		case "CNT":
			r.Bits = nil
		case "CCR5":
			for _, b := range r.Bits {
				if b == nil {
					continue
				}
				if b.Name == "CCR5" {
					b.Name = "CCR5V"
					break
				}
			}
		}
	}
}

func fsmcBank1_(p *Periph) {
	var btcr *Register
	for i, r := range p.Regs {
		switch r.Name {
		case "BTCR":
			r.Len /= 2
			btcr = r
		case "BCR1":
			r.Name = "BCR"
			btcr.SubRegs = append(btcr.SubRegs, r)
			p.Regs[i] = nil
		case "BTR1":
			r.Name = "BTR"
			btcr.SubRegs = append(btcr.SubRegs, r)
			p.Regs[i] = nil
		}
	}
}

func fsmcBank1e(p *Periph) {
	for _, r := range p.Regs {
		for _, bit := range r.Bits {
			bit.Name = "E" + bit.Name // Temporary workaround.
		}
	}
}

func fsmcBank5_6(p *Periph) {
	for _, r := range p.Regs {
		for _, bit := range r.Bits {
			if strings.HasPrefix(bit.Name, "MWID") {
				bit.Name = "SD" + bit.Name
			}
		}
	}
}

func adc(p *Periph) {
	regs := make([]*Register, 0, len(p.Regs))
	var ofr, jdr *Register
	for _, r := range p.Regs {
		switch {
		case r.Name == "ISR":
			for _, bit := range r.Bits {
				if bit.Name == "ADRD" {
					bit.Name = "ADRDY"
					break
				}
			}
		case r.Name == "IER":
			for _, bit := range r.Bits {
				if bit.Name == "RDY" {
					bit.Name = "ADRDY"
				}
				bit.Name += "IE"
			}
		case r.Name == "OFR1":
			ofr = r
			ofr.Name = "OFR"
			ofr.Len = 1
			ofr.Descr = "Offset registers"
		case strings.HasPrefix(r.Name, "OFR"):
			ofr.Len++
			continue
		case r.Name == "JDR1":
			jdr = r
			jdr.Name = "JDR"
			jdr.Len = 1
			jdr.Descr = "Injected data registers"
		case strings.HasPrefix(r.Name, "JDR"):
			jdr.Len++
			continue
		case r.Name == "DIFSEL":
			r.Bits = nil
		}
		regs = append(regs, r)
	}
	p.Regs = regs
}

func adcCommon(p *Periph) {
	for _, r := range p.Regs {
		switch r.Name {
		case "CCR":
			for _, bit := range r.Bits {
				if bit.Name == "DMACFG" {
					bit.Name = "MDMACFG"
					break
				}
			}
		}
	}
}

func pwr(p *Periph) {
	for _, r := range p.Regs {
		switch {
		case strings.HasPrefix(r.Name, "PUCR"):
			for _, bit := range r.Bits {
				bit.Name = "PU" + bit.Name
			}
		case strings.HasPrefix(r.Name, "PDCR"):
			for _, bit := range r.Bits {
				bit.Name = "PD" + bit.Name
			}
		}
	}
}

func dmaRequest(p *Periph) {
	p.Insts = append(
		p.Insts,
		&Instance{Name: "DMA1_Request", Base: "mmap.DMA1_CSELR_BASE"},
		&Instance{Name: "DMA2_Request", Base: "mmap.DMA2_CSELR_BASE"},
	)
}
