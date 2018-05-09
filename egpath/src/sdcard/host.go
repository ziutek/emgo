package sdcard

type Host interface {
	Cmd(cmd Command, arg uint32) (Response, error)
}
