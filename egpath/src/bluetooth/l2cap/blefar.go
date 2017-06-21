package l2cap

import (
	"errors"

	"bluetooth/ble"
)

// BLEFAR implements Fragmentation And Recombination component of L2CAP for BLE
// controller.
type BLEFAR struct {
	hci  ble.HCI
	rpay []byte
	tord uint16
}

func (far *BLEFAR) SetHCI(hci ble.HCI) {
	far.hci = hci
}

var ErrShortFrame = errors.New("l2cap: short frame")

// Skips all data until reading next frame header. After reading frame header.
// Read method can be used to read length bytes of payload.
func (far *BLEFAR) ReadHeader() (length, cid int, err error) {
	for {
		pdu, err := far.hci.Recv()
		if err != nil {
			return 0, 0, err
		}
		if pdu.Header()&ble.LLID == ble.L2CAPStart {
			if pdu.PayLen() < 4 {
				// BUG? Can header span more than one BLE PDU?
				return 0, 0, ErrShortFrame
			}
			p := pdu.Payload()
			length = int(p[0]) | int(p[1])<<8
			cid = int(p[2]) | int(p[3])<<8
			far.rpay = p[4:]
			far.tord = uint16(length)
			return length, cid, nil
		}
	}
}

// Read can be used to read payload of L2CAP frame.
func (far *BLEFAR) Read(s []byte) (n int, err error) {
	for n < len(s) && far.tord > 0 {
		if len(far.rpay) == 0 {
			pdu, err := far.hci.Recv()
			if err != nil {
				return n, err
			}
			if pdu.Header()&ble.LLID != ble.L2CAPCont {
				// BUG? Probably valid L2CAPStart PDU is lost and as a result a
				// whole next frame too.
				return n, ErrShortFrame
			}
			far.rpay = pdu.Payload()
			if len(far.rpay) > int(far.tord) {
				far.rpay = far.rpay[:far.tord]
			}
		}
		m := copy(s[n:], far.rpay)
		n += m
		far.tord = uint16(int(far.tord) - m)
		far.rpay = far.rpay[m:]
	}
	return
}

// WriteHeader writes frame header.
func (far *BLEFAR) WriteHeader(length, cid int) error {
	return nil
}
