package onewire

import "errors"

var (
	ErrNoResponse = errors.New("no response")
	ErrBusFault   = errors.New("bus fault")
	ErrCRC        = errors.New("bad CRC")
	ErrDevType    = errors.New("bad device type")
)

type Driver interface {
	// Reset sends reset pulse. If there is no presence pulse received it
	// returns ErrNoResponse error.
	Reset() error

	// ReadBit generates read time slot on the bus. It returns received bit
	// value (0 or 1) or error.
	ReadBit() (bit byte, err error)

	// WriteBit generates write slot on the bus. It sends 0 if bit == 0 or 1
	// otherwise.
	WriteBit(bit byte) error

	// ReadByte receives a byte by generating 8 read slots on the bus. It
	// returns read byte or error.
	ReadByte() (byte, error)

	// WriteByte sends b by generating 8 write slots on the bus.
	WriteByte(b byte) error
}
