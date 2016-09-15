package main

import (
	"fmt"

	"hdc/hdcfb"
	"onewire"
	"sync/atomic"

	"stm32/hal/raw/tim"
)

type MenuItem struct {
	Name   string
	Action func(fbs *hdcfb.Slice)
	Period int // ms
}

//emgo:const
var menuItems = [...]MenuItem{
	{Name: "Status", Action: printStatus, Period: 1000},
	{Name: "Data i czas"},
	{Name: "Czujnik temp. wody", Action: setWaterTempSensor},
	{Name: "Czujnik temp. otocz.", Action: setEnvTempSensor},
}

type Menu struct {
	curItem int
	timer   *tim.TIM_Periph
	timeout chan struct{}
}

var (
	menu            Menu
	waterTempSensor onewire.Dev
	envTempSensor   onewire.Dev
	searchResp      = make(chan onewire.Dev, 1)
)

func (m *Menu) Setup(t *tim.TIM_Periph, pclk uint) {
	m.timeout = make(chan struct{}, 1)
	m.timer = t
	t.PSC.U16.Store(uint16(pclk/1000 - 1)) // 1 ms
	t.CR1.StoreBits(
		tim.OPM|tim.URS,
		1<<tim.OPMn|1<<tim.URSn,
	)
	t.DIER.Store(tim.UIE)
}

func (m *Menu) setTimeout(ms int) {
	t := m.timer
	t.ARR.U32.Store(uint32(ms))
	t.CEN().Set()
}

func (m *Menu) Loop() {
	fbs := lcd.NewSlice(0, 80)
	var cmi int
	for {
		item := menuItems[cmi]
		fbs.WriteString(item.Name)
		fbs.Fill(20-len(item.Name), ' ')
		var es EncState
		if item.Period > 0 {
		loop:
			for {
				item.Action(fbs)
				m.setTimeout(item.Period)
				select {
				case es = <-encoder.State:
					break loop
				case <-m.timeout:
				}
			}
		} else {
			fbs.Fill(fbs.Remain(), ' ')
			fbs.Flush(0)
			es = <-encoder.State
			if es.Btn {
				for es.Btn {
					es = <-encoder.State
				}
				item.Action(fbs)
				encoder.SetCnt(int16(m.curItem))
			}
		}
		milen := len(menuItems)
		cmi = (milen + int(es.Cnt)%milen) % milen
		atomic.StoreInt(&m.curItem, cmi)
	}
}

func menuISR() {
	menu.timer.SR.Store(0)
	select {
	case menu.timeout <- struct{}{}:
	default:
	}
}

func printStatus(fbs *hdcfb.Slice) {
	waterT := readTemp(waterTempSensor)
	envT := readTemp(envTempSensor)
	dt := readRTC()
	fbs.SetPos(12)
	fmt.Fprintf(
		fbs, "%02d:%02d:%02d",
		dt.Hour(), dt.Minute(), dt.Second(),
	)
	fbs.Fill(20, ' ')
	fmt.Fprintf(fbs, " Woda:      %4.1f\xdfC  ", waterT)
	fmt.Fprintf(fbs, " Otoczenie: %4.1f\xdfC  ", envT)
	fbs.Flush(0)
}

func setTempSensor(fbs *hdcfb.Slice, d *onewire.Dev) {
	owd.Cmd <- SearchCmd{Typ: onewire.DS18B20, Resp: searchResp}
	var (
		devs [3]onewire.Dev
		i    int
	)
	for dev := range searchResp {
		if dev.Type() == 0 {
			break
		}
		devs[i] = dev
		i++
	}
	fbs.SetPos(20)
	for _, dev := range devs {
		if dev.Type() == 0 {
			break
		}
		var c byte = ' '
		if dev == *d {
			c = 0x7e // '->'
		}
		fmt.Fprintf(fbs, "%c%v", c, dev)
	}
	fbs.Flush(0)
	<-encoder.State
}

func setWaterTempSensor(fbs *hdcfb.Slice) {
	setTempSensor(fbs, &waterTempSensor)
}

func setEnvTempSensor(fbs *hdcfb.Slice) {
	setTempSensor(fbs, &envTempSensor)
}
