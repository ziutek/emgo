// +build f303xe

package tim

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	// Advanced-control timers.
	TIM1  = (*Periph)(unsafe.Pointer(mmap.TIM1_BASE))
	TIM8  = (*Periph)(unsafe.Pointer(mmap.TIM8_BASE))
	TIM20 = (*Periph)(unsafe.Pointer(mmap.TIM20_BASE))

	// General-purpose timers.
	TIM2  = (*Periph)(unsafe.Pointer(mmap.TIM2_BASE)) // 32-bit
	TIM3  = (*Periph)(unsafe.Pointer(mmap.TIM3_BASE))
	TIM4  = (*Periph)(unsafe.Pointer(mmap.TIM4_BASE))

	// Basic timers.
	TIM6 = (*Periph)(unsafe.Pointer(mmap.TIM6_BASE))
	TIM7 = (*Periph)(unsafe.Pointer(mmap.TIM7_BASE))

	// General-purpose timers (1-channel).
	TIM16 = (*Periph)(unsafe.Pointer(mmap.TIM15_BASE))
	TIM17 = (*Periph)(unsafe.Pointer(mmap.TIM17_BASE))

	// General-purpose timers (2-channel).
	TIM15 = (*Periph)(unsafe.Pointer(mmap.TIM15_BASE))
)
