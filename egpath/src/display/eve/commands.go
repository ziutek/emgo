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
	PALETTED     = 8 // EVE1
	TEXT8X8      = 9
	TEXTVGA      = 10
	BARGRAPH     = 11
	PALETTED565  = 14 // EVE2
	PALETTED4444 = 15 // EVE2
	PALETTED8    = 16 // EVE2
	L2           = 17 // EVE2
	ARGB8        = 32 // EVE2 (CMD_SNAPSHOT2)
)

const DEFAULT = 0

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
	CMD_DLSTART      = 0xFFFFFF00
	CMD_SWAP         = 0xFFFFFF01
	CMD_COLDSTART    = 0xFFFFFF32
	CMD_INTERRUPT    = 0xFFFFFF02
	CMD_APPEND       = 0xFFFFFF1E // Arg: addr, num
	CMD_REGREAD      = 0xFFFFFF19 // Arg: addr
	CMD_MEMWRITE     = 0xFFFFFF1A // Arg: addr, num, ...
	CMD_INFLATE      = 0xFFFFFF22 // Arg: addr, ...
	CMD_LOADIMAGE    = 0xFFFFFF24 // Arg: addr, options, ...
	CMD_MEDIAFIFO    = 0xFFFFFF39 // Arg: addr, size (EVE2)
	CMD_PLAYVIDEO    = 0xFFFFFF3A // Arg: options, ... (EVE2)
	CMD_VIDEOSTART   = 0xFFFFFF40 // (EVE2)
	CMD_VIDEOFRAME   = 0xFFFFFF41 // Arg: dst, ptr (EVE2)
	CMD_MEMCRC       = 0xFFFFFF18 // Arg: addr, num
	CMD_MEMZERO      = 0xFFFFFF1C // Arg: addr, num
	CMD_MEMSET       = 0xFFFFFF1B // Arg: addr, val, num
	CMD_MEMCPY       = 0xFFFFFF1D // Arg: dst, src, num
	CMD_BUTTON       = 0xFFFFFF0D // Arg: x, y, w, h, font, options, ..., 0
	CMD_CLOCK        = 0xFFFFFF14 // Arg: x, y, r, options, h, m, s, ms
	CMD_FGCOLOR      = 0xFFFFFF0A // Arg: rgb
	CMD_BGCOLOR      = 0xFFFFFF09 // Arg: rgb
	CMD_GRADCOLOR    = 0xFFFFFF34 // Arg: rgb
	CMD_GAUGE        = 0xFFFFFF13 // Arg: x, y, r, optns, major, minor, val, max
	CMD_GRADIENT     = 0xFFFFFF0B // Arg: x0, y0, rgb0, x1, y1, rgb1
	CMD_KEYS         = 0xFFFFFF0E // Arg: x, y, w, h, font, options, ..., 0
	CMD_PROGRESS     = 0xFFFFFF0F // Arg: x, y, w, h, options, val, max
	CMD_SCROLLBAR    = 0xFFFFFF11 // Arg: x, y, w, h, options, val, size, max
	CMD_SLIDER       = 0xFFFFFF10 // Arg: x, y, w, h, options, val, max
	CMD_DIAL         = 0xFFFFFF2d // Arg: x, y, r, options, val
	CMD_TOGGLE       = 0xFFFFFF12 // Arg: x, y, w, font, options, state, ..., 0
	CMD_TEXT         = 0xFFFFFF0C // Arg: x, y, font, options, ..., 0
	CMD_SETBASE      = 0xFFFFFF38 // Arg: base (EVE2)
	CMD_NUMBER       = 0xFFFFFF2E // Arg: x, y int, font, options, n
	CMD_LOADIDENTITY = 0xFFFFFF26
	CMD_SETMATRIX    = 0xFFFFFF2A // Arg: a, b, c, d, e, f
	CMD_GETMATRIX    = 0xFFFFFF33
	CMD_GETPTR       = 0xFFFFFF23
	CMD_GETPROPS     = 0xFFFFFF25
	CMD_SCALE        = 0xFFFFFF28 // Arg: sx, sy
	CMD_ROTATE       = 0xFFFFFF29 // Arg: a
	CMD_TRANSLATE    = 0xFFFFFF27 // Arg: tx, ty
	CMD_CALIBRATE    = 0xFFFFFF15
	CMD_SETROTATE    = 0xFFFFFF36 // Arg: r (EVE2)
	CMD_SPINNER      = 0xFFFFFF16 // Arg: x, y, style, scale
	CMD_SCREENSAVER  = 0xFFFFFF2F
	CMD_SKETCH       = 0xFFFFFF30 // Arg: x, y, w, h, addr, format
	CMD_STOP         = 0xFFFFFF17
	CMD_SETFONT      = 0xFFFFFF2B // Arg: font, addr
	CMD_SETFONT2     = 0xFFFFFF3B // Arg: font, addr, firstchar (EVE2)
	CMD_SETSCRATCH   = 0xFFFFFF3C // Arg: handle (EVE2)
	CMD_ROMFONT      = 0xFFFFFF3F // Arg: font, romslot (EVE2)
	CMD_TRACK        = 0xFFFFFF2C // Arg: x, y, w, h, tag
	CMD_SNAPSHOT     = 0xFFFFFF1F // Arg: addr
	CMD_SNAPSHOT2    = 0xFFFFFF37 // Arg: format, addr, x, y, w, h (EVE2)
	CMD_SETBITMAP    = 0xFFFFFF43 // Arg: addr, format, w, h (EVE2)
	CMD_LOGO         = 0xFFFFFF31
	CMD_CSKETCH      = 0xFFFFFF35 // Arg: x, y, w, h, addr, format, freq (FT801)
)

// REG_DLSWAP values.
const (
	DLSWAP_DONE  = 0
	DLSWAP_LINE  = 1
	DLSWAP_FRAME = 2
)

// Interrupt flags.
const (
	INT_SWAP         = 1 << 0 // Display list swap occurred.
	INT_TOUCH        = 1 << 1 // Touch detected.
	INT_TAG          = 1 << 2 // Touch-screen tag value change.
	INT_SOUND        = 1 << 3 // Sound effect ended.
	INT_PLAYBACK     = 1 << 4 // Audio playback ended.
	INT_CMDEMPTY     = 1 << 5 // Command FIFO empty.
	INT_CMDFLAG      = 1 << 6 // Command FIFO flag.
	INT_CONVCOMPLETE = 1 << 7 // Touch-screen conversions completed.
)

// Image/video options (CMD_LOAD_IMAGE, CMD_PLAYVIDEO).
const (
	OPT_RGB565     = 0
	OPT_MONO       = 1 << 0
	OPT_NODL       = 1 << 1
	OPT_NOTEAR     = 1 << 2 // EVE2
	OPT_FULLSCREEN = 1 << 3 // EVE2
	OPT_MEDIAFIFO  = 1 << 4 // EVE2
	OPT_SOUND      = 1 << 5 // EVE2
)

// Widget options.
const (
	OPT_FLAT    = 1 << 8
	OPT_SIGNED  = 1 << 8
	OPT_CENTERX = 1 << 9
	OPT_CENTERY = 1 << 10
	OPT_RIGHTX  = 1 << 11
	OPT_CENTER  = OPT_CENTERX | OPT_CENTERY
)

// Clock, gauge options.
const (
	OPT_NOBACK    = 1 << 12
	OPT_NOTICKS   = 1 << 13
	OPT_NOPOINTER = 1 << 14
	OPT_NOHM      = 1 << 14
	OPT_NOSECS    = 1 << 15
	OPT_NOHANDS   = OPT_NOHM | OPT_NOSECS
)

const (
	DECR_WRAP            = 7
	INCR_WRAP            = 6
	LINEAR_SAMPLES       = 0
	PLAYCOLOR            = 0x00a0a080
	REPEAT               = 1
	TOUCHMODE_CONTINUOUS = 3
	TOUCHMODE_FRAME      = 2
	TOUCHMODE_OFF        = 0
	TOUCHMODE_ONESHOT    = 1
	AW_SAMPLES           = 1
)

func MakeRGB(r, g, b int) uint32 {
	return uint32(r&255<<16 | g&255<<8 | b&255)
}
