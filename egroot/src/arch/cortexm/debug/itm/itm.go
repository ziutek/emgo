// Package itm provides interface to Instrumentation Trace Macrocell
package itm

import (
	"mmio"
	"unsafe"
)

type regs struct {
	stim [256]mmio.U32
	_    [640]mmio.U32
	te   [8]mmio.U32
	_    [8]mmio.U32
	tp   mmio.U32
	_    [15]mmio.U32
	tc   mmio.U32
}

var irs = (*regs)(unsafe.Pointer(uintptr(0xe0000000)))

type Control uint32

// Flags
const (
	ITMEna Control = 1 << iota
	TSEna
	SyncEna
	TxEna
	SWOEna

	Busy Control = 1 << 23
)

func (c Control) TSPresc() int {
	return int(c>>8) & 3
}

func (c *Control) SetTSPresc(f int) {
	*c = *c&^(3<<8) | Control(f&3)<<8
}

func (c Control) GTSFreq() int {
	return int(c>>10) & 3
}

func (c *Control) SetGTSFreq(f int) {
	*c = *c&^(3<<10) | Control(f&3)<<10
}

// Ctrl returns value of ITM Trace Control Register.
func Ctrl() Control {
	return Control(irs.tc.Load())
}

func SetCtrl(c Control) {
	irs.tc.Store(uint32(c))
}

// PrivMask returns conten of Trace Privilege Register. Every bit in returned
// value corresponds to eight stimulus ports. If bit is set then the
// corresponding ports can be accessed by privileged code only.
//
// Bits for unimplemented ports are always returned as 0. To determine the
// number of implemnted ports call SetPrivMask(0xffffffff) and next call
// PrivMask().
func PrivMask() uint32 {
	return irs.tp.Load()
}

// SetPrivMask writes mask to Trace Privilege Register.
func SetPrivMask(mask uint32) {
	irs.tp.Store(mask)
}

type Port int

// Enabled returns true if port is enabled.
func (p Port) Enabled() bool {
	bit := int(p & 31)
	p >>= 5
	return irs.te[p].Bit(bit) != 0
}

// Enable enables stimulus por.
func (p Port) Enable() {
	bit := int(p & 31)
	p >>= 5
	irs.te[p].SetBit(bit)
}

// Disable enables stimulus por.
func (p Port) Disable() {
	bit := int(p & 31)
	p >>= 5
	irs.te[p].ClearBit(bit)
}

// Ready returns true if port can accept data.
func (p Port) Ready() bool {
	return irs.stim[p].Bit(0) != 0
}

// Store8 stores byte into p.
func (p Port) Store8(b byte) {
	mmio.PtrU8(unsafe.Pointer(irs.stim[p].Addr())).Store(b)
}

// Store16 stores half-word into p.
func (p Port) Store16(h uint16) {
	mmio.PtrU16(unsafe.Pointer(irs.stim[p].Addr())).Store(h)
}

// Store32 stores word into p.
func (p Port) Store32(w uint32) {
	irs.stim[p].Store(w)
}

// WriteString implements io.StringWriter interface. Use p < 0 to disable
// writtening (useful for temporary disable debug messages).
func (p Port) WriteString(s string) (int, error) {
	if p < 0 {
		return len(s), nil
	}
	n := len(s)
	i := 0
loop:
	for n > 0 {
		for !p.Ready() {
			if !p.Enabled() || Ctrl()&ITMEna == 0 {
				// Silently discard data.
				break loop
			}
		}
		switch {
		case n >= 4:
			p.Store32(uint32(s[i]) + uint32(s[i+1])<<8 + uint32(s[i+2])<<16 +
				uint32(s[i+3])<<24)
			n -= 4
			i += 4
		case n >= 2:
			p.Store16(uint16(s[i]) + uint16(s[i+1])<<8)
			n -= 2
			i += 2
		default:
			p.Store8(s[i])
			n--
			i++
		}
	}
	return n + i, nil
}

// Write implements io.Writer interface. Use p < 0 to disable writtening
// (useful for temporary disable debug messages).
func (p Port) Write(b []byte) (int, error) {
	return p.WriteString(*(*string)(unsafe.Pointer(&b)))
}

func (p Port) WriteByte(b byte) error {
	p.Store8(b)
	return nil
}
