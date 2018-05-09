package sdcard

type Command byte

const (
	NoResp    Command = 0 << 6
	ShortResp Command = 1 << 6
	LongResp  Command = 3 << 6

	CMD0 = 0 | NoResp
	CMD8 = 8 | ShortResp

	GO_IDLE_STATE = CMD0
	SEND_IF_COND  = CMD8
)
