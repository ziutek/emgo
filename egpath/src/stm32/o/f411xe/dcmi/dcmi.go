// +build f411xe

// Peripheral: DCMI_Periph  DCMI.
// Instances:
//  DCMI  mmap.DCMI_BASE
// Registers:
//  0x00 32  CR      Control register 1.
//  0x04 32  SR      Status register.
//  0x08 32  RISR    Raw interrupt status register.
//  0x0C 32  IER     Interrupt enable register.
//  0x10 32  MISR    Masked interrupt status register.
//  0x14 32  ICR     Interrupt clear register.
//  0x18 32  ESCR    Embedded synchronization code register.
//  0x1C 32  ESUR    Embedded synchronization unmask register.
//  0x20 32  CWSTRTR Crop window start.
//  0x24 32  CWSIZER Crop window size.
//  0x28 32  DR      Data register.
// Import:
//  stm32/o/f411xe/mmap
package dcmi

const (
	CAPTURE CR_Bits = 0x01 << 0  //+
	CM      CR_Bits = 0x01 << 1  //+
	CROP    CR_Bits = 0x01 << 2  //+
	JPEG    CR_Bits = 0x01 << 3  //+
	ESS     CR_Bits = 0x01 << 4  //+
	PCKPOL  CR_Bits = 0x01 << 5  //+
	HSPOL   CR_Bits = 0x01 << 6  //+
	VSPOL   CR_Bits = 0x01 << 7  //+
	FCRC_0  CR_Bits = 0x01 << 8  //+
	FCRC_1  CR_Bits = 0x01 << 9  //+
	EDM_0   CR_Bits = 0x01 << 10 //+
	EDM_1   CR_Bits = 0x01 << 11 //+
	CRE     CR_Bits = 0x01 << 12 //+
	ENABLE  CR_Bits = 0x01 << 14 //+
)

const (
	HSYNC SR_Bits = 0x01 << 0 //+
	VSYNC SR_Bits = 0x01 << 1 //+
	FNE   SR_Bits = 0x01 << 2 //+
)

const (
	FRAME_RIS RISR_Bits = 0x01 << 0 //+
	OVF_RIS   RISR_Bits = 0x01 << 1 //+
	ERR_RIS   RISR_Bits = 0x01 << 2 //+
	VSYNC_RIS RISR_Bits = 0x01 << 3 //+
	LINE_RIS  RISR_Bits = 0x01 << 4 //+
)

const (
	FRAME_IE IER_Bits = 0x01 << 0 //+
	OVF_IE   IER_Bits = 0x01 << 1 //+
	ERR_IE   IER_Bits = 0x01 << 2 //+
	VSYNC_IE IER_Bits = 0x01 << 3 //+
	LINE_IE  IER_Bits = 0x01 << 4 //+
)

const (
	FRAME_MIS MISR_Bits = 0x01 << 0 //+
	OVF_MIS   MISR_Bits = 0x01 << 1 //+
	ERR_MIS   MISR_Bits = 0x01 << 2 //+
	VSYNC_MIS MISR_Bits = 0x01 << 3 //+
	LINE_MIS  MISR_Bits = 0x01 << 4 //+
)

const (
	FRAME_ISC ICR_Bits = 0x01 << 0 //+
	OVF_ISC   ICR_Bits = 0x01 << 1 //+
	ERR_ISC   ICR_Bits = 0x01 << 2 //+
	VSYNC_ISC ICR_Bits = 0x01 << 3 //+
	LINE_ISC  ICR_Bits = 0x01 << 4 //+
)
