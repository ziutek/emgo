// Peripheral: PWR_Periph  Power Control.
// Instances:
//  PWR  mmap.PWR_BASE
// Registers:
//  0x00 32  CR  Power control register.
//  0x04 32  CSR Power control/status register.
// Import:
//  stm32/o/l1xx_md/mmap
package pwr

const (
	LPSDSR   CR_Bits = 0x01 << 0  //+ Low-power deepsleep/sleep/low power run.
	PDDS     CR_Bits = 0x01 << 1  //+ Power Down Deepsleep.
	CWUF     CR_Bits = 0x01 << 2  //+ Clear Wakeup Flag.
	CSBF     CR_Bits = 0x01 << 3  //+ Clear Standby Flag.
	PVDE     CR_Bits = 0x01 << 4  //+ Power Voltage Detector Enable.
	PLS      CR_Bits = 0x07 << 5  //+ PLS[2:0] bits (PVD Level Selection).
	PLS_0    CR_Bits = 0x01 << 5  //  Bit 0.
	PLS_1    CR_Bits = 0x02 << 5  //  Bit 1.
	PLS_2    CR_Bits = 0x04 << 5  //  Bit 2.
	PLS_LEV0 CR_Bits = 0x00 << 5  //  PVD level 0.
	PLS_LEV1 CR_Bits = 0x01 << 5  //  PVD level 1.
	PLS_LEV2 CR_Bits = 0x02 << 5  //  PVD level 2.
	PLS_LEV3 CR_Bits = 0x03 << 5  //  PVD level 3.
	PLS_LEV4 CR_Bits = 0x04 << 5  //  PVD level 4.
	PLS_LEV5 CR_Bits = 0x05 << 5  //  PVD level 5.
	PLS_LEV6 CR_Bits = 0x06 << 5  //  PVD level 6.
	PLS_LEV7 CR_Bits = 0x07 << 5  //  PVD level 7.
	DBP      CR_Bits = 0x01 << 8  //+ Disable Backup Domain write protection.
	ULP      CR_Bits = 0x01 << 9  //+ Ultra Low Power mode.
	FWU      CR_Bits = 0x01 << 10 //+ Fast wakeup.
	VOS      CR_Bits = 0x03 << 11 //+ VOS[1:0] bits (Voltage scaling range selection).
	VOS_0    CR_Bits = 0x01 << 11 //  Bit 0.
	VOS_1    CR_Bits = 0x02 << 11 //  Bit 1.
	LPRUN    CR_Bits = 0x01 << 14 //+ Low power run mode.
)

const (
	WUF         CSR_Bits = 0x01 << 0  //+ Wakeup Flag.
	SBF         CSR_Bits = 0x01 << 1  //+ Standby Flag.
	PVDO        CSR_Bits = 0x01 << 2  //+ PVD Output.
	VREFINTRDYF CSR_Bits = 0x01 << 3  //+ Internal voltage reference (VREFINT) ready flag.
	VOSF        CSR_Bits = 0x01 << 4  //+ Voltage Scaling select flag.
	REGLPF      CSR_Bits = 0x01 << 5  //+ Regulator LP flag.
	EWUP1       CSR_Bits = 0x01 << 8  //+ Enable WKUP pin 1.
	EWUP2       CSR_Bits = 0x01 << 9  //+ Enable WKUP pin 2.
	EWUP3       CSR_Bits = 0x01 << 10 //+ Enable WKUP pin 3.
)
