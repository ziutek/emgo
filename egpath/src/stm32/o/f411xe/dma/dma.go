// +build f411xe

// Peripheral: DMA_Periph  DMA Controller.
// Instances:
//  DMA1  mmap.DMA1_BASE
//  DMA2  mmap.DMA2_BASE
// Registers:
//  0x00 32  LISR  Low interrupt status register.
//  0x04 32  HISR  High interrupt status register.
//  0x08 32  LIFCR Low interrupt flag clear register.
//  0x0C 32  HIFCR High interrupt flag clear register.
// Import:
//  stm32/o/f411xe/mmap
package dma

const (
	TCIF3  LISR_Bits = 0x01 << 27 //+
	HTIF3  LISR_Bits = 0x01 << 26 //+
	TEIF3  LISR_Bits = 0x01 << 25 //+
	DMEIF3 LISR_Bits = 0x01 << 24 //+
	FEIF3  LISR_Bits = 0x01 << 22 //+
	TCIF2  LISR_Bits = 0x01 << 21 //+
	HTIF2  LISR_Bits = 0x01 << 20 //+
	TEIF2  LISR_Bits = 0x01 << 19 //+
	DMEIF2 LISR_Bits = 0x01 << 18 //+
	FEIF2  LISR_Bits = 0x01 << 16 //+
	TCIF1  LISR_Bits = 0x01 << 11 //+
	HTIF1  LISR_Bits = 0x01 << 10 //+
	TEIF1  LISR_Bits = 0x01 << 9  //+
	DMEIF1 LISR_Bits = 0x01 << 8  //+
	FEIF1  LISR_Bits = 0x01 << 6  //+
	TCIF0  LISR_Bits = 0x01 << 5  //+
	HTIF0  LISR_Bits = 0x01 << 4  //+
	TEIF0  LISR_Bits = 0x01 << 3  //+
	DMEIF0 LISR_Bits = 0x01 << 2  //+
	FEIF0  LISR_Bits = 0x01 << 0  //+
)

const (
	TCIF7  HISR_Bits = 0x01 << 27 //+
	HTIF7  HISR_Bits = 0x01 << 26 //+
	TEIF7  HISR_Bits = 0x01 << 25 //+
	DMEIF7 HISR_Bits = 0x01 << 24 //+
	FEIF7  HISR_Bits = 0x01 << 22 //+
	TCIF6  HISR_Bits = 0x01 << 21 //+
	HTIF6  HISR_Bits = 0x01 << 20 //+
	TEIF6  HISR_Bits = 0x01 << 19 //+
	DMEIF6 HISR_Bits = 0x01 << 18 //+
	FEIF6  HISR_Bits = 0x01 << 16 //+
	TCIF5  HISR_Bits = 0x01 << 11 //+
	HTIF5  HISR_Bits = 0x01 << 10 //+
	TEIF5  HISR_Bits = 0x01 << 9  //+
	DMEIF5 HISR_Bits = 0x01 << 8  //+
	FEIF5  HISR_Bits = 0x01 << 6  //+
	TCIF4  HISR_Bits = 0x01 << 5  //+
	HTIF4  HISR_Bits = 0x01 << 4  //+
	TEIF4  HISR_Bits = 0x01 << 3  //+
	DMEIF4 HISR_Bits = 0x01 << 2  //+
	FEIF4  HISR_Bits = 0x01 << 0  //+
)

const (
	CTCIF3  LIFCR_Bits = 0x01 << 27 //+
	CHTIF3  LIFCR_Bits = 0x01 << 26 //+
	CTEIF3  LIFCR_Bits = 0x01 << 25 //+
	CDMEIF3 LIFCR_Bits = 0x01 << 24 //+
	CFEIF3  LIFCR_Bits = 0x01 << 22 //+
	CTCIF2  LIFCR_Bits = 0x01 << 21 //+
	CHTIF2  LIFCR_Bits = 0x01 << 20 //+
	CTEIF2  LIFCR_Bits = 0x01 << 19 //+
	CDMEIF2 LIFCR_Bits = 0x01 << 18 //+
	CFEIF2  LIFCR_Bits = 0x01 << 16 //+
	CTCIF1  LIFCR_Bits = 0x01 << 11 //+
	CHTIF1  LIFCR_Bits = 0x01 << 10 //+
	CTEIF1  LIFCR_Bits = 0x01 << 9  //+
	CDMEIF1 LIFCR_Bits = 0x01 << 8  //+
	CFEIF1  LIFCR_Bits = 0x01 << 6  //+
	CTCIF0  LIFCR_Bits = 0x01 << 5  //+
	CHTIF0  LIFCR_Bits = 0x01 << 4  //+
	CTEIF0  LIFCR_Bits = 0x01 << 3  //+
	CDMEIF0 LIFCR_Bits = 0x01 << 2  //+
	CFEIF0  LIFCR_Bits = 0x01 << 0  //+
)

const (
	CTCIF7  HIFCR_Bits = 0x01 << 27 //+
	CHTIF7  HIFCR_Bits = 0x01 << 26 //+
	CTEIF7  HIFCR_Bits = 0x01 << 25 //+
	CDMEIF7 HIFCR_Bits = 0x01 << 24 //+
	CFEIF7  HIFCR_Bits = 0x01 << 22 //+
	CTCIF6  HIFCR_Bits = 0x01 << 21 //+
	CHTIF6  HIFCR_Bits = 0x01 << 20 //+
	CTEIF6  HIFCR_Bits = 0x01 << 19 //+
	CDMEIF6 HIFCR_Bits = 0x01 << 18 //+
	CFEIF6  HIFCR_Bits = 0x01 << 16 //+
	CTCIF5  HIFCR_Bits = 0x01 << 11 //+
	CHTIF5  HIFCR_Bits = 0x01 << 10 //+
	CTEIF5  HIFCR_Bits = 0x01 << 9  //+
	CDMEIF5 HIFCR_Bits = 0x01 << 8  //+
	CFEIF5  HIFCR_Bits = 0x01 << 6  //+
	CTCIF4  HIFCR_Bits = 0x01 << 5  //+
	CHTIF4  HIFCR_Bits = 0x01 << 4  //+
	CTEIF4  HIFCR_Bits = 0x01 << 3  //+
	CDMEIF4 HIFCR_Bits = 0x01 << 2  //+
	CFEIF4  HIFCR_Bits = 0x01 << 0  //+
)
