package eve

// DL allows to write Display List commands.
type DL Writer

func (dl DL) wr(cmd uint32) {
	Writer(dl).wr32(cmd)
}

// AlphaFunc sets the alpha test function.
func (dl DL) AlphaFunc(fun byte, ref int) {
	dl.wr(ALPHA_FUNC | uint32(fun)<<8 | uint32(ref&0xFF))
}

// Begin begins drawing a graphics primitive.
func (dl DL) Begin(prim byte) {
	dl.wr(BEGIN | uint32(prim))
}

// BitmapHandle selscts the bitmap handle.
func (dl DL) BitmapHandle(handle byte) {
	dl.wr(BITMAP_HANDLE | uint32(handle))
}

// BitmapLayout sets the bitmap memory format and layout for the current handle.
func (dl DL) BitmapLayout(format byte, linestride, height int) {
	a := uint32(linestride) & 0x3FF
	b := uint32(height) & 0x1FF
	dl.wr(BITMAP_LAYOUT | uint32(format)<<19 | a<<9 | b)
	if linestride > 1023 || height > 512 {
		a = uint32(linestride) >> 10 & 3
		b = uint32(height) >> 9 & 3
		dl.wr(BITMAP_LAYOUT_H | a<<2 | b)
	}
}
