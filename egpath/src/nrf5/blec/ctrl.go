package blec

import (
	"encoding/binary/le"
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

type Controller struct {
	rnd     rand.XorShift64
	devAddr int64

	advIntervalRTC uint32
	connInterval   uint32
	winSize        uint32
	sca            fixed19

	advChi   byte
	flags    ble.Header // NESN, SN, MD
	noReq    uint16
	noReqMax uint16
	chm      chmap
	connUpd  connUpdate

	advPDU ble.AdvPDU
	rxaPDU ble.AdvPDU
	txcPDU ble.DataPDU
	rxcRef ble.DataPDU
	rspRef ble.PDU

	recv pduChan
	send pduChan

	isr func(c *Controller)

	radio *radio.Periph
	rtc0  *rtc.Periph
	tim0  *timer.Periph

	LEDs *[5]gpio.Pin
}

// NewController returns controller that supports payloads of maxpay length
// (counts MIC if used). Use 31 in case of BLE 4.0, 4.1 and any value from 31 to
// 255 in case of BLE 4.2+. nRF51 can not support BLE 4.2+ full length payloads
// (radio hardware limits max. payload length to 253 bytes). Rxcap and txcap
// sets capacities of internal Rx and Tx queues (channels).
func NewController(maxpay, rxcap, txcap int) *Controller {
	c := new(Controller)
	c.devAddr = getDevAddr()
	c.advPDU = ble.MakeAdvPDU(6 + 3)
	c.advPDU.SetType(ble.AdvInd)
	c.advPDU.SetTxAdd(c.devAddr < 0)
	c.advPDU.AppendAddr(c.devAddr)
	c.advPDU.AppendBytes(ble.Flags, ble.GeneralDisc|ble.OnlyLE)
	c.rxaPDU = ble.MakeAdvPDU(ble.MaxAdvPay)
	c.txcPDU = ble.MakeDataPDU(27)
	c.recv = makePDUChan(maxpay, rxcap)
	c.send = makePDUChan(maxpay, txcap)
	c.radio = radio.RADIO
	c.rtc0 = rtc.RTC0
	c.tim0 = timer.TIMER0
	c.rnd.Seed(rtos.Nanosec())
	return c
}

func (c *Controller) Recv() (ble.DataPDU, error) {
	return <-c.recv.Ch, nil
}

func (c *Controller) GetSend() ble.DataPDU {
	return c.send.Get()
}

func (c *Controller) Send() error {
	c.send.Ch <- c.send.Get()
	c.send.Next()
	return nil
}

func (c *Controller) SetTxPwr(dBm int) {
	c.radio.StoreTXPOWER(dBm)
}

func (c *Controller) DevAddr() int64 {
	return c.devAddr
}

func (c *Controller) ConnEventCnt() uint16 {
	return c.chm.ConnEventCnt()
}

// InitPhy initialises physical layer: radio, timers, PPI.
func (c *Controller) InitPhy() {
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

func (c *Controller) Advertise(rspPDU ble.AdvPDU, advIntervalms int) {
	r := c.radio
	r.StoreCRCINIT(0x555555)
	r.Event(radio.DISABLED).EnableIRQ()
	radioSetAA(r, 0x8E89BED6)
	radioSetMaxLen(r, ble.MaxAdvPay)
	radioSetChi(r, 37)

	c.advChi = 37
	c.rspRef = rspPDU.PDU
	c.advIntervalRTC = uint32(advIntervalms*32768+500) / 1000
	c.scanDisabledISR()

	// TIMER0_COMPARE1 is used to implement Rx timeout. ADDRESS event disables
	// timeout (sets timeout to >1h).
	ppm := ppi.TIMER0_COMPARE1__RADIO_DISABLE.Mask() |
		ppi.RADIO_ADDRESS__TIMER0_CAPTURE1.Mask() |
		ppi.RADIO_END__TIMER0_CAPTURE2.Mask()
	ppm.Enable()

	/*p := ppi.Chan(9)
	p.SetEEP(r.Event(radio.READY))
	p.SetTEP(c.tim0.Task(timer.CAPTURE(3)))
	p.Enable()*/

	// RTC0_COMPARE0 is used to implement periodical Tx events.
	rt := c.rtc0
	rt.Event(rtc.COMPARE(0)).EnablePPI()
	rt.Task(rtc.CLEAR).Trigger()
	rt.StoreCC(0, 1)
	ppi.RTC0_COMPARE0__RADIO_TXEN.Enable()

	fence.W()
	rt.Task(rtc.START).Trigger()
}

func (c *Controller) RadioISR() {
	r := c.radio
	if ev := r.Event(radio.DISABLED); ev.IRQEnabled() && ev.IsSet() {
		ev.Clear()
		c.isr(c)
		return
	}

	if ev := r.Event(radio.ADDRESS); ev.IRQEnabled() && ev.IsSet() {
		// Receiving PDU in connection event.
		ev.DisableIRQ()
		r.StoreSHORTS(radio.READY_START | radio.END_DISABLE |
			radio.DISABLED_TXEN)
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
		if c.advChi != 37 {
			c.advChi--
		} else {
			c.advChi = 39
		}
		radioSetChi(r, int(c.advChi))
		r.StoreSHORTS(radio.READY_START | radio.END_DISABLE |
			radio.DISABLED_TXEN)
		c.isr = (*Controller).scanReqRxDisabledISR
		c.LEDs[3].Set()
	case pdu.Type() == ble.ConnectReq && pdu.PayLen() == 6+6+22:
		// Setup for connection state.
		r.StoreSHORTS(radio.READY_START | radio.END_DISABLE)
		c.isr = (*Controller).connectReqRxDisabledISR
		c.LEDs[2].Set()
	}
}

// scanDisabledISR handles DISABLED->(TXEN/NOP) between RxTxScan* and TxAdvInd.
func (c *Controller) scanDisabledISR() {
	c.LEDs[0].Clear()

	r := c.radio

	// Must be before TxAdvInd.START.
	radioSetPDU(r, c.advPDU.PDU)

	// Must be before TxAdvInd.PAYLOAD
	r.Event(radio.PAYLOAD).DisableIRQ()

	// Must be before TxAdvInd.DISABLED.
	r.StoreSHORTS(radio.READY_START | radio.END_DISABLE | radio.DISABLED_RXEN)
	c.isr = (*Controller).advIndTxDisabledISR

	c.tim0.Task(timer.STOP).Trigger()

	if c.advChi == 39 {
		rt := c.rtc0
		rt.StoreCC(0, rt.LoadCC(0)+c.advIntervalRTC+c.rnd.Uint32()&255)
	}
}

// advIndTxDisabledISR handles DISABLED->RXEN between TxAdvInd and RxScanReq.
func (c *Controller) advIndTxDisabledISR() {
	c.LEDs[0].Set()

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
func (c *Controller) scanReqRxDisabledISR() {
	// Must be before TxScanRsp.START.
	radioSetPDU(c.radio, c.rspRef)

	// Must be before TxScanRsp.DISABLED.
	c.setupTxAdvInd()

	c.LEDs[3].Clear()
}

func (c *Controller) setupTxAdvInd() {
	// Calculate next channel index. Setup shorts and wakeup time.
	shorts := radio.READY_START | radio.END_DISABLE | radio.DISABLED_TXEN
	if c.advChi != 39 {
		c.advChi++
	} else {
		c.advChi = 37
		shorts &^= radio.DISABLED_TXEN
	}

	// Must be before RxTxScan*.DISABLED
	r := c.radio
	r.StoreSHORTS(shorts)
	radioSetChi(r, int(c.advChi))

	c.isr = (*Controller).scanDisabledISR
}

const (
	rxRU    = 138    // RADIO DISABLED to Rx READY (Rx ramp-up) for BLE (µs).
	rtcRA   = 31     // RTC read accuracy (µs): readRTC <= RTC <= readRTC+rtcRA.
	aaDelay = 40 + 1 // Delay between start of packet and ADDRESS event (µs).
)

// SetRxTimers setups RTC0 and TIMER0 to trigger RXEN event or timeout DISABLE
// event using start and window (µs). TIMER0 still counts from some base event.
// Start contains the value of TIMER0 at which radio should start Rx.
func (c *Controller) setRxTimers(start, window uint32) {
	start -= rxRU + 1
	window += rxRU + aaDelay + rtcRA + 2

	rt := c.rtc0
	t := c.tim0

	// Capture RTC0 (takes 5 CPU cycles) and TIMER0 (just after RTC0).
	rtcBase := rt.LoadCOUNTER()
	t.Task(timer.CAPTURE(0)).Trigger()
	t.Task(timer.STOP).Trigger()
	t.Task(timer.CLEAR).Trigger()

	delay := start - t.LoadCC(0)
	if delay>>31 != 0 {
		c.LEDs[3].Set()
	}
	rtcTick := delay * 67 / 2048         // 67/2048 < 32768/1e6 == 512/15625
	timTick := delay - rtcTick*15625/512 // µs

	// rt.CC0 with t.CC0 controlls RXEN, t.CC1 controlls Rx timeout.

	rt.StoreCC(0, rtcBase+rtcTick)
	t.StoreCC(0, timTick)
	t.StoreCC(1, timTick+window)
	ppm := ppi.RTC0_COMPARE0__RADIO_TXEN.Mask() |
		ppi.RTC0_COMPARE0__RADIO_RXEN.Mask() |
		ppi.TIMER0_COMPARE0__RADIO_RXEN.Mask() |
		ppi.RADIO_ADDRESS__TIMER0_CAPTURE1.Mask() |
		ppi.RADIO_END__TIMER0_CAPTURE2.Mask()
	ppm.Disable()
	if timTick < 2 {
		ppm = ppi.RTC0_COMPARE0__TIMER0_START.Mask() |
			ppi.RTC0_COMPARE0__RADIO_RXEN.Mask()
	} else {
		// Use RTC+TIMER for better accuracy of RXEN.
		ppm = ppi.RTC0_COMPARE0__TIMER0_START.Mask() |
			ppi.TIMER0_COMPARE0__RADIO_RXEN.Mask()
	}
	ppm.Enable()
}

func (c *Controller) setupNoReq(timeout uint32) {
	c.noReq = 0
	c.noReqMax = 10
	return

	if timeout > 32e6 {
		timeout = 32e6 // 32 s
	}
	widen := c.sca.MulUp(c.connInterval)
	maxWin := c.connInterval - rxRU - rtcRA
	noReqMax := timeout / c.connInterval
	if noReqMax*2*widen >= maxWin {
		noReqMax = maxWin/(2*widen) - 1
	}
	c.noReq = 0
	c.noReqMax = uint16(noReqMax)
}

func (c *Controller) connectReqRxDisabledISR() {
	r := c.radio
	if !r.LoadCRCSTATUS() {
		// Return to advertising.
		r.Task(radio.TXEN).Trigger()
		c.scanDisabledISR()
		c.LEDs[2].Clear()
		return
	}
	// Both timers (RTC0, TIMER0) are running. TIMER0.CC2 contains time of END
	// event (end of received ConnectReq packet).

	d := connReqLLData(c.rxaPDU.Payload()[6+6:])
	c.chm.Init(d.ChMapL(), d.ChMapH(), d.Hop())
	rxPDU := c.recv.Get().PDU

	r.StoreCRCINIT(d.CRCInit())
	radioSetMaxLen(r, rxPDU.PayLen())
	radioSetAA(r, d.AA())
	radioSetChi(r, c.chm.NextChi())
	radioSetPDU(r, rxPDU)
	r.StoreSHORTS(radio.READY_START | radio.END_DISABLE)
	ev := r.Event(radio.ADDRESS)
	ev.Clear()
	ev.EnableIRQ()

	c.winSize = d.WinSize()
	c.sca = d.SCA() + ppmToFixedUp(100) // Assume 100 ppm local SCA.
	c.connInterval = d.Interval()
	c.setupNoReq(d.Timeout())

	winOffset := d.WinOffset()
	widen := c.sca.MulUp(winOffset + c.winSize) // Widen at end of window (µs).

	// Setup first anchor point.
	c.setRxTimers(c.tim0.LoadCC(2)+winOffset-widen, c.winSize+2*widen)
	c.isr = (*Controller).connRxDisabledISR

	// Reenable ADDRESS time capture. See setRxTimer and connTxDisabledISR.
	ppi.RADIO_ADDRESS__TIMER0_CAPTURE1.Enable()

	c.txcPDU.SetHeader(ble.L2CAPCont | c.flags&(ble.NESN|ble.SN))
	c.txcPDU.SetPayLen(0)
	c.rspRef = c.txcPDU.PDU

	for _, led := range c.LEDs {
		led.Clear()
	}
}

func (c *Controller) connRxDisabledISR() {
	r := c.radio
	if !r.Event(radio.ADDRESS).IsSet() {
		// Timeout.
		if c.winSize != 0 {
			// Still in Connection Setup state.
			c.LEDs[0].Set()
			for {
				// BUG:
			}
		}
		if c.noReq++; c.noReq == c.noReqMax {
			// Connection timed-out.
			c.LEDs[1].Set()
			for {
				// BUG:
			}
		}
		led := c.LEDs[2]
		led.Store(^led.Load())

		c.flags &^= ble.MD
		radioSetChi(r, c.chm.NextChi())

		c.tim0.Task(timer.CAPTURE(0)).Trigger()
		c.setRxTimers(c.tim0.LoadCC(0)+c.connInterval/2, c.connInterval)
		/*
			// TIMER0 still counts. TIMER0.CC0 contains time of previous Rx start.
			widen := c.sca.MulUp(c.connInterval)
			c.setRxTimers(
				c.tim0.LoadCC(0)+c.connInterval-widen,
				2*(uint32(c.noReq)+1)*widen,
			)
		*/

		// Need to record time of Rx ADDRESS event. At the same time, this
		// avoids Rx timeout (sets it to >1h).
		ppi.RADIO_ADDRESS__TIMER0_CAPTURE1.Enable() // setRxTimers disabled it.

		return
	}
	c.noReq = 0
	c.winSize = 0

	/*
		// Test safety margin to START event.
		r.Event(radio.READY).Clear()
	*/

	// Both timers (RTC0, TIMER0) are running. TIMER0.CC1 contains time of
	// ADDRESS event (end of address field in received data packet),

	// BUG: In case of MD (more data) Rx timers will be broken.
	if c.connUpd.CheckInstant(c.chm.ConnEventCnt() + 1) {
		c.winSize = c.connUpd.WinSize()
		winOffset := c.connUpd.WinOffset() + c.connInterval
		c.connInterval = c.connUpd.Interval()
		c.setupNoReq(c.connUpd.Timeout())
		c.connUpd.SetDone()
		widen := c.sca.MulUp(winOffset + c.winSize) // Abs. SCA at end o window.
		c.setRxTimers(
			c.tim0.LoadCC(1)+winOffset-widen-aaDelay,
			c.winSize+2*widen,
		)
	} else {
		widen := c.sca.MulUp(c.connInterval) // Abs. SCA after connInterval (µs).
		c.setRxTimers(c.tim0.LoadCC(1)+c.connInterval-widen-aaDelay, 2*widen)
	}

	if r.LoadCRCSTATUS() {
		rxPDU := c.recv.Get()
		header := rxPDU.Header()
		//c.flags |= header&ble.MD // BUG: MD not supported.

		/*
			c.LEDs[0].Store(int(header) >> 3)
			c.LEDs[1].Store(int(c.flags) >> 3)
			c.LEDs[2].Store(int(header) >> 2)
			c.LEDs[3].Store(int(c.flags) >> 2)
		*/

		if header&ble.SN == c.flags&ble.NESN<<1 {
			// New PDU received.
			llid := header & ble.LLID
			switch {
			case llid == ble.LLControl:
				// LL Control PDU. Pass it for further processing if the
				// previous one was done.
				if c.rxcRef.IsZero() {
					c.rxcRef = rxPDU
					c.flags ^= ble.NESN
				}
			case llid == ble.L2CAPCont && rxPDU.PayLen() == 0:
				// Empty L2CAP PDU.
				c.flags ^= ble.NESN
			default:
				// Non-empty L2CAP PDU.
				select {
				case c.recv.Ch <- rxPDU:
					c.recv.Next()
					c.flags ^= ble.NESN
				default:
				}
			}
		}
		if header&ble.NESN != c.flags&ble.SN>>1 {
			// Previous packet ACKed. Can send new one.
			var rspPDU ble.DataPDU
			if !c.rxcRef.IsZero() {
				// Process last controll PDU.
				header = ble.LLControl
				rspPDU = c.txcPDU
				req := c.rxcRef.Payload()
				opcode, req := req[0], req[1:]
				switch opcode {
				case llConnUpdateReq:
					for len(req) != 11 {
						// BUG:
					}
					c.connUpd.Init(req)
					for c.connUpd.Instant()-c.chm.ConnEventCnt() >= 32767 {
						// BUG:
					}
					rspPDU = ble.DataPDU{} // There is no llConnUpdateRsp.
					c.recv.Ch <- rxPDU
				case llChanMapReq:
					for len(req) != 7 {
						// BUG:
					}
					instant := le.Decode16(req[5:])
					for instant-c.chm.ConnEventCnt() >= 32767 {
						// BUG:
					}
					c.chm.Update(le.Decode32(req), req[4], instant)
					rspPDU = ble.DataPDU{} // There is no llChanMapRsp.
					c.recv.Ch <- rxPDU
				case llFeatureReq:
					rspPDU.SetPayLen(9)
					pay := rspPDU.Payload()
					pay[0] = llFeatureRsp
					pay[1] = 0
					pay[2] = 0
					pay[3] = 0
					pay[4] = 0
					pay[5] = 0
					pay[6] = 0
					pay[7] = 0
					pay[8] = 0
				case llVersionInd:
					rspPDU.SetPayLen(6)
					pay := rspPDU.Payload()
					pay[0] = llVersionInd
					pay[1] = 6    // BLE version: 6: 4.0, 7: 4.1, 8: 4.2, 9: 5.
					pay[2] = 0xFF // Company ID (2 octets).
					pay[3] = 0xFF // Using 0xFFFF: tests / not assigned.
					pay[4] = 0    // Subversion (2 octets). Unique for each
					pay[5] = 0    // implementation or revision of controller.
				default:
					rspPDU.SetPayLen(2)
					pay := rspPDU.Payload()
					pay[0] = llUnknownRsp
					pay[1] = opcode
					c.recv.Ch <- rxPDU
				}
				c.rxcRef = ble.DataPDU{}
			}
			if rspPDU.IsZero() {
				// Send data PDU from send queue or empty PDU.
				select {
				case rspPDU = <-c.send.Ch:
					header = rspPDU.Header()
				default:
					header = ble.L2CAPCont
					rspPDU = c.txcPDU
					rspPDU.SetPayLen(0)
				}
			}
			c.flags ^= ble.SN
			rspPDU.SetHeader(header | c.flags&(ble.NESN|ble.SN))
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
	c.isr = (*Controller).connTxDisabledISR
}

func (c *Controller) connTxDisabledISR() {
	r := c.radio
	if c.flags&ble.MD == 0 {
		radioSetChi(r, c.chm.NextChi())
	}
	ev := r.Event(radio.ADDRESS)
	ev.Clear()
	ev.EnableIRQ()

	// Need to record time of Rx ADDRESS event. At the same time, this avoids
	// Rx timeout (sets it to >1h).
	ppi.RADIO_ADDRESS__TIMER0_CAPTURE1.Enable() // setRxTimers disabled it.

	radioSetPDU(r, c.recv.Get().PDU)

	c.isr = (*Controller).connRxDisabledISR
}
