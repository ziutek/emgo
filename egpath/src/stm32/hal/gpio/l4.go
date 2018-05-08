// +build l476xx

package gpio

import (
	"stm32/hal/raw/rcc"
)

const (
	MCO = AF0

	TIM2_AF1 = AF1
	TIM1_AF1 = AF1
	IR_OUT   = AF1
	LPTIM1   = AF1
	TIM8_AF1 = AF1

	TIM5     = AF2
	TIM2_AF2 = AF2
	TIM3     = AF2
	TIM1_AF2 = AF2
	TIM4     = AF2

	TIM8_AF3 = AF3
	TIM1_AF3 = AF3

	I2C3 = AF4
	I2C1 = AF4
	I2C2 = AF4

	SPI1 = AF5
	SPI2 = AF5

	SPI3  = AF6
	DFSDM = AF6

	USART2 = AF7
	USART3 = AF7
	USART1 = AF7

	UART4   = AF8
	UART5   = AF8
	LPUART1 = AF8

	CAN1 = AF9
	TSC  = AF9

	QUADSPI = AF10
	OTG_FS  = AF10

	LCD = AF11

	TIM1_AF12 = AF12
	COMP1     = AF12
	TIM8_AF12 = AF12
	FMC       = AF12
	SDMMC1    = AF12
	SWPMI1    = AF12

	SAI1      = AF13
	SAI2      = AF13
	TIM8_AF13 = AF13

	TIM2_AF14 = AF14
	TIM15     = AF14
	LPTIM2    = AF14
	TIM16     = AF14
	TIM17     = AF14
	TIM8_AF14 = AF14
)

const (
	veryLow  = -1 // Not supported.
	low      = -1 // v MHz (CL = 50 pF, VDD > 2.7 V)
	_        = 0  // x MHz (CL = 50 pF, VDD > 2.7 V)
	high     = 1  // y MHz (CL = 40 pF, VDD > 2.7 V)
	veryHigh = 2  // z MHz (CL = 30 pF, VDD > 2.7 V)
)

func enreg() *rcc.RAHB2ENR   { return &rcc.RCC.AHB2ENR }
func rstreg() *rcc.RAHB2RSTR { return &rcc.RCC.AHB2RSTR }

func lpenaclk(pnum uint) {
	rcc.RCC.AHB2SMENR.AtomicSetBits(rcc.GPIOASMEN << pnum)
}
func lpdisclk(pnum uint) {
	rcc.RCC.AHB2SMENR.AtomicClearBits(rcc.GPIOASMEN << pnum)
}
