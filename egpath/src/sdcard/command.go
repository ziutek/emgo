package sdcard

import (
	"bits"
)

type Command uint16

const (
	// Command fields
	CmdIdx   Command = 63 << 0 // Command index.
	HasResp  Command = 1 << 6  // Response expected.
	LongResp Command = 1 << 7  // Long response.
	RespIdx  Command = 63 << 8 // Response index.
	BusyCmd  Command = 1 << 14 // Command can set D0 low to signal busy state.
	AppCmd   Command = 1 << 15 // Application command (hint, APP_CMD required).

	// Response types
	RespType = RespIdx | HasResp | LongResp
	R1       = 1<<10 | HasResp
	R2       = 2<<10 | HasResp | LongResp
	R3       = 3<<10 | HasResp
	R4       = 4<<10 | HasResp
	R5       = 5<<10 | HasResp
	R6       = 6<<10 | HasResp
	R7       = 7<<10 | HasResp

	cmd0  = 0                 // GO_IDLE_STATE
	cmd2  = 2 | R2            // ALL_SEND_CID
	cmd3  = 3 | R6            // SEND_RELATIVE_ADDR
	cmd4  = 4                 // SET_DSR
	cmd5  = 5 | R4            // IO_SEND_OP_COND
	cmd6  = 6 | R1            // SWITCH_FUNC
	cmd7  = 7 | R1 | BusyCmd  // SELECT_CARD/DESELECT_CARD
	cmd8  = 8 | R7            // SEND_IF_COND
	cmd9  = 9 | R2            // SEND_CSD
	cmd10 = 10 | R2           // SEND_CID
	cmd11 = 11 | R1           // VOLTAGE_SWITCH
	cmd12 = 12 | R1 | BusyCmd // STOP_TRANSMISSION
	cmd13 = 13 | R1           // SEND_STATUS/SEND_TASK_STATUS
	cmd15 = 15                // GO_INACTIVE_STATE
	cmd16 = 16 | R1           // SET_BLOCKLEN
	cmd17 = 17 | R1           // READ_SINGLE_BLOCK
	cmd18 = 18 | R1           // READ_MULTIPLE_BLOCK
	cmd19 = 19 | R1           // SEND_TUNING_BLOCK
	cmd20 = 20 | R1 | BusyCmd // SPEED_CLASS_CONTROL
	cmd23 = 23 | R1           // SET_BLOCK_COUNT
	cmd24 = 24 | R1           // WRITE_BLOCK
	cmd25 = 25 | R1           // WRITE_MULTIPLE_BLOCK
	cmd27 = 27 | R1           // PROGRAM_CSD
	cmd28 = 28 | R1 | BusyCmd // SET_WRITE_PROT
	cmd29 = 29 | R1 | BusyCmd // CLR_WRITE_PROT
	cmd30 = 30 | R1           // SEND_WRITE_PROT
	cmd32 = 30 | R1           // ERASE_WR_BLK_START
	cmd33 = 33 | R1           // ERASE_WR_BLK_END
	cmd38 = 38 | R1 | BusyCmd // ERASE
	cmd40 = 40 | R1           // TODO: See DPS spec.
	cmd42 = 42 | R1           // LOCK_UNLOCK
	cmd52 = 52 | R5           // IO_RW_DIRECT
	cmd53 = 53 | R5           // IO_RW_EXTENDED
	cmd55 = 55 | R1           // APP_CMD
	cmd56 = 56 | R1           // GEN_CMD

	acmd6  = AppCmd | 6 | R1  // SET_BUS_WIDTH
	acmd13 = AppCmd | 13 | R1 // SD_STATUS
	acmd22 = AppCmd | 22 | R1 // SEND_NUM_WR_BLOCKS
	acmd23 = AppCmd | 23 | R1 // SET_WR_BLK_ERASE_COUNT
	acmd41 = AppCmd | 41 | R3 // SD_SEND_OP_COND
	acmd42 = AppCmd | 42 | R1 // SET_CLR_CARD_DETECT
	acmd51 = AppCmd | 51 | R1 // SEND_SCR
)

// Use command numbers (eg. CMD17) as function names instead of abbreviated
// command names (eg. READ_SINGLE_BLOCK) because the SD Card protocol
// specification uses numbers instead of names when describes state machines
// and algorithms. It turned out to be difficult to follow the source code side
// by side with the protocol specification when code uses command names.

// CMD0 (GO_IDLE_STATE) performs software reset and sets the card into Idle
// State.
func CMD0() (Command, uint32) {
	return cmd0, 0
}

// CMD2 (ALL_SEND_CID, R2) gets Card Identification Data.
func CMD2() (Command, uint32) {
	return cmd2, 0
}

// CMD3 (SEND_RELATIVE_ADDR, R6) asks the card to publishets a new Relative
// Card Address (RCA) and Card Status bits 23,22,19,12:0
func CMD3() (Command, uint32) {
	return cmd3, 0
}

type SwitchFunc uint32

const (
	AccessMode   SwitchFunc = 0x00000F // Access mode (keep current).
	DefaultSpeed SwitchFunc = 0x000000 // Default Speed or SDR12.
	HighSpeed    SwitchFunc = 0x000001 // High Speed or SDR25.
	SDR50        SwitchFunc = 0x000002 // SDR50.
	SDR104       SwitchFunc = 0x000003 // SDR104.
	DDR50        SwitchFunc = 0x000004 // DDR50.

	CommandSystem SwitchFunc = 0x0000F0 // Command system (keep current).
	DefaultSystem SwitchFunc = 0x000000 // Default Command System.
	OTP           SwitchFunc = 0x000030
	ASSD          SwitchFunc = 0x000040
	VendorSpec    SwitchFunc = 0x0000E0

	Driver       SwitchFunc = 0x000F00 // UHS-I driver strength (keep current).
	DefaultTypeB SwitchFunc = 0x000000 // Default Type B driver.
	TypeA        SwitchFunc = 0x000100 // Type A driver.
	TypeC        SwitchFunc = 0x000200 // Type C driver.
	TypeD        SwitchFunc = 0x000300 // Type D driver.

	PowerLimit SwitchFunc = 0x00F000 // Power limit (keep current).
	Default720 SwitchFunc = 0x000000 // Default limit: 720 mW.
	Power1440  SwitchFunc = 0x001000 // Limit: 1440 mW.
	Power2160  SwitchFunc = 0x002000 // Limit: 2160 mW.
	Power2880  SwitchFunc = 0x003000 // Limit: 2880 mW.
	Power1800  SwitchFunc = 0x004000 // Limit: 1800 mW.

	ModeCheck  SwitchFunc = 0 << 31
	ModeSwitch SwitchFunc = 1 << 31
)

// CMD6 (SWITCH_FUNC, R1) switches or expands memory card functions.
func CMD6(sf SwitchFunc) (Command, uint32) {
	return cmd6, uint32(sf)
}

// CMD7 (SELECT_CARD/DESELECT_CARD, R1b) selects card with rca address (puts
// into Transfer State) and deselects all other (puts into Stand-by State).
func CMD7(rca uint16) (Command, uint32) {
	return cmd7, uint32(rca) << 16
}

type VHS byte

const (
	V27_36 VHS = 1 << 0
	LVR    VHS = 1 << 1
)

// CMD8 (SEND_IF_COND, R7) initializes SD Memory Cards compliant to the Physical
// Layer Specification Version 2.00 or later.
func CMD8(vhs VHS, checkPattern byte) (Command, uint32) {
	return cmd8, uint32(vhs)&0xF<<8 | uint32(checkPattern)
}

// CMD9 (SEND_CSD, R2) gets Card Specific Data from card indentified by rca.
func CMD9(rca uint16) (Command, uint32) {
	return cmd9, uint32(rca) << 16
}

// CMD12 (STOP_TRANSMISSION, R1b) forces the card to stop transmission in
// Multiple Block Read Operation.
func CMD12() (Command, uint32) {
	return cmd12, 0
}

// CMD16 (SET_BLOCKLEN, R1) sets the block length (in bytes) for block commands.
func CMD16(blen int) (Command, uint32) {
	return cmd16, uint32(blen)
}

// CMD17 (READ_SINGLE_BLOCK, R1) reads a block of the size selected by CMD16.
func CMD17(addr uint) (Command, uint32) {
	return cmd17, uint32(addr)
}

// CMD18 (READ_MULTIPLE_BLOCK, R1) works like CMD17 but does not stop the
// transmision after first data block. Instead the card continuously transfers
// data blocks until it receives CMD12 (STOP_TRANSMISSION) command.
func CMD18(addr uint) (Command, uint32) {
	return cmd18, uint32(addr)
}

// CMD55 (APP_CMD, R1) indicates to the card that the next command is an
// application specific command.
func CMD55(rca uint16) (Command, uint32) {
	return cmd55, uint32(rca) << 16
}

type BusWidth byte

const (
	Bus1bit BusWidth = 0
	Bus4bit BusWidth = 2
)

// ACMD6 (SET_BUS_WIDTH, R1) sets the data bus width.
func ACMD6(bw BusWidth) (Command, uint32) {
	return acmd6, uint32(bw)
}

// ACMD41 (SD_SEND_OP_COND, R3) starts initialization/identification process.
func ACMD41(ocr OCR) (Command, uint32) {
	return acmd41, uint32(ocr)
}

// ACMD42 (SET_CLR_CARD_DETECT, R1) enables/disables pull-up resistor on D3/CD.
func ACMD42(pullUp bool) (Command, uint32) {
	return acmd42, uint32(bits.One(pullUp))
}

// ACMD51 (SEND_SCR, R1) reads SD Configuration Register.
func ACMD51() (Command, uint32) {
	return acmd51, 0
}
