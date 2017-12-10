package ft80

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
