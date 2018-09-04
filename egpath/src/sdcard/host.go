package sdcard

import (
	"errors"
)

// BusWidth describes SD data bus width.
type BusWidth byte

const (
	Bus1 BusWidth = 0 // 1-lane SD data bus.
	Bus4 BusWidth = 2 // 4-lane SD data bus.
	Bus8 BusWidth = 3 // 8-lane SD data bus.
)

// DataMode describes data transfer mode.
type DataMode uint16

// All DataMode constants are defined to be friendly to use with ARM PrimeCell
// Multimedia Card Interface (used by STM32, LPC and probably more MCUs). Do
// not add, delete, modify these constants without checking the stm32/hal/sdmmc
// and lpc/hal/sdmmc.
const (
	Send       DataMode = 0 << 1  // Send data to a card.
	Recv       DataMode = 1 << 1  // Receive data from a card.
	Stream     DataMode = 1 << 2  // Stream or SDIO multibyte data transfer.
	Block1     DataMode = 0 << 4  // Block data transfer, block size: 1 B.
	Block2     DataMode = 1 << 4  // Block data transfer, block size: 2 B.
	Block4     DataMode = 2 << 4  // Block data transfer, block size: 4 B.
	Block8     DataMode = 3 << 4  // Block data transfer, block size: 8 B.
	Block16    DataMode = 4 << 4  // Block data transfer, block size: 16 B.
	Block32    DataMode = 5 << 4  // Block data transfer, block size: 32 B.
	Block64    DataMode = 6 << 4  // Block data transfer, block size: 64 B.
	Block128   DataMode = 7 << 4  // Block data transfer, block size: 128 B.
	Block256   DataMode = 8 << 4  // Block data transfer, block size: 256 B.
	Block512   DataMode = 9 << 4  // Block data transfer, block size: 512 B.
	Block1K    DataMode = 10 << 4 // Block data transfer, block size: 1 KiB.
	Block2K    DataMode = 11 << 4 // Block data transfer, block size: 2 KiB.
	Block4K    DataMode = 12 << 4 // Block data transfer, block size: 4 KiB.
	Block8K    DataMode = 13 << 4 // Block data transfer, block size: 8 KiB.
	Block16K   DataMode = 14 << 4 // Block data transfer, block size: 16 KiB.
	RWaitStart DataMode = 1 << 8  // Read wait start.
	RWaitStop  DataMode = 1 << 9  // Read wait stop.
	RWaitCK    DataMode = 1 << 10 // Read wait control using CK instead od D2.
)

// ErrCmdTimeout is returned by Host in case of command response timeout.
var ErrCmdTimeout = errors.New("sdio: cmd timeout")

type Host interface {
	// SetClock sets SD/SPI clock frequency. If pwrsave is true the host can
	// disable clock in case of inactive bus.
	SetClock(freqhz int, pwrsave bool)

	// SetBusWidth sets the SD bus width. Returns supported bus widths. SD host
	// returns combination of SDBus1 and SDBus4. SPI host returns SPIBus.
	SetBusWidth(width BusWidth) BusWidths

	// SendCmd sends the cmd to the card and receives its response, if any.
	// Short response is returned in r[0], long is returned in r[0:3] (r[0]
	// contains the least significant bits, r[3] contains the most significant
	// bits). If preceded by SetupData, SendCmd performs the data transfer.
	SendCmd(cmd Command, arg uint32) (r Response)

	// SetupData setups the data transfer for subsequent command.
	SetupData(mode DataMode, buf []uint64, nbytes int)

	// Wait waits for deassertion of busy signal on DATA0 line (READY_FOR_DATA
	// state). It returns false if the deadline has passed. Wait can not be used
	// while transfer is in progress.
	Wait(deadline int64) bool

	// Err returns and clears the host internal error. The internal error, if
	// not nil, prevents any subsequent operations. Host should convert its
	// internal command timeout error to ErrCmdTimeout.
	Err(clear bool) error
}
