// This example shows how to use DMA for memory to memory transfers. In case of
// STM32F2xx/4xx only DMA2 supports MTM transfer.
package main

import (
	"delay"
	"fmt"
	"rtos"
	"unsafe"

	"stm32/hal/dma"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var (
	ch  *dma.Channel
	tce rtos.EventFlag
)

func init() {
	system.Setup96(8)
	systick.Setup()

	d := dma.DMA2
	d.EnableClock(true)
	ch = d.Channel(0, 0)
	rtos.IRQ(irq.DMA2_Stream0).Enable()
}

var (
	P = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, -1, -2, -3, -4, -5, -6, -7, -8, -9}
	M [len(P)]int
)

func main() {
	ch.Setup(dma.MTM | dma.IncP | dma.IncM | dma.FIFO_4_4)
	ch.SetWordSize(unsafe.Sizeof(P[0]), unsafe.Sizeof(M[0]))
	ch.SetLen(len(P))
	ch.SetAddrP(unsafe.Pointer(&P[0]))
	ch.SetAddrM(unsafe.Pointer(&M[0]))
	ch.EnableIRQ(dma.Complete, dma.ErrAll)
	ch.Enable()
	tce.Wait(1, 0)

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
		tce.Signal(1)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.DMA2_Stream0: dmaISR,
}
