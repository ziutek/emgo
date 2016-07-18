package onewire

import "errors"

var (
	ErrNoResponse = errors.New("no response")
	ErrBusFault   = errors.New("bus fault")
	ErrCRC        = errors.New("bad CRC")
	ErrDevType    = errors.New("bad device type")
)

type Driver interface {
	Reset() error
	ReadBit() (bit byte, err error)
	WriteBit(bit byte) error
	ReadByte() (byte, error)
	WriteByte(b byte) error
}
