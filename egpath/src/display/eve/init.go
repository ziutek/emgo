package eve

import (
	"delay"
	"errors"
)

const chipidAddr = 0x0C0000

type mmap struct {
	ramdl       int
	ramcmd      int
	regdlswap   int
	regintflags int
	regcmdwrite int
	regtracker  int
}

//emgo:const
var eve1 = mmap{
	ramdl:       0x100000,
	ramcmd:      0x108000, // 264 * 4096
	regdlswap:   0x102450,
	regintflags: 0x102498,
	regcmdwrite: 0x1024e8,
	regtracker:  0x109000,
}

//emgo:const
var eve2 = mmap{
	ramdl:       0x300000,
	ramcmd:      0x308000, // 776 * 4096
	regdlswap:   0x302054,
	regintflags: 0x3020a8,
	regcmdwrite: 0x3020fc,
	regtracker:  0x309000,
}

// Register offsets relative to REG_DLSWAP.
const (
	ohcycle  = -40
	orotate  = 4
	opclkpol = 24
	opclk    = 28
	ogpio    = 64
)

// Register offsets relative to REG_INT_FLAGS.
const (
	ointen   = 4
	ointmask = 8
	opwmduty = 44
)

// Register offsets relative to REG_CMD_WRITE.
const (
	ocmdread       = -4
	ocmddl         = 4
	otouchscreenxy = 40
	otouchtagxy    = 44
	otouchtag      = 48
)

// EVE2 bulk write registers.
const (
	regcmdbspace = 0x302574
	regcmdbwrite = 0x302578
)

// Host commands for initialisation.
const (
	cmdActive = 0
	cmdClkExt = 0x44
)

// REG_ID adresses
const (
	ramreg  = 0x102400
	ramreg2 = 0x302000
)

// DisplayConfig contains LCD timing parameters. It seems to be an error in
// datasheet that describes HOFFSET as non-visible part of line (Thf+Thp+Thb)
// and VOFFSET as number of non-visible lines (Tvf+Tvp+Tvb).
type DisplayConfig struct {
	Hcycle  uint16 // Total number of clocks per line.(Th)
	Hsize   uint16 // Active width of LCD display.....(Thd)
	Hsync0  byte   // Start of horizontal sync pulse..(Thf)
	Hsync1  byte   // End of horizontal sync pulse....(Thf+Thp)
	Hoffset byte   // Start of active line............(Thp+Thb)
	ClkPol  byte   // Define active edge of pixel clock.
	Vcycle  uint16 // Total number of lines per scree.(Tv)
	Vsize   uint16 // Active height of LCD display....(Tvd)
	Vsync0  byte   // Start of vertical sync pulse....(Tvf)
	Vsync1  byte   // End of vertical sync pulse......(Tvf+Tvp)
	Voffset byte   // Start of active screen..........(Tvp+Tvb)
	ClkMHz  byte   // Pixel Clock MHz.................(Fclk)
}

//emgo:const
var (
	Default320x240 = DisplayConfig{
		Hcycle: 408, Hsize: 320, Hsync0: 0, Hsync1: 10, Hoffset: 70,
		Vcycle: 263, Vsize: 240, Vsync0: 0, Vsync1: 2, Voffset: 13,
		ClkPol: 0, ClkMHz: 6,
	}
	Default480x272 = DisplayConfig{
		Hcycle: 548, Hsize: 480, Hsync0: 0, Hsync1: 41, Hoffset: 43,
		Vcycle: 292, Vsize: 272, Vsync0: 0, Vsync1: 10, Voffset: 12,
		ClkPol: 1, ClkMHz: 9,
	}
	Default800x480 = DisplayConfig{
		Hcycle: 928, Hsize: 800, Hsync0: 40, Hsync1: 40 + 48, Hoffset: 88,
		Vcycle: 525, Vsize: 480, Vsync0: 13, Vsync1: 13 + 3, Voffset: 32,
		ClkPol: 1, ClkMHz: 30, // KD50G21-40NT-A1
	}
)

type Config struct {
	OutBits uint16 // Bits 8-0 set number of red, green, blue output signals.
	Rotate  byte   // Screen rotation controll.
	Dither  byte   // Dithering controll (note that 0 overrides default 1).
	Swizzle byte   // Control the arrangement of output RGB pins.
	Spread  byte   // Control the color signals spread (for reducing EM noise).
}

// Init initializes EVE and writes first display list. Dcf describes display
// configuration. Cfg describes EVE configuration and can be nil to use reset
// defaults.
func (d *Driver) Init(dcf *DisplayConfig, cfg *Config) error {
	d.dci.SetPDN(0)
	delay.Millisec(20)
	d.dci.SetPDN(1)
	delay.Millisec(20) // Wait 20 ms for internal oscilator and PLL.

	d.width = dcf.Hsize
	d.height = dcf.Vsize

	d.HostCmd(cmdClkExt, 0) // Select external 12 MHz oscilator as clock source.
	d.HostCmd(cmdActive, 0)

	// Read both possible REG_ID locations for max. 300 ms, then check CHIPID.
	for i := 0; i < 30; i++ {
		if d.ReadByte(ramreg) == 0x7C {
			break
		}
		if d.ReadByte(ramreg2) == 0x7C {
			break
		}
		delay.Millisec(10)
	}
	if err := d.Err(true); err != nil {
		return err
	}
	chipid := d.ReadUint32(chipidAddr)
	switch {
	case chipid == 0x10008:
		d.mmap = &eve1
	case 0x11008 <= chipid && chipid <= 0x111308:
		d.mmap = &eve2
	default:
		return errors.New("eve: unknown controller")
	}

	/*
		// Simple triming algorithm if internal oscilator is used.
		for trim := uint32(0); trim <= 31; trim++ {
			d.W(ft80.REG_TRIM).W32(trim)
			if f := curFreq(d); f > 47040000 {
				lcd.W(ft80.REG_FREQUENCY).W32(f)
				break
			}
		}
	*/

	d.SetBacklight(0)
	d.SetIntMask(0)
	d.WriteByte(d.mmap.regintflags+ointen, 1)
	if cfg != nil {
		w := d.W(d.mmap.regdlswap + orotate)
		w.Write32(
			uint32(cfg.Rotate),
			uint32(cfg.OutBits),
			uint32(cfg.Dither),
			uint32(cfg.Swizzle),
			uint32(cfg.Spread),
			uint32(dcf.ClkPol),
			0, // REG_PCLK
		)
	} else {
		w := d.W(d.mmap.regdlswap + opclkpol)
		w.Write32(
			uint32(dcf.ClkPol),
			0, // REG_PCLK
		)
	}
	w := d.W(d.mmap.regdlswap + ohcycle)
	w.Write32(
		uint32(dcf.Hcycle),
		uint32(dcf.Hoffset),
		uint32(dcf.Hsize),
		uint32(dcf.Hsync0),
		uint32(dcf.Hsync1),
		uint32(dcf.Vcycle),
		uint32(dcf.Voffset),
		uint32(dcf.Vsize),
		uint32(dcf.Vsync0),
		uint32(dcf.Vsync1),
	)
	w = d.W(d.mmap.ramdl)
	w.Write32(
		CLEAR|CST,
		DISPLAY,
	)
	d.SwapDL()
	b := d.ReadByte(d.mmap.regdlswap + ogpio)
	d.WriteByte(d.mmap.regdlswap+ogpio, b|0x80) // Set DISP high.
	// Calculate prescaler. +1 causes that the half-way cases are rounded up.
	var presc int
	if d.mmap == &eve1 {
		presc = (48*2 + 1 + int(dcf.ClkMHz)) / (int(dcf.ClkMHz) * 2)
	} else {
		presc = (60*2 + 1 + int(dcf.ClkMHz)) / (int(dcf.ClkMHz) * 2)
	}
	d.WriteByte(d.mmap.regdlswap+opclk, byte(presc)) // Enable PCLK.

	delay.Millisec(20) // Wait for new main clock.

	return d.Err(true)
}
