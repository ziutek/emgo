package ble

// AdvPDU represents Advertising Channel PDU. See PDU for more information.
type AdvPDU struct {
	PDU
}

// MakeAdvPDU is more readable counterpart of AdvPDU{MakePDU(maxpay)}.
func MakeAdvPDU(maxpay int) AdvPDU {
	return AdvPDU{MakePDU(maxpay)}
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

func (pdu AdvPDU) UpdateAddr(offset int, addr int64) {
	data := pdu.Payload()[offset:]
	data[0] = byte(addr)
	data[1] = byte(addr >> 8)
	data[2] = byte(addr >> 16)
	data[3] = byte(addr >> 24)
	data[4] = byte(addr >> 32)
	data[5] = byte(addr >> 40)
}

func (pdu AdvPDU) AppendAddr(addr int64) (offset int) {
	offset = pdu.PayLen()
	pdu.SetPayLen(offset + 6)
	pdu.UpdateAddr(offset, addr)
	return
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
	r := pdu.Remain()
	pdu.SetPayLen(pdu.PayLen() + 2 + len(s))
	r[0] = byte(1 + len(s))
	r[1] = byte(typ)
	copy(r[2:], s)
}

func (pdu AdvPDU) AppendString(typ ADType, s string) {
	r := pdu.Remain()
	pdu.SetPayLen(pdu.PayLen() + 2 + len(s))
	r[0] = byte(1 + len(s))
	r[1] = byte(typ)
	copy(r[2:], s)
}

func (pdu AdvPDU) AppendWords16(typ ADType, s ...uint16) {
	r := pdu.Remain()
	pdu.SetPayLen(pdu.PayLen() + 2 + len(s)*2)
	r[0] = byte(1 + len(s)*2)
	r[1] = byte(typ)
	for i, u := range s {
		r[2+i*2] = byte(u)
		r[3+i*2] = byte(u >> 8)
	}
}

func (pdu AdvPDU) AppendWords32(typ ADType, s ...uint32) {
	r := pdu.Remain()
	pdu.SetPayLen(pdu.PayLen() + 2 + len(s)*4)
	r[0] = byte(1 + len(s)*4)
	r[1] = byte(typ)
	for i, u := range s {
		r[2+i*4] = byte(u)
		r[3+i*4] = byte(u >> 8)
		r[4+i*4] = byte(u >> 16)
		r[5+i*4] = byte(u >> 24)
	}
}
