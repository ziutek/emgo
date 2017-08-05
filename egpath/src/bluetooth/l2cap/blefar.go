package l2cap

import (
	"errors"
	"io"

	"bluetooth/ble"
)

// BLEFAR implements Fragmentation And Recombination component of L2CAP for BLE
// controller. Reading and writing operations are mutually independent (can be
// performed concurently).
type BLEFAR struct {
	hci   ble.HCI
	rxpay []byte
	tord  uint16
	towr  uint16
}

func (f *BLEFAR) Init(hci ble.HCI) {
	f.hci = hci
}

var ErrShortFrame = errors.New("l2cap: short frame")

// ReadHeader skips all data until reading next header of L2CAP frame. After
// that Read can be used to read length bytes of payload or ReadHeader can be
// used to skip current payload and read next header.
func (f *BLEFAR) ReadHeader() (cid int, err error) {
	for {
		pdu, err := f.hci.Recv()
		if err != nil {
			return 0, err
		}
		if pdu.Header()&ble.LLID == ble.L2CAPStart {
			if pdu.PayLen() < 4 {
				// BUG? Can header span more than one BLE PDU?
				return 0, ErrShortFrame
			}
			p := pdu.Payload()
			f.tord = uint16(p[0]) | uint16(p[1])<<8
			cid = int(p[2]) | int(p[3])<<8
			f.rxpay = p[4:]
			return cid, nil
		}
	}
}

// Len returns number of bytes of the unread portion of payload.
func (f *BLEFAR) Len() int {
	return int(f.tord)
}

// Read reads payload of L2CAP frame. It always reads full s if there are enough
// unread data in payload. Otherwise Read returns n < len(s) and error, which is
// io.EOF in case of successfull read of whole payload.
func (f *BLEFAR) Read(s []byte) (n int, err error) {
	for {
		if n == len(s) {
			return n, nil
		}
		if f.tord == 0 {
			return n, io.EOF
		}
		if len(f.rxpay) == 0 {
			pdu, err := f.hci.Recv()
			if err != nil {
				return n, err
			}
			if pdu.Header()&ble.LLID != ble.L2CAPCont {
				// BUG? Probably valid L2CAPStart PDU is lost and as a result a
				// whole next frame too.
				return n, ErrShortFrame
			}
			f.rxpay = pdu.Payload()
			if len(f.rxpay) > int(f.tord) {
				f.rxpay = f.rxpay[:f.tord]
			}
		}
		m := copy(s[n:], f.rxpay)
		n += m
		f.tord = uint16(int(f.tord) - m)
		f.rxpay = f.rxpay[m:]
	}
}

// WriteHeader writes header of L2CAP frame to the internal buffer. It panics if
// the previous frame was not completed.
func (f *BLEFAR) WriteHeader(length, cid int) {
	if f.towr != 0 {
		panic("l2cap: buried frame")
	}
	f.towr = uint16(length)
	pdu := f.hci.GetSend()
	pdu.SetHeader(ble.L2CAPStart)
	pdu.SetPayLen(4)
	pay := pdu.Payload()
	pay[0] = byte(length)
	pay[1] = byte(length >> 8)
	pay[2] = byte(cid)
	pay[3] = byte(cid >> 8)
}

// WriteString writes s as payload of current L2CAP frame (the frame started by
// WriteHeader). It returns (0, io.ErrShortWrite) if there is no space for s.
func (f *BLEFAR) WriteString(s string) (n int, err error) {
	if int(f.towr) < len(s) {
		return 0, io.ErrShortWrite
	}
	pdu := f.hci.GetSend()
	for n < len(s) {
		if pdu.PayLen() == pdu.MaxPay() {
			if err = f.hci.Send(); err != nil {
				return n, err
			}
			pdu = f.hci.GetSend()
			pdu.SetHeader(ble.L2CAPCont)
			pdu.SetPayLen(0)
		}
		m := copy(pdu.Remain(), s[n:])
		n += m
		pdu.SetPayLen(pdu.PayLen() + m)
		f.towr = uint16(int(f.towr) - m)
	}
	if f.towr == 0 {
		return n, f.hci.Send()
	}
	return n, nil
}
