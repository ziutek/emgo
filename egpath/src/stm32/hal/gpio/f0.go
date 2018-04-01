// +build f030x6

package gpio

import (
	"unsafe"

	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

const (
	EVENTOUT_AF0 = AF0
	TIM15_AF0    = AF0
	SPI1         = AF0
	MCO_AF0      = AF0
	TIM17_AF0    = AF0
	SWDIO        = AF0
	SWCLK        = AF0
	TIM14_AF0    = AF0
	USART1_AF0   = AF0
	IR_OUT       = AF0
	SPI2_AF0     = AF0
	TIM3_AF0     = AF0
	USART4_AF0   = AF0
	TIM3_ETR     = AF0

	USART1_AF1   = AF1
	USART2       = AF1
	TIM3_AF1     = AF1
	I2C1_AF1     = AF1
	I2C2_AF1     = AF1
	EVENTOUT_AF1 = AF1
	TIM15_AF1    = AF1
	SPI2_AF1     = AF1
	USART3_AF1   = AF1

	TIM16_AF2  = AF2
	TIM17_AF2  = AF2
	TIM1       = AF2
	USART6_AF2 = AF2
	USART5_AF2 = AF2

	EVENTOUT_AF3 = AF3
	I2C1_AF3     = AF3
	TIM15_AF3    = AF3

	USART4_AF4 = AF4
	TIM14_AF4  = AF4
	USART3_AF4 = AF4
	I2C1L_AF4  = AF4
	USART5_AF4 = AF4

	TIM15_AF5  = AF5
	USART6_AF5 = AF5
	TIM16_AF5  = AF5
	TIM17_AF5  = AF5
	MCO_AF5    = AF5
	TIM17_AF5  = AF5
	SPI2_AF5   = AF5
	I2C2_AF5   = AF5

	EVENTOUT_AF6 = AF6
)

const (
	veryLow  = -1 // Not supported.
	low      = -1 // 2 MHz (CL = 50 pF, VDD >= 2.4 V)
	_        = 0  // 10 MHz (CL = 50 pF, VDD >= 2.4 V)
	high     = 2  // 30 MHz (CL = 50 pF, VDD >= 2.7 V)
	veryHigh = 2  // Not supported.
)

//emgo:const
var (
	A = (*Port)(unsafe.Pointer(mmap.GPIOA_BASE))
	B = (*Port)(unsafe.Pointer(mmap.GPIOB_BASE))
	C = (*Port)(unsafe.Pointer(mmap.GPIOC_BASE))
	D = (*Port)(unsafe.Pointer(mmap.GPIOD_BASE))
	F = (*Port)(unsafe.Pointer(mmap.GPIOF_BASE))
)

func enreg() *rcc.RAHBENR   { return &rcc.RCC.AHBENR }
func rstreg() *rcc.RAHBRSTR { return &rcc.RCC.AHBRSTR }

func lpenaclk(pnum uint) {}
func lpdisclk(pnum uint) {}
