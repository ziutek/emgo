// Package mpu allows configure Cortex-M memory protection unit.
package mpu

// An Attr is bitfield that describes attributes of memory region.
type Attr uint32

const (
	Enable     Attr = 1 << 0
	Bufferable Attr = 1 << 16
	Cacheable  Attr = 1 << 17
	Shareable  Attr = 1 << 18
)
