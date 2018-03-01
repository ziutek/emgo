// Peripheral: AUX_Periph
// Instances:
//  AUX  0xe000e008
// Registers:
//  0x00 32  ACTLR Auxiliary Control Register
package scb

const (
	DISMCYCINT ACTLR = 1 << 0 //+
	DISDEFWBUF ACTLR = 1 << 1 //+
	DISFOLD    ACTLR = 1 << 2 //+
	DISFPCA    ACTLR = 1 << 8 //+
	DISOOFP    ACTLR = 1 << 9 //+
)
