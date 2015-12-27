// Package irq provides list of all defined external interrupts.
package irq

import "arch/cortexm/nvic"

const (
	WWDG               nvic.IRQ = 0  // Window WatchDog Interrupt.
	PVD                nvic.IRQ = 1  // PVD through EXTI Line detection Interrupt.
	TAMP_STAMP         nvic.IRQ = 2  // Tamper and TimeStamp interrupts through the EXTI line.
	RTC_WKUP           nvic.IRQ = 3  // RTC Wakeup interrupt through the EXTI line.
	FLASH              nvic.IRQ = 4  // FLASH global Interrupt.
	RCC                nvic.IRQ = 5  // RCC global Interrupt.
	EXTI0              nvic.IRQ = 6  // EXTI Line0 Interrupt.
	EXTI1              nvic.IRQ = 7  // EXTI Line1 Interrupt.
	EXTI2              nvic.IRQ = 8  // EXTI Line2 Interrupt.
	EXTI3              nvic.IRQ = 9  // EXTI Line3 Interrupt.
	EXTI4              nvic.IRQ = 10 // EXTI Line4 Interrupt.
	DMA1_Stream0       nvic.IRQ = 11 // DMA1 Stream 0 global Interrupt.
	DMA1_Stream1       nvic.IRQ = 12 // DMA1 Stream 1 global Interrupt.
	DMA1_Stream2       nvic.IRQ = 13 // DMA1 Stream 2 global Interrupt.
	DMA1_Stream3       nvic.IRQ = 14 // DMA1 Stream 3 global Interrupt.
	DMA1_Stream4       nvic.IRQ = 15 // DMA1 Stream 4 global Interrupt.
	DMA1_Stream5       nvic.IRQ = 16 // DMA1 Stream 5 global Interrupt.
	DMA1_Stream6       nvic.IRQ = 17 // DMA1 Stream 6 global Interrupt.
	ADC                nvic.IRQ = 18 // ADC1, ADC2 and ADC3 global Interrupts.
	EXTI9_5            nvic.IRQ = 23 // External Line[9:5] Interrupts.
	TIM1_BRK_TIM9      nvic.IRQ = 24 // TIM1 Break interrupt and TIM9 global interrupt.
	TIM1_UP_TIM10      nvic.IRQ = 25 // TIM1 Update Interrupt and TIM10 global interrupt.
	TIM1_TRG_COM_TIM11 nvic.IRQ = 26 // TIM1 Trigger and Commutation Interrupt and TIM11 global interrupt.
	TIM1_CC            nvic.IRQ = 27 // TIM1 Capture Compare Interrupt.
	TIM2               nvic.IRQ = 28 // TIM2 global Interrupt.
	TIM3               nvic.IRQ = 29 // TIM3 global Interrupt.
	TIM4               nvic.IRQ = 30 // TIM4 global Interrupt.
	I2C1_EV            nvic.IRQ = 31 // I2C1 Event Interrupt.
	I2C1_ER            nvic.IRQ = 32 // I2C1 Error Interrupt.
	I2C2_EV            nvic.IRQ = 33 // I2C2 Event Interrupt.
	I2C2_ER            nvic.IRQ = 34 // I2C2 Error Interrupt.
	SPI1               nvic.IRQ = 35 // SPI1 global Interrupt.
	SPI2               nvic.IRQ = 36 // SPI2 global Interrupt.
	USART1             nvic.IRQ = 37 // USART1 global Interrupt.
	USART2             nvic.IRQ = 38 // USART2 global Interrupt.
	EXTI15_10          nvic.IRQ = 40 // External Line[15:10] Interrupts.
	RTC_Alarm          nvic.IRQ = 41 // RTC Alarm (A and B) through EXTI Line Interrupt.
	OTG_FS_WKUP        nvic.IRQ = 42 // USB OTG FS Wakeup through EXTI line interrupt.
	DMA1_Stream7       nvic.IRQ = 47 // DMA1 Stream7 Interrupt.
	SDIO               nvic.IRQ = 49 // SDIO global Interrupt.
	TIM5               nvic.IRQ = 50 // TIM5 global Interrupt.
	SPI3               nvic.IRQ = 51 // SPI3 global Interrupt.
	DMA2_Stream0       nvic.IRQ = 56 // DMA2 Stream 0 global Interrupt.
	DMA2_Stream1       nvic.IRQ = 57 // DMA2 Stream 1 global Interrupt.
	DMA2_Stream2       nvic.IRQ = 58 // DMA2 Stream 2 global Interrupt.
	DMA2_Stream3       nvic.IRQ = 59 // DMA2 Stream 3 global Interrupt.
	DMA2_Stream4       nvic.IRQ = 60 // DMA2 Stream 4 global Interrupt.
	OTG_FS             nvic.IRQ = 67 // USB OTG FS global Interrupt.
	DMA2_Stream5       nvic.IRQ = 68 // DMA2 Stream 5 global interrupt.
	DMA2_Stream6       nvic.IRQ = 69 // DMA2 Stream 6 global interrupt.
	DMA2_Stream7       nvic.IRQ = 70 // DMA2 Stream 7 global interrupt.
	USART6             nvic.IRQ = 71 // USART6 global interrupt.
	I2C3_EV            nvic.IRQ = 72 // I2C3 event interrupt.
	I2C3_ER            nvic.IRQ = 73 // I2C3 error interrupt.
	FPU                nvic.IRQ = 81 // FPU global interrupt.
	SPI4               nvic.IRQ = 84 // SPI4 global Interrupt.
	SPI5               nvic.IRQ = 85 // SPI5 global Interrupt.
)
