// +build l1xx_md

// Package irq provides list of external interrupts.
package irq

import "arch/cortexm/nvic"

const (
	WWDG          nvic.IRQ = 0  // Window WatchDog Interrupt.
	PVD           nvic.IRQ = 1  // PVD through EXTI Line detection Interrupt.
	TAMPER_STAMP  nvic.IRQ = 2  // Tamper and Time Stamp through EXTI Line Interrupts.
	RTC_WKUP      nvic.IRQ = 3  // RTC Wakeup Timer through EXTI Line Interrupt.
	FLASH         nvic.IRQ = 4  // FLASH global Interrupt.
	RCC           nvic.IRQ = 5  // RCC global Interrupt.
	EXTI0         nvic.IRQ = 6  // EXTI Line0 Interrupt.
	EXTI1         nvic.IRQ = 7  // EXTI Line1 Interrupt.
	EXTI2         nvic.IRQ = 8  // EXTI Line2 Interrupt.
	EXTI3         nvic.IRQ = 9  // EXTI Line3 Interrupt.
	EXTI4         nvic.IRQ = 10 // EXTI Line4 Interrupt.
	DMA1_Channel1 nvic.IRQ = 11 // DMA1 Channel 1 global Interrupt.
	DMA1_Channel2 nvic.IRQ = 12 // DMA1 Channel 2 global Interrupt.
	DMA1_Channel3 nvic.IRQ = 13 // DMA1 Channel 3 global Interrupt.
	DMA1_Channel4 nvic.IRQ = 14 // DMA1 Channel 4 global Interrupt.
	DMA1_Channel5 nvic.IRQ = 15 // DMA1 Channel 5 global Interrupt.
	DMA1_Channel6 nvic.IRQ = 16 // DMA1 Channel 6 global Interrupt.
	DMA1_Channel7 nvic.IRQ = 17 // DMA1 Channel 7 global Interrupt.
	ADC1          nvic.IRQ = 18 // ADC1 global Interrupt.
	USB_HP        nvic.IRQ = 19 // USB High Priority Interrupt.
	USB_LP        nvic.IRQ = 20 // USB Low Priority Interrupt.
	DAC           nvic.IRQ = 21 // DAC Interrupt.
	COMP          nvic.IRQ = 22 // Comparator through EXTI Line Interrupt.
	EXTI9_5       nvic.IRQ = 23 // External Line[9:5] Interrupts.
	LCD           nvic.IRQ = 24 // LCD Interrupt.
	TIM9          nvic.IRQ = 25 // TIM9 global Interrupt.
	TIM10         nvic.IRQ = 26 // TIM10 global Interrupt.
	TIM11         nvic.IRQ = 27 // TIM11 global Interrupt.
	TIM2          nvic.IRQ = 28 // TIM2 global Interrupt.
	TIM3          nvic.IRQ = 29 // TIM3 global Interrupt.
	TIM4          nvic.IRQ = 30 // TIM4 global Interrupt.
	I2C1_EV       nvic.IRQ = 31 // I2C1 Event Interrupt.
	I2C1_ER       nvic.IRQ = 32 // I2C1 Error Interrupt.
	I2C2_EV       nvic.IRQ = 33 // I2C2 Event Interrupt.
	I2C2_ER       nvic.IRQ = 34 // I2C2 Error Interrupt.
	SPI1          nvic.IRQ = 35 // SPI1 global Interrupt.
	SPI2          nvic.IRQ = 36 // SPI2 global Interrupt.
	USART1        nvic.IRQ = 37 // USART1 global Interrupt.
	USART2        nvic.IRQ = 38 // USART2 global Interrupt.
	USART3        nvic.IRQ = 39 // USART3 global Interrupt.
	EXTI15_10     nvic.IRQ = 40 // External Line[15:10] Interrupts.
	RTC_Alarm     nvic.IRQ = 41 // RTC Alarm through EXTI Line Interrupt.
	USB_FS_WKUP   nvic.IRQ = 42 // USB FS WakeUp from suspend through EXTI Line Interrupt.
	TIM6          nvic.IRQ = 43 // TIM6 global Interrupt.
	TIM7          nvic.IRQ = 44 // TIM7 global Interrupt.
)
