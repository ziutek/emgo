// Package power provides interface to power managemnt peripheral.
package power

import (
	"mmio"
	"unsafe"

	"nrf5/hal/internal/mmap"
	"nrf5/hal/te"
)

type Periph struct {
	te.Regs

	resetreas mmio.U32     // 0x400
	_         [9]mmio.U32  //
	ramstatus mmio.U32     // 0x428
	_         [53]mmio.U32 //
	systemoff mmio.U32     // 0x500
	_         [3]mmio.U32  //
	pofcon    mmio.U32     // 0x510
	_         [2]mmio.U32  //
	gpregret  [2]mmio.U32  // 0x51C
	ramon     mmio.U32     // 0x524
	_         [7]mmio.U32  //
	reset     mmio.U32     // 0x544
	_         [3]mmio.U32  //
	ramonb    mmio.U32     // 0x554
	_         [8]mmio.U32  //
	dcdcen    mmio.U32     // 0x578
	_         [225]mmio.U32
	ram       [8]struct{ power, powerset, powerclr mmio.U32 }
}

//emgo:const
var POWER = (*Periph)(unsafe.Pointer(mmap.APB_BASE))

type Task byte

const (
	CONSTLAT Task = 0x78 // Enable constant latency mode.
	LOWPWR   Task = 0x7C // Enable low power mode (variable latency).
)

type Event byte

const (
	POFWARN    Event = 0x08 // Power failure warning.
	SLEEPENTER Event = 0x14 // CPU entered WFI/WFE sleep (nRF52).
	SLEEPEXIT  Event = 0x18 // CPU exited WFI/WFE sleep (nRF52).
)

func (p *Periph) Task(t Task) *te.Task    { return p.Regs.Task(int(t)) }
func (p *Periph) Event(e Event) *te.Event { return p.Regs.Event(int(e)) }

// ResetReas is a bitfield that describes reset reason.
type ResetReas uint32

const (
	RESETPIN ResetReas = 1 << 0  // Reset from pin-reset.
	DOG      ResetReas = 1 << 1  // Reset from watchdog.
	SREQ     ResetReas = 1 << 2  // Reset from AIRCR.SYSRESETREQ.
	LOCKUP   ResetReas = 1 << 3  // Reset from CPU lock-up.
	OFF      ResetReas = 1 << 16 // Wake up from OFF mode by GPIO DETECT.
	LPCOMP   ResetReas = 1 << 17 // Wake up from OFF mode by LPCOMP ANADETECT.
	DIF      ResetReas = 1 << 18 // Wake up from OFF mode by debug interface.
	NFC      ResetReas = 1 << 19 // Wake up from OFF mode by NFC.
)

// LoadRESETREAS returns reset reason bits.
func (p *Periph) LoadRESETREAS() ResetReas {
	return ResetReas(p.resetreas.Load())
}

// ClearRESETREAS clears reset reason bits specified by mask.
func (p *Periph) ClearRESETREAS(mask ResetReas) {
	p.resetreas.Store(uint32(mask))
}

// RAMStatus is a bitfield that describes RAM blocks status.
type RAMStatus byte

const (
	RAMBLOCK0 RAMStatus = 1 << 0 // RAM block 0 is powered up.
	RAMBLOCK1 RAMStatus = 1 << 1 // RAM block 1 is powered up.
	RAMBLOCK2 RAMStatus = 1 << 2 // RAM block 2 is powered up.
	RAMBLOCK3 RAMStatus = 1 << 3 // RAM block 3 is powered up.
)

// LoadRAMSTATUS returns bitfield that describes status of RAM blocks.
func (p *Periph) LoadRAMSTATUS() RAMStatus {
	return RAMStatus(p.ramstatus.Load())
}

// SetSYSTEMOFF sets system into OFF state.
func (p *Periph) SetSYSTEMOFF() {
	p.systemoff.Store(1)
}

// POFCon is power failure comparator configuration.
type POFCon byte

const (
	POF       POFCon = 1 << 0  // Set if power failure comparoator is enabled.
	THRESHOLD POFCon = 15 << 1 // Power failure comparator threshold mask.

	// Power failure comparoator thresholds.

	V2_1 POFCon = 0 << 1 // Threshold: 2.1 V (nrF51).
	V2_3 POFCon = 1 << 1 // Threshold: 2.3 V (nrF51).
	V2_5 POFCon = 2 << 1 // Threshold: 2.5 V (nrF51).
	V2_7 POFCon = 3 << 1 // Threshold: 2.5 V (nrF51).

	V17 POFCon = 4 << 1  // Threshold: 1.7 V (nRF52).
	V18 POFCon = 5 << 1  // Threshold: 1.8 V (nRF52).
	V19 POFCon = 6 << 1  // Threshold: 1.9 V (nRF52).
	V20 POFCon = 7 << 1  // Threshold: 2.0 V (nRF52).
	V21 POFCon = 8 << 1  // Threshold: 2.1 V (nRF52).
	V22 POFCon = 9 << 1  // Threshold: 2.2 V (nRF52).
	V23 POFCon = 10 << 1 // Threshold: 2.3 V (nRF52).
	V24 POFCon = 11 << 1 // Threshold: 2.4 V (nRF52).
	V25 POFCon = 12 << 1 // Threshold: 2.5 V (nRF52).
	V26 POFCon = 13 << 1 // Threshold: 2.6 V (nRF52).
	V27 POFCon = 14 << 1 // Threshold: 2.7 V (nRF52).
	V28 POFCon = 15 << 1 // Threshold: 2.8 V (nRF52).
)

// LoadPOFCON returns power failure comparoator configuration.
func (p *Periph) LoadPOFCON() POFCon {
	return POFCon(p.pofcon.Load())
}

// StorePOFCON sets power failure comparoator configuration.
func (p *Periph) StorePOFCON(pofcon POFCon) {
	p.pofcon.Store(uint32(pofcon))
}

// GPREGRET returns pointer to n-th general purpose retention register. nRF51
// supports one, nRF52 supports two. Only lowest 8 bits can be used.
func (p *Periph) GPREGRET(n int) *mmio.U32 {
	return &p.gpregret[n]
}

/*
// LoadRAMON
func (p *Periph) LoadRAMON() RAMBlocks {
	return RAMBlocks(p.ramon.Load())
}

func (p *Periph) StoreRAMON(ramon RAMBlocks) {
	p.ramon.Store(uint32(ramon))
}
*/
