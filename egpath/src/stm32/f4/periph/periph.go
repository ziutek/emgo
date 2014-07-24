package periph

import "unsafe"

type regs struct {
	ahb1rst uint32 `C:"volatile"`
	ahb2rst uint32 `C:"volatile"`
	ahb3rst uint32 `C:"volatile"`
	_       uint32 `C:"volatile"`
	apb1rst uint32 `C:"volatile"`
	apb2rst uint32 `C:"volatile"`
	_       uint32 `C:"volatile"`
	_       uint32 `C:"volatile"`

	ahb1en uint32 `C:"volatile"`
	ahb2en uint32 `C:"volatile"`
	ahb3en uint32 `C:"volatile"`
	_      uint32 `C:"volatile"`
	apb1en uint32 `C:"volatile"`
	apb2en uint32 `C:"volatile"`
	_      uint32 `C:"volatile"`
	_      uint32 `C:"volatile"`

	ahb1lpen uint32 `C:"volatile"`
	ahb2lpen uint32 `C:"volatile"`
	ahb3lpen uint32 `C:"volatile"`
	_        uint32 `C:"volatile"`
	apb1lpen uint32 `C:"volatile"`
	apb2lpen uint32 `C:"volatile"`
}

var p = (*regs)(unsafe.Pointer(uintptr(0x40023810)))

type AHB1Dev uint32

const (
	GPIOA AHB1Dev = 1 << iota
	GPIOB
	GPIOC
	GPIOD
	GPIOE
	GPIOF
	GPIOG
	GPIOH

	GPIOI
	GPIOJ
	GPIOK
	_
	CRC
	_
	_
	_

	_
	_
	BkpSRAM
	_
	CCMDataRAM
	DMA1
	DMA2
	DMA2D

	_
	EthMAC
	EthMACTx
	EthMACRx
	EthMACPTP
	OTGHS
	OTGHSULPI
)

func AHB1Reset(d AHB1Dev) {
	p.ahb1rst |= uint32(d)
	p.ahb1rst &^= uint32(d)
}

func AHB1ClockEnable(d AHB1Dev) {
	p.ahb1en |= uint32(d)
}

func AHB1ClockDisable(d AHB1Dev) {
	p.ahb1en &^= uint32(d)
}

type AHB2Dev uint32

const (
	DCMI AHB2Dev = 1 << iota
	_
	_
	_
	CRYP
	HASH
	RNG
	OTGFS
)

func AHB2Reset(d AHB2Dev) {
	p.ahb2rst |= uint32(d)
	p.ahb2rst &^= uint32(d)
}

func AHB2ClockEnable(d AHB2Dev) {
	p.ahb2en |= uint32(d)
}

func AHB2ClockDisable(d AHB2Dev) {
	p.ahb2en &^= uint32(d)
}

type AHB3Dev uint32

const (
	FMC AHB3Dev = 1
)

func AHB3Reset(d AHB3Dev) {
	p.ahb3rst |= uint32(d)
	p.ahb3rst &^= uint32(d)
}

func AHB3ClockEnable(d AHB3Dev) {
	p.ahb3en |= uint32(d)
}

func AHB3ClockDisable(d AHB3Dev) {
	p.ahb3en &^= uint32(d)
}

type APB1Dev uint32

const (
	TIM2 APB1Dev = 1 << iota
	TIM3
	TIM4
	TIM5
	TIM6
	TIM7
	TIM12
	TIM13

	TIM14
	_
	_
	WWDG
	_
	_
	SPI2
	SPI3

	_
	USART2
	USART3
	UART4
	UART5
	I2C1
	I2C2
	I2C3

	_
	CAN1
	CAN2
	_
	PWR
	DAC
	UART7
	UART8
)

func APB1Reset(d APB1Dev) {
	p.apb1rst |= uint32(d)
	p.apb1rst &^= uint32(d)
}

func APB1ClockEnable(d APB1Dev) {
	p.apb1en |= uint32(d)
}

func APB1ClockDisable(d APB1Dev) {
	p.apb1en &^= uint32(d)
}

type APB2Dev uint32

const (
	TIM1 APB2Dev = 1 << iota
	TIM8
	_
	_
	USART1
	USART6
	_
	_

	ADC
	_
	_
	SDIO
	SPI1
	SPI4
	SYSCFG
	_

	TIM9
	TIM10
	TIM11
	_
	SPI5
	SPI6
	SAI1
	_

	_
	_
	LTDC
)

func APB2Reset(d APB2Dev) {
	p.apb2rst |= uint32(d)
	p.apb2rst &^= uint32(d)
}

func APB2ClockEnable(d APB2Dev) {
	p.apb2en |= uint32(d)
}

func APB2ClockDisable(d APB2Dev) {
	p.apb2en &^= uint32(d)
}
