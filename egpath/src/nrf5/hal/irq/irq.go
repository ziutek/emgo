package irq

import "arch/cortexm/nvic"

const (
	CLK_PWR_BPR nvic.IRQ = 0  // Clock control/Power control/Block protect
	RADIO       nvic.IRQ = 1  // 2.4 GHz Radio
	UART0       nvic.IRQ = 2  // UART  / UART with EasyDMA
	SPI0_TWI0   nvic.IRQ = 3  // SPI 0 / Two-Wire (I2C) 0
	SPI1_TWI1   nvic.IRQ = 4  // SPI 1 / Two-Wire (I2C) 1 / SPI Slave 1
	GPIOTE      nvic.IRQ = 6  // GPIO Tasks and Events
	ADC         nvic.IRQ = 7  // Analog-to-digital converter
	TIMER0      nvic.IRQ = 8  // Timer/counter 0
	TIMER1      nvic.IRQ = 9  // Timer/counter 1
	TIMER2      nvic.IRQ = 10 // Timer/counter 2
	RTC0        nvic.IRQ = 11 // Real time counter 0
	TEMP        nvic.IRQ = 12 // Temperature sensor
	RNG         nvic.IRQ = 13 // Random number generator
	ECB         nvic.IRQ = 14 // Random number generator
	CCM_AAR     nvic.IRQ = 15 // AES CCM mode encrypt./accelerated address resolver
	WDT         nvic.IRQ = 16 // Watchdog timer
	RTC1        nvic.IRQ = 17 // Real time counter 1
	QDEC        nvic.IRQ = 18 //  Quadrature decoder
	LPCOMP      nvic.IRQ = 19 // Low power comparator
	SWI0        nvic.IRQ = 20 // Software interrupt 0
	SWI1        nvic.IRQ = 21 // Software interrupt 1
	SWI2        nvic.IRQ = 22 // Software interrupt 2
	SWI3        nvic.IRQ = 23 // Software interrupt 3
	SWI4        nvic.IRQ = 24 // Software interrupt 4
	SWI5        nvic.IRQ = 25 // Software interrupt 5
	NVMC        nvic.IRQ = 30 // Non Volatile Memory Controller
	PPI         nvic.IRQ = 31 // PPI controller
)
