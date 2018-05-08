// +build f303xe

package gpio

import (
	"stm32/hal/raw/rcc"
)

const (
	RTC_REFIN     = AF0
	MCO           = AF0
	SWDIO_FTNS    = AF0
	SWCLK_JTCK    = AF0
	JTDI          = AF0
	JTDO_TRACESWO = AF0
	JRST          = AF0
	TRACECK       = AF0
	TRACED0       = AF0
	TRACED1       = AF0
	TRACED2       = AF0
	TRACED3       = AF0

	TIM2_AF1  = AF1
	TIM15_AF1 = AF1
	TIM16_AF1 = AF1
	TIM17_AF1 = AF1
	EVENTOUT  = AF1

	I2C3_AF2  = AF2
	TIM1_AF2  = AF2
	TIM2_AF2  = AF2
	TIM3_AF2  = AF2
	TIM4_AF2  = AF2
	TIM8_AF2  = AF2
	TIM20_AF2 = AF2
	TIM15_AF2 = AF2
	COMP1_AF2 = AF2

	I2C3_AF3  = AF3
	TIM8_AF3  = AF3
	TIM20_AF3 = AF3
	TIM15_AF3 = AF3
	COMP7     = AF3
	TSC       = AF3

	I2C1      = AF4
	I2C2      = AF4
	TIM1_AF4  = AF4
	TIM8_AF4  = AF4
	TIM16_AF4 = AF4
	TIM17_AF4 = AF4

	SPI1     = AF5
	SPI2_AF5 = AF5
	I2S2_AF5 = AF5
	SPI3_AF5 = AF5
	I2S3_AF5 = AF5
	SPI4     = AF5
	UART4    = AF5
	UART5    = AF5
	TIM8_AF5 = AF5
	IR_AF5   = AF5

	SPI2_AF6  = AF6
	I2S2_AF6  = AF6
	SPI3_AF6  = AF6
	I2S3_AF6  = AF6
	TIM1_AF6  = AF6
	TIM8_AF6  = AF6
	TIM20_AF6 = AF6
	IR_AF6    = AF6

	USART1  = AF7
	USART2  = AF7
	USART3  = AF7
	CAN_AF7 = AF7
	COMP3   = AF7
	COMP5   = AF7
	COMP6   = AF7

	I2C3_AF8  = AF8
	COMP1_AF8 = AF8
	COMP2     = AF8
	COMP3_AF8 = AF8
	COMP4     = AF8
	COMP5_AF8 = AF8
	COMP6_AF8 = AF8

	CAN_AF9   = AF9
	TIM1_AF9  = AF9
	TIM8_AF9  = AF9
	TIM15_AF9 = AF9

	TIM2_AF10  = AF10
	TIM3_AF10  = AF10
	TIM4_AF10  = AF10
	TIM8_AF10  = AF10
	TIM17_AF10 = AF10

	TIM1_AF11 = AF11
	TIM8_AF11 = AF11

	FSMC      = AF12
	TIM1_AF12 = AF12
)

const (
	veryLow  = -1 // Not supported.
	low      = -1 // 2 MHz ???
	_        = 0  // 18 MHz ???
	high     = 2  // 36 MHz (CL = 47 pF, VDD = 3.3 V)
	veryHigh = 2  // Not supported.
)

func enreg() *rcc.RAHBENR   { return &rcc.RCC.AHBENR }
func rstreg() *rcc.RAHBRSTR { return &rcc.RCC.AHBRSTR }

func lpenaclk(pnum uint) {}
func lpdisclk(pnum uint) {}
