// +build cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package bits

//c:static inline
func reverse32(u uint32) uint32

//c:static inline
func reverse64(u uint64) uint64

//c:static inline
func reversePtr(u uintptr) uintptr