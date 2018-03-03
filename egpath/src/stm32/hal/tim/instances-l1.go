// +build l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package tim

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	// General-purpose timers.
	TIM2 = (*Periph)(unsafe.Pointer(mmap.TIM2_BASE))
	TIM3 = (*Periph)(unsafe.Pointer(mmap.TIM3_BASE))
	TIM4 = (*Periph)(unsafe.Pointer(mmap.TIM4_BASE))
	TIM5 = (*Periph)(unsafe.Pointer(mmap.TIM5_BASE))

	// Basic timers.
	TIM6 = (*Periph)(unsafe.Pointer(mmap.TIM6_BASE))
	TIM7 = (*Periph)(unsafe.Pointer(mmap.TIM7_BASE))

	// General-purpose timers (1-channel).
	TIM10 = (*Periph)(unsafe.Pointer(mmap.TIM10_BASE))
	TIM11 = (*Periph)(unsafe.Pointer(mmap.TIM11_BASE))

	// General-purpose timers (2-channel).
	TIM9 = (*Periph)(unsafe.Pointer(mmap.TIM9_BASE))
)
