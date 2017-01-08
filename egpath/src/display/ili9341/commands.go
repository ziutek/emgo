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

// Reset invokes Software Reset command (8-bit).
func (d *Display) Reset() {
	d.dci.Cmd(SWRESET)
}

// SlpIn invokes Enter Sleep Mode command (8-bit).
func (d *Display) SlpIn() {
	d.dci.Cmd(SPLIN)
}

// SlpOut invokes Sleep Out command (8-bit). SlpOut usually requires 120 ms
// delay before next command.
func (d *Display) SlpOut() {
	d.dci.Cmd(SLPOUT)
}

// DispOn invokes Display ON command (8-bit).
func (d *Display) DispOn() {
	d.dci.Cmd(DISPON)
}

// MAD is a bitmask that describes memory access direction.
type MAD byte

const (
	MH  MAD = 1 << 2 // Horizontal refresh order.
	BGR MAD = 1 << 3 // RGB-BGR order.
	ML  MAD = 1 << 4 // Vertical refresh order.
	MV  MAD = 1 << 5 // Row/column exchange.
	MX  MAD = 1 << 6 // Column address order.
	MY  MAD = 1 << 7 // Row address order.
)

// MADCtl invokes Memory Access Control command (8-bit).
func (d *Display) MADCtl(mad MAD) {
	d.dci.Cmd(MADCTL)
	d.dci.Byte(byte(mad))
}

// PixFmt describes pixel format.
type PixFmt byte

const (
	PF16 PixFmt = 0x55 // 16-bit 565 pixel format.
	PF18 PixFmt = 0x66 // 18-bit 666 pixel format.
)

// PixSet invokes Pixel Format Set command (8-bit).
func (d *Display) PixSet(pf PixFmt) {
	d.dci.Cmd(PIXSET)
	d.dci.Byte(byte(pf))
}
