// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

package i2c

import (
	"stm32/hal/raw/i2c"
)

var (
	I2C1 = Periph{i2c.I2C1}
	I2C2 = Periph{i2c.I2C2}
)
