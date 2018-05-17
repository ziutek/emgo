package sdcard

import (
	"errors"
)

var ErrTimeout = errors.New("sdio: timeout")

type Host interface {
	// Cmd sends a command to the card and receives its response, if any. Short
	// response is returned in r[0]. Long is returned in r[0:3] (r[0] contains
	// the least significant bits, r[3] contains the most significant bits).
	Cmd(cmd Command, arg uint32) (r Response)

	// Err returns and clears internal error status. The internal error status,
	// if set,  prevents any subsequent operations on the card. the ErrTimeout
	// should be returned if the internal error represents the lack of reponse
	// from the card.
	Err(clear bool) error
}
