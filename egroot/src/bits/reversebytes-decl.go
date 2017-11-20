// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package bits

//c:inline
func reverseBytes16(u uint16) uint16

//c:inline
func reverseBytes32(u uint32) uint32

//c:inline
func reverseBytes64(u uint64) uint64
