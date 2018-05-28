package sdcard

import (
	"errors"
)

// DataMode describes data transfer mode.
type DataMode byte

const (
	Send     DataMode = 0 << 1  // Send data to a card.
	Recv     DataMode = 1 << 1  // Receive data from a card.
	Stream   DataMode = 1 << 2  // Stream or SDIO multibyte data transfer.
	Block1   DataMode = 0 << 4  // Block data transfer, block size: 1 B.
	Block2   DataMode = 1 << 4  // Block data transfer, block size: 2 B.
	Block4   DataMode = 2 << 4  // Block data transfer, block size: 4 B.
	Block8   DataMode = 3 << 4  // Block data transfer, block size: 8 B.
	Block16  DataMode = 4 << 4  // Block data transfer, block size: 16 B.
	Block32  DataMode = 5 << 4  // Block data transfer, block size: 32 B.
	Block62  DataMode = 6 << 4  // Block data transfer, block size: 64 B.
	Block128 DataMode = 7 << 4  // Block data transfer, block size: 128 B.
	Block256 DataMode = 8 << 4  // Block data transfer, block size: 256 B.
	Block512 DataMode = 9 << 4  // Block data transfer, block size: 512 B.
	Block1K  DataMode = 10 << 4 // Block data transfer, block size: 1 KiB.
	Block2K  DataMode = 11 << 4 // Block data transfer, block size: 2 KiB.
	Block4K  DataMode = 12 << 4 // Block data transfer, block size: 4 KiB.
	Block8K  DataMode = 13 << 4 // Block data transfer, block size: 8 KiB.
	Block16K DataMode = 14 << 4 // Block data transfer, block size: 16 KiB.
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
