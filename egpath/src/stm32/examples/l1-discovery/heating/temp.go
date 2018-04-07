package main

import (
	"delay"
	"fmt"

	"onewire"

	"stm32/hal/dma"
	"stm32/hal/usart"
	"stm32/onedci"
)

type OneWireDaemon struct {
	m   onewire.Master
	Cmd chan interface{}
}

func (d *OneWireDaemon) Start(u *usart.Periph, txdma, rxdma *dma.Channel) {
	drv := usart.NewDriver(u, txdma, rxdma, make([]byte, 16))
	drv.Periph().EnableClock(true)
	drv.Periph().SetMode(usart.HalfDuplex | usart.OneBit)
	drv.Periph().Enable()
	drv.EnableRx()
	drv.EnableTx()
	d.m.DCI = onedci.USART{drv}
	d.Cmd = make(chan interface{}, 1)
	go d.loop()
}

func log1wireErr(err error) bool {
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

type ConfigureCmd struct {
	Dev  onewire.Dev
	Cfg  byte
	Resp chan onewire.Dev
}

type TempCmd struct {
	Dev  onewire.Dev
	Resp chan int
}

const InvalidTemp = -300 * 16

func (d *OneWireDaemon) loop() {
	for cmd := range d.Cmd {
		switch c := cmd.(type) {
		case SearchCmd:
			s := onewire.MakeSearch(c.Typ, false)
			for d.m.SearchNext(&s) {
				c.Resp <- s.Dev()
			}
			c.Resp <- onewire.Dev{}
			log1wireErr(s.Err())
		case ConfigureCmd:
			{
				if c.Dev.Type() == 0 {
					goto abortConfigureCmd
				}
				if log1wireErr(d.m.MatchROM(c.Dev)) {
					goto abortConfigureCmd
				}
				if log1wireErr(d.m.WriteScratchpad(127, -128, c.Cfg)) {
					goto abortConfigureCmd
				}
				if log1wireErr(d.m.MatchROM(c.Dev)) {
					goto abortConfigureCmd
				}
				if log1wireErr(d.m.CopyScratchpad()) {
					goto abortConfigureCmd
				}
				c.Resp <- c.Dev
				break
			}
		abortConfigureCmd:
			c.Resp <- onewire.Dev{}
		case TempCmd:
			{
				if c.Dev.Type() == 0 {
					goto abortTempCmd
				}
				if log1wireErr(d.m.MatchROM(c.Dev)) {
					goto abortTempCmd
				}
				if log1wireErr(d.m.ConvertT()) {
					goto abortTempCmd
				}
				//delay.Millisec(200)

				for i := 0; i < 750/50; i++ {
					delay.Millisec(50)
					b, err := d.m.ReadBit()
					if log1wireErr(err) {
						goto abortTempCmd
					}
					if b != 0 {
						break
					}
				}

				if log1wireErr(d.m.MatchROM(c.Dev)) {
					goto abortTempCmd
				}
				s, err := d.m.ReadScratchpad()
				if log1wireErr(err) {
					goto abortTempCmd
				}
				t, err := s.Temp16(c.Dev.Type())
				if log1wireErr(err) {
					goto abortTempCmd
				}
				c.Resp <- t
				break
			}
		abortTempCmd:
			c.Resp <- InvalidTemp
		}
	}
}

var owd OneWireDaemon

func owdUSARTISR() {
	owd.m.DCI.(onedci.USART).Drv.ISR()
}

func owdRxDMAISR() {
	owd.m.DCI.(onedci.USART).Drv.RxDMAISR()
}

func owdTxDMAISR() {
	owd.m.DCI.(onedci.USART).Drv.TxDMAISR()
}
