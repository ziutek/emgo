package main

import (
	"delay"
	"fmt"

	"onewire"

	"stm32/hal/dma"
	"stm32/hal/usart"
	"stm32/onedrv"
)

type OneWireDaemon struct {
	m   onewire.Master
	Cmd chan interface{}
}

func (d *OneWireDaemon) Start(u *usart.Periph, rxdma, txdma *dma.Channel) {
	drv := usart.NewDriver(u, rxdma, txdma, make([]byte, 16))
	drv.EnableClock(true)
	drv.SetBaudRate(115200)
	drv.SetMode(usart.HalfDuplex)
	drv.Enable()
	drv.EnableRx()
	drv.EnableTx()
	d.m.Driver = onedrv.USARTDriver{drv}
	d.Cmd = make(chan interface{}, 1)
	go d.loop()
}

func checkPrintErr(err error) bool {
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

type SearchCmd struct {
	Typ  onewire.Type
	Resp chan onewire.Dev
}

type TempCmd struct {
	Dev  onewire.Dev
	Resp chan float32
}

func (d *OneWireDaemon) loop() {
	for cmd := range d.Cmd {
		switch c := cmd.(type) {
		case SearchCmd:
			s := onewire.MakeSearch(c.Typ, false)
			for d.m.SearchNext(&s) {
				c.Resp <- s.Dev()
			}
			c.Resp <- onewire.Dev{}
			checkPrintErr(s.Err())
		case TempCmd:
			for {
				if checkPrintErr(d.m.MatchROM(c.Dev)) {
					continue
				}
				if checkPrintErr(d.m.WriteScratchpad(
					127, -128, onewire.T10bit,
				)) {
					continue
				}
				if checkPrintErr(d.m.MatchROM(c.Dev)) {
					continue
				}
				if checkPrintErr(d.m.ConvertT()) {
					continue
				}
				delay.Millisec(200)
				if checkPrintErr(d.m.MatchROM(c.Dev)) {
					continue
				}
				s, err := d.m.ReadScratchpad()
				if checkPrintErr(err) {
					continue
				}
				t, err := s.Temp(c.Dev.Type())
				if checkPrintErr(err) {
					continue
				}
				c.Resp <- t
				break
			}
		}
	}
}

var owd OneWireDaemon

func owdUSARTISR() {
	owd.m.Driver.(onedrv.USARTDriver).USART.ISR()
}

func owdRxDMAISR() {
	owd.m.Driver.(onedrv.USARTDriver).USART.RxDMAISR()
}

func owdTxDMAISR() {
	owd.m.Driver.(onedrv.USARTDriver).USART.TxDMAISR()
}
