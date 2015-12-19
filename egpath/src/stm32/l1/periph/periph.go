package periph

import "unsafe"

type regs struct {
	ahbrstr  uint32
	apb2rstr uint32
	apb1rstr uint32

	ahbenr  uint32
	apb2enr uint32
	apb1enr uint32

	ahblpenr  uint32
	apb2lpenr uint32
	apb1lpenr uint32
} //c:volatile

var p = (*regs)(unsafe.Pointer(uintptr(0x40023810)))

type AHBDev uint32

const (
	GPIOA AHBDev = 1 << iota
	GPIOB
	GPIOC
	GPIOD
	GPIOE
	GPIOH
	GPIOF
	GPIOG
	_
	_
	_
	_
	CRC
	_
	_
	FlashIf
	_
	_
	_
	_
	_
	_
	_
	_
	DMA1
	DMA2
	_
	AES
	_
	_
	FSMC
)

func AHBReset(d AHBDev) {
	p.ahbrstr |= uint32(d)
	p.ahbrstr &^= uint32(d)
}

func AHBClockEnable(d AHBDev) {
	p.ahbenr |= uint32(d)
}

func AHBClockDisable(d AHBDev) {
	p.ahbenr &^= uint32(d)
}

type APB1Dev uint32

const (
	Tim2 APB1Dev = 1 << iota
	Tim3
	Tim4
	Tim5
	Tim6
	Tim7
	_
	_
	_
	LCD
	_
	WWdg
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
	USB
	_
	_
	_
	_
	PWR
	DAC
	_
	Comp
)

func APB1Reset(d APB1Dev) {
	p.apb1rstr |= uint32(d)
	p.apb1rstr &^= uint32(d)
}

func APB1ClockEnable(d APB1Dev) {
	p.apb1enr |= uint32(d)
}

func APB1ClockDisable(d APB1Dev) {
	p.apb1enr &^= uint32(d)
}

type APB2Dev uint32

const (
	SysCfg APB2Dev = 1 << iota
	_
	Tim9
	Tim10
	Tim11
	_
	_
	_
	_
	ADC1
	_
	SDIO
	SPI1
	_
	USART1
)

func APB2Reset(d APB2Dev) {
	p.apb2rstr |= uint32(d)
	p.apb2rstr &^= uint32(d)
}

func APB2ClockEnable(d APB2Dev) {
	p.apb2enr |= uint32(d)
}

func APB2ClockDisable(d APB2Dev) {
	p.apb2enr &^= uint32(d)
}
