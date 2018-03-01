// +build f303xe l476xx

package exti

import (
	"stm32/hal/raw/exti"
)

type lines uint64

func risiTrigEnabled() Lines {
	return Lines(exti.EXTI.RTSR1.Load()) | Lines(exti.EXTI.RTSR2.Load())<<32
}

func (li Lines) enableRisiTrig() {
	if m := exti.RTSR1(li); m != 0 {
		exti.EXTI.RTSR1.AtomicSetBits(m)
	}
	if m := exti.RTSR2(li >> 32); m != 0 {
		exti.EXTI.RTSR2.AtomicSetBits(m)
	}
}

func (li Lines) disableRisiTrig() {
	if m := exti.RTSR1(li); m != 0 {
		exti.EXTI.RTSR1.AtomicClearBits(m)
	}
	if m := exti.RTSR2(li >> 32); m != 0 {
		exti.EXTI.RTSR2.AtomicClearBits(m)
	}
}

func fallTrigEnabled() Lines {
	return Lines(exti.EXTI.FTSR1.Load()) | Lines(exti.EXTI.FTSR2.Load())<<32
}

func (li Lines) enableFallTrig() {
	if m := exti.FTSR1(li); m != 0 {
		exti.EXTI.FTSR1.AtomicSetBits(m)
	}
	if m := exti.FTSR2(li >> 32); m != 0 {
		exti.EXTI.FTSR2.AtomicSetBits(m)
	}
}

func (li Lines) disableFallTrig() {
	if m := exti.FTSR1(li); m != 0 {
		exti.EXTI.FTSR1.AtomicClearBits(m)
	}
	if m := exti.FTSR2(li >> 32); m != 0 {
		exti.EXTI.FTSR2.AtomicClearBits(m)
	}
}

func (li Lines) trigger() {
	if m := exti.SWIER1(li); m != 0 {
		exti.EXTI.SWIER1.Store(m)
	}
	if m := exti.SWIER2(li >> 32); m != 0 {
		exti.EXTI.SWIER2.Store(m)
	}
}

func irqEnabled() Lines {
	return Lines(exti.EXTI.IMR1.Load()) | Lines(exti.EXTI.IMR2.Load())<<32
}

func (li Lines) enableIRQ() {
	if m := exti.IMR1(li); m != 0 {
		exti.EXTI.IMR1.AtomicSetBits(m)
	}
	if m := exti.IMR2(li >> 32); m != 0 {
		exti.EXTI.IMR2.AtomicSetBits(m)
	}
}

func (li Lines) disableIRQ() {
	if m := exti.IMR1(li); m != 0 {
		exti.EXTI.IMR1.AtomicClearBits(m)
	}
	if m := exti.IMR2(li >> 32); m != 0 {
		exti.EXTI.IMR2.AtomicClearBits(m)
	}
}

func eventEnabled() Lines {
	return Lines(exti.EXTI.EMR1.Load()) | Lines(exti.EXTI.EMR2.Load())<<32
}

func (li Lines) enableEvent() {
	if m := exti.EMR1(li); m != 0 {
		exti.EXTI.EMR1.AtomicSetBits(m)
	}
	if m := exti.EMR2(li >> 32); m != 0 {
		exti.EXTI.EMR2.AtomicSetBits(m)
	}
}

func (li Lines) disableEvent() {
	if m := exti.EMR1(li); m != 0 {
		exti.EXTI.EMR1.AtomicClearBits(m)
	}
	if m := exti.EMR2(li >> 32); m != 0 {
		exti.EXTI.EMR2.AtomicClearBits(m)
	}
}

func pending() Lines {
	return Lines(exti.EXTI.PR1.Load()) | Lines(exti.EXTI.PR2.Load())<<32
}

func (li Lines) clearPending() {
	if m := exti.PR1(li); m != 0 {
		exti.EXTI.PR1.Store(m)
	}
	if m := exti.PR2(li >> 32); m != 0 {
		exti.EXTI.PR2.Store(m)
	}
}
