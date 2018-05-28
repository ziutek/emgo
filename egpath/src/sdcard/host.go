package sdcard

import (
	"errors"
)

// DataMode describes data transfer mode.
type DataMode byte

const (
	Send   DataMode = 0 << 1 // Send data to a card.
	Recv   DataMode = 1 << 1 // Receive data from a card.
	Block  DataMode = 0 << 2 // Block data transfer.
	Stream DataMode = 1 << 2 // Stream or SDIO multibyte data transfer.
)

var ErrCmdTimeout = errors.New("sdio: cmd timeout")

type Host interface {
	// SetFreq sets the SD/SPI clock frequency to freqhz. Host can implement
	// disabling clock output if the bus is idle and pwrsave is set to true.
	SetFreq(freqhz int, pwrsave bool)

	// SetBusWidth allow to change the the host data bus width.
	SetBusWidth(width int)

	// Cmd sends the cmd to the card and receives its response, if any. Short
	// response is returned in r[0]. Long is returned in r[0:3] (r[0] contains
	// the least significant bits, r[3] contains the most significant bits).
	// If preceded by Data, Cmd performs data transfer.
	Cmd(cmd Command, arg uint32) (r Response)

	// Data prepares the data transfer for subsequent command.
	Data(mode DataMode, buf Data)

	// Err returns and clears the host internal error. The internal error, if
	// not nil, prevents any subsequent operations on the card. Host should
	// convert its internal command timeout error to ErrCmdTimeout.
	Err(clear bool) error
}
