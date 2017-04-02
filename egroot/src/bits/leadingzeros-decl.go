// +build cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package bits

//c:inline
func leadingZeros32(u uint32) uint

//c:inline
func leadingZeros64(u uint64) uint

//c:inline
func leadingZerosPtr(u uintptr) uint