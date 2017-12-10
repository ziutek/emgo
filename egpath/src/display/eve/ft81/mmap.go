package ft81

// Memory map.
const (
	RAM_G         = 0x000000 // 1024 KB
	RAM_CHIPID    = 0x0C0000 //    4 B
	ROM_FONT      = 0x1E0000 // 1152 KB
	ROM_FONT_ADDR = 0x2FFFFC //    4 B
	RAM_DL        = 0x300000 //    8 KB
	RAM_REG       = 0x302000 //    4 KB
	RAM_CMD       = 0x308000 //    4 KB
)
