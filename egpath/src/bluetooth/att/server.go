package att

import (
	"errors"
	"io"

	"bluetooth/l2cap"
)

var (
	ErrBadPDU    = errors.New("att: bad PDU")
	ErrBadMethod = errors.New("att: bad method")
)

type Request struct {
	Method    Method
	Cmd       bool
	Handle    uint16
	EndHandle uint16
	Other     uint16 // Offset, MTU, Flags.
	UUID      UUID

	far *l2cap.BLEFAR
}

func (r *Request) parseExchangeMTU(arg1 uint16, buf *[18]byte, n int) {
	r.Other = arg1 // MTU
}

func (r *Request) parseFindInformation(arg1 uint16, buf *[18]byte, n int) {
	r.Handle = arg1
	r.EndHandle = Decode16(buf[0:2])
}

func (r *Request) parseFindByTypeValue(arg1 uint16, buf *[18]byte, n int) {
	r.Handle = arg1
	r.EndHandle = Decode16(buf[0:2])
	r.UUID = DecodeUUID16(buf[2:4]).Full()
}

func (r *Request) parseReadByType(arg1 uint16, buf *[18]byte, n int) {
	r.Handle = arg1
	r.EndHandle = Decode16(buf[0:2])
	if n == 4 {
		r.UUID = DecodeUUID16(buf[2:4]).Full()
	} else {
		r.UUID = DecodeUUID(buf[2:18])
	}
}

func (r *Request) parseRead(arg1 uint16, buf *[18]byte, n int) {
	r.Handle = arg1
}

func (r *Request) parseReadBlob(arg1 uint16, buf *[18]byte, n int) {
	r.Handle = arg1
	r.Other = Decode16(buf[0:2]) // Value Offset
}

type reqParser struct {
	n1 uint16 // Number of bytes to read.
	n2 uint16 // Alternate number of bytes to read.
	f  func(r *Request, arg1 uint16, buf *[18]byte, n int)
}

//emgo:const
var reqParsers = [...]reqParser{
	ExchangeMTU:     {0, 0, (*Request).parseExchangeMTU},
	FindInformation: {2, 2, (*Request).parseExchangeMTU},
	FindByTypeValue: {4, 0xFFFF, (*Request).parseFindByTypeValue},
	ReadByType:      {18, 4, (*Request).parseReadByType},
	Read:            {0, 0, (*Request).parseRead},
	ReadBlob:        {2, 2, (*Request).parseReadBlob},
}

func (r *Request) readAndParse(length int) error {
	var buf [18]byte
	n, err := r.far.Read(buf[:3])
	if err != nil && err != io.EOF {
		return err
	}
	if n != 3 {
		return ErrBadPDU
	}

	r.Method = Method(buf[0] & 0x3F >> 1)
	r.Cmd = buf[0]&0x40 != 0
	// BUG: Authentication Signature Flag unhandled.
	arg1 := Decode16(buf[1:3])

	if int(r.Method) >= len(reqParsers) {
		return ErrBadMethod
	}
	p := reqParsers[r.Method]
	if p.f == nil {
		return ErrBadMethod
	}

	if p.n1 != 0 {
		n, err = r.far.Read(buf[:p.n1])
		if err != nil && err != io.EOF {
			return err
		}
		if n != int(p.n1) && n != int(p.n2) {
			return ErrBadPDU
		}
	}
	if p.n2 != 0xFFFF && r.far.Len() != 0 {
		return ErrBadPDU
	}

	p.f(r, arg1, &buf, n)
	return nil
}

// Len returns number of bytes of the unparsed portion of request.
func (r *Request) Len() int {
	return r.far.Len()
}

// Read can be used to read the unparsed portion of request (eg: attribute value
// in case of Find By Type Value request).
func (r *Request) Read(s []byte) (int, error) {
	return r.far.Read(s)
}

type Handler interface {
	ServeATT(*Request)
}
