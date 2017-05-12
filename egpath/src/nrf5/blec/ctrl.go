package blec

import (
	"math/rand"
	"rtos"
	"sync/fence"

	"ble"

	"nrf5/hal/gpio"
	"nrf5/hal/ppi"
	"nrf5/hal/radio"
	"nrf5/hal/rtc"
	"nrf5/hal/timer"
)

type Ctrl struct {
	rnd rand.XorShift64

	devAddr int64
	advPDU  ble.AdvPDU
	rxd     []byte
	txd     []byte

	advInterval  uint32
	connInterval uint32
	hop          byte
	chi          byte

	isr func(c *Ctrl)

	radio *radio.Periph
	rtc0  *rtc.Periph
	tim0  *timer.Periph
	ppi   ppi.Chan

	LEDs *[5]gpio.Pin
}

func (c *Ctrl) RxD() []byte {
	return c.rxd
}

func NewCtrl(p ppi.Chan) *Ctrl {
	c := new(Ctrl)
	c.devAddr = getDevAddr()
	c.advPDU = ble.MakeAdvPDU(make([]byte, 2+6+3))
	c.advPDU.SetType(ble.AdvInd)
	c.advPDU.SetTxAdd(c.devAddr < 0)
	c.advPDU.AppendAddr(c.devAddr)
	c.advPDU.AppendBytes(ble.Flags, ble.GeneralDisc|ble.OnlyLE)
	c.rxd = make([]byte, 39)
	c.radio = radio.RADIO
	c.rtc0 = rtc.RTC0
	c.tim0 = timer.TIMER0
	c.ppi = p
	c.rnd.Seed(rtos.Nanosec())
	return c
}

func (c *Ctrl) SetTxPwr(dBm int) {
	c.radio.StoreTXPOWER(dBm)
}

func (c *Ctrl) DevAddr() int64 {
	return c.devAddr
}

// InitPhy initialises physical layer: radio, timers, PPI.
func (c *Ctrl) InitPhy() {
	radioInit(c.radio, 39)
	rtcInit(c.rtc0)
	timerInit(c.tim0)
	pp := ppi.RADIO_BCMATCH__AAR_START.Mask() |
		ppi.RADIO_READY__CCM_KSGEN.Mask() |
		ppi.RADIO_ADDRESS__CCM_CRYPT.Mask() |
		ppi.RADIO_ADDRESS__TIMER0_CAPTURE1.Mask() |
		ppi.RADIO_END__TIMER0_CAPTURE2.Mask() |
		ppi.RTC0_COMPARE0__RADIO_TXEN.Mask() |
		ppi.RTC0_COMPARE0__RADIO_RXEN.Mask() |
		ppi.RTC0_COMPARE0__TIMER0_CLEAR.Mask() |
		ppi.RTC0_COMPARE0__TIMER0_START.Mask() |
		ppi.TIMER0_COMPARE0__RADIO_TXEN.Mask() |
		ppi.TIMER0_COMPARE0__RADIO_RXEN.Mask() |
		ppi.TIMER0_COMPARE1__RADIO_DISABLE.Mask()
	pp.Disable()
}

func (c *Ctrl) Advertise(respPDU ble.AdvPDU, advIntervalms int) {
	r := c.radio
	r.StoreCRCINIT(0x555555)
	r.Event(radio.DISABLED).EnableIRQ()
	radioSetAA(r, 0x8E89BED6)
	radioSetMaxLen(r, 39-2)
	radioSetChi(r, 37)

	c.chi = 37
	c.txd = respPDU.Bytes()
	c.advInterval = uint32(advIntervalms*32768+500) / 1000
	c.scanDisabledISR()

	// TIMER0_COMPARE1 is used to implement Rx timeout.
	t := c.tim0
	t.StoreCC(1, 200) // > t_IFS+t_Preamble+t_AA = 190 µs.
	ppi.TIMER0_COMPARE1__RADIO_DISABLE.Enable()
	p := c.ppi
	p.SetEEP(r.Event(radio.ADDRESS))
	p.SetTEP(t.Task(timer.STOP))
	p.Enable()

	// RTC0_COMPARE0 is used to implement periodical events.
	rt := c.rtc0
	rt.Event(rtc.COMPARE(0)).EnablePPI()
	rt.Task(rtc.CLEAR).Trigger()
	rt.StoreCC(0, 1)
	ppi.RTC0_COMPARE0__RADIO_TXEN.Enable()

	fence.W()
	rt.Task(rtc.START).Trigger()
}

/*pc := c.ppic
pc.SetEEP(t.Event(rtc.COMPARE(1)))
pc.SetTEP(r.Task(radio.DISABLE))
pc.Enable()*/

func (c *Ctrl) RadioISR() {
	r := c.radio
	if ev := r.Event(radio.DISABLED); ev.IsSet() {
		ev.Clear()
		c.isr(c)
		return
	}

	r.Event(radio.PAYLOAD).DisableIRQ()

	pdu := ble.AsAdvPDU(c.rxd)
	if len(pdu.Payload()) < 12 ||
		decodeDevAddr(pdu.Payload()[6:], pdu.RxAdd()) != c.devAddr {
		return
	}
	switch {
	case pdu.Type() == ble.ScanReq && len(pdu.Payload()) == 6+6:
		// Setup for TxScanRsp.
		if c.chi != 37 {
			c.chi--
		} else {
			c.chi = 39
		}
		radioSetChi(r, c.chi)
		r.StoreSHORTS(radio.READY_START | radio.END_DISABLE |
			radio.DISABLED_TXEN)
		c.isr = (*Ctrl).scanReqRxDisabledISR
		c.LEDs[1].Set()
	case pdu.Type() == ble.ConnectReq && len(pdu.Payload()) == 6+6+22:
		// Setup for connection state.
		c.connInterval = c.rtc0.LoadCOUNTER()
		r.StoreSHORTS(radio.READY_START | radio.END_DISABLE)
		c.isr = (*Ctrl).connectReqRxDisabledISR
		c.LEDs[2].Set()
	}
}

// scanDisabledISR handles DISABLED->(TXEN/NOP) between RxTxScan* and TxAdvInd.
func (c *Ctrl) scanDisabledISR() {
	c.LEDs[4].Clear()

	r := c.radio

	// Must be before TxAdvInd.START.
	radioSetPDU(r, c.advPDU.Bytes())

	// Must be before TxAdvInd.PAYLOAD
	r.Event(radio.PAYLOAD).DisableIRQ()

	// Must be before TxAdvInd.DISABLED.
	r.StoreSHORTS(radio.READY_START | radio.END_DISABLE | radio.DISABLED_RXEN)
	c.isr = (*Ctrl).advIndTxDisabledISR

	if c.chi == 39 {
		rt := c.rtc0
		rt.StoreCC(0, rt.LoadCC(0)+c.advInterval+c.rnd.Uint32()&255)
	}
}

// advIndTxDisabledISR handles DISABLED->RXEN between TxAdvInd and RxScanReq.
func (c *Ctrl) advIndTxDisabledISR() {
	c.LEDs[4].Set()

	r := c.radio

	// Must be before RxScanReq.START.
	radioSetPDU(r, c.rxd)

	// Must be before RxScanReq.PAYLOAD
	ev := r.Event(radio.PAYLOAD)
	ev.Clear()
	ev.EnableIRQ()

	// Enable RxScanReq timeout.
	t := c.tim0
	t.Task(timer.CLEAR).Trigger()
	t.Task(timer.START).Trigger()

	c.setupTxAdvInd()
}

// scanReqRxDisabledISR handles DISABLED->TXEN between RxScanReq and TxScanRsp.
func (c *Ctrl) scanReqRxDisabledISR() {
	// Must be before TxScanRsp.START.
	radioSetPDU(c.radio, c.txd)

	// Must be before TxScanRsp.DISABLED.
	c.setupTxAdvInd()

	c.LEDs[1].Clear()
}

func (c *Ctrl) setupTxAdvInd() {
	// Calculate next channel index. Setup shorts and wakeup time.
	shorts := radio.READY_START | radio.END_DISABLE | radio.DISABLED_TXEN
	if c.chi != 39 {
		c.chi++
	} else {
		c.chi = 37
		shorts &^= radio.DISABLED_TXEN
	}

	// Must be before RxTxScan*.DISABLED
	r := c.radio
	r.StoreSHORTS(shorts)
	radioSetChi(r, c.chi)

	c.isr = (*Ctrl).scanDisabledISR
}

func (c *Ctrl) connectReqRxDisabledISR() {
	r := c.radio
	if !r.LoadCRCSTATUS() {
		// Return to advertising.
		r.Task(radio.TXEN).Trigger()
		c.scanDisabledISR()
		c.LEDs[2].Clear()
		return
	}
	ppi.RTC0_COMPARE0__RADIO_TXEN.Disable()

	d := llData(c.rxd[2+6+6:])

	r.StoreCRCINIT(d.CRCInit())
	radioSetAA(r, d.AA())
	radioSetMaxLen(r, 253-2)
	//radioSetChi()

	// Calculate timing parameters as RTC ticks, rounding to safe direction.
	const (
		scale = 41943 // 32768*1.25*1024/1000 = 41943.04 (error ~= 1ppm).
		tcrc  = (24*32768 + 500) / 1000
	)
	winSize := (d.WinSize()*scale + 2<<10 - 1) >> 10
	winOffset := (d.WinOffset() + 1) * scale >> 10
	ssca := d.SSCA() + (100<<19+999999)/1000000
	tsca := (winOffset*ssca + 1<<19 - 1) >> 19
	winOffset -= tsca
	winSize += 2 * tsca

	// c.connInterval contains time of ConnectReq.PAYLOAD event.
	rxstart := c.connInterval + tcrc + winOffset
	rxstop := rxstart + winSize + ((150+2080)*32768+999)/1000
	rt := c.rtc0

	rt.Task(rtc.STOP).Trigger()

	rt.StoreCC(0, rxstart)
	rt.StoreCC(1, rxstop)
	c.connInterval = d.Interval() * scale >> 10

	ppi.RTC0_COMPARE0__RADIO_RXEN.Enable()
	r.StoreSHORTS(radio.READY_START | radio.END_DISABLE | radio.DISABLED_TXEN)
	//c.isr = (*Ctrl).connRxDisabledISR

}
