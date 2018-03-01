// +build f746xx

package gpio

import (
	"unsafe"

	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

const (
	TIM1 = AF1
	TIM2 = AF1

	TIM3 = AF2
	TIM4 = AF2
	TIM5 = AF2

	TIM8    = AF3
	TIM9    = AF3
	TIM10   = AF3
	TIM11   = AF3
	LPTIM1  = AF3
	CEC_AF3 = AF3

	I2C1    = AF4
	I2C2    = AF4
	I2C3    = AF4
	I2C4    = AF4
	CEC_AF4 = AF4

	SPI1     = AF5
	SPI2_AF5 = AF5
	SPI4     = AF5
	SPI5     = AF5
	SPI6     = AF5

	SPI3_AF6 = AF6
	SAI1     = AF6

	SPI2_AF7    = AF7
	SPI3_AF7    = AF7
	USART1      = AF7
	USART2      = AF7
	USART3      = AF7
	UART5_AF7   = AF7
	SPDIFRX_AF7 = AF7

	SAI2_AF8    = AF8
	USART6      = AF8
	UART4       = AF8
	UART55_AF8  = AF8
	UART7       = AF8
	UART8       = AF8
	SPDIFRX_AF8 = AF8

	CAN1        = AF9
	CAN2        = AF9
	TIM12       = AF9
	TIM13       = AF9
	TIM14       = AF9
	QUADSPI_AF9 = AF9
	LCD_AF9     = AF9

	SAI2_AF10    = AF10
	QUADSPI_AF10 = AF10
	OTG2HS       = AF10
	OTG1FS_AF10  = AF10

	ETH = AF11
	OTG1FS_AF11

	FMC    = AF12
	SDMMS1 = AF12
	OTG2FS = AF12

	DCMI = AF13

	LCD_AF14 = AF14
)

const (
	veryLow  = -1 // Not supported.
	low      = -1 //   2 MHz (CL = 50 pF, VDD > 2.7 V)
	_        = 0  //  25 MHz (CL = 50 pF, VDD > 2.7 V)
	high     = 1  //  50 MHz (CL = 40 pF, VDD > 2.7 V)
	veryHigh = 2  // 100 MHz (CL = 30 pF, VDD > 2.7 V)
)

//emgo:const
var (
	A = (*Port)(unsafe.Pointer(mmap.GPIOA_BASE))
	B = (*Port)(unsafe.Pointer(mmap.GPIOB_BASE))
	C = (*Port)(unsafe.Pointer(mmap.GPIOC_BASE))
	D = (*Port)(unsafe.Pointer(mmap.GPIOD_BASE))
	E = (*Port)(unsafe.Pointer(mmap.GPIOE_BASE))
	F = (*Port)(unsafe.Pointer(mmap.GPIOF_BASE))
	G = (*Port)(unsafe.Pointer(mmap.GPIOG_BASE))
	H = (*Port)(unsafe.Pointer(mmap.GPIOH_BASE))
	I = (*Port)(unsafe.Pointer(mmap.GPIOI_BASE))
	J = (*Port)(unsafe.Pointer(mmap.GPIOJ_BASE))
	K = (*Port)(unsafe.Pointer(mmap.GPIOK_BASE))
)

func enreg() *rcc.RAHB1ENR   { return &rcc.RCC.AHB1ENR }
func rstreg() *rcc.RAHB1RSTR { return &rcc.RCC.AHB1RSTR }

func lpenaclk(pnum uint) {
	rcc.RCC.AHB1LPENR.AtomicSetBits(rcc.GPIOALPEN << pnum)
}
func lpdisclk(pnum uint) {
	rcc.RCC.AHB1LPENR.AtomicClearBits(rcc.GPIOALPEN << pnum)
}
