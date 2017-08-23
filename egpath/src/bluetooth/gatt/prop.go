package gatt

const (
	Broadcast        = 0x01
	Read             = 0x02
	WriteWithoutResp = 0x04
	Write            = 0x08
	Notify           = 0x10
	Indicate         = 0x20
	AuthSignWrites   = 0x40
	Extended         = 0x80
)
