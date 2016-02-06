package system

type Bus int

const (
	Core Bus = iota
	AHB
	APB1
	APB2
	busN
)

var clock [busN]uint

// Clock returns clock frequency [Hz] for bus.
func (b Bus) Clock() uint {
	return clock[b]
}
