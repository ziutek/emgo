package sdmmc

import (
	"unsafe"

	"stm32/hal/internal"
	"stm32/hal/system"
)

type Periph periph

// Bus returns a bus to which p is connected to.
func (p *Periph) Bus() system.Bus {
	return internal.Bus(unsafe.Pointer(p))
}

// EnableClock enables clock for p.
// lp determines whether the clock remains on in low power (sleep) mode.
func (p *Periph) EnableClock(lp bool) {
	addr := unsafe.Pointer(p)
	internal.APB_SetLPEnabled(addr, lp)
	internal.APB_SetEnabled(addr, true)
}

// DisableClock disables clock for p.
func (p *Periph) DisableClock() {
	internal.APB_SetEnabled(unsafe.Pointer(p), false)
}

// Reset resets p.
func (p *Periph) Reset() {
	internal.APB_Reset(unsafe.Pointer(p))
}

// Enabled reports whether the p is enabled.
func (p *Periph) Enabled() bool {
	return p.raw.PWRCTRL().Load() == 3
}

// Enable enables p. At least 7 PCLK clock periods are needed between any
// Enable or Disable. At least 3 SDMMCCLK clock periods plus 3 PCLK clock
// periods are needed between any Enable or Disable.
func (p *Periph) Enable() {
	p.raw.POWER.Store(3)
}

// Disable disables gp. At least 7 PCLK clock periods are needed between any
// Enable or Disable. At least 3 SDMMCCLK clock periods plus 3 PCLK clock
// periods are needed between any Enable or Disable.
func (p *Periph) Disable() {
	p.raw.POWER.Store(0)
}

type BusClock byte

const (
	ClkEna   BusClock = 1 << 0 // Enables bus clock.
	PwrSave  BusClock = 1 << 1 // Enables power saving mode.
	ClkByp   BusClock = 1 << 2 // Pass SDMMC clock directly to CK pin.
	BusWidth BusClock = 3 << 3 // Describes data bus width.
	Bus1     BusClock = 0 << 3 // Single data bus line.
	Bus4     BusClock = 1 << 3 // Four data bus lines.
	Bus8     BusClock = 2 << 3 // Eight data bus lines.
	NegEdge  BusClock = 1 << 5 // Command and data changed on CK falling edge.
	FlowCtrl BusClock = 1 << 6 // Enables hardware flow controll.
)

// BusClock returns the current configuration of SDMMC bus.
func (p *Periph) BusClock() (cfg BusClock, clkdiv int) {
	clkcr := p.raw.CLKCR.Load()
	return BusClock(clkcr >> 8), int(clkcr & 255)
}

// SetBusClock configures the SDMMC bus and clock. If ClkByp is set bus clock
// frequency is equal to SDMMCCLK, otherwise it is SDMMCCLK / (clkdiv+2).
//
// Note (Errata Sheet DocID027036 Rev 2):
// 2.7.1 Don't use hardware flow controll (FlowCtrl).
// 2.7.3 Don't use clock dephasing (NegEdge).
// 2.7.5 Ensure that:
//
//	3*period(PCLK)+3*period(SDMMCCLK) < 32/BusWidth*period(SDMMC_CK)
//  (always met for: PCLK > 28.8 MHz).
//
// At least 3 SDMMCCLK clock periods plus 2 PCLK clock periods are needed
// between subsequent SetBusClock.
func (p *Periph) SetBusClock(cfg BusClock, clkdiv int) {
	if uint(clkdiv) > 255 {
		panic("sdio: bad clkdiv")
	}
	p.raw.CLKCR.U32.Store(uint32(cfg)<<8 | uint32(clkdiv))
}

// Arg returns value of command argument.
func (p *Periph) Arg() uint32 {
	return p.raw.ARG.U32.Load()
}

// SetArg sets command argument.
func (p *Periph) SetArg(arg uint32) {
	p.raw.ARG.U32.Store(arg)
}

type Command uint16

const (
	CmdIdx   Command = 63 << 0 // Command index.
	HasResp  Command = 1 << 6  // Response expected.
	LongResp Command = 1 << 7  // Long response.
	WaitInt  Command = 1 << 8  // Disable command timeout and wait for IRQ.
	WaitPend Command = 1 << 9  // Wait for end of stream data transfer.
	CmdEna   Command = 1 << 10 // Enable CPSM (send command) / CPSM is enabled.
	SuspIO   Command = 1 << 11 // SDIO suspend command.
)

// Cmd returns the last command and configuration/state of the Command Path
// State Machine (CPSM).
func (p *Periph) Cmd() Command {
	return Command(p.raw.CMD.Load())
}

// SetCmd passes command and sets configuration of the Command Path State
// Machine (CPSM). At least 3 SDMMCCLK clock periods plus 3 PCLK clock
// periods are needed between subsequent SetCmd.
func (p *Periph) SetCmd(cmd Command) {
	p.raw.CMD.U32.Store(uint32(cmd))
}

// RespCmd returns the command index field of the last command response.
func (p *Periph) RespCmd() Command {
	return Command(p.raw.RESPCMD.Load()) & CmdIdx
}

// Resp returns n-th 32-bit word of the last received response. Resp(0) returns
// the most significant bits, Resp(3) returns the least significant bits.
func (p *Periph) Resp(n int) uint32 {
	return p.raw.RESP[n].U32.Load()
}

// DataTimeout returns data timeout period as number of card bus clock periods.
func (p *Periph) DataTimeout() uint {
	return uint(p.raw.DTIMER.Load())
}

// SetDataTimeout sets data timeout period as number of card bus clock periods.
func (p *Periph) SetDataTimeout(ck uint) {
	p.raw.DTIMER.U32.Store(uint32(ck))
}

// DataLen returns the number of data bytes to be transfered.
func (p *Periph) DataLen() int {
	return int(p.raw.DLEN.Load())
}

// SetDataLen sets the number of data bytes to be transfered.
func (p *Periph) SetDataLen(dlen int) {
	if uint(dlen) > 1<<24-1 {
		panic("sdio: bad data len")
	}
	p.raw.DLEN.U32.Store(uint32(dlen))
}

// DataCtrl represents Data Path State Machine (DPMS) configuration.
type DataCtrl uint16

const (
	DTEna    DataCtrl = 1 << 0  // Enable data transfer.
	Send     DataCtrl = 0 << 1  // Send data to card.
	Recv     DataCtrl = 1 << 1  // Receive data from card.
	Stream   DataCtrl = 1 << 2  // Stream or SDIO multibyte data transfer.
	UseDMA   DataCtrl = 1 << 3  // Use DMA.
	Block1   DataCtrl = 0 << 4  // Block data transfer, block size: 1 B.
	Block2   DataCtrl = 1 << 4  // Block data transfer, block size: 2 B.
	Block4   DataCtrl = 2 << 4  // Block data transfer, block size: 4 B.
	Block8   DataCtrl = 3 << 4  // Block data transfer, block size: 8 B.
	Block16  DataCtrl = 4 << 4  // Block data transfer, block size: 16 B.
	Block32  DataCtrl = 5 << 4  // Block data transfer, block size: 32 B.
	Block64  DataCtrl = 6 << 4  // Block data transfer, block size: 64 B.
	Block128 DataCtrl = 7 << 4  // Block data transfer, block size: 128 B.
	Block256 DataCtrl = 8 << 4  // Block data transfer, block size: 256 B.
	Block512 DataCtrl = 9 << 4  // Block data transfer, block size: 512 B.
	Block1K  DataCtrl = 10 << 4 // Block data transfer, block size: 1 KiB.
	Block2K  DataCtrl = 11 << 4 // Block data transfer, block size: 2 KiB.
	Block4K  DataCtrl = 12 << 4 // Block data transfer, block size: 4 KiB.
	Block8K  DataCtrl = 13 << 4 // Block data transfer, block size: 8 KiB.
	Block16K DataCtrl = 14 << 4 // Block data transfer, block size: 16 KiB.
	RWStart  DataCtrl = 1 << 8  // Read wait start.
	RWStop   DataCtrl = 1 << 9  // Read wait stop.
	RWCK     DataCtrl = 1 << 10 // Read wait constrol using CK instead of D2.
	IO       DataCtrl = 1 << 11 // SDIO specific operation.
)

// DataCtrl returns current state/configuration of Data Path State Machine
// (DPMS).
func (p *Periph) DataCtrl() DataCtrl {
	return DataCtrl(p.raw.DCTRL.Load())
}

// SetDataCtrl controls Data Path State Machine (DPMS).
func (p *Periph) SetDataCtrl(cfg DataCtrl) {
	p.raw.DCTRL.U32.Store(uint32(cfg))
}

// RemainBytes returns the number of remaining data bytes to be transfered.
func (p *Periph) RemainBytes() int {
	return int(p.raw.DCOUNT.Load())
}

type Event uint32

const (
	CmdRespOK   Event = 1 << 6  // Command response received, CRC OK.
	CmdSent     Event = 1 << 7  // Command sent (no response required).
	DataEnd     Event = 1 << 8  // DataCount() == 0.
	DataBlkEnd  Event = 1 << 10 // Data block sent/received, CRC OK.
	CmdAct      Event = 1 << 11 // Command transfer in progress.
	TxAct       Event = 1 << 12 // Data transmit in progress.
	RxAct       Event = 1 << 13 // Data receive in progress.
	TxHalfEmpty Event = 1 << 14 // Tx FIFO half empty.
	RxHalfFull  Event = 1 << 15 // Rx FIFO half full.
	TxFull      Event = 1 << 16 // Tx FIFO full.
	RxFull      Event = 1 << 17 // Rx FIFO full.
	TxEmpty     Event = 1 << 18 // Tx FIFO empty.
	RxEmpty     Event = 1 << 19 // Rx FIFO empty.
	TxNotEmpty  Event = 1 << 20 // Tx FIFO not empty.
	RxNotEmpty  Event = 1 << 21 // Rx FIFO not empty.
	IOIRQ       Event = 1 << 22 // SDIO interrupt request.
	EvAll       Event = 0x7FFDC0
)

type Error byte

const (
	ErrCmdCRC      Error = 1 << 0 // Command response received, CRC failed.
	ErrDataCRC     Error = 1 << 1 // Data response receifed, CRC failed.
	ErrCmdTimeout  Error = 1 << 2 // Command response timeout.
	ErrDataTimeout Error = 1 << 3 // Data response timeout.
	ErrTxUnderrun  Error = 1 << 4 // Tx FIFO underrun.
	ErrRxOverrun   Error = 1 << 5 // Rx FIFO overrun.
	ErrAll         Error = 0x3F
)

func (err Error) Error() string {
	var (
		s string
		d Error
	)
	switch {
	case err&ErrCmdCRC != 0:
		d = ErrCmdCRC
		s = "sdmmc: cmd CRC+"
	case err&ErrDataCRC != 0:
		d = ErrDataCRC
		s = "sdmmc: data CRC+"
	case err&ErrCmdTimeout != 0:
		d = ErrCmdTimeout
		s = "sdmmc: cmd timeout+"
	case err&ErrDataTimeout != 0:
		d = ErrDataTimeout
		s = "sdmmc: data timeout+"
	case err&ErrTxUnderrun != 0:
		d = ErrTxUnderrun
		s = "sdmmc: Tx underrun+"
	case err&ErrRxOverrun != 0:
		d = ErrRxOverrun
		s = "sdmmc: Rx overrun+"
	default:
		return ""
	}
	if err == d {
		s = s[:len(s)-1]
	}
	return s
}

// Status returns the status bits: events and errors.
//
// Note (Errata Sheet DocID027036 Rev 2 2.7.2): Ignore ErrCmdCRC for R3
// and R4 responses.
//
func (p *Periph) Status() (Event, Error) {
	sta := p.raw.STA.Load()
	return Event(sta) & EvAll, Error(sta) & ErrAll
}

// Clear clears specified events and errors. All errors and CmdRespOK, CmdSent,
// DataEnd, DataBlkEnd, IOIRQ events can be cleared this way.
func (p *Periph) Clear(ev Event, err Error) {
	p.raw.ICR.U32.Store(uint32(ev) | uint32(err))
}

// EnableIRQ enables generating of IRQ by specified events and errors.
func (p *Periph) EnableIRQ(ev Event, err Error) {
	p.raw.MASK.U32.SetBits(uint32(ev) | uint32(err))
}

// DisableIRQ disables generating of IRQ by specified events and errors.
func (p *Periph) DisableIRQ(ev Event, err Error) {
	p.raw.MASK.U32.ClearBits(uint32(ev) | uint32(err))
}

// RemainWords returns the number of remaining words that can be read from /
// written to the FIFO: (RemainBytes() + 3) / 4.
func (p *Periph) RemainWords() int {
	return int(p.raw.FIFOCNT.Load())
}

// Load loads one word from FIFO.
func (p *Periph) Load() uint32 {
	return p.raw.FIFO.U32.Load()
}

// Store stores one word into FIFO.
func (p *Periph) Store(w uint32) {
	p.raw.FIFO.U32.Store(w)
}

//emgo:inline
func burstCopyPTM(p, m uintptr) uintptr

//emgo:inline
func burstCopyMTP(m, p uintptr) uintptr

func panicShortBuf() {
	panic("sdio: buf too short")
}

// BurstLoad reads 16 words from FIFO using LDM and STM instructions.
func (p *Periph) BurstLoad(buf []uint32) {
	if len(buf) < 16 {
		panicShortBuf()
	}
	burstCopyPTM(p.raw.FIFO.Addr(), uintptr(unsafe.Pointer(&buf[0])))
}

// BurstStore stores 16 words into FIFO using LDM and STM instructions.
func (p *Periph) BurstStore(buf []uint32) {
	if len(buf) < 16 {
		panicShortBuf()
	}
	burstCopyMTP(uintptr(unsafe.Pointer(&buf[0])), p.raw.FIFO.Addr())
}
