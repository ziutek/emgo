package eve

// Display list commands.
const (
	ALPHA_FUNC         = 0x09000000 // Arg: func<<8 | ref
	BEGIN              = 0x1F000000 // Arg: prim
	BITMAP_HANDLE      = 0x05000000 // Arg: handle
	BITMAP_LAYOUT      = 0x07000000 // Arg: linestride<<9 | height
	BITMAP_LAYOUT_H    = 0x28000000 // Arg: linestride<<2 | height (EVE2)
	BITMAP_SIZE        = 0x08000000 // Arg: opt<<18 | width<<9 | height
	BITMAP_SIZE_H      = 0x29000000 // Arg: width<<2 | height (EVE2)
	BITMAP_SOURCE      = 0x01000000 // Arg: addr
	BITMAP_TRANSFORM_A = 0x15000000 // Arg: a
	BITMAP_TRANSFORM_B = 0x16000000 // Arg: b
	BITMAP_TRANSFORM_C = 0x17000000 // Arg: c
	BITMAP_TRANSFORM_D = 0x18000000 // Arg: d
	BITMAP_TRANSFORM_E = 0x19000000 // Arg: e
	BITMAP_TRANSFORM_F = 0x1A000000 // Arg: f
	BLEND_FUNC         = 0x0B000000 // Arg: src<<3 | dst
	CALL               = 0x1D000000 // Arg: dest
	CELL               = 0x06000000 // Arg: cell
	CLEAR              = 0x26000000 // Arg: cst
	CLEAR_COLOR_A      = 0x0F000000 // Arg: alpha
	CLEAR_COLOR_RGB    = 0x02000000 // Arg: red<<16 | blue<<8 | green
	CLEAR_STENCIL      = 0x11000000 // Arg: s
	CLEAR_TAG          = 0x12000000 // Arg: t
	COLOR_A            = 0x10000000 // Arg: alpha
	COLOR_MASK         = 0x20000000 // Arg: rgba
	COLOR_RGB          = 0x04000000 // Arg: red<<16 | blue<<8 | green
	DISPLAY            = 0x00000000
	END                = 0x21000000
	JUMP               = 0x1E000000 // Arg: dest
	LINE_WIDTH         = 0x0E000000 // Arg: width
	MACRO              = 0x25000000 // Arg: m
	NOP                = 0x2D000000
	PALETTE_SOURCE     = 0x2A000000 // Arg: addr (EVE2)
	POINT_SIZE         = 0x0D000000 // Arg: size
	RESTORE_CONTEXT    = 0x23000000
	RETURN             = 0x24000000
	SAVE_CONTEXT       = 0x22000000
	SCISSOR_SIZE       = 0x1C000000 // Arg: width<<12 | height
	SCISSOR_XY         = 0x1B000000 // Arg: x<<11 | y
	STENCIL_FUNC       = 0x0A000000 // Arg: func<<16 | ref<<8 | mask
	STENCIL_MASK       = 0x13000000 // Arg: mask
	STENCIL_OP         = 0x0C000000 // Arg: sfail<<3 | spass
	TAG                = 0x03000000 // Arg: t
	TAG_MASK           = 0x14000000 // Arg: mask
	VERTEX2F           = 0x40000000 // Arg: x<<15 | y
	VERTEX2II          = 0x80000000 // Arg: x<<21 | y<<12 | handle<<7 | cell
	VERTEX_FORMAT      = 0x27000000 // Arg: frac (EVE2)
)

// Alpha/stencil function (ALPHA_FUNC, STENCIL_FUNC).
const (
	NEVER    = 0
	LESS     = 1
	LEQUAL   = 2
	GREATER  = 3
	GEQUAL   = 4
	EQUAL    = 5
	NOTEQUAL = 6
	ALWAYS   = 7
)

// Graphics primitive (BEGIN).
const (
	BITMAPS      = 1
	POINTS       = 2
	LINES        = 3
	LINE_STRIP   = 4
	EDGE_STRIP_R = 5
	EDGE_STRIP_L = 6
	EDGE_STRIP_A = 7
	EDGE_STRIP_B = 8
	RECTS        = 9
)

// Bitmap formats (BITMAP_LAYOUT).
const (
	ARGB1555     = 0
	L1           = 1
	L4           = 2
	L8           = 3
	RGB332       = 4
	ARGB2        = 5
	ARGB4        = 6
	RGB565       = 7
	PALETTED     = 8 // FT80x
	TEXT8X8      = 9
	TEXTVGA      = 10
	BARGRAPH     = 11
	PALETTED565  = 14 // FT81x
	PALETTED4444 = 15 // FT81x
	PALETTED8    = 16 // FT81x
	L2           = 17 // FT81x
)

// Bitmap options (BITMAP_SIZE).
const (
	BILINEAR = 1 << 2
	REPEATX  = 1 << 1
	REPEATY  = 1 << 0
)

// Blending options (BLEND_FUNC).
const (
	ZERO                = 0
	ONE                 = 1
	SRC_ALPHA           = 2
	DST_ALPHA           = 3
	ONE_MINUS_SRC_ALPHA = 4
	ONE_MINUS_DST_ALPHA = 5
)

// Clearing options: cst (CLEAR).
const (
	T   = 1 << 0
	S   = 1 << 1
	C   = 1 << 2
	CS  = C | S
	ST  = S | T
	CT  = C | T
	CST = C | S | T
)

// Color mask: rgba (COLOR_MASK).
const (
	A   = 1 << 0
	B   = 1 << 1
	G   = 1 << 2
	R   = 1 << 3
	RG  = R | G
	GB  = G | B
	RB  = R | B
	RGB = R | G | B
)

// Stencil test actions: sfail, spass.
const (
	//ZERO  = 0 // Alredy defined in blending options.
	KEEP    = 1
	REPLACE = 2
	INCR    = 3
	DECR    = 4
	INVERT  = 5
)

// Graphics Engine Commands.
const (
	CMD_DLSTART    = 0xffffff00
	CMD_SWAP       = 0xffffff01
	CMD_COLDSTART  = 0xffffff32
	CMD_INTERRUPT  = 0xffffff02
	CMD_APPEND     = 0xffffff1e // Arg: addr, num
	CMD_REGREAD    = 0xffffff19 // Arg: addr
	CMD_MEMWRITE   = 0xffffff1a // Arg: addr, num, ...
	CMD_INFLATE    = 0xffffff22 // Arg: addr, ...
	CMD_LOADIMAGE  = 0xffffff24 // Arg: addr, options, ...
	CMD_MEDIAFIFO  = 0xffffff39 // Arg: addr, size (EVE2)
	CMD_PLAYVIDEO  = 0xffffff3a // Arg: options, ... (EVE2)
	CMD_VIDEOSTART = 0xffffff40 // (EVE2)
	CMD_VIDEOFRAME = 0xffffff41 // Arg: dst, ptr (EVE2)
	CMD_MEMCRC     = 0xffffff18 // Arg: addr, num
	CMD_MEMZERO    = 0xffffff1c // Arg: addr, num
	CMD_MEMSET     = 0xffffff1b // Arg: addr, val, num
	CMD_MEMCPY     = 0xffffff1d // Arg: dst, src, num
	CMD_BUTTON     = 0xffffff0d // Arg: x, y, w, h, font, options, ..., 0

	CMD_BGCOLOR      = 0xffffff09
	CMD_CALIBRATE    = 0xffffff15
	CMD_CLOCK        = 0xffffff14
	CMD_DIAL         = 0xffffff2d
	CMD_FGCOLOR      = 0xffffff0a
	CMD_GAUGE        = 0xffffff13
	CMD_GETMATRIX    = 0xffffff33
	CMD_GETPTR       = 0xffffff23
	CMD_GRADCOLOR    = 0xffffff34
	CMD_GRADIENT     = 0xffffff0b
	CMD_KEYS         = 0xffffff0e
	CMD_LOADIDENTITY = 0xffffff26
	CMD_LOGO         = 0xffffff31
	CMD_NUMBER       = 0xffffff2e
	CMD_PROGRESS     = 0xffffff0f
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
	CMD_TEXT         = 0xffffff0c
	CMD_TOGGLE       = 0xffffff12
	CMD_TRACK        = 0xffffff2c
	CMD_TRANSLATE    = 0xffffff27
)

// REG_DLSWAP values.
const (
	DLSWAP_DONE  = 0
	DLSWAP_LINE  = 1
	DLSWAP_FRAME = 2
)

// Interrupt flags.
const (
	INT_SWAP         = 1
	INT_TOUCH        = 2
	INT_TAG          = 4
	INT_SOUND        = 8
	INT_PLAYBACK     = 16
	INT_CMDEMPTY     = 32
	INT_CMDFLAG      = 64
	INT_CONVCOMPLETE = 128
)

// Image/video options (CMD_LOAD_IMAGE, CMD_PLAYVIDEO).
const (
	OPT_MONO       = 1
	OPT_NODL       = 2
	OPT_NOTEAR     = 4  // EVE2
	OPT_FULLSCREEN = 8  // EVE2
	OPT_MEDIAFIFO  = 16 // EVE2
	OPT_SOUND      = 32 // EVE2
)

const (
	DECR_WRAP            = 7
	INCR_WRAP            = 6
	LINEAR_SAMPLES       = 0
	OPT_CENTER           = 1536  // = 0x6000
	OPT_CENTERX          = 512   // = 0x0200
	OPT_CENTERY          = 1024  // = 0x0400
	OPT_FLAT             = 256   // = 0x0100
	OPT_NOBACK           = 4096  // = 0x1000
	OPT_NOHANDS          = 49152 // = 0xC168
	OPT_NOHM             = 16384 // = 0x4000
	OPT_NOPOINTER        = 16384 // = 0x4000
	OPT_NOSECS           = 32768 // = 0x8000
	OPT_NOTICKS          = 8192  // = 0x2000
	OPT_RIGHTX           = 2048  // = 0x0800
	OPT_SIGNED           = 256   // = 0x0100
	PLAYCOLOR            = 0x00a0a080
	REPEAT               = 1
	TOUCHMODE_CONTINUOUS = 3
	TOUCHMODE_FRAME      = 2
	TOUCHMODE_OFF        = 0
	TOUCHMODE_ONESHOT    = 1
	AW_SAMPLES           = 1
)
