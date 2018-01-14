package input

import (
	"nrf5/hal/gpio"
)

// Event contains 32-bit identifier of source of event and 32-bit event
// value.
type Event struct {
	src uint32
	val int32
}

// MakeIntEvent creates new event with int value.
func MakeIntEvent(src uint32, val int) Event {
	return Event{src, int32(val)}
}

// MakePinEvent creates new event with gpio.Pins value.
func MakePinEvent(src uint32, pins gpio.Pins) Event {
	return Event{src, int32(pins)}
}

func (e Event) Src() uint32 {
	return e.src
}

func (e Event) Int() int {
	return int(e.val)
}

func (e Event) Pins() gpio.Pins {
	return gpio.Pins(e.val)
}
