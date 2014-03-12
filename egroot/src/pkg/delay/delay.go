package delay

import "sync/barrier"


// Loop can be used to perform short active delay.
func Loop(n int) {
	for n > 0 {
		n--
		barrier.Compiler()
	}
}
