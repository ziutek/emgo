// +build cortexm3 cortexm4 cortexm4f

package bits

//c:static inline
func leadingZeros32(u uint32) uint

//c:static inline
func leadingZeros64(u uint64) uint

//c:static inline
func leadingZerosPtr(u uintptr) uint