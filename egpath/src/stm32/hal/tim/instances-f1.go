// +build f10x_ld f10x_ld_vl f10x_md f10x_md_vl f10x_hd f10x_hd_vl f10x_xl f10x_cl

package tim

import (
	"unsafe"

	"stm32/hal/raw/mmap"
)

//emgo:const
var (
	// Advanced-control timers.
	TIM1 = (*Periph)(unsafe.Pointer(mmap.TIM1_BASE))
	TIM8 = (*Periph)(unsafe.Pointer(mmap.TIM8_BASE))

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
	TIM13 = (*Periph)(unsafe.Pointer(mmap.TIM13_BASE))
	TIM14 = (*Periph)(unsafe.Pointer(mmap.TIM14_BASE))

	// General-purpose timers (2-channel).
	TIM9  = (*Periph)(unsafe.Pointer(mmap.TIM9_BASE))
	TIM12 = (*Periph)(unsafe.Pointer(mmap.TIM12_BASE))

	// General-purpose timers (2-channel with complementary output).
	TIM15 = (*Periph)(unsafe.Pointer(mmap.TIM15_BASE))

	// General-purpose timers (1-channel with complementary output).
	TIM16 = (*Periph)(unsafe.Pointer(mmap.TIM16_BASE))
	TIM17 = (*Periph)(unsafe.Pointer(mmap.TIM17_BASE))
)
