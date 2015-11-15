package irqs

import "arch/cortexm/exce"

const irq0 = exce.IRQ0

const (
	PwrClk  = irq0 + 0  // Power control / clock control
	Radio   = irq0 + 1  // 2.4 GHz Radio
	UART0   = irq0 + 2  // UART 0
	SPITWI0 = irq0 + 3  // SPI 0 / Two-Wire (I2C) 0
	SPITWI1 = irq0 + 4  // SPI 1 / Two-Wire (I2C) 1 / SPI Slave 1
	GPIO    = irq0 + 6  // GPIO
	ADC     = irq0 + 7  // Analog-to-digital converter
	Timer0  = irq0 + 8  // Timer/counter 0
	Timer1  = irq0 + 9  // Timer/counter 1
	Timer2  = irq0 + 10 // Timer/counter 2
	RTC0    = irq0 + 11 // Real time counter 0
	Temp    = irq0 + 12 // Temperature sensor
	RNG     = irq0 + 13 // Random number generator
	ECB     = irq0 + 14 // Random number generator
	CCMAAR  = irq0 + 15 // AES CCM mode encrypt. / accelerated address resolver
	WDT     = irq0 + 16 // Watchdog timer
	RTC1    = irq0 + 17 // Real time counter 1
	QDec    = irq0 + 18 //  Quadrature decoder
	LPComp  = irq0 + 19 // Low power comparator
)
