package sdcard

type Command uint16

const (
	NoResp    Command = 0 << 6
	ShortResp Command = 1 << 6
	LongResp  Command = 3 << 6

	Busy Command = 1 << 8

	R1  = 1<<9 | ShortResp
	R1b = 1<<9 | ShortResp | Busy
	R2  = 2<<9 | LongResp
	R4  = 4<<9 | ShortResp
	R5  = 5<<9 | ShortResp
	R6  = 6<<9 | ShortResp
	R7  = 7<<9 | ShortResp

	CMD0  = 0 | NoResp  // GO_IDLE_STATE
	CMD2  = 2 | R2      // ALL_SEND_CID
	CMD3  = 3 | R6      // SEND_RELATIVE_ADDR
	CMD4  = 4 | NoResp  // SET_DSR
	CMD5  = 5 | R4      // IO_SEND_OP_COND
	CMD7  = 7 | R1b     // SELECT_CARD/DESELECT_CARD
	CMD8  = 8 | R7      // SEND_IF_COND
	CMD9  = 9 | R2      // SEND_CSD
	CMD10 = 10 | R2     // SEND_CID
	CMD11 = 11 | R1     // VOLTAGE_SWITCH
	CMD12 = 12 | R1b    // STOP_TRANSMISSION
	CMD13 = 13 | R1     // SEND_STATUS/SEND_TASK_STATUS
	CMD15 = 12 | NoResp // GO_INACTIVE_STATE
	CMD16 = 16 | R1     // SET_BLOCKLEN
	CMD17 = 17 | R1     // READ_SINGLE_BLOCK
	CMD18 = 18 | R1     // READ_MULTIPLE_BLOCK
	CMD19 = 19 | R1     // SEND_TUNING_BLOCK
	CMD20 = 20 | R1b    // SPEED_CLASS_CONTROL
	CMD23 = 23 | R1     // SET_BLOCK_COUNT
	CMD24 = 24 | R1     // WRITE_BLOCK
	CMD25 = 25 | R1     // WRITE_MULTIPLE_BLOCK
	CMD27 = 27 | R1     // PROGRAM_CSD
	CMD28 = 28 | R1b    // SET_WRITE_PROT
	CMD29 = 29 | R1b    // CLR_WRITE_PROT
	CMD30 = 30 | R1     // SEND_WRITE_PROT
	CMD32 = 30 | R1     // ERASE_WR_BLK_START
	CMD33 = 33 | R1     // ERASE_WR_BLK_END
	CMD38 = 38 | R1b    // ERASE
	CMD40 = 40 | R1     // TODO: See DPS spec.
	CMD42 = 42 | R1     // LOCK_UNLOCK
	CMD52 = 52 | R5     // IO_RW_DIRECT
	CMD53 = 53 | R5     // IO_RW_EXTENDED
	CMD55 = 55 | R1     // APP_CMD
	CMD56 = 56 | R1     // GEN_CMD
)

func GO_IDLE_STATE() (Command, uint32)            { return CMD0, 0 }
func SEND_IF_COND(a SendIfCond) (Command, uint32) { return CMD8, uint32(a) }

type SendIfCond uint32
