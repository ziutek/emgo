package blec

import (
	"rtos"
	"unsafe"

	"nrf5/hal/ficr"
	"nrf5/hal/radio"
	"nrf5/hal/te"
)

func radioInit(r *radio.Periph, maxPayLen int) {
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
	r.StoreRXADDRESSES(1)
	r.StoreTXADDRESS(0)
	r.StoreCRCCNF(3, true)
	r.StoreCRCPOLY(0x100065B)
	r.DisableIRQ(te.EvAll)
	irq := rtos.IRQ(r.IRQ())
	irq.SetPrio(rtos.IRQPrioHighest)
	irq.Enable()
}

func radioSetMaxLen(r *radio.Periph, maxPayLen int) {
	r.StorePCNF1(radio.WhiteEna | radio.MaxLen(maxPayLen) | radio.BALen(3))
}

func radioSetPDU(r *radio.Periph, pdu []byte) {
	r.StorePACKETPTR(unsafe.Pointer(&pdu[0]))
}

func radioSetAA(r *radio.Periph, addr uint32) {
	r.StoreBASE(0, addr<<8)
	r.StorePREFIX(0, addr>>24)
}

func radioSetChi(r *radio.Periph, chi byte) {
	r.StoreDATAWHITEIV(uint32(chi))
	var ch int
	switch {
	case chi <= 10:
		ch = int(chi*2 + 4)
	case chi <= 36:
		ch = int(chi*2 + 6)
	case chi == 37:
		ch = 2
	case chi == 38:
		ch = 26
	case chi == 39:
		ch = 80
	default:
		panic("ble: bad ch.index")
	}
	r.StoreFREQUENCY(radio.Channel(ch))
}

func radioSetAddrMatch(r *radio.Periph, addr int64) {
	r.StoreDAB(0, uint32(addr))
	msw := uint32(addr >> 32)
	r.StoreDAP(0, uint16(msw))
	r.StoreDACNF(1, byte(addr>>31))
}
