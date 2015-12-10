package delay

import (
	"sync"
)

// Loop can be used to perform short active delay.
func Loop(n int) {
	for n > 0 {
		n--
		sync.Fence()
	}
}

// Millisec can by used to perform delays of the order from few milliseconds
// to houres or days. For small values it can be very inaccurate. This function
// need some support from runtime or OS and can panic if there is no such
// support (eg: noos target and MaxTasks == 0).
func Millisec(ms int) {
	millisec(ms)
}
