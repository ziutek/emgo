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
	ds  *dma.Stream
	tce rtos.EventFlag
)

func init() {
	system.Setup168(8)
	systick.Setup()

	DMA := dma.DMA2
	DMA.EnableClock(false)
	ds = DMA.Stream(0)
	rtos.IRQ(irq.DMA2_Stream0).Enable()
}

var (
	P = [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, -1, -2, -3, -4, -5, -6, -7, -8, -9}
	M [len(P)]int
)

func main() {
	ds.Setup(dma.MTM|dma.IncP|dma.IncM, 0)
	ds.SetWordSize(unsafe.Sizeof(P[0]), unsafe.Sizeof(M[0]))
	ds.SetNum(len(P))
	ds.SetAddrP(unsafe.Pointer(&P[0]))
	ds.SetAddrM(unsafe.Pointer(&M[0]))
	ds.EnableInt(dma.TCE) // Simplified (should handle dma.ERR too).
	ds.Enable()
	tce.Wait(0)

	delay.Millisec(250) // Wait for OpenOCD (press reset if you see nothing).
	fmt.Println(M[:])
}

func dmaISR() {
	if ds.Events()&dma.TCE != 0 {
		ds.ClearEvents(dma.TCE)
		tce.Set()
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.DMA2_Stream0: dmaISR,
}
