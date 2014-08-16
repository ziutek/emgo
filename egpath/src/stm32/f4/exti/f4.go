package exti

import (
	"unsafe"
	"mmio"

	"stm32/f4/gpio"
)

const (
	L0 Lines = 1 << iota
	L1
	L2
	L3
	L4
	L5
	L6
	L7
	L8
	L9
	L10
	L11
	L12
	L13
	L14
	L15
	LPVD
	LRTCAlarm
	LOTGFS
	LEth
	LOTGHS
	LTampStamp
	LRTCWkup
	LAll Lines = (1 << iota) - 1
)

var (
	regs1 = (*extiRegs)(unsafe.Pointer(uintptr(0x40013C00)))
	regs2 = (*[4]mmio.Reg32)(unsafe.Pointer(uintptr(0x40013808)))
)

// Connect connects port to exti lines. periph.SysCfg should
// be enabled before use this method.
func (l Lines) Connect(port *gpio.Port) {
	l.connect(uint32(port.Number()))
}
