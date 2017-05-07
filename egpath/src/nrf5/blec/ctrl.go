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
)

type Ctrl struct {
	rnd rand.XorShift64

	addr   int64
	advPDU ble.AdvPDU
	rxd    []byte
	txd    []byte

	period uint32
	chi    byte

	isr func(c *Ctrl)

	radio *radio.Periph
	rtc0  *rtc.Periph
	ppic  ppi.Chan

	LEDs *[5]gpio.Pin
}

func (c *Ctrl) RxD() []byte {
	return c.rxd
}

func NewCtrl(ppic ppi.Chan) *Ctrl {
	c := new(Ctrl)
	c.addr = getDevAddr()
	c.advPDU = ble.MakeAdvPDU(make([]byte, 2+6+3))
	c.advPDU.SetType(ble.AdvInd)
	c.advPDU.SetTxAdd(c.addr < 0)
	c.advPDU.AppendAddr(c.addr)
	c.advPDU.AppendBytes(ble.Flags, ble.GeneralDisc|ble.OnlyLE)
	c.rxd = make([]byte, 39)
	c.radio = radio.RADIO
	c.rtc0 = rtc.RTC0
	c.ppic = ppic
	c.rnd.Seed(rtos.Nanosec())
	return c
}

func (c *Ctrl) SetTxPwr(dBm int) {
	c.radio.StoreTXPOWER(dBm)
}

func (c *Ctrl) DevAddr() int64 {
	return c.addr
}

// InitPhy initialises physical layer: radio, timers, PPI.
func (c *Ctrl) InitPhy() {
	radioInit(c.radio, 39)
	rtcInit(c.rtc0)
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

func (c *Ctrl) Advertise(respPDU ble.AdvPDU, periodms int) {
	r := c.radio
	r.StoreCRCINIT(0x555555)
	r.Event(radio.DISABLED).EnableIRQ()
	radioSetAA(r, 0x8E89BED6)
	radioSetAddrMatch(r, c.addr)
	radioSetMaxLen(r, 39-2)
	radioSetChi(r, 37)

	c.chi = 37
	c.txd = respPDU.Bytes()
	c.period = uint32(periodms*32768+500) / 1000
	c.scanDisabledISR()

	t := c.rtc0
	t.Event(rtc.COMPARE(0)).EnablePPI()
	t.Task(rtc.CLEAR).Trigger()
	t.StoreCC(0, 1)
	pc := c.ppic
	pc.SetEEP(t.Event(rtc.COMPARE(1)))
	pc.SetTEP(r.Task(radio.DISABLE))
	pc.Enable()
	ppi.RTC0_COMPARE0__RADIO_TXEN.Enable()

	fence.W()
	t.Task(rtc.START).Trigger()
}

func (c *Ctrl) RadioISR() {
	c.rtc0.Event(rtc.COMPARE(1)).DisablePPI()

	r := c.radio
	if ev := r.Event(radio.DISABLED); ev.IsSet() {
		ev.Clear()
		c.isr(c)
		return
	}

	r.Event(radio.PAYLOAD).DisableIRQ()

	pdu := ble.AsAdvPDU(c.rxd)
	if len(pdu.Payload()) < 12 ||
		decodeDevAddr(pdu.Payload()[6:], pdu.RxAdd()) != c.addr {
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
		r.StoreSHORTS(radio.READY_START | radio.END_DISABLE)
		c.isr = (*Ctrl).connectReqRxDisabledISR
		c.LEDs[2].Set()
	}
}

// scanDisabledISR handles DISABLED/RTC->TXEN between RxTxScan* and TxAdvInd.
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
		t := c.rtc0
		t.StoreCC(0, t.LoadCC(0)+c.period+c.rnd.Uint32()&255)
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

	// Setup for RxScanReq timeout.
	t := c.rtc0
	t.StoreCC(1, t.LoadCOUNTER()+16) // 488 µs
	t.Event(rtc.COMPARE(1)).EnablePPI()
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

	c.rtc0.Task(rtc.STOP).Trigger()
}

/*func (c *Ctrl) WaitConn(deadline int64) {

}*/
