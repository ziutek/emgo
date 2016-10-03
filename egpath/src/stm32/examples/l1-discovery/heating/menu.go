package main

import (
	"fmt"
	"rtos"
	"strconv"

	"hdc/hdcfb"
	"onewire"

	"stm32/hal/raw/tim"
)

type MenuItem struct {
	Status func(fbs *hdcfb.SyncSlice)
	Period int // Status refersh in ms (0 means show once).
	Action func(fbs *hdcfb.SyncSlice)
}

//emgo:const
var menuItems = [...]MenuItem{
	{Status: showStatus, Period: 1000},
	{Status: showWaterTempSensor, Action: setWaterTempSensor},
	{Status: showEnvTempSensor, Action: setEnvTempSensor},
	{Status: showDesiredWaterTemp, Action: setDesiredWaterTemp},
	{Status: showDateTime, Period: 1000, Action: setDateTime},
	{Status: showDisplOffTimeout, Action: setDisplOffTimeout},
}

type Menu struct {
	curItem         int
	timer           *tim.TIM_Periph
	timeout         chan struct{}
	tempResp        chan int
	devResp         chan onewire.Dev
	displOffTimeout int
}

var (
	menu          Menu
	envTempSensor onewire.Dev
)

func (m *Menu) Setup(t *tim.TIM_Periph, pclk uint) {
	m.timeout = make(chan struct{}, 1)
	m.timer = t
	m.displOffTimeout = 20
	m.tempResp = make(chan int, 1)
	m.devResp = make(chan onewire.Dev, 1)

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
	t.EGR.Store(tim.UG)
	t.CEN().Set()
}

func (m *Menu) stopTimeout() {
	m.timer.CEN().Clear()
}

func btnPreRel(es EncoderState) bool {
	if !es.Btn() {
		return false
	}
	for es.Btn() {
		es = <-encoder.State
	}
	return true
}

func (m *Menu) waitEncoder(deadline1, deadline2 int64) (EncoderState, int) {
	d := 1
	if deadline1 <= 0 || deadline2 > 0 && deadline2 < deadline1 {
		deadline1 = deadline2
		d = 2
	}
	if deadline1 > 0 {
		m.setTimeout(int((deadline1 - rtos.Nanosec() + 0.5e6) / 1e6))
	}
	select {
	case es := <-encoder.State:
		m.stopTimeout()
		select {
		case <-m.timeout:
		default:
		}
		return es, 0
	case <-m.timeout:
		return 0, d
	}

}

func (m *Menu) Loop() {
	fbs := lcd.NewSyncSlice(0, 80)
	disp := lcd.Display()
	for {
		item := menuItems[m.curItem]
		now := rtos.Nanosec()
		var displOff, next int64
		if m.displOffTimeout > 0 {
			displOff = now + int64(m.displOffTimeout)*1e9
		}
		if item.Period > 0 {
			next = now + int64(item.Period)*1e6
		}

	printStatus:
		item.Status(fbs)

	waitEncoder:
		es, dn := m.waitEncoder(displOff, next)

		switch dn {
		case 1:
			logLCDErr(disp.ClearAUX()) // Backlight off.
			displOff = -1
			next = -1
			goto waitEncoder
		case 2:
			next += int64(item.Period) * 1e6
			goto printStatus
		}
		if displOff == -1 {
			encoder.SetCnt(m.curItem)
			logLCDErr(disp.SetAUX()) // Backlight on.
			continue
		}
		if item.Action != nil && btnPreRel(es) {
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

func printTemp(fbs *hdcfb.SyncSlice, d onewire.Dev) {
	if d.Type() == 0 {
		fbs.WriteString(" ---- ")
		return
	}
	owd.Cmd <- TempCmd{Dev: d, Resp: menu.tempResp}
	t := <-menu.tempResp
	if t == InvalidTemp {
		fbs.WriteString(" blad ")
		return
	}
	strconv.WriteFloat(fbs, float64(t)*0.0625, 'f', 4, 1, 32)
	fbs.WriteString("\xdfC")
}

func showStatus(fbs *hdcfb.SyncSlice) {
	dt := readRTC()
	fbs.SetPos(0)
	fmt.Fprintf(
		fbs, "Status      %02d:%02d:%02d",
		dt.Hour(), dt.Minute(), dt.Second(),
	)
	fmt.Fprintf(fbs, " Woda: (%2dkW) ", water.LastPower())
	printTemp(fbs, water.TempSensor)
	fbs.WriteString(" Otoczenie:   ")
	printTemp(fbs, envTempSensor)
	fbs.Fill(fbs.Remain(), ' ')
	lcdDraw()
}

func showTempSensor(fbs *hdcfb.SyncSlice, name string, d *onewire.Dev) {
	tempsensor := "Czujnik temp. "
	fbs.SetPos(0)
	fbs.WriteString(tempsensor)
	fbs.WriteString(name)
	fbs.Fill(41-len(tempsensor)-len(name), ' ')
	if d.Type() == 0 {
		fbs.WriteString("   nie  wybrano")
	} else {
		fmt.Fprint(fbs, *d)
	}
	fbs.Fill(fbs.Remain(), ' ')
	lcdDraw()
}

func showWaterTempSensor(fbs *hdcfb.SyncSlice) {
	showTempSensor(fbs, "wody", &water.TempSensor)
}

func showEnvTempSensor(fbs *hdcfb.SyncSlice) {
	showTempSensor(fbs, "otocz.", &envTempSensor)
}

func setTempSensor(fbs *hdcfb.SyncSlice, d *onewire.Dev) {
	owd.Cmd <- SearchCmd{Typ: onewire.DS18B20, Resp: menu.devResp}
	var (
		devs [4]onewire.Dev
		n    int
	)
	for dev := range menu.devResp {
		if dev.Type() == 0 {
			break
		}
		if n < len(devs) {
			devs[n] = dev
			n++
		}
	}
	if n == 0 {
		fbs.SetPos(0)
		fbs.WriteString("   nie znaleziono")
		fbs.Fill(fbs.Remain(), ' ')
		lcdDraw()
		for es := range encoder.State {
			if btnPreRel(es) {
				return
			}
		}
	}
	encoder.SetCnt(0)
	sel := -1
	for {
		fbs.SetPos(0)
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
		lcdDraw()
		es := <-encoder.State
		if btnPreRel(es) {
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
				Dev: dev, Cfg: onewire.T10bit, Resp: menu.devResp,
			}
			dev = <-menu.devResp
			if dev.Type() == 0 {
				*d = onewire.Dev{}
				sel = -1
			}
			break
		}
		sel = es.ModCnt(n)
	}
}

func setWaterTempSensor(fbs *hdcfb.SyncSlice) {
	setTempSensor(fbs, &water.TempSensor)
}

func setEnvTempSensor(fbs *hdcfb.SyncSlice) {
	setTempSensor(fbs, &envTempSensor)
}

func showDesiredWaterTemp(fbs *hdcfb.SyncSlice) {
	fbs.SetPos(0)
	fbs.WriteString("Zadana temp. wody")
	fbs.Fill(30, ' ')
	fmt.Fprintf(fbs, "%2d\xdfC", water.DesiredTemp())
	fbs.Fill(fbs.Remain(), ' ')
	lcdDraw()
}

func setDesiredWaterTemp(fbs *hdcfb.SyncSlice) {
	const (
		min = 30 // °C
		max = 60 // °C
	)
	encoder.SetCnt(water.DesiredTemp())
	for es := range encoder.State {
		if btnPreRel(es) {
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
		water.SetDesiredTemp(temp)
		showDesiredWaterTemp(fbs)
	}
}

func printDateTime(fbs *hdcfb.SyncSlice, dt DateTime) {
	fbs.SetPos(0)
	fbs.WriteString("Data i czas")
	fbs.Fill(29, ' ')
	fmt.Fprintf(fbs, "%04d-%02d-%02d", 2000+dt.Year(), dt.Month(), dt.Day())
	fmt.Fprintf(fbs, "  %02d:%02d:%02d", dt.Hour(), dt.Minute(), dt.Second())
	wd := dt.Weekday().String()
	fbs.Fill(10-len(wd)/2, ' ')
	fbs.WriteString(wd)
	fbs.Fill(fbs.Remain(), ' ')
	lcdDraw()
}

func showDateTime(fbs *hdcfb.SyncSlice) {
	printDateTime(fbs, readRTC())
}

func updateDateTime(fbs *hdcfb.SyncSlice, dt *DateTime, get func(DateTime) int,
	set func(*DateTime, int), offs, mod int) {

	encoder.SetCnt(get(*dt))
	for es := range encoder.State {
		if btnPreRel(es) {
			break
		}
		set(dt, offs+es.ModCnt(mod))
		printDateTime(fbs, *dt)
	}
}

func setDateTime(fbs *hdcfb.SyncSlice) {
	dt := readRTC()
	updateDateTime(fbs, &dt, (DateTime).Year, (*DateTime).SetYear, 0, 100)
	updateDateTime(fbs, &dt, (DateTime).Month, (*DateTime).SetMonth, 1, 12)
	updateDateTime(fbs, &dt, (DateTime).Day, (*DateTime).SetDay, 1, 31)
	dt.SetWeekday(dayofweek(2000+dt.Year(), dt.Month(), dt.Day()))
	printDateTime(fbs, dt)
	updateDateTime(fbs, &dt, (DateTime).Hour, (*DateTime).SetHour, 0, 24)
	updateDateTime(fbs, &dt, (DateTime).Minute, (*DateTime).SetMinute, 0, 60)
	updateDateTime(fbs, &dt, (DateTime).Second, (*DateTime).SetSecond, 0, 60)
	setRTC(dt)
}

func showDisplOffTimeout(fbs *hdcfb.SyncSlice) {
	fbs.SetPos(0)
	fbs.WriteString("Wygaszanie ekranu po")
	fbs.Fill(26, ' ')
	fmt.Fprintf(fbs, "%4d s", menu.displOffTimeout)
	fbs.Fill(fbs.Remain(), ' ')
	lcdDraw()
}

func setDisplOffTimeout(fbs *hdcfb.SyncSlice) {
	encoder.SetCnt(menu.displOffTimeout)
	for es := range encoder.State {
		if btnPreRel(es) {
			break
		}
		cnt := es.Cnt()
		if cnt < 0 {
			cnt = 0
			encoder.SetCnt(cnt)
		} else if cnt > 65 {
			cnt = 65 // == 65000 ms (can not be more than 65535 ms).
			encoder.SetCnt(cnt)
		}
		menu.displOffTimeout = cnt
		showDisplOffTimeout(fbs)
	}
}
