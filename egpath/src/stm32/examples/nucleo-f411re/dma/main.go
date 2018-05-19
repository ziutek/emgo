// This example tests different ways of coping memory. It also shows how to use
// DMA for memory to memory transfers. In case of STM32F2xx/4xx only DMA2
// supports MTM transfer.
package main

import (
	"delay"
	"fmt"
	"rtos"
	"sync/fence"
	"unsafe"

	"stm32/hal/dma"
	"stm32/hal/irq"
	"stm32/hal/system"
	"stm32/hal/system/timer/systick"
)

var (
	ch     *dma.Channel
	dmaErr dma.Error
	tce    rtos.EventFlag
)

func init() {
	system.Setup96(8)
	systick.Setup(2e6)

	d := dma.DMA2
	d.EnableClock(true)
	ch = d.Channel(0, 0)
	ch.EnableIRQ(dma.Complete, dma.ErrAll)
	rtos.IRQ(irq.DMA2_Stream0).Enable()
}

const n = 40 * 1024 / 4

var (
	src = make([]uint32, n)
	dst = make([]uint32, n)
)

func printSpeed(t int64, check bool) {
	t1 := rtos.Nanosec()
	t2 := rtos.Nanosec()
	dt := (t1 - t) - (t2 - t1)
	if check {
		for i := range dst {
			if dst[i] != uint32(i) {
				fmt.Printf(" dst != src\n")
				return
			}
			dst[i] = 0
		}
	}
	fmt.Printf(" %6d kB/s\n", (int64(n*unsafe.Sizeof(dst[0]))*1e6+dt/2)/dt)
}

func copyDMA(mode dma.Mode) {
	ch.Setup(dma.MTM | dma.IncP | dma.IncM | mode)
	ch.SetWordSize(unsafe.Sizeof(src[0]), unsafe.Sizeof(dst[0]))
	ch.SetLen(n)
	ch.SetAddrP(unsafe.Pointer(&src[0]))
	ch.SetAddrM(unsafe.Pointer(&dst[0]))
	tce.Reset(0)
	fence.W()
	ch.Enable()
	tce.Wait(1, 0)
	if dmaErr != 0 {
		fmt.Println(dmaErr)
	}
}

func main() {
	delay.Millisec(250) // Wait for SWO (press reset if you see nothing).

	fmt.Printf("Initialize src                        ")
	t := rtos.Nanosec()
	for i := range src {
		src[i] = uint32(i)
	}
	printSpeed(t, false)

	fmt.Printf("for i := range src { dst[i] = src[i] }")
	t = rtos.Nanosec()
	for i := range src {
		dst[i] = src[i]
	}
	printSpeed(t, true)

	fmt.Printf("copy(dst, src)                        ")
	t = rtos.Nanosec()
	copy(dst, src)
	printSpeed(t, true)

	fmt.Printf("DMA                                   ")
	t = rtos.Nanosec()
	copyDMA(0)
	printSpeed(t, true)

	fmt.Printf("DMA FT1                               ")
	t = rtos.Nanosec()
	copyDMA(dma.FT1)
	printSpeed(t, true)

	fmt.Printf("DMA FT2                               ")
	t = rtos.Nanosec()
	copyDMA(dma.FT2)
	printSpeed(t, true)

	fmt.Printf("DMA FT3                               ")
	t = rtos.Nanosec()
	copyDMA(dma.FT3)
	printSpeed(t, true)

	fmt.Printf("DMA FT4                               ")
	t = rtos.Nanosec()
	copyDMA(dma.FT4)
	printSpeed(t, true)

	fmt.Printf("DMA FT4 PB4 MB4                       ")
	t = rtos.Nanosec()
	copyDMA(dma.FT4 | dma.PB4 | dma.MB4)
	printSpeed(t, true)
}

func dmaISR() {
	ev, err := ch.Status()
	ch.Clear(ev, err)
	if ev&dma.Complete != 0 || err != 0 {
		dmaErr = err
		tce.Signal(1)
	}
}

//emgo:const
//c:__attribute__((section(".ISRs")))
var ISRs = [...]func(){
	irq.DMA2_Stream0: dmaISR,
}
