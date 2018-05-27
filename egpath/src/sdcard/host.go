package sdcard

import (
	"errors"
)

var ErrTimeout = errors.New("sdio: timeout")

type Host interface {
	// SetFreq sets the CLK/SCLK clock frequency to freqhz. Host can implement
	// disabling clock output if the bus is idle and pwrsave is set to true.
	SetFreq(freqhz int, pwrsave bool)

	// SetBusWidth allow to change the the host data bus width.
	SetBusWidth(width int)

	// Cmd sends the cmd to the card and receives its response, if any. Short
	// response is returned in r[0]. Long is returned in r[0:3] (r[0] contains
	// the least significant bits, r[3] contains the most significant bits).
	Cmd(cmd Command, arg uint32) (r Response)

	// Err returns and clears the host internal error. The internal error, if
	// not nil, prevents any subsequent operations on the card.
	Err(clear bool) error
}
