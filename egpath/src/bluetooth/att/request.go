package att

import (
	"encoding/binary/le"
	"errors"
	"io"

	"bluetooth/l2cap"
	"bluetooth/uuid"
)

var (
	ErrBadPDU    = errors.New("att: bad PDU")
	ErrBadMethod = errors.New("att: bad method")
)

type Request struct {
	Method    Method
	Cmd       bool
	Handle    uint16 // Handle, Start Handle.
	EndHandle uint16
	Other     uint16    // Offset, MTU, Flags.
	UUID      uuid.Long // Attribute type.

	far *l2cap.BLEFAR
}

// BUG: Authentication Signature unhandled.
func (r *Request) readAndParse(far *l2cap.BLEFAR) error {
	var buf [18]byte
	r.far = far
	n, err := r.Read(buf[:3])
	if err != nil && err != io.EOF {
		return err
	}
	if n < 2 {
		return ErrBadPDU
	}
	r.Method = Method(buf[0] & 0x3F)
	r.Cmd = buf[0]&0x40 != 0
	r.Handle = 0
	r.EndHandle = 0
	r.Other = 0
	r.UUID = uuid.Long{}
	if n == 2 {
		if r.Method != ExecuteWrite {
			return ErrBadPDU
		}
		r.Other = uint16(buf[1]) // Flags
		return nil
	}
	m := int(r.Method)>>1 - 1
	if r.Method == unusedMethod || m >= len(reqDecoders) {
		return ErrBadMethod
	}
	r.Handle = le.Decode16(buf[1:3])
	dec := reqDecoders[m]
	if dec.n1 != 0 {
		n, err = r.Read(buf[:dec.n1])
		if err != nil && err != io.EOF {
			return err
		}
		if n != int(dec.n1) && n != int(dec.n2) {
			return ErrBadPDU
		}
	}
	if dec.n2 >= 0 && r.far.Len() != 0 {
		return ErrBadPDU
	}
	if dec.f != nil {
		dec.f(r, &buf, n)
	}
	return nil
}

func (r *Request) decodeMTU(buf *[18]byte, n int) {
	r.Other = r.Handle // MTU
	r.Handle = 0
}

func (r *Request) decodeEndHandle(buf *[18]byte, n int) {
	r.EndHandle = le.Decode16(buf[0:2])
}

func (r *Request) decodeEndHandleUUID(buf *[18]byte, n int) {
	r.EndHandle = le.Decode16(buf[0:2])
	if n == 4 {
		r.UUID = uuid.DecodeShort(buf[2:4]).Long()
	} else {
		r.UUID = uuid.DecodeLong(buf[2:18])
	}
}

func (r *Request) decodeOffset(buf *[18]byte, n int) {
	r.Other = le.Decode16(buf[0:2]) // Value Offset
}

type reqDecder struct {
	n1 int8 // Number of bytes to read.
	n2 int8 // Alternate number of bytes to read.
	f  func(r *Request, buf *[18]byte, n int)
}

//emgo:const
var reqDecoders = [...]reqDecder{
	ExchangeMTU>>1 - 1:     {0, 0, (*Request).decodeMTU},
	FindInformation>>1 - 1: {2, 2, (*Request).decodeEndHandle},
	FindByTypeValue>>1 - 1: {4, -1, (*Request).decodeEndHandleUUID},
	ReadByType>>1 - 1:      {18, 4, (*Request).decodeEndHandleUUID},
	Read>>1 - 1:            {0, 0, nil},
	ReadBlob>>1 - 1:        {2, 2, (*Request).decodeOffset},
	ReadMultiple>>1 - 1:    {0, 0, nil},
	ReadByGroupType>>1 - 1: {18, 4, (*Request).decodeEndHandleUUID},
	Write>>1 - 1:           {0, -1, nil},
	PrepareWrite>>1 - 1:    {2, -1, (*Request).decodeOffset},
}

// Len returns number of bytes of the unparsed portion of request.
func (r *Request) Len() int {
	return r.far.Len()
}

// Read can be used to read the unparsed portion of request (eg: attribute value
// in case of Find By Type Value, Write, PrepareWrite requests).
func (r *Request) Read(s []byte) (int, error) {
	return r.far.Read(s)
}

// ReadNextHandle can be used to read more attribute handles in case of
// ReadMultiple request.
func (r *Request) ReadNextHandle() error {
	var buf [2]byte
	_, err := r.Read(buf[:])
	r.Handle = le.Decode16(buf[:])
	return err
}
