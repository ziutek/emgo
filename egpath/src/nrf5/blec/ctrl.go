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

	LEDs *[5]gpio.Pin
}

func NewCtrl() *Ctrl {
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
	c.rnd.Seed(rtos.Nanosec())
	return c
}

func (c *Ctrl) SetTxPwr(dBm int) {
	c.radio.StoreTXPOWER(dBm)
}

func (c *Ctrl) DevAddr() int64 {
	return c.addr
}

func (c *Ctrl) Chi() int {
	return int(c.chi)
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

	ppi.RTC0_COMPARE0__RADIO_TXEN.Enable()
	rt := c.rtc0
	rt.Event(rtc.COMPARE(0)).EnablePPI()
	rt.Task(rtc.CLEAR).Trigger()
	rt.StoreCC(0, 1)

	fence.W()
	rt.Task(rtc.START).Trigger()
}

func (c *Ctrl) RadioISR() {
	r := c.radio
	if ev := r.Event(radio.DISABLED); ev.IsSet() {
		ev.Clear()
		c.isr(c)
		return
	}

	r.Event(radio.PAYLOAD).DisableIRQ()
	c.rtc0.Event(rtc.COMPARE(1)).DisableIRQ()

	pdu := ble.AsAdvPDU(c.rxd)
	if pdu.Type() != ble.ScanReq ||
		decodeDevAddr(pdu.Payload()[6:], pdu.RxAdd()) != c.addr {

		c.setupTxAdvInd()
		return
	}

	// Setup for TxScanRsp.
	r.StoreSHORTS(radio.READY_START | radio.END_DISABLE | radio.DISABLED_TXEN)
	c.isr = (*Ctrl).scanReqRxDisabledISR
	c.LEDs[1].Set()
}

func (c *Ctrl) RTCISR() {
	rt := c.rtc0
	rt.Event(rtc.COMPARE(1)).DisableIRQ()
	rt.Event(rtc.COMPARE(1)).Clear()
	c.setupTxAdvInd()
	r := c.radio
	r.Task(radio.DISABLE).Trigger()
	r.Event(radio.PAYLOAD).DisableIRQ()
}

// scanDisabledISR handles DISABLED/RTC->TXEN between RxTxScan* and TxAdvInd.
func (c *Ctrl) scanDisabledISR() {
	r := c.radio

	// Must be before TxAdvInd.START.
	radioSetPDU(r, c.advPDU.Bytes())

	// Must be before TxAdvInd.DISABLED.
	r.StoreSHORTS(radio.READY_START | radio.END_DISABLE | radio.DISABLED_RXEN)
	c.isr = (*Ctrl).advIndTxDisabledISR
}

// advIndTxDisabledISR handles DISABLED->RXEN between TxAdvInd and RxScanReq.
func (c *Ctrl) advIndTxDisabledISR() {
	r := c.radio

	// Must be before RxScanReq.START.
	radioSetPDU(r, c.rxd)

	// Must be before RxScanReq.PAYLOAD
	ev := r.Event(radio.PAYLOAD)
	ev.Clear()
	ev.EnableIRQ()

	// Setup for RxScanReq timeout.
	const timeout = 4 // ms
	rt := c.rtc0
	rt.StoreCC(1, rt.LoadCOUNTER()+(timeout*32768+999)/1000)
	rt.Event(rtc.COMPARE(1)).EnableIRQ()
	c.isr = (*Ctrl).invalidISR
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
	switch c.chi {
	case 37:
		c.LEDs[4].Set()
		c.chi = 38
	case 38:
		c.chi = 39
	default:
		c.chi = 37
		shorts &^= radio.DISABLED_TXEN
		rt := c.rtc0
		rt.StoreCC(0, rt.LoadCC(0)+c.period+c.rnd.Uint32()&255)
		c.LEDs[4].Clear()
	}
	r := c.radio

	// Must be before RxTxScan*.DISABLED
	r.StoreSHORTS(shorts)
	radioSetChi(r, c.chi)

	c.isr = (*Ctrl).scanDisabledISR
}

func (c *Ctrl) setupTxScanRsp() {

}

func (c *Ctrl) invalidISR() {
	for {
		_ = c
	}
}

/*func (c *Ctrl) WaitConn(deadline int64) {

}*/
