package ble

import "unsafe"

const (
	MaxAdvPay      = 37  // Max. advertising payload length in bytes.
	MaxDataPay     = 31  // Max. data payload length in bytes (BLE 4.0, 4.1).
	MaxLongDataPay = 255 // Max. data payload length in bytes (BLE 4.2+).
)

type pduBytes [MaxLongDataPay + 2]byte

// PDU represents BLE Protocol Data Unit. It has reference semantic (like
// slice), so it should be initialized before use.
type PDU struct {
	b      *pduBytes
	maxpay int
}

func checkPayLen(n, max int) {
	if uint(n) > uint(max) {
		panic("ble: bad payload len.")
	}
}

// MakePDU returns ready to use PDU of capacity maxpay.
func MakePDU(maxpay int) PDU {
	checkPayLen(maxpay, MaxLongDataPay)
	pdu := make([]byte, maxpay+2)
	return PDU{(*pduBytes)(unsafe.Pointer(&pdu[0])), maxpay}
}

// IsZero reports whether value of pdu is zero.
func (pdu PDU) IsZero() bool {
	return pdu.b == nil
}

// MaxPay returns number of bytes that can be stored in payload part of pdu.
func (pdu PDU) MaxPay() int {
	return pdu.maxpay
}

// PayLen returns payload length.
func (pdu PDU) PayLen() int {
	return int(pdu.b[1])
}

// SetPayLen sets payload length to n.
func (pdu PDU) SetPayLen(n int) {
	checkPayLen(n, pdu.maxpay)
	pdu.b[1] = byte(n)
}

// Payload returns used payload part of PDU.
func (pdu PDU) Payload() []byte {
	return pdu.b[2 : 2+pdu.PayLen()]
}

// Remain returns unused payload part of PDU.
func (pdu PDU) Remain() []byte {
	return pdu.b[2+pdu.PayLen() : 2+pdu.maxpay]
}

// Bytes returns underlying bytes of header and payload part of pdu.
func (pdu PDU) Bytes() []byte {
	return pdu.b[:2+pdu.PayLen()]
}
