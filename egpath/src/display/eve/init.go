package eve

const chipid = 0x0C0000

type mmap struct {
	dl     int
	cmd    int
	dlswap int
}

//emgo:const
var eve1 = mmap{
	dl:     0x100000,
	cmd:    0x108000,
	dlswap: 0x102450,
}

//emgo:const
var eve1 = mmap{
	dl:     0x300000,
	cmd:    0x308000,
	dlswap: 0x302054,
}

// Register offsets relative to REG_DLSWAP.
const (
	hcycleOffset  = -40
	cspreadOffset = 20
	pclkOffset    = 28
)

type DisplayConfig struct {
	Hcycle  uint16 // Total number of clocks per line.
	Hsize   uint16 // Active width of LCD display.
	Hoffset uint16  // Start of active line.
	Hsync0  byte   // Start of horizontal sync pulse.
	Hsync1  byte   // End of horizontal sync pulse.
	Vcycle  uint16 // Total number of lines per screen.
	Vsize   uint16 // Active height of LCD display.
	Voffset uint16 // Start of active screen.
	Vsync0  byte   // Start of vertical sync pulse.
	Vsync1  byte   // End of vertical sync pulse.
	PclkDiv byte   // Pixel Clock divider.
	PclkPol byte   // Define active edge of pixel clock.
	CSpreed byte   // Color signals spread, reduces EM noise
}
