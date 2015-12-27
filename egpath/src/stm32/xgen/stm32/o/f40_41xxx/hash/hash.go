// Peripheral: HASH_Periph  HASH.
// Instances:
//  HASH  mmap.HASH_BASE
// Registers:
//  0x00 32  CR      Control register.
//  0x04 32  DIN     Data input register.
//  0x08 32  STR     Start register.
//  0x0C 32  HR[5]   Digest registers.
//  0x20 32  IMR     Interrupt enable register.
//  0x24 32  SR      Status register.
//  0xF8 32  CSR[54] Context swap registers.
// Import:
//  stm32/o/f40_41xxx/mmap
package hash

const (
	INIT       CR_Bits = 0x01 << 2 //+
	DMAE       CR_Bits = 0x01 << 3 //+
	DATATYPE   CR_Bits = 0x03 << 4 //+
	DATATYPE_0 CR_Bits = 0x01 << 4
	DATATYPE_1 CR_Bits = 0x02 << 4
	MODE       CR_Bits = 0x01 << 6  //+
	ALGO       CR_Bits = 0x801 << 7 //+
	ALGO_0     CR_Bits = 0x01 << 7
	ALGO_1     CR_Bits = 0x800 << 7
	NBW        CR_Bits = 0x0F << 8 //+
	NBW_0      CR_Bits = 0x01 << 8
	NBW_1      CR_Bits = 0x02 << 8
	NBW_2      CR_Bits = 0x04 << 8
	NBW_3      CR_Bits = 0x08 << 8
	DINNE      CR_Bits = 0x01 << 12 //+
	MDMAT      CR_Bits = 0x01 << 13 //+
	LKEY       CR_Bits = 0x01 << 16 //+
)

const (
	NBW   STR_Bits = 0x1F << 0 //+
	NBW_0 STR_Bits = 0x01 << 0
	NBW_1 STR_Bits = 0x02 << 0
	NBW_2 STR_Bits = 0x04 << 0
	NBW_3 STR_Bits = 0x08 << 0
	NBW_4 STR_Bits = 0x10 << 0
	DCAL  STR_Bits = 0x01 << 8 //+
)

const (
	DINIM IMR_Bits = 0x01 << 0 //+
	DCIM  IMR_Bits = 0x01 << 1 //+
)

const (
	DINIS SR_Bits = 0x01 << 0 //+
	DCIS  SR_Bits = 0x01 << 1 //+
	DMAS  SR_Bits = 0x01 << 2 //+
	BUSY  SR_Bits = 0x01 << 3 //+
)
