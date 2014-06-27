// +build cortexm3 cortexm4 cortexm4f

package bits

func leadingZeros32(u uint32) uint
func leadingZeros64(u uint64) uint
func leadingZerosPtr(u uintptr) uint