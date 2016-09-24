package main

import (
	"fmt"
	"math"
	"rtos"
	"strconv"
	"sync/atomic"

	"hdc/hdcfb"
	"onewire"

	"stm32/hal/raw/tim"
)

type MenuItem struct {
	Status func(fbs *hdcfb.Slice) // Should not call fbs.Flush.
	Period int                    // Status refersh in ms (0 means show once).
	Action func(fbs *hdcfb.Slice) // Must call fbs.Flush(0) before return.
}

//emgo:const
var menuItems = [...]MenuItem{
	{Status: showStatus, Period: 1000},
	{Status: showWaterTempSensor, Action: setWaterTempSensor},
	{Status: showEnvTempSensor, Action: setEnvTempSensor},
	{Status: showDesiredWaterTemp, Action: setDesiredWaterTemp},
}

type Menu struct {
	curItem int
	timer   *tim.TIM_Periph
	timeout chan struct{}
}

var (
	menu             Menu
	envTempSensor    onewire.Dev
	waterTempSensor  onewire.Dev
	devResp          = make(chan onewire.Dev, 1)
	desiredWaterTemp = 450 // d°C
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
	if ms < 10 {
		ms = 10
	}
	t := m.timer
	t.ARR.U32.Store(uint32(ms))
	t.CEN().Set()
}

func (m *Menu) Loop() {
	fbs := lcd.NewSlice(0, 80)
	for {
		item := menuItems[m.curItem]
		var (
			es   EncState
			next int64
		)
		if item.Period > 0 {
			next = rtos.Nanosec()
		}
	status:
		for {
			item.Status(fbs)
			fbs.Flush(0)
			if item.Period > 0 {
				next += int64(item.Period * 1e6)
				m.setTimeout(int((next - rtos.Nanosec() + 0.5e6) / 1e6))
			}
			select {
			case es = <-encoder.State:
				break status
			case <-m.timeout:
			}
		}
		if es.Btn() && item.Action != nil {
			for es.Btn() {
				es = <-encoder.State
			}
			item.Action(fbs)
			encoder.SetCnt(m.curItem)
		}
		m.curItem = es.ModCnt(len(menuItems))
	}
}

func menuISR() {
	menu.timer.SR.Store(0)
	select {
	case menu.timeout <- struct{}{}:
	default:
	}
}

func printTemp(fbs *hdcfb.Slice, d onewire.Dev) {
	if d.Type() == 0 {
		fbs.WriteString(" ---- ")
		return
	}
	owd.Cmd <- TempCmd{Dev: d, Resp: tempResp}
	t := <-tempResp
	if t == math.MaxFloat32 {
		fbs.WriteString(" blad ")
		return
	}
	strconv.WriteFloat(fbs, float64(t), 'f', 4, 1, 32)
	fbs.WriteString("\xdfC")
}

func showStatus(fbs *hdcfb.Slice) {
	dt := readRTC()
	fmt.Fprintf(
		fbs, "Status      %02d:%02d:%02d",
		dt.Hour(), dt.Minute(), dt.Second(),
	)
	fbs.WriteString(" Woda:      ")
	printTemp(fbs, waterTempSensor)
	fbs.WriteString("   Otoczenie: ")
	printTemp(fbs, envTempSensor)
	fbs.Fill(fbs.Remain(), ' ')
}

func showTempSensor(fbs *hdcfb.Slice, name string, d *onewire.Dev) {
	tempsensor := "Czujnik temp. "
	fbs.WriteString(tempsensor)
	fbs.WriteString(name)
	fbs.Fill(41-len(tempsensor)-len(name), ' ')
	if d.Type() == 0 {
		fbs.WriteString("   nie  wybrano")
	} else {
		fmt.Fprint(fbs, *d)
	}
	fbs.Fill(fbs.Remain(), ' ')
}

func showWaterTempSensor(fbs *hdcfb.Slice) {
	showTempSensor(fbs, "wody", &waterTempSensor)
}

func showEnvTempSensor(fbs *hdcfb.Slice) {
	showTempSensor(fbs, "otocz.", &envTempSensor)
}

func setTempSensor(fbs *hdcfb.Slice, d *onewire.Dev) {
	owd.Cmd <- SearchCmd{Typ: onewire.DS18B20, Resp: devResp}
	var (
		devs [4]onewire.Dev
		n    int
	)
	for dev := range devResp {
		if dev.Type() == 0 {
			break
		}
		if n < len(devs) {
			devs[n] = dev
			n++
		}
	}
	sel := -1
	for {
		for i := 0; i < n; i++ {
			dev := devs[i]
			if sel == -1 && dev == *d {
				sel = i
				encoder.SetCnt(i)
			}
			var c byte = ' '
			if sel == i {
				sel = i
				c = 0x7e // '->'
			}
			fmt.Fprintf(fbs, "%c%v", c, dev)
		}
		fbs.Fill(fbs.Remain(), ' ')
		fbs.Flush(0)
		es := <-encoder.State
		if es.Btn() {
			for es.Btn() {
				es = <-encoder.State
			}
			if sel == -1 {
				break
			}
			dev := devs[sel]
			if *d == dev {
				break
			}

			// BUG: Not atomic operation. Works, because reading from 1-wire
			// sensor is generally unrealiable so code that uses it muse be
			// fault-tolerant.
			*d = dev

			owd.Cmd <- ConfigureCmd{
				Dev: dev, Cfg: onewire.T10bit, Resp: devResp,
			}
			dev = <-devResp
			if dev.Type() == 0 {
				*d = onewire.Dev{}
				sel = -1
			}
			break
		}
		sel = es.ModCnt(n)
	}
}

func setWaterTempSensor(fbs *hdcfb.Slice) {
	setTempSensor(fbs, &waterTempSensor)
}

func setEnvTempSensor(fbs *hdcfb.Slice) {
	setTempSensor(fbs, &envTempSensor)
}

func showDesiredWaterTemp(fbs *hdcfb.Slice) {
	fbs.WriteString("Zadana temp. wody")
	fbs.Fill(30, ' ')
	fmt.Fprintf(fbs, "%2d.%d\xdfC", desiredWaterTemp/10, desiredWaterTemp%10)
	fbs.Fill(fbs.Remain(), ' ')
}

func setDesiredWaterTemp(fbs *hdcfb.Slice) {
	const (
		min = 300 // d°C
		max = 500 // d°C
	)
	encoder.SetCnt(desiredWaterTemp)
	for es := range encoder.State {
		if es.Btn() {
			for es.Btn() {
				es = <-encoder.State
			}
			break
		}
		temp := es.Cnt()
		if temp < min {
			temp = min
			encoder.SetCnt(min)
		} else if temp > max {
			temp = max
			encoder.SetCnt(max)
		}
		atomic.StoreInt(&desiredWaterTemp, temp)
		showDesiredWaterTemp(fbs)
		fbs.Flush(0)
	}
}
