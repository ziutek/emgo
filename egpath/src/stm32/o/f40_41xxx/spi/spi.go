// +build f40_41xxx

// Peripheral: SPI_Periph  Serial Peripheral Interface.
// Instances:
//  I2S2ext  mmap.I2S2ext_BASE
//  SPI2     mmap.SPI2_BASE
//  SPI3     mmap.SPI3_BASE
//  I2S3ext  mmap.I2S3ext_BASE
//  SPI1     mmap.SPI1_BASE
//  SPI4     mmap.SPI4_BASE
//  SPI5     mmap.SPI5_BASE
//  SPI6     mmap.SPI6_BASE
// Registers:
//  0x00 16  CR1     Control register 1 (not used in I2S mode).
//  0x04 16  CR2     Control register 2.
//  0x08 16  SR      Status register.
//  0x0C 16  DR      Data register.
//  0x10 16  CRCPR   CRC polynomial register (not used in I2S mode).
//  0x14 16  RXCRCR  RX CRC register (not used in I2S mode).
//  0x18 16  TXCRCR  TX CRC register (not used in I2S mode).
//  0x1C 16  I2SCFGR SPI_I2S configuration register.
//  0x20 16  I2SPR   SPI_I2S prescaler register.
// Import:
//  stm32/o/f40_41xxx/mmap
package spi

const (
	CPHA     CR1_Bits = 0x01 << 0  //+ Clock Phase.
	CPOL     CR1_Bits = 0x01 << 1  //+ Clock Polarity.
	MSTR     CR1_Bits = 0x01 << 2  //+ Master Selection.
	BR       CR1_Bits = 0x07 << 3  //+ BR[2:0] bits (Baud Rate Control).
	BR_0     CR1_Bits = 0x01 << 3  //  Bit 0.
	BR_1     CR1_Bits = 0x02 << 3  //  Bit 1.
	BR_2     CR1_Bits = 0x04 << 3  //  Bit 2.
	SPE      CR1_Bits = 0x01 << 6  //+ SPI Enable.
	LSBFIRST CR1_Bits = 0x01 << 7  //+ Frame Format.
	SSI      CR1_Bits = 0x01 << 8  //+ Internal slave select.
	SSM      CR1_Bits = 0x01 << 9  //+ Software slave management.
	RXONLY   CR1_Bits = 0x01 << 10 //+ Receive only.
	DFF      CR1_Bits = 0x01 << 11 //+ Data Frame Format.
	CRCNEXT  CR1_Bits = 0x01 << 12 //+ Transmit CRC next.
	CRCEN    CR1_Bits = 0x01 << 13 //+ Hardware CRC calculation enable.
	BIDIOE   CR1_Bits = 0x01 << 14 //+ Output enable in bidirectional mode.
	BIDIMODE CR1_Bits = 0x01 << 15 //+ Bidirectional data mode enable.
)

const (
	RXDMAEN CR2_Bits = 0x01 << 0 //+ Rx Buffer DMA Enable.
	TXDMAEN CR2_Bits = 0x01 << 1 //+ Tx Buffer DMA Enable.
	SSOE    CR2_Bits = 0x01 << 2 //+ SS Output Enable.
	ERRIE   CR2_Bits = 0x01 << 5 //+ Error Interrupt Enable.
	RXNEIE  CR2_Bits = 0x01 << 6 //+ RX buffer Not Empty Interrupt Enable.
	TXEIE   CR2_Bits = 0x01 << 7 //+ Tx buffer Empty Interrupt Enable.
)

const (
	RXNE   SR_Bits = 0x01 << 0 //+ Receive buffer Not Empty.
	TXE    SR_Bits = 0x01 << 1 //+ Transmit buffer Empty.
	CHSIDE SR_Bits = 0x01 << 2 //+ Channel side.
	UDR    SR_Bits = 0x01 << 3 //+ Underrun flag.
	CRCERR SR_Bits = 0x01 << 4 //+ CRC Error flag.
	MODF   SR_Bits = 0x01 << 5 //+ Mode fault.
	OVR    SR_Bits = 0x01 << 6 //+ Overrun flag.
	BSY    SR_Bits = 0x01 << 7 //+ Busy flag.
)

const ()

const (
	CRCPOLY CRCPR_Bits = 0xFFFF << 0 //+ CRC polynomial register.
)

const (
	RXCRC RXCRCR_Bits = 0xFFFF << 0 //+ Rx CRC Register.
)

const (
	TXCRC TXCRCR_Bits = 0xFFFF << 0 //+ Tx CRC Register.
)

const (
	CHLEN    I2SCFGR_Bits = 0x01 << 0  //+ Channel length (number of bits per audio channel).
	DATLEN   I2SCFGR_Bits = 0x03 << 1  //+ DATLEN[1:0] bits (Data length to be transferred).
	DATLEN_0 I2SCFGR_Bits = 0x01 << 1  //  Bit 0.
	DATLEN_1 I2SCFGR_Bits = 0x02 << 1  //  Bit 1.
	CKPOL    I2SCFGR_Bits = 0x01 << 3  //+ steady state clock polarity.
	I2SSTD   I2SCFGR_Bits = 0x03 << 4  //+ I2SSTD[1:0] bits (I2S standard selection).
	I2SSTD_0 I2SCFGR_Bits = 0x01 << 4  //  Bit 0.
	I2SSTD_1 I2SCFGR_Bits = 0x02 << 4  //  Bit 1.
	PCMSYNC  I2SCFGR_Bits = 0x01 << 7  //+ PCM frame synchronization.
	I2SCFG   I2SCFGR_Bits = 0x03 << 8  //+ I2SCFG[1:0] bits (I2S configuration mode).
	I2SCFG_0 I2SCFGR_Bits = 0x01 << 8  //  Bit 0.
	I2SCFG_1 I2SCFGR_Bits = 0x02 << 8  //  Bit 1.
	I2SE     I2SCFGR_Bits = 0x01 << 10 //+ I2S Enable.
	I2SMOD   I2SCFGR_Bits = 0x01 << 11 //+ I2S mode selection.
)

const (
	I2SDIV I2SPR_Bits = 0xFF << 0 //+ I2S Linear prescaler.
	ODD    I2SPR_Bits = 0x01 << 8 //+ Odd factor for the prescaler.
	MCKOE  I2SPR_Bits = 0x01 << 9 //+ Master Clock Output Enable.
)
