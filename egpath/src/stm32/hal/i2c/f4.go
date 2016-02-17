// +build f40_41xxx f411xe

package i2c

import (
	"stm32/hal/raw/i2c"
)

var (
	I2C1 = Periph{i2c.I2C1}
	I2C2 = Periph{i2c.I2C2}
	I2C3 = Periph{i2c.I2C3}
)
