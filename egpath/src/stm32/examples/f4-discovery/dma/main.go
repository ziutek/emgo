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
	ch  dma.Channel
	tce rtos.EventFlag
)

func init() {
	system.Setup168(8)
	systick.Setup()

	d := dma.DMA2
	d.EnableClock(false)
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
	ch.EnableInt(dma.TRCE) // Simplified (should handle dma.ERR too).
	ch.Enable()
	tce.Wait(0)

	delay.Millisec(250) // Wait for OpenOCD (press reset if you see nothing).
	fmt.Println(M[:])
}

func dmaISR() {
	if ch.Events()&dma.TRCE != 0 {
		ch.ClearEvents(dma.TRCE)
		tce.Set()
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.DMA2_Stream0: dmaISR,
}
