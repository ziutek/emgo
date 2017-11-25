package eve

import (
	"io"
)

type EVE struct {
	DCI DCI
}

type HostCmd byte

// Cmd invokes host command. Arg is command argument. It must be zero for
// commands that do not require arguments..
func (e *EVE) Cmd(cmd HostCmd, arg byte) {
	dci := e.DCI
	dci.Begin()
	dci.Write([]byte{byte(cmd), arg, 0})
	dci.End()
}

func checkAddr(addr int) {
	if uint(addr) >= 1<<22 {
		panic("eve: bad addr")
	}
}

type writeCloser struct {
	e *EVE
}

func (wc writeCloser) Write(s []byte) (int, error) {
	dci := wc.e.DCI
	dci.Write(s)
	if err := dci.Err(); err != nil {
		return 0, err
	}
	return len(s), nil
}

func (wc writeCloser) WriteString(s string) (int, error) {
	dci := wc.e.DCI
	dci.WriteString(s)
	if err := dci.Err(); err != nil {
		return 0, err
	}
	return len(s), nil
}

func (wc writeCloser) Close() error {
	dci := wc.e.DCI
	dci.End()
	return dci.Err()
}

func (e *EVE) StartWrite(addr int) io.WriteCloser {
	checkAddr(addr)
	dci := e.DCI
	dci.Begin()
	dci.Write([]byte{1<<7 | byte(addr>>16), byte(addr >> 8), byte(addr)})
	return writeCloser{e}
}

func (e *EVE) Write(addr int, s []byte) {
	checkAddr(addr)
	dci := e.DCI
	dci.Begin()
	dci.Write([]byte{1<<7 | byte(addr>>16), byte(addr >> 8), byte(addr)})
	dci.Write(s)
	dci.End()
}

type readCloser struct {
	e *EVE
}

func (rc readCloser) Read(s []byte) (int, error) {
	dci := rc.e.DCI
	return dci.Read(s), dci.Err()
}

func (rc readCloser) Close() error {
	dci := rc.e.DCI
	dci.End()
	return dci.Err()
}

func (e *EVE) StartRead(addr int) io.ReadCloser {
	checkAddr(addr)
	dci := e.DCI
	dci.Begin()
	buf := []byte{byte(addr >> 16), byte(addr >> 8), byte(addr)}
	dci.Write(buf)
	dci.Read(buf[:1]) // Switch to input for dummy byte to support FT81x QSPI.
	return readCloser{e}
}

func (e *EVE) Read(addr int, s []byte) {
	checkAddr(addr)
	dci := e.DCI
	dci.Begin()
	buf := []byte{byte(addr >> 16), byte(addr >> 8), byte(addr)}
	dci.Write(buf)
	dci.Read(buf[:1]) // Switch to input for dummy byte to support FT81x QSPI.
	dci.Read(s)
	dci.End()
}

func (e *EVE) ReadByte(addr int) byte {
	checkAddr(addr)
	dci := e.DCI
	dci.Begin()
	buf := []byte{byte(addr >> 16), byte(addr >> 8), byte(addr)}
	dci.Write(buf)
	dci.Read(buf[:2])
	dci.End()
	return buf[1]
}

func (e *EVE) ReadWord16(addr int) uint16 {
	checkAddr(addr)
	dci := e.DCI
	dci.Begin()
	buf := []byte{byte(addr >> 16), byte(addr >> 8), byte(addr)}
	dci.Write(buf)
	dci.Read(buf)
	dci.End()
	return uint16(buf[1]) | uint16(buf[2])<<8
}

func (e *EVE) ReadWord32(addr int) uint32 {
	checkAddr(addr)
	dci := e.DCI
	dci.Begin()
	buf := []byte{byte(addr >> 16), byte(addr >> 8), byte(addr), 0, 0}
	dci.Write(buf[:3])
	dci.Read(buf)
	dci.End()
	return uint32(buf[1]) | uint32(buf[2])<<8 | uint32(buf[3])<<16 |
		uint32(buf[4])<<24
}
