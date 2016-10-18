package system

type Bus int8

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

//emgo:const
var busStr = [...]string{
	Core: "Core",
	AHB:  "AHB",
	APB1: "APB1",
	APB2: "APB2",
}

func (b Bus) String() string {
	if uint(b) >= uint(busN) {
		return ""
	}
	return busStr[b]
}
