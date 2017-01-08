package ili9341

// These commands can be used to directly interract with display controller
// using DCI.
const (
	NOP     = 0x00
	SWRESET = 0x01
	SPLIN   = 0x10
	SLPOUT  = 0x11
	DISPOFF = 0x28
	DISPON  = 0x29
	RAMWR   = 0x2C
	MADCTL  = 0x36
	PIXSET  = 0x3A
	CASET   = 0x2A
	PASET   = 0x2B
)

// Reset invokes Software Reset, 8-bit command.
func (d *Display) Reset() {
	d.dci.Cmd(SWRESET)
}

// SlpIn invokes Enter Sleep Mode, 8-bit command.
func (d *Display) SlpIn() {
	d.dci.Cmd(SPLIN)
}

// SlpOut invokes Sleep Out, 8-bit command. SlpOut usually requires 120 ms
// delay before next command.
func (d *Display) SlpOut() {
	d.dci.Cmd(SLPOUT)
}

// DispOn invokes Display ON, 8-bit command.
func (d *Display) DispOn() {
	d.dci.Cmd(DISPON)
}
