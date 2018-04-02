// +build f030x6 f030x8

package tim

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	// Advanced-control timers.
	TIM1 = (*Periph)(unsafe.Pointer(mmap.TIM1_BASE))

	// General-purpose timers.
	TIM3 = (*Periph)(unsafe.Pointer(mmap.TIM3_BASE))

	// General-purpose timers (1-channel).
	TIM14 = (*Periph)(unsafe.Pointer(mmap.TIM14_BASE))

	// General-purpose timers (1-channel with complementary output).
	TIM16 = (*Periph)(unsafe.Pointer(mmap.TIM16_BASE))
	TIM17 = (*Periph)(unsafe.Pointer(mmap.TIM17_BASE))
)
