package main

import (
	"fmt"
	"rtos"
	"strconv"
	"sync/fence"
	"time"
	"time/tz"

	"display/hdc/hdcfb"
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
	{Status: showRoomTempSensor, Action: setEnvTempSensor},
	{Status: showDesiredWaterTemp, Action: setDesiredWaterTemp},
	{Status: showDesiredRoomTemp, Action: setDesiredRoomTemp},
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

type Sensor struct {
	dev onewire.Dev
}

func (s *Sensor) Load() onewire.Dev {
	fence.Compiler()
	return s.dev
}

func (s *Sensor) Store(dev onewire.Dev) {
	fence.Compiler()
	s.dev = dev
}

var menu Menu

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

func printTemp(fbs *hdcfb.SyncSlice, sensor *Sensor) {
	dev := sensor.Load()
	if dev.Type() == 0 {
		fbs.WriteString(" ---- ")
		return
	}
	owd.Cmd <- TempCmd{Dev: dev, Resp: menu.tempResp}
	t := <-menu.tempResp
	if t == InvalidTemp {
		fbs.WriteString(" blad ")
		return
	}
	strconv.WriteFloat(fbs, float64(t)*0.0625, 'f', 1, 32, 4, ' ')
	fbs.WriteString("\xdfC")
}

func showStatus(fbs *hdcfb.SyncSlice) {
	t := time.Now()
	fbs.SetPos(0)
	fmt.Fprintf(
		fbs, "Status      %02d:%02d:%02d",
		t.Hour(), t.Minute(), t.Second(),
	)
	fmt.Fprintf(fbs, " Woda: (%2dkW) ", water.LastPower())
	printTemp(fbs, water.TempSensor())
	fbs.WriteString(" Otoczenie:   ")
	printTemp(fbs, room.TempSensor())
	fbs.Fill(fbs.Remain(), ' ')
	lcdDraw()
}

func showTempSensor(fbs *hdcfb.SyncSlice, name string, sensor *Sensor) {
	tempsensor := "Czujnik temp. "
	fbs.SetPos(0)
	fbs.WriteString(tempsensor)
	fbs.WriteString(name)
	fbs.Fill(41-len(tempsensor)-len(name), ' ')
	dev := sensor.Load()
	if dev.Type() == 0 {
		fbs.WriteString("   nie  wybrano")
	} else {
		fmt.Fprint(fbs, dev)
	}
	fbs.Fill(fbs.Remain(), ' ')
	lcdDraw()
}

func showWaterTempSensor(fbs *hdcfb.SyncSlice) {
	showTempSensor(fbs, "wody", water.TempSensor())
}

func showRoomTempSensor(fbs *hdcfb.SyncSlice) {
	showTempSensor(fbs, "otocz.", room.TempSensor())
}

func setTempSensor(fbs *hdcfb.SyncSlice, sensor *Sensor) {
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
	d := sensor.Load()
	for {
		fbs.SetPos(0)
		for i := 0; i < n; i++ {
			dev := devs[i]
			if sel == -1 && dev == d {
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
			if d == dev {
				break
			}

			// BUG: Not atomic operation. Works, because reading from 1-wire
			// sensor is generally unrealiable so code that uses it muse be
			// fault-tolerant.
			sensor.Store(dev)

			owd.Cmd <- ConfigureCmd{
				Dev: dev, Cfg: onewire.T10bit, Resp: menu.devResp,
			}
			dev = <-menu.devResp
			if dev.Type() == 0 {
				sensor.Store(onewire.Dev{})
				sel = -1
			}
			break
		}
		sel = es.ModCnt(n)
	}
}

func setWaterTempSensor(fbs *hdcfb.SyncSlice) {
	setTempSensor(fbs, water.TempSensor())
}

func setEnvTempSensor(fbs *hdcfb.SyncSlice) {
	setTempSensor(fbs, room.TempSensor())
}

func showDesiredTemp(fbs *hdcfb.SyncSlice, s string, desiredTemp16 int) {
	fbs.SetPos(0)
	fbs.WriteString("Zadana temp. ")
	fbs.WriteString(s)
	fbs.Fill(35-len(s), ' ')
	fmt.Fprintf(fbs, "%2.1f\xdfC", float32(desiredTemp16)/16)
	fbs.Fill(fbs.Remain(), ' ')
	lcdDraw()
}

func showDesiredWaterTemp(fbs *hdcfb.SyncSlice) {
	showDesiredTemp(fbs, "wody", water.DesiredTemp16())
}

func showDesiredRoomTemp(fbs *hdcfb.SyncSlice) {
	showDesiredTemp(fbs, "otocz.", room.DesiredTemp16())
}

func setDesiredTemp(fbs *hdcfb.SyncSlice, set func(int), show func(*hdcfb.SyncSlice), cur16, min16, max16 int) {
	cur2, min2, max2 := cur16/8, min16/8, max16/8 // °C/2
	encoder.SetCnt(cur2)
	for es := range encoder.State {
		if btnPreRel(es) {
			break
		}
		temp2 := es.Cnt()
		if temp2 < min2 {
			temp2 = min2
			encoder.SetCnt(min2)
		} else if temp2 > max2 {
			temp2 = max2
			encoder.SetCnt(max2)
		}
		set(temp2 * 8)
		show(fbs)
	}
}

func setDesiredWaterTemp(fbs *hdcfb.SyncSlice) {
	setDesiredTemp(
		fbs, water.SetDesiredTemp16, showDesiredWaterTemp,
		water.DesiredTemp16(), 30*16, 60*16,
	)
}

func setDesiredRoomTemp(fbs *hdcfb.SyncSlice) {
	setDesiredTemp(
		fbs, room.SetDesiredTemp16, showDesiredRoomTemp,
		room.DesiredTemp16(), 10*16, 26*16,
	)
}

//emgo:const
var wdaysPL = [7]string{
	"Niedziela",    // Sunday
	"Poniedzialek", // Monday
	"Wtorek",       // Tuesday
	"Sroda",        // Wednesday
	"Czwartek",     // Thursday
	"Piatek",       // Friday
	"Sobota",       // Saturday
}

func printDateTime(fbs *hdcfb.SyncSlice, t time.Time) {
	fbs.SetPos(0)
	zone, _ := t.Zone()
	fmt.Fprintf(fbs, "Data i czas     %4s", zone)
	fbs.Fill(20, ' ')
	y, mo, d := t.Date()
	h, mi, s := t.Clock()
	fmt.Fprintf(fbs, "%04d-%02d-%02d", y, mo, d)
	fmt.Fprintf(fbs, "  %02d:%02d:%02d", h, mi, s)
	wd := wdaysPL[t.Weekday()]
	fbs.Fill(10-len(wd)/2, ' ')
	fbs.WriteString(wd)
	fbs.Fill(fbs.Remain(), ' ')
	lcdDraw()
}

func showDateTime(fbs *hdcfb.SyncSlice) {
	printDateTime(fbs, time.Now())
}

type dateTime struct {
	year, month, day int
	hour, min, sec   int
}

func (dt *dateTime) time() time.Time {
	return time.Date(
		dt.year, time.Month(dt.month), dt.day,
		dt.hour, dt.min, dt.sec, 0,
		&tz.EuropeWarsaw,
	)
}

func updateDateTime(fbs *hdcfb.SyncSlice, dt *dateTime, field *int, offs, mod int) {
	encoder.SetCnt(*field - offs)
	for es := range encoder.State {
		if btnPreRel(es) {
			break
		}
		*field = offs + es.ModCnt(mod)
		printDateTime(fbs, dt.time())
	}
}

func setDateTime(fbs *hdcfb.SyncSlice) {
	var dt dateTime
	t := time.Now()
	dt.year = t.Year()
	if dt.year > 2000 {
		dt.month = int(t.Month())
		dt.day = t.Day()
		dt.hour, dt.min, dt.sec = t.Clock()
	} else {
		dt.year = 2000
		dt.month = 1
		dt.day = 1
	}
	updateDateTime(fbs, &dt, &dt.year, 0, 3000)
	updateDateTime(fbs, &dt, &dt.month, 1, 12)
	updateDateTime(fbs, &dt, &dt.day, 1, 31)
	printDateTime(fbs, dt.time())
	updateDateTime(fbs, &dt, &dt.hour, 0, 24)
	updateDateTime(fbs, &dt, &dt.min, 0, 60)
	updateDateTime(fbs, &dt, &dt.sec, 0, 60)
	time.Set(dt.time(), rtos.Nanosec())
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
