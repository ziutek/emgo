package main

import (
	"delay"
	"fmt"

	"stm32/hal/raw/pwr"
	"stm32/hal/raw/rcc"
	"stm32/hal/raw/rtc"
)

func initRTC() {
	const (
		mask = rcc.LSEON | rcc.RTCSEL | rcc.RTCEN
		cfg  = rcc.LSEON | rcc.RTCSEL_LSE | rcc.RTCEN
	)
	RCC := rcc.RCC
	PWR := pwr.PWR
	if RCC.CSR.Bits(mask) != cfg {
		RCC.PWREN().Set()
		RCC.PWREN().Load()
		PWR.DBP().Set()
		RCC.RTCRST().Set()
		RCC.RTCRST().Clear()
		RCC.CSR.StoreBits(mask, cfg)
		for RCC.LSERDY().Load() == 0 {
		}
		PWR.DBP().Clear()
		RCC.PWREN().Clear()
		fmt.Println("Done.")
	}
}

func setRTC(t DateTime) {
	RCC := rcc.RCC
	RTC := rtc.RTC
	PWR := pwr.PWR
	RCC.PWREN().Set()
	RCC.PWREN().Load()
	PWR.DBP().Set()
	RTC.WPR.Store(0xca)
	RTC.WPR.Store(0x53)
	RTC.INIT().Set()
	for RTC.INITF().Load() == 0 {
	}
	//const prer = (x-1)<<16 + (y-1)
	//RTC.PRER.Store(prer)
	//RTC.PRER.Store(prer)
	RTC.DR.U32.Store(t.dr)
	RTC.TR.U32.Store(t.tr)
	RTC.INIT().Clear()
	RTC.WPR.Store(0xff)
	PWR.DBP().Clear()
	RCC.PWREN().Clear()
	for RTC.RSF().Load() == 0 {
		delay.Millisec(50)
	}
}

type Weekday byte

const (
	Monday Weekday = iota + 1
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

//emgo:const
var wdayStr = [8]string{
	"--",
	"Poniedzialek", // Monday
	"Wtorek",       // Tuesday
	"Sroda",        // Wednesday
	"Czwartek",     // Thursday
	"Piatek",       // Friday
	"Sobota",       // Saturday
	"Niedziela",    // Sunday
}

func (wd Weekday) String() string {
	return wdayStr[wd]
}

type DateTime struct {
	tr uint32
	dr uint32
}

func readRTC() DateTime {
	var t DateTime
	if rtc.RTC.INITS().Load() == 0 {
		return t
	}
	RTC := rtc.RTC
	t.dr = RTC.DR.U32.Load()
	for {
		t.tr = RTC.TR.U32.Load()
		dr := RTC.DR.U32.Load()
		if dr == t.dr {
			return t
		}
		t.dr = dr
	}
}

func (t DateTime) IsValid() bool {
	return t.dr&0xff0000 != 0 // Year 0 means: RTC is not set.
}

func (t DateTime) Year() int {
	return int(t.dr>>20&0xf*10 + t.dr>>16&0xf)
}

func (t *DateTime) SetYear(d int) {
	bcd := uint32(d/10<<4 + d%10)
	t.dr = t.dr&^(0xff<<16) | bcd<<16
}

func (t DateTime) Month() int {
	return int(t.dr>>12&1*10 + t.dr>>8&0xf)
}

func (t *DateTime) SetMonth(d int) {
	bcd := uint32(d/10<<4 + d%10)
	t.dr = t.dr&^(0x1f<<8) | bcd<<8
}

func (t DateTime) Day() int {
	return int(t.dr>>4&3*10 + t.dr&0xf)
}

func (t *DateTime) SetDay(d int) {
	bcd := uint32(d/10<<4 + d%10)
	t.dr = t.dr&^0x3f | bcd
}

func (t DateTime) Weekday() Weekday {
	return Weekday(t.dr >> 13 & 7)
}

func (t *DateTime) SetWeekday(d Weekday) {
	t.dr = t.dr&^(7<<13) | uint32(d)&7<<13
}

func (t DateTime) Hour() int {
	return int(t.tr>>20&3*10 + t.tr>>16&0xf)
}

func (t *DateTime) SetHour(d int) {
	bcd := uint32(d/10<<4 + d%10)
	t.tr = t.tr&^(0x3f<<16) | bcd<<16
}

func (t DateTime) Minute() int {
	return int(t.tr>>12&7*10 + t.tr>>8&0xf)
}

func (t *DateTime) SetMinute(d int) {
	bcd := uint32(d/10<<4 + d%10)
	t.tr = t.tr&^(0x7f<<8) | bcd<<8
}

func (t DateTime) Second() int {
	return int(t.tr>>4&7*10 + t.tr&0xf)
}

func (t *DateTime) SetSecond(d int) {
	bcd := uint32(d/10<<4 + d%10)
	t.tr = t.tr&^0x7f | bcd
}

func makeDateTime(Y, M, D, h, m, s int, wd Weekday) (t DateTime) {
	Y -= 2000
	t.dr = uint32(Y/10<<20 + Y%10<<16 + M/10<<12 + M%10<<8 + D/10<<4 + D%10 + int(wd)<<13)
	t.tr = uint32(h/10<<20 + h%10<<16 + m/10<<12 + m%10<<8 + s/10<<4 + s%10)
	return
}

func dayofweek(y, m, d int) Weekday {
	if y <= 1752 || m < 1 || m > 12 || d < 1 || d > 31 {
		return 0
	}
	t := [...]byte{0, 3, 2, 5, 0, 3, 5, 1, 4, 6, 2, 4}
	if m < 3 {
		y--
	}
	w := (y + y/4 - y/100 + y/400 + int(t[m-1]) + d) % 7
	if w == 0 {
		w = 7
	}
	return Weekday(w)
}
