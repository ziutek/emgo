package periph

import "unsafe"

type regs struct {
	ahb1rstr uint32 `C:"volatile"`
	ahb2rstr uint32 `C:"volatile"`
	ahb3rstr uint32 `C:"volatile"`
	_        uint32 `C:"volatile"`
	apb1rstr uint32 `C:"volatile"`
	apb2rstr uint32 `C:"volatile"`
	_        uint32 `C:"volatile"`
	_        uint32 `C:"volatile"`

	ahb1enr uint32 `C:"volatile"`
	ahb2enr uint32 `C:"volatile"`
	ahb3enr uint32 `C:"volatile"`
	_       uint32 `C:"volatile"`
	apb1enr uint32 `C:"volatile"`
	apb2enr uint32 `C:"volatile"`
	_       uint32 `C:"volatile"`
	_       uint32 `C:"volatile"`

	ahb1lpenr uint32 `C:"volatile"`
	ahb2lpenr uint32 `C:"volatile"`
	ahb3lpenr uint32 `C:"volatile"`
	_         uint32 `C:"volatile"`
	apb1lpenr uint32 `C:"volatile"`
	apb2lpenr uint32 `C:"volatile"`
}

const base uintptr = 0x40023810

var p = (*regs)(unsafe.Pointer(base))

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
	p.ahb1rstr |= uint32(d)
	p.ahb1rstr &^= uint32(d)
}

func AHB1ClockEnable(d AHB1Dev) {
	p.ahb1enr |= uint32(d)
}

func AHB1ClockDisable(d AHB1Dev) {
	p.ahb1enr &^= uint32(d)
}
