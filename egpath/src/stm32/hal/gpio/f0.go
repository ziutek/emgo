// +build f030x6

package gpio

import (
	"unsafe"

	"stm32/hal/raw/mmap"
	"stm32/hal/raw/rcc"
)

const (
	EVENTOUT_AF0   = AF0
	TIM15_CH1_AF0  = AF0
	TIM15_CH2_AF0  = AF0
	SPI1_NSS       = AF0
	SPI1_SCK       = AF0
	SPI1_MISO      = AF0
	SPI1_MOSI      = AF0
	MCO_AF0        = AF0
	TIM15_BKIN     = AF0
	TIM17_BKIN_AF0 = AF0
	SWDIO          = AF0
	SWCLK          = AF0
	TIM14_CH1_AF0  = AF0
	USART1_TX_AF0  = AF0
	USART1_RX_AF0  = AF0
	IR_OUT         = AF0
	SPI2_SCK_AF0   = AF0
	SPI2_MISO_AF0  = AF0
	SPI2_MOSI_AF0  = AF0
	TIM3_CH1_AF0   = AF0
	TIM3_CH2_AF0   = AF0
	TIM3_CH3_AF0   = AF0
	TIM3_CH4_AF0   = AF0
	USART4_TX_AF0  = AF0
	USART4_RX_AF0  = AF0
	USART4_CK      = AF0
	TIM3_ETR       = AF0

	USART1_CTS     = AF1
	USART2_CTS     = AF1
	USART1_RTS     = AF1
	USART2_RTS     = AF1
	USART1_TX_AF1  = AF1
	USART2_TX      = AF1
	USART1_RX_AF1  = AF1
	USART2_RX      = AF1
	USART1_CK      = AF1
	USART2_CK      = AF1
	TIM3_CH1_AF1   = AF1
	TIM3_CH2_AF1   = AF1
	TIM3_CH3_AF1   = AF1
	TIM3_CH4_AF1   = AF1
	I2C1_SCL_AF1   = AF1
	I2C1_SDA_AF1   = AF1
	I2C2_SCL_AF1   = AF1
	I2C2_SDA_AF1   = AF1
	EVENTOUT_AF1   = AF1
	TIM15_CH1_AF1  = AF1
	TIM15_CH2_AF1  = AF1
	SPI2_MISO_AF1  = AF1
	SPI2_MOSI_AF1  = AF1
	USART3_TX_AF1  = AF1
	USART3_RX_AF1  = AF1
	USART3_CK_AF1  = AF1
	USART3_RTS_AF1 = AF1

	TIM16_CH1_AF2 = AF2
	TIM17_CH1_AF2 = AF2
	TIM1_BKIN     = AF2
	TIM1_CH1N     = AF2
	TIM1_CH2N     = AF2
	TIM1_CH3N     = AF2
	USART6_TX_AF2 = AF2
	USART6_RX_AF2 = AF2
	USART5_TX_AF2 = AF2
	USART5_RX_AF2 = AF2

	EVENTOUT_AF3   = AF3
	I2C1_SMBA      = AF3
	TIM15_CH1N_AF3 = AF3

	USART4_TX_AF4  = AF4
	USART4_RX_AF4  = AF4
	TIM14_CH1_AF4  = AF4
	USART3_CTS     = AF4
	I2C1_SCL_AF4   = AF4
	I2C1_SDA_AF4   = AF4
	USART4_RTS     = AF4
	USART3_CK_AF4  = AF4
	USART3_RTS_AF4 = AF4
	USART5_TX_AF4  = AF4
	USART5_RX_AF4  = AF4
	SART5_CK_RTS   = AF4
	USART4_CTS     = AF4
	USART3_TX_AF4  = AF4
	USART3_RX_AF4  = AF4

	IM15_CH1N_AF5  = AF5
	USART6_TX_AF5  = AF5
	USART6_RX_AF5  = AF5
	TIM16_CH1_AF5  = AF5
	TIM17_CH1_AF5  = AF5
	MCO_AF5        = AF5
	SCL            = AF5
	SDA            = AF5
	TIM17_BKIN_AF5 = AF5
	SPI2_NSS       = AF5
	SPI2_SCK_AF5   = AF5
	TIM15          = AF5
	I2C2_SCL_AF5   = AF5
	I2C2_SDA_AF5   = AF5

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
