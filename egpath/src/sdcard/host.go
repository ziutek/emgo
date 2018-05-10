package sdcard

type Host interface {
	Cmd(cmd Command, arg uint32) Response
	Err(clear bool) error // Err returns and clears internal error status.
}
