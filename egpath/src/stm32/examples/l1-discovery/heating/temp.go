package main

import (
	"delay"
	"fmt"

	"onewire"

	"stm32/hal/dma"
	"stm32/hal/usart"
	"stm32/onedrv"
)

var tempd tempDaemon

type tempDaemon struct {
	m onewire.Master
}

func (td *tempDaemon) Init(u *usart.Periph, rxdma, txdma *dma.Channel) {
	d := usart.NewDriver(u, rxdma, txdma, make([]byte, 16))
	d.EnableClock(true)
	d.SetBaudRate(115200)
	d.SetMode(usart.HalfDuplex)
	d.Enable()
	d.EnableRx()
	d.EnableTx()
	td.m.Driver = onedrv.USARTDriver{d}
}

func printErr(err error) bool {
	if err == nil {
		return false
	}
	fmt.Printf("1-wire: %v\n", err)
	for i := 0; i < 5; i++ {
		ledBlue.Set()
		delay.Millisec(100)
		ledBlue.Clear()
		delay.Millisec(100)
	}
	return true
}

func printOK() {
	fmt.Printf("OK.\n")
}

//emgo:const
var dstypes = [...]onewire.Type{onewire.DS18B20}

func (td *tempDaemon) Loop() {
start:
	for {
		if printErr(td.m.SkipROM()) {
			continue start
		}
		if printErr(td.m.WriteScratchpad(127, -128, onewire.T10bit)) {
			continue start
		}
		if printErr(td.m.SkipROM()) {
			continue start
		}
		if printErr(td.m.ConvertT()) {
			continue start
		}
		for {
			delay.Millisec(50)
			b, err := td.m.ReadBit()
			if printErr(err) {
				continue start
			}
			if b != 0 {
				break
			}
		}
		for _, typ := range dstypes {
			s := onewire.MakeSearch(onewire.DS18B20, false)
			for td.m.SearchNext(&s) {
				d := s.Dev()
				fmt.Printf("%v : ", d)
				if printErr(td.m.MatchROM(d)) {
					continue start
				}
				s, err := td.m.ReadScratchpad()
				if printErr(err) {
					continue start
				}
				t, err := s.Temp(typ)
				if printErr(err) {
					continue start
				}
				fmt.Printf("%6.2f C\n", t)
			}
			if printErr(s.Err()) {
				continue start
			}
		}
		delay.Millisec(4e3)
	}
}

func tempdUSARTISR() {
	tempd.m.Driver.(onedrv.USARTDriver).USART.ISR()
}

func tempdRxDMAISR() {
	tempd.m.Driver.(onedrv.USARTDriver).USART.RxDMAISR()
}

func tempdTxDMAISR() {
	tempd.m.Driver.(onedrv.USARTDriver).USART.TxDMAISR()
}
