// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl f40_41xxx f411xe l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package usart

import (
	"stm32/hal/raw/usart"
)

const (
	lbd   = usart.LBD
	lbdie = usart.LBDIE
)
