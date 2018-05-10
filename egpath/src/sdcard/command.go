package sdcard

type Command uint16

const (
	// Command fields
	CMD     Command = 0x3F << 0
	RespLen Command = 0x03 << 6
	ACMD    Command = 0x01 << 8
	Busy    Command = 0x01 << 9
	R       Command = 0x3F << 10

	NoResp    Command = 0 << 6
	ShortResp Command = 1 << 6
	LongResp  Command = 3 << 6

	R1  = 1<<10 | ShortResp
	R1b = 1<<10 | ShortResp | Busy
	R2  = 2<<10 | LongResp
	R3  = 3<<10 | ShortResp
	R4  = 4<<10 | ShortResp
	R5  = 5<<10 | ShortResp
	R6  = 6<<10 | ShortResp
	R7  = 7<<10 | ShortResp

	cmd0  = 0 | NoResp  // GO_IDLE_STATE
	cmd2  = 2 | R2      // ALL_SEND_CID
	cmd3  = 3 | R6      // SEND_RELATIVE_ADDR
	cmd4  = 4 | NoResp  // SET_DSR
	cmd5  = 5 | R4      // IO_SEND_OP_COND
	cmd7  = 7 | R1b     // SELECT_CARD/DESELECT_CARD
	cmd8  = 8 | R7      // SEND_IF_COND
	cmd9  = 9 | R2      // SEND_CSD
	cmd10 = 10 | R2     // SEND_CID
	cmd11 = 11 | R1     // VOLTAGE_SWITCH
	cmd12 = 12 | R1b    // STOP_TRANSMISSION
	cmd13 = 13 | R1     // SEND_STATUS/SEND_TASK_STATUS
	cmd15 = 15 | NoResp // GO_INACTIVE_STATE
	cmd16 = 16 | R1     // SET_BLOCKLEN
	cmd17 = 17 | R1     // READ_SINGLE_BLOCK
	cmd18 = 18 | R1     // READ_MULTIPLE_BLOCK
	cmd19 = 19 | R1     // SEND_TUNING_BLOCK
	cmd20 = 20 | R1b    // SPEED_CLASS_CONTROL
	cmd23 = 23 | R1     // SET_BLOCK_COUNT
	cmd24 = 24 | R1     // WRITE_BLOCK
	cmd25 = 25 | R1     // WRITE_MULTIPLE_BLOCK
	cmd27 = 27 | R1     // PROGRAM_CSD
	cmd28 = 28 | R1b    // SET_WRITE_PROT
	cmd29 = 29 | R1b    // CLR_WRITE_PROT
	cmd30 = 30 | R1     // SEND_WRITE_PROT
	cmd32 = 30 | R1     // ERASE_WR_BLK_START
	cmd33 = 33 | R1     // ERASE_WR_BLK_END
	cmd38 = 38 | R1b    // ERASE
	cmd40 = 40 | R1     // TODO: See DPS spec.
	cmd42 = 42 | R1     // LOCK_UNLOCK
	cmd52 = 52 | R5     // IO_RW_DIRECT
	cmd53 = 53 | R5     // IO_RW_EXTENDED
	cmd55 = 55 | R1     // APP_CMD
	cmd56 = 56 | R1     // GEN_CMD

	acmd6  = ACMD | 6 | R1  // SET_BUS_WIDTH
	acmd13 = ACMD | 13 | R1 // SD_STATUS
	acmd22 = ACMD | 22 | R1 // SEND_NUM_WR_BLOCKS
	acmd23 = ACMD | 23 | R1 // SET_WR_BLK_ERASE_COUNT
	acmd41 = ACMD | 41 | R3 // SD_SEND_OP_COND
	acmd42 = ACMD | 42 | R1 // SET_CLR_CARD_DETECT
	acmd51 = ACMD | 51 | R1 // SEND_SCR
)

// CMD0 (GO_IDLE_STATE, -) performs software reset and sets a card into Idle
// State.
func CMD0() (Command, uint32) {
	return cmd0, 0
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

// CMD55 (APP_CMD, R1) indicates to the card that the next command is an
// application specific command.
func CMD55(rca uint16) (Command, uint32) {
	return cmd55, uint32(rca) << 16
}
