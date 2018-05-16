package sdcard

import (
	"errors"
)

var ErrTimeout = errors.New("sdio: timeout")

type Host interface {
	// Cmd sends a command to the card and receives its response, if any.
	Cmd(cmd Command, arg uint32) Response

	// Err returns and clears internal error status. The internal error status,
	// if set,  prevents any subsequent operations on the card. the ErrTimeout
	// should be returned if the internal error represents the lack of reponse
	// from the card.
	Err(clear bool) error

	// Timeout check if the internal error status is set and represents the lack
	// of reponse from the card. In such case it clears the error and returns
	// true.
	Timeout() bool
}
