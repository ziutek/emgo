package blec

import (
	"math/rand"
	"rtos"
	"sync/fence"

	"bluetooth/ble"

	"nrf5/hal/gpio"
	"nrf5/hal/ppi"
	"nrf5/hal/radio"
	"nrf5/hal/rtc"
	"nrf5/hal/timer"
)

type Ctrl struct {
	rnd     rand.XorShift64
	devAddr int64

	advIntervalRTC uint32
	connDelay      uint32
	connWindow     uint16
	sn, nesn       byte

	advPDU  ble.AdvPDU
	rxaPDU  ble.AdvPDU
	txcPDU  ble.DataPDU
	rxdPDUs []ble.DataPDU
	rxcRef  ble.DataPDU
	rspRef  ble.PDU
	rxn     byte
	md      bool

	chi byte
	chm chmap

	isr func(c *Ctrl)

	radio *radio.Periph
	rtc0  *rtc.Periph
	tim0  *timer.Periph

	recv chan ble.DataPDU
	send chan ble.DataPDU

	Iter int
	LEDs *[5]gpio.Pin
}

// NewCtrl returns controller that supports payloads of maxpay length (counts
// MIC if used). Use 31 in case of BLE 4.0, 4.1 and any value from 31 to 255 in
// case of BLE 4.2+. nRF51 can not support BLE 4.2+ full length payloads (radio
// hardware limits max. payload length to 253 bytes). Rxcap and txcap sets
// capacities of internal Rx and Tx queues (channels).
func NewCtrl(maxpay, rxcap, txcap int) *Ctrl {
	c := new(Ctrl)
	c.devAddr = getDevAddr()
	c.advPDU = ble.MakeAdvPDU(6 + 3)
	c.advPDU.SetType(ble.AdvInd)
	c.advPDU.SetTxAdd(c.devAddr < 0)
	c.advPDU.AppendAddr(c.devAddr)
	c.advPDU.AppendBytes(ble.Flags, ble.GeneralDisc|ble.OnlyLE)
	c.rxaPDU = ble.MakeAdvPDU(ble.MaxAdvPay)
	c.txcPDU = ble.MakeDataPDU(27)
	c.rxdPDUs = make([]ble.DataPDU, rxcap+2)
	for i := range c.rxdPDUs {
		c.rxdPDUs[i] = ble.MakeDataPDU(maxpay)
	}
	c.recv = make(chan ble.DataPDU, rxcap)
	c.send = make(chan ble.DataPDU, txcap)
	c.radio = radio.RADIO
	c.rtc0 = rtc.RTC0
	c.tim0 = timer.TIMER0
	c.rnd.Seed(rtos.Nanosec())
	return c
}

func (c *Ctrl) Recv() (ble.DataPDU, error) {
	return <-c.recv, nil
}

func (c *Ctrl) Send(pdu ble.DataPDU) error {
	c.send <- pdu
	return nil
}

func (c *Ctrl) SetTxPwr(dBm int) {
	c.radio.StoreTXPOWER(dBm)
}

func (c *Ctrl) DevAddr() int64 {
	return c.devAddr
}

// InitPhy initialises physical layer: radio, timers, PPI.
func (c *Ctrl) InitPhy() {
	radioInit(c.radio)
	rtcInit(c.rtc0)
	timerInit(c.tim0)
	ppm := ppi.RADIO_BCMATCH__AAR_START.Mask() |
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
	ppm.Disable()
}

func (c *Ctrl) Advertise(rspPDU ble.AdvPDU, advIntervalms int) {
	r := c.radio
	r.StoreCRCINIT(0x555555)
	r.Event(radio.DISABLED).EnableIRQ()
	radioSetAA(r, 0x8E89BED6)
	radioSetMaxLen(r, ble.MaxAdvPay)
	radioSetChi(r, 37)

	c.chi = 37
	c.rspRef = rspPDU.PDU
	c.advIntervalRTC = uint32(advIntervalms*32768+500) / 1000
	c.scanDisabledISR()

	// TIMER0_COMPARE1 is used to implement Rx timeout during advertising
	// and periodical events during connection.
	ppi.TIMER0_COMPARE1__RADIO_DISABLE.Enable()
	ppi.RADIO_ADDRESS__TIMER0_CAPTURE1.Enable()
	ppi.RADIO_END__TIMER0_CAPTURE2.Enable()

	// RTC0_COMPARE0 is used to implement periodical Tx events.
	rt := c.rtc0
	rt.Event(rtc.COMPARE(0)).EnablePPI()
	rt.Task(rtc.CLEAR).Trigger()
	rt.StoreCC(0, 1)
	ppi.RTC0_COMPARE0__RADIO_TXEN.Enable()

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

	pdu := c.rxaPDU
	if len(pdu.Payload()) < 12 ||
		decodeDevAddr(pdu.Payload()[6:], pdu.RxAdd()) != c.devAddr {
		return
	}
	switch {
	case pdu.Type() == ble.ScanReq && pdu.PayLen() == 6+6:
		// Setup for TxScanRsp.
		if c.chi != 37 {
			c.chi--
		} else {
			c.chi = 39
		}
		radioSetChi(r, int(c.chi))
		r.StoreSHORTS(radio.READY_START | radio.END_DISABLE |
			radio.DISABLED_TXEN)
		c.isr = (*Ctrl).scanReqRxDisabledISR
		c.LEDs[3].Set()
	case pdu.Type() == ble.ConnectReq && pdu.PayLen() == 6+6+22:
		// Setup for connection state.
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
	radioSetPDU(r, c.advPDU.PDU)

	// Must be before TxAdvInd.PAYLOAD
	r.Event(radio.PAYLOAD).DisableIRQ()

	// Must be before TxAdvInd.DISABLED.
	r.StoreSHORTS(radio.READY_START | radio.END_DISABLE | radio.DISABLED_RXEN)
	c.isr = (*Ctrl).advIndTxDisabledISR

	c.tim0.Task(timer.STOP).Trigger()

	if c.chi == 39 {
		rt := c.rtc0
		rt.StoreCC(0, rt.LoadCC(0)+c.advIntervalRTC+c.rnd.Uint32()&255)
	}
}

// advIndTxDisabledISR handles DISABLED->RXEN between TxAdvInd and RxScanReq.
func (c *Ctrl) advIndTxDisabledISR() {
	c.LEDs[4].Set()

	r := c.radio

	// Must be before RxScanReq.START.
	radioSetPDU(r, c.rxaPDU.PDU)

	// Must be before RxScanReq.PAYLOAD
	ev := r.Event(radio.PAYLOAD)
	ev.Clear()
	ev.EnableIRQ()

	// Enable RxScanReq timeout.
	t := c.tim0
	t.Task(timer.CLEAR).Trigger()
	t.StoreCC(1, 200) // > t_IFS+t_Preamble+t_AA = 190 µs.
	t.Task(timer.START).Trigger()

	c.setupTxAdvInd()
}

// scanReqRxDisabledISR handles DISABLED->TXEN between RxScanReq and TxScanRsp.
func (c *Ctrl) scanReqRxDisabledISR() {
	// Must be before TxScanRsp.START.
	radioSetPDU(c.radio, c.rspRef)

	// Must be before TxScanRsp.DISABLED.
	c.setupTxAdvInd()

	c.LEDs[3].Clear()
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
	radioSetChi(r, int(c.chi))

	c.isr = (*Ctrl).scanDisabledISR
}

// setRxTimers setups RTC0 and Timer0 to trigger RXEN event or timeout DISABLE
// event using base, delay, window in µs.
func (c *Ctrl) setRxTimers(base, delay, window uint32) {
	t := c.tim0
	rt := c.rtc0
	t.Task(timer.CAPTURE(0)).Trigger()
	baseRTC := rt.LoadCOUNTER()
	t.Task(timer.STOP).Trigger()
	t.Task(timer.CLEAR).Trigger()

	delay -= t.LoadCC(0) - base
	rtcTick := delay * 67 / 2048         // 67/2048 < 32768/1e6 == 512/15625
	timTick := delay - rtcTick*15625/512 // µs

	rt.StoreCC(0, baseRTC+rtcTick)
	if timTick < 4 {
		t.StoreCC(1, timTick+window)
		ppm := ppi.RTC0_COMPARE0__RADIO_TXEN.Mask() |
			ppi.TIMER0_COMPARE0__RADIO_RXEN.Mask()
		ppm.Disable()
		ppm = ppi.RTC0_COMPARE0__RADIO_RXEN.Mask() |
			ppi.RTC0_COMPARE0__TIMER0_START.Mask()
		ppm.Enable()
	} else {
		t.StoreCC(0, timTick)
		t.StoreCC(1, timTick+window)
		ppm := ppi.RTC0_COMPARE0__RADIO_TXEN.Mask() |
			ppi.RTC0_COMPARE0__RADIO_RXEN.Mask()
		ppm.Disable()
		ppm = ppi.RTC0_COMPARE0__TIMER0_START.Mask() |
			ppi.TIMER0_COMPARE0__RADIO_RXEN.Mask()
		ppm.Enable()
	}
}

func (c *Ctrl) connectReqRxDisabledISR() {
	const (
		rxRU    = 130 // RADIO Rx ramp up (µs).
		rtcRA   = 30  // RTC read accuracy (µs): read <= real <= read+rtcRA.
		aaDelay = 40  // Delay between start of packet and ADDRESS event (µs).
	)

	r := c.radio
	if !r.LoadCRCSTATUS() {
		// Return to advertising.
		r.Task(radio.TXEN).Trigger()
		c.scanDisabledISR()
		c.LEDs[2].Clear()
		return
	}
	// Both timers (RTC0, TIMER0) are running. TIMER0.CC2 contains time of
	// END event (end of ConnectReq packet).

	d := llData(c.rxaPDU.Payload()[6+6:])
	c.chm = d.ChM()
	rxPDU := c.rxdPDUs[c.rxn].PDU

	r.Event(radio.ADDRESS).Clear()
	r.StoreCRCINIT(d.CRCInit())
	radioSetMaxLen(r, rxPDU.PayLen())
	radioSetAA(r, d.AA())
	radioSetChi(r, c.chm.NextChi())
	radioSetPDU(r, rxPDU)
	r.StoreSHORTS(radio.READY_START | radio.END_DISABLE | radio.DISABLED_TXEN)

	rsca := d.SCA() + (100<<19+999999)/1000000 // Assume 100 ppm local SCA.

	winOffset := d.WinOffset()
	sca := (winOffset*rsca + 1<<19 - 1) >> 19 // Absolute SCA for (µs).

	// Setup first anchor point.
	c.setRxTimers(
		c.tim0.LoadCC(2),
		winOffset-rxRU-sca,
		d.WinSize()+rxRU+2*sca+rtcRA,
	)
	c.isr = (*Ctrl).connRxDisabledISR

	connInterval := d.Interval()
	sca = (connInterval*rsca + 1<<19 - 1) >> 19
	c.connDelay = connInterval - rxRU - sca - aaDelay
	c.connWindow = uint16(rxRU + 2*sca + rtcRA)

	c.txcPDU.SetHeader(ble.L2CAPCont | ble.Header(c.nesn&1<<2|c.sn&1<<3))
	c.txcPDU.SetPayLen(0)
	c.rspRef = c.txcPDU.PDU
}

func (c *Ctrl) connRxDisabledISR() {
	r := c.radio
	if !r.Event(radio.ADDRESS).IsSet() {
		c.tim0.Task(timer.STOP).Trigger()
		c.rtc0.Task(rtc.STOP).Trigger()
		c.LEDs[2].Clear()
		return
	}
	// Both timers (RTC0, TIMER0) are running. TIMER0.CC1 contains time of
	// ADDRESS event (end of address field in data packet).

	/*
		// Test safety margin to START event.
		r.Event(radio.READY).Clear()
	*/

	if r.LoadCRCSTATUS() {
		rxPDU := c.rxdPDUs[c.rxn]
		header := rxPDU.Header()
		c.md = header&ble.MD != 0 // BUG: fix this.
		nesn := byte(header) >> 2 & 1
		sn := byte(header) >> 3 & 1
		c.LEDs[0].Store(int(sn))
		c.LEDs[1].Store(int(c.nesn))
		c.LEDs[2].Store(int(nesn))
		c.LEDs[3].Store(int(c.sn))
		if sn == c.nesn&1 {
			// New PDU received.
			llid := header & ble.LLID
			switch {
			case llid == ble.LLControl:
				// LL Control PDU. Pass it for further processing if the
				// previous one was done.
				if c.rxcRef.IsZero() {
					c.rxcRef = rxPDU
					c.nesn++
				}
			case llid == ble.L2CAPCont && rxPDU.PayLen() == 0:
				// Empty L2CAP PDU.
				c.nesn++
			default:
				// Non-empty L2CAP PDU.
				select {
				case c.recv <- rxPDU:
					if c.rxn++; int(c.rxn) == len(c.rxdPDUs) {
						c.rxn = 0
					}
					c.nesn++
				default:
				}
			}
		}
		if nesn != c.sn&1 {
			// Previous packet ACKed. Can send new one.
			var rspPDU ble.DataPDU
			if !c.rxcRef.IsZero() {
				// Process last controll PDU.
				header = ble.LLControl
				req := c.rxcRef.Payload()
				switch req[0] {
				case llFeatureReq:
					c.txcPDU.SetPayLen(9)
					rsp := c.txcPDU.Payload()
					rsp[0] = llFeatureRsp
					rsp[1] = 0
					rsp[2] = 0
					rsp[3] = 0
					rsp[4] = 0
					rsp[5] = 0
					rsp[6] = 0
					rsp[7] = 0
					rsp[8] = 0
				case llVersionInd:
					c.txcPDU.SetPayLen(6)
					rsp := c.txcPDU.Payload()
					rsp[0] = llVersionInd
					rsp[1] = 6    // BLE version: 6: 4.0, 7: 4.1, 8: 4.2, 9: 5.
					rsp[2] = 0xFF // Company ID (2 octets).
					rsp[3] = 0xFF // Using 0xFFFF: tests / not assigned.
					rsp[4] = 0    // Subversion (2 octets). Unique for each
					rsp[5] = 0    // implementation or revision of controller.
				default:
					c.txcPDU.SetPayLen(2)
					rsp := c.txcPDU.Payload()
					rsp[0] = llUnknownRsp
					rsp[1] = req[0]
					c.LEDs[4].Clear()
					c.recv <- rxPDU
				}
				c.rxcRef = ble.DataPDU{}
				rspPDU = c.txcPDU
			} else {
				// Send data PDU from send queue or empty PDU.
				select {
				case rspPDU = <-c.send:
					header = rspPDU.Header()
				default:
					header = ble.L2CAPCont
					c.txcPDU.SetPayLen(0)
					rspPDU = c.txcPDU
				}
			}
			c.sn++
			rspPDU.SetHeader(header | ble.Header(c.nesn&1<<2|c.sn&1<<3))
			c.rspRef = rspPDU.PDU
		}
	}

	// Must be before ConnTx.START.
	radioSetPDU(r, c.rspRef)

	/*
		// Test safety margin to START event.
		delay.Loop(320)
		if r.Event(radio.READY).IsSet() {
			c.LEDs[0].Set()
		}
	*/

	// Must be before ConnTx.DISABLED
	r.StoreSHORTS(radio.READY_START | radio.END_DISABLE)
	c.isr = (*Ctrl).connTxDisabledISR

	c.setRxTimers(c.tim0.LoadCC(1), c.connDelay, uint32(c.connWindow))
}

func (c *Ctrl) connTxDisabledISR() {
	r := c.radio
	if !c.md {
		radioSetChi(r, c.chm.NextChi())
	}
	radioSetPDU(r, c.rxdPDUs[c.rxn].PDU)
	r.StoreSHORTS(radio.READY_START | radio.END_DISABLE | radio.DISABLED_TXEN)
	r.Event(radio.ADDRESS).Clear()

	c.isr = (*Ctrl).connRxDisabledISR
	c.Iter++
}
