package sdcard

type Command byte

const (
	NoResp    Command = 0 << 6
	ShortResp Command = 1 << 6
	LongResp  Command = 3 << 6

	CMD0  = 0 | NoResp     // GO_IDLE_STATE                -
	CMD2  = 2 | LongResp   // ALL_SEND_CID                 R2
	CMD3  = 3 | ShortResp  // SEND_RELATIVE_ADDR           R6
	CMD4  = 4 | NoResp     // SET_DSR                      -
	CMD5  = 5 | ShortResp  // IO_SEND_OP_COND              R4
	CMD7  = 7 | ShortResp  // SELECT_CARD/DESELECT_CARD    R1b
	CMD8  = 8 | ShortResp  // SEND_IF_COND                 R7
	CMD9  = 9 | LongResp   // SEND_CSD                     R2
	CMD10 = 10 | LongResp  // SEND_CID                     R2
	CMD11 = 11 | ShortResp // VOLTAGE_SWITCH               R1
	CMD12 = 12 | ShortResp // STOP_TRANSMISSION            R1b
	CMD13 = 13 | ShortResp // SEND_STATUS/SEND_TASK_STATUS R1
	CMD15 = 12 | NoResp    // GO_INACTIVE_STATE            -
	CMD16 = 16 | ShortResp // SET_BLOCKLEN                 R1
	CMD17 = 17 | ShortResp // READ_SINGLE_BLOCK            R1
	CMD18 = 18 | ShortResp // READ_MULTIPLE_BLOCK          R1
	CMD19 = 19 | ShortResp // SEND_TUNING_BLOCK            R1
	CMD20 = 20 | ShortResp // SPEED_CLASS_CONTROL          R1b
	CMD23 = 23 | ShortResp // SET_BLOCK_COUNT              R1
	CMD24 = 24 | ShortResp // WRITE_BLOCK                  R1
	CMD25 = 25 | ShortResp // WRITE_MULTIPLE_BLOCK         R1
	CMD27 = 27 | ShortResp // PROGRAM_CSD                  R1
	CMD28 = 28 | ShortResp // SET_WRITE_PROT               R1b
	CMD29 = 29 | ShortResp // CLR_WRITE_PROT               R1b
	CMD30 = 30 | ShortResp // SEND_WRITE_PROT              R1
	CMD32 = 30 | ShortResp // ERASE_WR_BLK_START           R1
	CMD33 = 33 | ShortResp // ERASE_WR_BLK_END             R1
	CMD38 = 38 | ShortResp // ERASE                        R1b
	CMD40 = 40 | ShortResp // TODO: See DPS spec.          R1
	CMD42 = 42 | ShortResp // LOCK_UNLOCK                  R1
	CMD52 = 52 | ShortResp // IO_RW_DIRECT                 R5
	CMD53 = 53 | ShortResp // IO_RW_EXTENDED               R5
	CMD55 = 55 | ShortResp // APP_CMD                      R1
	CMD56 = 56 | ShortResp // GEN_CMD                      R1
)

func GO_IDLE_STATE() (Command, uint32)            { return CMD0, 0 }
func SEND_IF_COND(a SendIfCond) (Command, uint32) { return CMD8, uint32(a) }

type SendIfCond uint32
