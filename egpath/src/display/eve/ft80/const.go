package ft80

import (
	"display/eve"
)

const (
	ACTIVE  eve.HostCmd = 0x00 // Switch mode to Active.
	STANDBY eve.HostCmd = 0x41 // Switch mode to Standby: PLL and Oscillator on.
	SLEEP   eve.HostCmd = 0x42 // Switch mode to Sleep: PLL and Oscillator off.
	PWRDOWN eve.HostCmd = 0x50 // Switch off LDO, Clock, PLL and Oscillator.

	CLKEXT eve.HostCmd = 0x44 // Select PLL external clock source.
	CLK48M eve.HostCmd = 0x62 // Switch PLL output to 48 MHz (default).
	CLK36M eve.HostCmd = 0x61 // Switch PLL output to 36 MHz.

	CORERST eve.HostCmd = 0x68 // Send reset pulse to FT80x core.
)

// Memory map.
const (
	RAM_G          = 0x000000 // 256 KB
	ROM_CHIPID     = 0x0C0000 //   4 B
	ROM_FONT       = 0x0BB23C // 275 KB
	ROM_FONT_ADDR  = 0x0FFFFC //   4 B
	RAM_DL         = 0x100000 //   8 KB
	RAM_PAL        = 0x102000 //   1 KB
	RAM_REG        = 0x102400 // 380 B
	RAM_CMD        = 0x108000 //   4 KB
	RAM_SCREENSHOT = 0x1C2000 //   2 KB
)

// Register Addresses.
const (
	REG_ID         = 0x102400
	REG_FRAMES     = 0x102404
	REG_CLOCK      = 0x102408
	REG_FREQUENCY  = 0x10240c
	REG_RENDERMODE = 0x102410
	REG_SNAPY      = 0x102414
	REG_SNAPSHOT   = 0x102418
	REG_CPURESET   = 0x10241c
	REG_TAP_CRC    = 0x102420
	REG_TAP_MASK   = 0x102424
	REG_HCYCLE     = 0x102428
	REG_HOFFSET    = 0x10242c
	REG_HSIZE      = 0x102430
	REG_HSYNC0     = 0x102434
	REG_HSYNC1     = 0x102438
	REG_VCYCLE     = 0x10243c
	REG_VOFFSET    = 0x102440
	REG_VSIZE      = 0x102444
	REG_VSYNC0     = 0x102448
	REG_VSYNC1     = 0x10244c
	REG_DLSWAP     = 0x102450
	REG_ROTATE     = 0x102454
	REG_OUTBITS    = 0x102458
	REG_DITHER     = 0x10245c
	REG_SWIZZLE    = 0x102460
	REG_CSPREAD    = 0x102464
	REG_PCLK_POL   = 0x102468
	REG_PCLK       = 0x10246c
	REG_TAG_X      = 0x102470
	REG_TAG_Y      = 0x102474
	REG_TAG        = 0x102478
	REG_VOL_PB     = 0x10247c
	REG_VOL_SOUND  = 0x102480
	REG_SOUND      = 0x102484
	REG_PLAY       = 0x102488
	REG_GPIO_DIR   = 0x10248c
	REG_GPIO       = 0x102490

	REG_INT_FLAGS        = 0x102498
	REG_INT_EN           = 0x10249c
	REG_INT_MASK         = 0x1024a0
	REG_PLAYBACK_START   = 0x1024a4
	REG_PLAYBACK_LENGTH  = 0x1024a8
	REG_PLAYBACK_READPTR = 0x1024ac
	REG_PLAYBACK_FREQ    = 0x1024b0
	REG_PLAYBACK_FORMAT  = 0x1024b4
	REG_PLAYBACK_LOOP    = 0x1024b8
	REG_PLAYBACK_PLAY    = 0x1024bc
	REG_PWM_HZ           = 0x1024c0
	REG_PWM_DUTY         = 0x1024c4
	REG_MACRO_0          = 0x1024c8
	REG_MACRO_1          = 0x1024cc

	REG_CMD_READ          = 0x1024e4
	REG_CMD_WRITE         = 0x1024e8
	REG_CMD_DL            = 0x1024ec
	REG_TOUCH_MODE        = 0x1024f0
	REG_TOUCH_ADC_MODE    = 0x1024f4
	REG_TOUCH_CHARGE      = 0x1024f8
	REG_TOUCH_SETTLE      = 0x1024fc
	REG_TOUCH_OVERSAMPLE  = 0x102500
	REG_TOUCH_RZTHRESH    = 0x102504
	REG_TOUCH_RAW_XY      = 0x102508
	REG_TOUCH_RZ          = 0x10250c
	REG_TOUCH_SCREEN_XY   = 0x102510
	REG_TOUCH_TAG_XY      = 0x102514
	REG_TOUCH_TAG         = 0x102518
	REG_TOUCH_TRANSFORM_A = 0x10251c
	REG_TOUCH_TRANSFORM_B = 0x102520
	REG_TOUCH_TRANSFORM_C = 0x102524
	REG_TOUCH_TRANSFORM_D = 0x102528
	REG_TOUCH_TRANSFORM_E = 0x10252c
	REG_TOUCH_TRANSFORM_F = 0x102530

	REG_TOUCH_DIRECT_XY   = 0x102574
	REG_TOUCH_DIRECT_Z1Z2 = 0x102578

	REG_TRACKER = 0x109000
)

// Graphics Engine Commands.
const (
	CMDBUF_SIZE      = 4096
	CMD_APPEND       = 0xffffff1e
	CMD_BGCOLOR      = 0xffffff09
	CMD_BUTTON       = 0xffffff0d
	CMD_CALIBRATE    = 0xffffff15
	CMD_CLOCK        = 0xffffff14
	CMD_COLDSTART    = 0xffffff32
	CMD_DIAL         = 0xffffff2d
	CMD_DLSTART      = 0xffffff00
	CMD_FGCOLOR      = 0xffffff0a
	CMD_GAUGE        = 0xffffff13
	CMD_GETMATRIX    = 0xffffff33
	CMD_GETPTR       = 0xffffff23
	CMD_GRADCOLOR    = 0xffffff34
	CMD_GRADIENT     = 0xffffff0b
	CMD_INFLATE      = 0xffffff22
	CMD_INTERRUPT    = 0xffffff02
	CMD_KEYS         = 0xffffff0e
	CMD_LOADIDENTITY = 0xffffff26
	CMD_LOADIMAGE    = 0xffffff24
	CMD_LOGO         = 0xffffff31
	CMD_MEMCPY       = 0xffffff1d
	CMD_MEMCRC       = 0xffffff18
	CMD_MEMSET       = 0xffffff1b
	CMD_MEMWRITE     = 0xffffff1a
	CMD_MEMZERO      = 0xffffff1c
	CMD_NUMBER       = 0xffffff2e
	CMD_PROGRESS     = 0xffffff0f
	CMD_REGREAD      = 0xffffff19
	CMD_ROTATE       = 0xffffff29
	CMD_SCALE        = 0xffffff28
	CMD_SCREENSAVER  = 0xffffff2f
	CMD_SCROLLBAR    = 0xffffff11
	CMD_SETFONT      = 0xffffff2b
	CMD_SETMATRIX    = 0xffffff2a
	CMD_SKETCH       = 0xffffff30
	CMD_SLIDER       = 0xffffff10
	CMD_SNAPSHOT     = 0xffffff1f
	CMD_SPINNER      = 0xffffff16
	CMD_STOP         = 0xffffff17
	CMD_SWAP         = 0xffffff01
	CMD_TEXT         = 0xffffff0c
	CMD_TOGGLE       = 0xffffff12
	CMD_TRACK        = 0xffffff2c
	CMD_TRANSLATE    = 0xffffff27
)

// Display list commands to be embedded in Graphics Processor.
const (
	DL_ALPHA_FUNC      = 0x09000000 // requires OR'd arguments
	DL_BITMAP_HANDLE   = 0x05000000 // requires OR'd arguments
	DL_BITMAP_LAYOUT   = 0x07000000 // requires OR'd arguments
	DL_BITMAP_SIZE     = 0x08000000 // requires OR'd arguments
	DL_BITMAP_SOURCE   = 0x01000000 // requires OR'd arguments
	DL_BITMAP_TFORM_A  = 0x15000000 // requires OR'd arguments
	DL_BITMAP_TFORM_B  = 0x16000000 // requires OR'd arguments
	DL_BITMAP_TFORM_C  = 0x17000000 // requires OR'd arguments
	DL_BITMAP_TFORM_D  = 0x18000000 // requires OR'd arguments
	DL_BITMAP_TFORM_E  = 0x19000000 // requires OR'd arguments
	DL_BITMAP_TFORM_F  = 0x1A000000 // requires OR'd arguments
	DL_BLEND_FUNC      = 0x0B000000 // requires OR'd arguments
	DL_BEGIN           = 0x1F000000 // requires OR'd arguments
	DL_CALL            = 0x1D000000 // requires OR'd arguments
	DL_CLEAR           = 0x26000000 // requires OR'd arguments
	DL_CELL            = 0x06000000 // requires OR'd arguments
	DL_CLEAR_RGB       = 0x02000000 // requires OR'd arguments
	DL_CLEAR_STENCIL   = 0x11000000 // requires OR'd arguments
	DL_CLEAR_TAG       = 0x12000000 // requires OR'd arguments
	DL_COLOR_A         = 0x0F000000 // requires OR'd arguments
	DL_COLOR_MASK      = 0x20000000 // requires OR'd arguments
	DL_COLOR_RGB       = 0x04000000 // requires OR'd arguments
	DL_DISPLAY         = 0x00000000
	DL_END             = 0x21000000
	DL_JUMP            = 0x1E000000 // requires OR'd arguments
	DL_LINE_WIDTH      = 0x0E000000 // requires OR'd arguments
	DL_MACRO           = 0x25000000 // requires OR'd arguments
	DL_POINT_SIZE      = 0x0D000000 // requires OR'd arguments
	DL_RESTORE_CONTEXT = 0x23000000
	DL_RETURN          = 0x24000000
	DL_SAVE_CONTEXT    = 0x22000000
	DL_SCISSOR_SIZE    = 0x1C000000 // requires OR'd arguments
	DL_SCISSOR_XY      = 0x1B000000 // requires OR'd arguments
	DL_STENCIL_FUNC    = 0x0A000000 // requires OR'd arguments
	DL_STENCIL_MASK    = 0x13000000 // requires OR'd arguments
	DL_STENCIL_OP      = 0x0C000000 // requires OR'd arguments
	DL_TAG             = 0x03000000 // requires OR'd arguments
	DL_TAG_MASK        = 0x14000000 // requires OR'd arguments
	DL_VERTEX2F        = 0x40000000 // requires OR'd arguments
	DL_VERTEX2II       = 0x02000000 // requires OR'd arguments
)

// Command and register value options.
const (
	CLR_COL              = 0x4
	CLR_STN              = 0x2
	CLR_TAG              = 0x1
	DECR                 = 4
	DECR_WRAP            = 7
	DLSWAP_DONE          = 0
	DLSWAP_FRAME         = 2
	DLSWAP_LINE          = 1
	DST_ALPHA            = 3
	EDGE_STRIP_A         = 7
	EDGE_STRIP_B         = 8
	EDGE_STRIP_L         = 6
	EDGE_STRIP_R         = 5
	EQUAL                = 5
	GEQUAL               = 4
	GREATER              = 3
	INCR                 = 3
	INCR_WRAP            = 6
	INT_CMDEMPTY         = 32
	INT_CMDFLAG          = 64
	INT_CONVCOMPLETE     = 128
	INT_PLAYBACK         = 16
	INT_SOUND            = 8
	INT_SWAP             = 1
	INT_TAG              = 4
	INT_TOUCH            = 2
	INVERT               = 5
	KEEP                 = 1
	L1                   = 1
	L4                   = 2
	L8                   = 3
	LEQUAL               = 2
	LESS                 = 1
	LINEAR_SAMPLES       = 0
	LINES                = 3
	LINE_STRIP           = 4
	NEAREST              = 0
	NEVER                = 0
	NOTEQUAL             = 6
	ONE                  = 1
	ONE_MINUS_DST_ALPHA  = 5
	ONE_MINUS_SRC_ALPHA  = 4
	OPT_CENTER           = 1536 // = 0x6000
	OPT_CENTERX          = 512  // = 0x0200
	OPT_CENTERY          = 1024 // = 0x0400
	OPT_FLAT             = 256  // = 0x0100
	OPT_MONO             = 1
	OPT_NOBACK           = 4096 // = 0x1000
	OPT_NODL             = 2
	OPT_NOHANDS          = 49152 // = 0xC168
	OPT_NOHM             = 16384 // = 0x4000
	OPT_NOPOINTER        = 16384 // = 0x4000
	OPT_NOSECS           = 32768 // = 0x8000
	OPT_NOTICKS          = 8192  // = 0x2000
	OPT_RIGHTX           = 2048  // = 0x0800
	OPT_SIGNED           = 256   // = 0x0100
	PALETTED             = 8
	PLAYCOLOR            = 0x00a0a080
	POINTS               = 2
	RECTS                = 9
	REPEAT               = 1
	REPLACE              = 2
	RGB332               = 4
	RGB565               = 7
	SRC_ALPHA            = 2
	TEXT8X8              = 9
	TEXTVGA              = 10
	TOUCHMODE_CONTINUOUS = 3
	TOUCHMODE_FRAME      = 2
	TOUCHMODE_OFF        = 0
	TOUCHMODE_ONESHOT    = 1
	AW_SAMPLES           = 1
	ZERO                 = 0
)
