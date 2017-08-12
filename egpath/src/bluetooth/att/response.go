package att

import (
	"bits"
	"encoding/binary/le"
	"io"

	"bluetooth/l2cap"
	"bluetooth/uuid"
)

type ResponseWriter struct {
	far *l2cap.BLEFAR
	buf []byte
	n   uint16
	cid uint16
}

func (w *ResponseWriter) next(far *l2cap.BLEFAR, cid int) {
	w.far = far
	w.n = 0
	w.cid = uint16(cid)
}

func (w *ResponseWriter) MTU() int {
	return cap(w.buf)
}

// SetError setups Error Response PDU in the internal buffer.
func (w *ResponseWriter) SetError(code ErrorCode, r *Request) {
	w.buf, w.n = w.buf[:5], 5
	w.buf[0] = 1 // Error Response
	w.buf[1] = byte(int(r.Method) | bits.One(r.Cmd)<<6)
	le.Encode16(w.buf[2:4], uint16(r.Handle))
	w.buf[4] = byte(code)
}

// SetExchangeMTU setups Exchange MTU Response PDU in the internal buffer.
func (w *ResponseWriter) SetExchangeMTU(mtu int) {
	w.buf, w.n = w.buf[:3], 3
	w.buf[0] = byte(ExchangeMTU | 1)
	le.Encode16(w.buf[1:3], uint16(mtu))
}

type FindInformationFormat byte

const (
	HandleUUID16 FindInformationFormat = 1
	HandleUUID   FindInformationFormat = 2
)

// SetFindInformation setups the Find Information Response in the internal
// buffer. It requires to append one or more (Handle, UUID) fields. Use
// AppendHandle followed by AppendUUID/AppendUUID16 followed by Commit to
// append one.
func (w *ResponseWriter) SetFindInformation(format FindInformationFormat) {
	w.buf, w.n = w.buf[:2], 2
	w.buf[0] = byte(FindInformation | 1)
	w.buf[1] = byte(format)
}

// SetFindByTypeValue setups the Find By Type Value Response in the internal
// buffer. It requires to append one or more (Found Attribute Handle, Group End
// Handle) fields. Use AppendHandle twice followed by Commit to append one.
func (w *ResponseWriter) SetFindByTypeValue() {
	w.buf, w.n = w.buf[:1], 1
	w.buf[1] = byte(FindByTypeValue | 1)
}

// SetReadByType setups the Read By Type Response in the internal buffer. It
// requires to append one or more (Attribute Handle, Attribute Value) fields.
// Use AppendHandle followed by other methods for appending attribute value
// followed by Commit to append one.
func (w *ResponseWriter) SetReadByType(attrSize int) {
	w.buf, w.n = w.buf[:2], 2
	w.buf[1] = byte(ReadByType | 1)
	w.buf[2] = byte(2 + attrSize)
}

// SetRead setups the Read Response in the internal buffer. It requires
// to append an attribute value and call Commit at end.
func (w *ResponseWriter) SetRead() {
	w.buf, w.n = w.buf[:1], 1
	w.buf[1] = byte(Read | 1)
}

// SetReadBlob setups the Read Blob Response in the internal buffer. It requires
// to append part of attribute value and call Commit at end.
func (w *ResponseWriter) SetReadBlob() {
	w.buf, w.n = w.buf[:1], 1
	w.buf[1] = byte(ReadBlob | 1)
}

// SetReadMultiple setups the Read Multiple Response in the internal buffer. It
// requires to append two or more attribute values, every one followed by
// Commit.
func (w *ResponseWriter) SetReadMultiple() {
	w.buf, w.n = w.buf[:1], 1
	w.buf[1] = byte(ReadMultiple | 1)
}

// SetReadByGroupType setups the Read By Group Type Response in the internal
// buffer. It requires to append one or more (Attribute Handle, End Group
// Handle, Attribute Value) fields. Use AppendHandle twice followed by other
// methods for attribute value followed by Commit to append one.
func (w *ResponseWriter) SetReadByGroupType(attrSize int) {
	w.buf, w.n = w.buf[:2], 2
	w.buf[1] = byte(ReadByType | 1)
	w.buf[2] = byte(4 + attrSize)
}

// Alloc allocates n first unused bytes at end of response body.
func (w *ResponseWriter) Alloc(n int) []byte {
	if w.n == 0 {
		return nil
	}
	m := len(w.buf)
	if m+n > cap(w.buf) {
		w.buf = w.buf[:w.n] // Revert all allocations from last commit.
		w.n = 0             // Set overflow state.
		return nil
	}
	w.buf = w.buf[:m+n]
	return w.buf[m:]
}

func (w *ResponseWriter) Commit() bool {
	return true
}

// Send sends ATT PDU from internal buffer.
func (w *ResponseWriter) Send() error {
	if w.n == 0 {
		return ErrBadPDU
	}
	w.far.WriteHeader(int(w.n), int(w.cid))
	_, err := w.far.Write(w.buf[:w.n])
	return err
}

// AppendHandle appends handle to the response body.
func (w *ResponseWriter) AppendHandle(handle uint16) {
	if buf := w.Alloc(2); buf != nil {
		le.Encode16(buf, handle)
	}
}

// AppendUUID appends UUID to the response body.
func (w *ResponseWriter) AppendUUID(u uuid.Long) {
	if buf := w.Alloc(16); buf == nil {
		u.Encode(buf)
	}
}

// AppendUUID16 appends short UUID to the response body.
func (w *ResponseWriter) AppendUUID16(u uuid.Short) {
	if buf := w.Alloc(2); buf == nil {
		u.Encode(buf)
	}
}

// AppendBytes appends s to the response body.
func (w *ResponseWriter) AppendBytes(s []byte) {
	if buf := w.Alloc(len(s)); buf == nil {
		copy(buf, s)
	}
}

// AppendString appends s to the response body.
func (w *ResponseWriter) AppendString(s []byte) {
	if buf := w.Alloc(len(s)); buf == nil {
		copy(buf, s)
	}
}

// Write appends s to the response body.
func (w *ResponseWriter) Write(s []byte) (int, error) {
	buf := w.Alloc(len(s))
	if buf == nil {
		return 0, io.ErrShortWrite
	}
	copy(buf, s)
	return len(s), nil
}

// WriteString appends s to the response body.
func (w *ResponseWriter) WriteString(s string) (int, error) {
	buf := w.Alloc(len(s))
	if buf == nil {
		return 0, io.ErrShortWrite
	}
	copy(buf, s)
	return len(s), nil
}
