package gatt

// Prop is a bitfield that describes properties of characteristic.
type Prop byte

const (
	Broadcast        Prop = 0x01
	Read             Prop = 0x02
	WriteWithoutResp Prop = 0x04
	Write            Prop = 0x08
	Notify           Prop = 0x10
	Indicate         Prop = 0x20
	AuthSignWrites   Prop = 0x40
	Extended         Prop = 0x80
)
