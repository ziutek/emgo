// This example shows how to use DMA for memory to memory transfers.
package main

import (
	"delay"
	"fmt"
	"rtos"
	"unsafe"

	"stm32/hal/dma"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/rtc"
)

var (
	ch  *dma.Channel
	tce rtos.EventFlag
)

func init() {
	system.Setup(8, 72/8, false)
	rtc.Setup(32768)

	d := dma.DMA1
	d.EnableClock(false)
	ch = d.Channel(1, 0)
	rtos.IRQ(irq.DMA1_Channel1).Enable()
}

// Try different element types (eg. P [...]uint16, M [...]byte) to see how DMA
// handles such asymetrical cases.
var (
	P = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, -1, -2, -3, -4, -5, -6, -7, -8, -9}
	M [len(P)]int
)

func main() {
	ch.Setup(dma.MTM | dma.IncP | dma.IncM)
	ch.SetWordSize(unsafe.Sizeof(P[0]), unsafe.Sizeof(M[0]))
	ch.SetLen(len(P))
	ch.SetAddrP(unsafe.Pointer(&P[0]))
	ch.SetAddrM(unsafe.Pointer(&M[0]))
	ch.EnableIRQ(dma.Complete, dma.ErrAll)
	ch.Enable()
	tce.Wait(0)

	delay.Millisec(250) // Wait for OpenOCD (press reset if you see nothing).

	if _, err := ch.Status(); err != 0 {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Println(M[:])
	}
}

func dmaISR() {
	ev, err := ch.Status()
	ch.Clear(ev, err)
	if ev&dma.Complete != 0 || err != 0 {
		tce.Set()
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.RTCAlarm:      rtc.ISR,
	irq.DMA1_Channel1: dmaISR,
}
