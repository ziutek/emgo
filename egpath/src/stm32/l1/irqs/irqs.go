package irqs

import "cortexm/exce"

const (
	WinWdg = exce.IRQ0 + iota
	PVD
	TampStamp
	RTCWkup
	Flash
	RCC
	Ext0
	Ext1
	Ext2
	Ext3
	Ext4
	DMA1Chan1
	DMA1Chan2
	DMA1Chan3
	DMA1Chan4
	DMA1Chan5
	DMA1Chan6
	DMA1Chan7
	ADC1
	USBHP
	USBLP
	DAC
	Comp
	Ext9_5
	LCD
	Tim9
	Tim10
	Tim11
	Tim2
	Tim3
	Tim4
	I2C1Ev
	I2C1Er
	I2C2Ev
	I2C2Er
	SPI1
	SPI2
	USART1
	USART2
	USART3
	Ext15_10
	RTCAlarm
	OTGFSWkup
	Tim6
	Tim7
	// BUG: following interrupt names are wrong in case of high density devices 
	SDIO
	Tim5
	SPI3
	UART4
	UART5
	DMA2Chan1
	DMA2Chan2
	DMA2Chan3
	DMA2Chan4
	DMA2Chan5
	AES
	CompACQ
)
