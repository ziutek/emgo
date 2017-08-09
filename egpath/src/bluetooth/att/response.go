package att

import (
	"bits"

	"bluetooth/l2cap"
)

type ResponseWriter struct {
	far *l2cap.BLEFAR
	cid int
}

func (w *ResponseWriter) writePDU(pdu []byte) error {
	w.far.WriteHeader(len(pdu), w.cid)
	_, err := w.far.Write(pdu)
	return err
}

// WriteError writes Error Response for requesr r.
func (w *ResponseWriter) WriteError(code ErrorCode, r *Request) error {
	var pdu [5]byte
	pdu[0] = 1 // Error Response
	pdu[1] = byte(int(r.Method) | bits.One(r.Cmd)<<6)
	Encode16(pdu[2:4], uint16(r.Handle))
	pdu[4] = byte(code)
	return w.writePDU(pdu[:])
}

//  WriteExchangeMTU writes Exchange MTU Response.
func (w *ResponseWriter) WriteExchangeMTU(mtu int) error {
	var pdu [3]byte
	pdu[0] = byte(ExchangeMTU | 1)
	Encode16(pdu[1:3], uint16(mtu))
	return w.writePDU(pdu[:])
}

type HU struct {
	Handle uint16
	UUID   UUID
}

// WriteFindInformation writes FindInformation Response that contains
// handle-UUID pairs.
func (w *ResponseWriter) WriteFindInformation(data []HU) error {
	var buf [2 + 16]byte
	buf[0] = byte(FindInformation | 1)
	buf[1] = 2 // 128-bit UUIDs
	w.far.WriteHeader(2+len(data)*len(buf), w.cid)
	if _, err := w.far.Write(buf[:2]); err != nil {
		return err
	}
	for _, hu := range data {
		Encode16(buf[:2], hu.Handle)
		hu.UUID.Encode(buf[2:])
		if _, err := w.far.Write(buf[:]); err != nil {
			return err
		}
	}
	return nil
}

type HU16 struct {
	Handle uint16
	UUID   UUID16
}

// WriteFindInformationShort writes FindInformation Response that contains
// handle-UUID16 pairs.
func (w *ResponseWriter) WriteFindInformationShort(data []HU16) error {
	var buf [2 + 2]byte
	buf[0] = byte(FindInformation | 1)
	buf[1] = 1 // 16-bit UUIDs
	w.far.WriteHeader(2+len(data)*len(buf), w.cid)
	if _, err := w.far.Write(buf[:2]); err != nil {
		return err
	}
	for _, hu := range data {
		Encode16(buf[:2], hu.Handle)
		hu.UUID.Encode(buf[2:])
		if _, err := w.far.Write(buf[:]); err != nil {
			return err
		}
	}
	return nil
}
