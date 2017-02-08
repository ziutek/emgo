// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl f40_41xxx f411xe f746xx l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package adc

import (
	"stm32/hal/internal"
)

func (p *Periph) enableClock(lp bool) {
	addr := unsafe.Pointer(p)
	internal.APB_SetLPEnabled(addr, lp)
	internal.APB_SetEnabled(addr, true)
}

func (p *Periph) disableClock() {
	internal.APB_SetEnabled(unsafe.Pointer(p), false)
}