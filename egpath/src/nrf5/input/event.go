package input

// Event contains 8-bit identifier of source of event and 24-bit event
// value.
type Event struct {
	v int32
}

// MakeEvent creates new event. Only 24 bits of val are used.
func MakeEvent(src byte, val int) Event {
	return Event{int32(val)<<8 | int32(src)}
}

func (e Event) Src() byte {
	return byte(e.v)
}

func (e Event) Val() int {
	return int(e.v >> 8)
}
