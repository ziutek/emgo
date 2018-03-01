// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl f40_41xxx f411xe f746xx l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package exti

import (
	"stm32/hal/raw/exti"
)

type lines uint32

func risiTrigEnabled() Lines {
	return Lines(exti.EXTI.RTSR.Load())
}

func (li Lines) enableRisiTrig() {
	exti.EXTI.RTSR.AtomicSetBits(exti.RTSR(li))
}

func (li Lines) disableRisiTrig() {
	exti.EXTI.RTSR.AtomicClearBits(exti.RTSR(li))
}

func fallTrigEnabled() Lines {
	return Lines(exti.EXTI.FTSR.Load())
}

func (li Lines) enableFallTrig() {
	exti.EXTI.FTSR.AtomicSetBits(exti.FTSR(li))
}

func (li Lines) disableFallTrig() {
	exti.EXTI.FTSR.AtomicClearBits(exti.FTSR(li))
}

func (li Lines) trigger() {
	exti.EXTI.SWIER.Store(exti.SWIER(li))
}

func irqEnabled() Lines {
	return Lines(exti.EXTI.IMR.Load())
}

func (li Lines) enableIRQ() {
	exti.EXTI.IMR.AtomicSetBits(exti.IMR(li))
}

func (li Lines) disableIRQ() {
	exti.EXTI.IMR.AtomicClearBits(exti.IMR(li))
}

func eventEnabled() Lines {
	return Lines(exti.EXTI.EMR.Load())
}

func (li Lines) enableEvent() {
	exti.EXTI.EMR.AtomicSetBits(exti.EMR(li))
}

func (li Lines) disableEvent() {
	exti.EXTI.EMR.AtomicClearBits(exti.EMR(li))
}

func pending() Lines {
	return Lines(exti.EXTI.PR.Load())
}

func (li Lines) clearPending() {
	exti.EXTI.PR.Store(exti.PR(li))
}
