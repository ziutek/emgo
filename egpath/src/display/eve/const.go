package eve

// Graphics Engine Commands.
const (
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

// Display list commands.
const (
	ALPHA_FUNC         = 0x09000000 // requires OR'd arguments
	BITMAP_HANDLE      = 0x05000000 // requires OR'd arguments
	BITMAP_LAYOUT      = 0x07000000 // requires OR'd arguments
	BITMAP_LAYOUT_H    = 0x28000000 // requires OR'd arguments, FT81x
	BITMAP_SIZE        = 0x08000000 // requires OR'd arguments
	BITMAP_SIZE_H      = 0x29000000 // requires OR'd arguments
	BITMAP_SOURCE      = 0x01000000 // requires OR'd arguments
	BITMAP_TRANSFORM_A = 0x15000000 // requires OR'd arguments
	BITMAP_TRANSFORM_B = 0x16000000 // requires OR'd arguments
	BITMAP_TRANSFORM_C = 0x17000000 // requires OR'd arguments
	BITMAP_TRANSFORM_D = 0x18000000 // requires OR'd arguments
	BITMAP_TRANSFORM_E = 0x19000000 // requires OR'd arguments
	BITMAP_TRANSFORM_F = 0x1A000000 // requires OR'd arguments
	BLEND_FUNC         = 0x0B000000 // requires OR'd arguments
	BEGIN              = 0x1F000000 // requires OR'd arguments
	CALL               = 0x1D000000 // requires OR'd arguments
	CLEAR              = 0x26000000 // requires OR'd arguments
	CELL               = 0x06000000 // requires OR'd arguments
	CLEAR_RGB          = 0x02000000 // requires OR'd arguments
	CLEAR_STENCIL      = 0x11000000 // requires OR'd arguments
	CLEAR_TAG          = 0x12000000 // requires OR'd arguments
	COLOR_A            = 0x0F000000 // requires OR'd arguments
	COLOR_MASK         = 0x20000000 // requires OR'd arguments
	COLOR_RGB          = 0x04000000 // requires OR'd arguments
	DISPLAY            = 0x00000000
	END                = 0x21000000
	JUMP               = 0x1E000000 // requires OR'd arguments
	LINE_WIDTH         = 0x0E000000 // requires OR'd arguments
	MACRO              = 0x25000000 // requires OR'd arguments
	POINT_SIZE         = 0x0D000000 // requires OR'd arguments
	RESTORE_CONTEXT    = 0x23000000
	RETURN             = 0x24000000
	SAVE_CONTEXT       = 0x22000000
	SCISSOR_SIZE       = 0x1C000000 // requires OR'd arguments
	SCISSOR_XY         = 0x1B000000 // requires OR'd arguments
	STENCIL_FUNC       = 0x0A000000 // requires OR'd arguments
	STENCIL_MASK       = 0x13000000 // requires OR'd arguments
	STENCIL_OP         = 0x0C000000 // requires OR'd arguments
	TAG                = 0x03000000 // requires OR'd arguments
	TAG_MASK           = 0x14000000 // requires OR'd arguments
	VERTEX2F           = 0x40000000 // requires OR'd arguments
	VERTEX2II          = 0x02000000 // requires OR'd arguments
)

// Alpha function (ALPHA_FUNC).
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
	WRAPX    = 1 << 1
	WRAPY    = 1 << 0
)

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
	LINEAR_SAMPLES       = 0
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
	PLAYCOLOR            = 0x00a0a080
	REPEAT               = 1
	REPLACE              = 2
	SRC_ALPHA            = 2
	TOUCHMODE_CONTINUOUS = 3
	TOUCHMODE_FRAME      = 2
	TOUCHMODE_OFF        = 0
	TOUCHMODE_ONESHOT    = 1
	AW_SAMPLES           = 1
	ZERO                 = 0
)
