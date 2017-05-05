package ble

// AdvPDU represents Advertising Channel PDU. It should be initialized before
// use.
type AdvPDU struct {
	b []byte
}

// MakeAdvPDU returns ready to use AdvPDU. If b is nil then MakeAdvPDU allocates
// 39 bytes that is enough to store any valid BLE4.x Advertising Channel PDU.
// Returned variable refers to the same memory as b or to the allocated storage.
func MakeAdvPDU(b []byte) AdvPDU {
	if b == nil {
		b = make([]byte, 39)
	}
	b[0] = 0
	b[1] = 0
	return AdvPDU{b}
}

func AsAdvPDU(b []byte) AdvPDU {
	return AdvPDU{b}
}

// IsNil returns true in case of uninitialized pdu variable.
func (pdu AdvPDU) IsNil() bool {
	return pdu.b == nil
}

func (pdu AdvPDU) length() int {
	return int(pdu.b[1]) + 2
}

func (pdu AdvPDU) setLength(n int) {
	if uint(n) > 39 {
		panic("ble: bad PDU length")
	}
	pdu.b[1] = byte(n - 2)
}

// Bytes returns underlying bytes of pdu.
func (pdu AdvPDU) Bytes() []byte {
	return pdu.b[:pdu.length()]
}

func (pdu AdvPDU) Payload() []byte {
	return pdu.b[2:pdu.length()]
}

// AdvPDUType represents Advertising Channel PDU type.
type AdvPDUType byte

const (
	AdvInd        AdvPDUType = 0
	AdvDirectInd  AdvPDUType = 1
	AdvNonconnInd AdvPDUType = 2
	ScanReq       AdvPDUType = 3
	ScanRsp       AdvPDUType = 4
	ConnectReq    AdvPDUType = 5
	AdvScanInd    AdvPDUType = 6
)

// Type returns pdu type from header.
func (pdu AdvPDU) Type() AdvPDUType {
	return AdvPDUType(pdu.b[0] & 0xf)
}

// SetType sets pdu type in header.
func (pdu AdvPDU) SetType(typ AdvPDUType) {
	pdu.b[0] = pdu.b[0]&^0xf | byte(typ)
	pdu.Reset()
}

// TxAdd returns TxAdd field from header.
func (pdu AdvPDU) TxAdd() bool {
	return pdu.b[0]>>6&1 != 0
}

// SetTxAdd sets TxAdd field in header.
func (pdu AdvPDU) SetTxAdd(rnda bool) {
	if rnda {
		pdu.b[0] |= 1 << 6
	} else {
		pdu.b[0] &^= 1 << 6
	}
}

// RxAdd returns RxAdd field from header.
func (pdu AdvPDU) RxAdd() bool {
	return pdu.b[0]>>7 != 0
}

// SetRxAdd sets RxAdd field in header.
func (pdu AdvPDU) SetRxAdd(rnda bool) {
	if rnda {
		pdu.b[0] |= 1 << 7
	} else {
		pdu.b[0] &^= 1 << 7
	}
}

// Reset resets length of pdu.
func (pdu AdvPDU) Reset() {
	pdu.setLength(2)
}

func (pdu AdvPDU) AppendAddr(addr int64) {
	n := pdu.length()
	pdu.setLength(n + 6)
	pdu.b[n] = byte(addr)
	pdu.b[n+1] = byte(addr >> 8)
	pdu.b[n+2] = byte(addr >> 16)
	pdu.b[n+3] = byte(addr >> 24)
	pdu.b[n+4] = byte(addr >> 32)
	pdu.b[n+5] = byte(addr >> 40)
}

// Flags
const (
	LimitedDisc  = 1 << 0 // LE Limited Discoverable mode.
	GeneralDisc  = 1 << 1 // LE General Discoverable mode.
	OnlyLE       = 1 << 2 // BR/EDR not supported.
	DualModeCtlr = 1 << 3 // LE and BR/EDR capable (controller).
	DualModeHost = 1 << 4 // LE and BR/EDR capable (host).
)

type ADType byte

const (
	Flags ADType = 0x01 // Flags

	MoreServices16 ADType = 0x02 // More 16-bit UUIDs available.
	Services16     ADType = 0x03 // Complete list of 16-bit UUIDs available.
	MoreServices32 ADType = 0x04 // More 32-bit UUIDs available.
	Services32     ADType = 0x05 // Complete list of 32-bit UUIDs available.
	MoreServices   ADType = 0x06 // More 128-bit UUIDs available.
	Services       ADType = 0x07 // Complete list of 128-bit UUIDs available.

	ShortLocalName ADType = 0x08 // Shortened local name.
	LocalName      ADType = 0x09 // Complete local name.

	TxPower ADType = 0x0A // Tx Power Level: -127 to +127 dBm.

	SlaveConnIntRange ADType = 0x12 // Slave Connection Interval Range.

	ManufSpecData ADType = 0xFF // Manufacturer Specific Data.
)

func (pdu AdvPDU) AppendBytes(typ ADType, s ...byte) {
	n := pdu.length()
	pdu.setLength(n + 2 + len(s))
	pdu.b[n] = byte(1 + len(s))
	pdu.b[n+1] = byte(typ)
	copy(pdu.b[n+2:], s)
}

func (pdu AdvPDU) AppendWords16(typ ADType, s ...uint16) {
	n := pdu.length()
	pdu.setLength(n + 2 + 2*len(s))
	pdu.b[n] = byte(1 + 2*len(s))
	pdu.b[n+1] = byte(typ)
	for i, u := range s {
		pdu.b[n+2+2*i] = byte(u)
		pdu.b[n+3+2*i] = byte(u >> 8)
	}
}

func (pdu AdvPDU) AppendWords32(typ ADType, s ...uint32) {
	n := pdu.length()
	pdu.setLength(n + 2 + 4*len(s))
	pdu.b[n] = byte(1 + 4*len(s))
	pdu.b[n+1] = byte(typ)
	for i, u := range s {
		pdu.b[n+2+4*i] = byte(u)
		pdu.b[n+3+4*i] = byte(u >> 8)
		pdu.b[n+4+4*i] = byte(u >> 16)
		pdu.b[n+5+4*i] = byte(u >> 24)
	}
}

func (pdu AdvPDU) AppendString(typ ADType, s string) {
	n := pdu.length()
	pdu.setLength(n + 2 + len(s))
	pdu.b[n] = byte(1 + len(s))
	pdu.b[n+1] = byte(typ)
	copy(pdu.b[n+2:], s)
}
