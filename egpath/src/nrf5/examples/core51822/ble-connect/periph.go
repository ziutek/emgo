package main

import (
	"math/rand"
	"rtos"
	"sync/fence"
	"unsafe"

	"ble"

	"nrf5/hal/ficr"
	"nrf5/hal/ppi"
	"nrf5/hal/radio"
	"nrf5/hal/rtc"
	"nrf5/hal/timer"
)

type Periph struct {
	rnd       rand.XorShift64
	radio     *radio.Periph
	rtc0      *rtc.Periph
	tim0      *timer.Periph
	ev        rtos.EventFlag
	rtc512ms  uint32
	rxPDU     ble.AdvPDU
	txPDU     ble.AdvPDU
	advPeriod uint32
	advFreq   [3]radio.Freq
	advCh     byte
	dirRx     bool
}

func NewPeriph(name string, txpwr int) *Periph {
	p := new(Periph)
	p.radio = radio.RADIO
	p.rtc0 = rtc.RTC0
	p.tim0 = timer.TIMER0
	p.init(name, txpwr)
	return p
}

func getDevAddr() int64 {
	FICR := ficr.FICR
	da0 := FICR.DEVICEADDR[0].Load()
	da1 := FICR.DEVICEADDR[1].Load()
	if FICR.DEVICEADDRTYPE.Load()&1 != 0 {
		da1 |= 0x8000C000
	}
	return int64(da1)<<32 | int64(da0)
}

func (p *Periph) init(name string, txpwr int) {
	r := p.radio
	if f := ficr.FICR; f.BLE_1MBIT_OK().Load() == 0 {
		r.StoreOVERRIDE(0, f.BLE_1MBIT[0].U32.Load())
		r.StoreOVERRIDE(1, f.BLE_1MBIT[1].U32.Load())
		r.StoreOVERRIDE(2, f.BLE_1MBIT[2].U32.Load())
		r.StoreOVERRIDE(3, f.BLE_1MBIT[3].U32.Load())
		r.StoreOVERRIDE(4, f.BLE_1MBIT[4].U32.Load()|1<<31)
	}
	r.StoreMODE(radio.BLE_1M)
	r.StoreTIFS(150)
	r.StorePCNF0(radio.S0_8b | radio.LenBitN(8))
	r.StorePCNF1(radio.WhiteEna | radio.MaxLen(39-2) | radio.BALen(3))
	r.StoreRXADDRESSES(1)
	r.StoreTXADDRESS(0)
	r.StoreCRCCNF(3, true)
	r.StoreCRCPOLY(0x100065B)
	r.StoreTXPOWER(txpwr)
	rtos.IRQ(r.IRQ()).Enable()

	af := &p.advFreq
	af[0] = radio.Channel(2)  // BLE 37, 2402 MHz
	af[1] = radio.Channel(26) // BLE 38, 2426 MHz
	af[2] = radio.Channel(80) // BLE 39, 2480 MHz

	pdu := ble.MakeAdvPDU(nil)
	pdu.SetType(ble.AdvNonconnInd)
	da := getDevAddr()
	pdu.SetTxAdd(da < 0)
	pdu.AppendAddr(da)
	pdu.AppendBytes(ble.Flags, ble.GeneralDisc|ble.OnlyLE)
	pdu.AppendString(ble.LocalName, name)
	pdu.AppendBytes(ble.TxPower, byte(txpwr))
	p.txPDU = pdu
	p.rxPDU = ble.MakeAdvPDU(nil)

	p.rnd.Seed(rtos.Nanosec())
}

func setAA(r *radio.Periph, addr uint32) {
	r.StoreBASE(0, addr<<8)
	r.StorePREFIX(0, addr>>24)
}

func (p *Periph) StartAdvert() {
	r := p.radio
	setAA(r, 0x8E89BED6)
	r.StoreCRCINIT(0x555555)
	r.StorePACKETPTR(unsafe.Pointer(&p.txPDU.Bytes()[0]))
	r.EnableIRQ(1<<radio.READY | 1<<radio.DISABLED)
	r.StoreFREQUENCY(p.advFreq[0])
	r.StoreDATAWHITEIV(37)
	r.StoreSHORTS(radio.READY_START | radio.END_DISABLE)
	p.advCh = 37

	rtc0 := p.rtc0
	rtc0.StoreCC(0, rtc0.LoadCOUNTER())
	rtc0.Event(rtc.COMPARE0).EnablePPI()
	p.rtc512ms = 512 * 32768 / (1e3 * (rtc0.LoadPRESCALER() + 1))

	ppi.RTC0_COMPARE0__RADIO_TXEN.Enable()

	tim0 := p.tim0
	tim0.StorePRESCALER(9) // 31250 Hz
	tim0.StoreCC(1, 256)   // 8.2 ms
	tim0.StoreSHORTS(timer.COMPARE1_STOP)
	ppi.TIMER0_COMPARE1__RADIO_DISABLE.Enable()

	fence.W()

	r.Task(radio.TXEN).Trigger()
}

func (p *Periph) msToRTC(ms int) uint32 {
	return uint32(ms) * p.rtc512ms >> 9
}

func rnd10(rnd *rand.XorShift64) int {
	r := rnd.Uint32()
	return int(r&7 + r>>3&3)
}

func (p *Periph) ISR() {
	r := p.radio
	if ev := r.Event(radio.READY); ev.IsSet() {
		ev.Clear()
		shorts := radio.READY_START | radio.END_DISABLE
		ch := uint32(p.advCh)
		if p.dirRx {
			// Now Rx. Setup Tx.
			leds[3].Set()
			if ch == 39 {
				// Now Rx on channel 39. Setup Tx on channel 37.
				ch = 37
				rtc0 := p.rtc0
				deadline := rtc0.LoadCC(0) + p.msToRTC(625+rnd10(&p.rnd))
				rtc0.StoreCC(0, deadline&0xFFFFFF)
			} else {
				// Setup Tx on next channel.
				ch++
				shorts |= radio.DISABLED_TXEN
				p.tim0.Task(timer.START).Trigger()
			}
			p.advCh = byte(ch)
			r.StoreDATAWHITEIV(ch)
			r.StoreFREQUENCY(p.advFreq[ch-37])
			r.StorePACKETPTR(unsafe.Pointer(&p.txPDU.Bytes()[0]))
			// Setup Rx timeout.
			tim0 := p.tim0
			tim0.Task(timer.CLEAR).Trigger()
			tim0.Task(timer.START).Trigger()
			p.dirRx = false
		} else {
			// Now Tx. Setup Rx on the same channel.
			leds[4].Set()
			shorts |= radio.DISABLED_RXEN
			r.StoreDATAWHITEIV(ch)
			r.StorePACKETPTR(unsafe.Pointer(&p.rxPDU.Bytes()[0]))
			p.dirRx = true
		}
		r.StoreSHORTS(shorts)
	}
	if ev := r.Event(radio.DISABLED); ev.IsSet() {
		ev.Clear()
		if !p.dirRx {
			leds[3].Clear()
		} else {
			leds[4].Clear()
		}
	}
}
