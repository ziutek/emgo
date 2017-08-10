package att

import (
	"bits"
	"io"

	"bluetooth/l2cap"
)

type ResponseWriter struct {
	buf []byte
	far *l2cap.BLEFAR
	cid int
}

func (w *ResponseWriter) MTU() int {
	return cap(w.buf)
}

// Send sends ATT PDU from internal buffer.
func (w *ResponseWriter) Send() error {
	w.far.WriteHeader(len(w.buf), w.cid)
	_, err := w.far.Write(w.buf)
	w.buf = w.buf[:0]
	return err
}

// SetError setups Error Response PDU in the internal buffer.
func (w *ResponseWriter) SetError(code ErrorCode, r *Request) {
	w.buf = w.buf[:5]
	w.buf[0] = 1 // Error Response
	w.buf[1] = byte(int(r.Method) | bits.One(r.Cmd)<<6)
	Encode16(w.buf[2:4], uint16(r.Handle))
	w.buf[4] = byte(code)
}

// SetExchangeMTU setups Exchange MTU Response PDU in the internal buffer.
func (w *ResponseWriter) SetExchangeMTU(mtu int) {
	w.buf = w.buf[:3]
	w.buf[0] = byte(ExchangeMTU | 1)
	Encode16(w.buf[1:3], uint16(mtu))
}

type FindInformationFormat byte

const (
	HandleUUID16 FindInformationFormat = 1
	HandleUUID   FindInformationFormat = 2
)

// SetFindInformation setups the Find Information Response in the internal
// buffer.
func (w *ResponseWriter) SetFindInformation(format FindInformationFormat) {
	w.buf = w.buf[:2]
	w.buf[0] = byte(FindInformation | 1)
	w.buf[1] = byte(format)
}

// SetFindByTypeValue setups the Find By Type Value Response in the internal
// buffer.
func (w *ResponseWriter) SetFindByTypeValue() {
	w.buf = w.buf[:1]
	w.buf[1] = byte(FindByTypeValue | 1)
}

// SetReadByType setups the Set Read By Type Response in the internal buffer.
func (w *ResponseWriter) SetReadByType(attrSize int) {
	w.buf = w.buf[:2]
	w.buf[1] = byte(ReadByType | 1)
	w.buf[2] = byte(2 + attrSize)
}

func (w *ResponseWriter) Alloc(n int) []byte {
	m := len(w.buf)
	if m+n > cap(w.buf) {
		return nil
	}
	w.buf = w.buf[:m+n]
	return w.buf[m:]
}

func (w *ResponseWriter) AppendHandle(handle uint16) bool {
	buf := w.Alloc(2)
	if buf == nil {
		return false
	}
	Encode16(buf, handle)
	return true
}

func (w *ResponseWriter) AppendUUID(uuid UUID) bool {
	buf := w.Alloc(16)
	if buf == nil {
		return false
	}
	uuid.Encode(buf)
	return true
}

func (w *ResponseWriter) AppendUUID16(uuid UUID16) bool {
	buf := w.Alloc(2)
	if buf == nil {
		return false
	}
	uuid.Encode(buf)
	return true
}

func (w *ResponseWriter) AppendBytes(s []byte) bool {
	buf := w.Alloc(len(s))
	if buf == nil {
		return false
	}
	copy(buf, s)
	return true
}

func (w *ResponseWriter) AppendString(s string) bool {
	buf := w.Alloc(len(s))
	if buf == nil {
		return false
	}
	copy(buf, s)
	return true
}

// Write wraps AppendBytes to implement io.Writer.
func (w *ResponseWriter) Write(s []byte) (int, error) {
	if w.AppendBytes(s) {
		return len(s), nil
	}
	return 0, io.ErrShortWrite
}

// WriteString wraps AppendString.
func (w *ResponseWriter) WriteString(s string) (int, error) {
	if w.AppendString(s) {
		return len(s), nil
	}
	return 0, io.ErrShortWrite
}
