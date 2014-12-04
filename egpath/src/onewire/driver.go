package onewire

import "errors"

var (
	ErrNoResponse = errors.New("no response")
	ErrBusFault   = errors.New("bus fault")
	ErrCRC        = errors.New("bad CRC")
)

type Driver interface {
	Reset() error
	SendBit(bit byte) error
	RecvBit() (bit byte, err error)
}
