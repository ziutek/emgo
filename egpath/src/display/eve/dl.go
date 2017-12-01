package eve

// DL provides convenient way to write Display List commands. Every command is
// a function call, so for better performance or lower RAM usage, use raw Writer
// with many Display List commands in array/slice.
type DL struct {
	Writer
}

// AlphaFunc sets the alpha test function.
func (dl DL) AlphaFunc(fun byte, ref int) {
	dl.wr32(ALPHA_FUNC | uint32(fun)<<8 | uint32(ref&0xFF))
}

// Begin begins drawing a graphics primitive.
func (dl DL) Begin(prim byte) {
	dl.wr32(BEGIN | uint32(prim))
}

// BitmapHandle selscts the bitmap handle.
func (dl DL) BitmapHandle(handle byte) {
	dl.wr32(BITMAP_HANDLE | uint32(handle))
}

// BitmapLayout sets the bitmap memory format and layout for the current handle.
func (dl DL) BitmapLayout(format byte, linestride, height int) {
	l := uint32(linestride) & 1023
	h := uint32(height) & 511
	dl.wr32(BITMAP_LAYOUT | uint32(format)<<19 | l<<9 | h)
	if linestride > 1023 || height > 511 {
		// BUG?: Does BITMAP_LAYOUT zeros bits set by previous BITMAP_LAYOUT_H?
		l = uint32(linestride) >> 10 & 3
		h = uint32(height) >> 9 & 3
		dl.wr32(BITMAP_LAYOUT_H | l<<2 | h)
	}
}

// BitmapSize sets the screen drawing of bitmaps for the current handle.
func (dl DL) BitmapSize(options byte, width, height int) {
	w := uint32(width) & 511
	h := uint32(height) & 511
	dl.wr32(BITMAP_SIZE | uint32(options)<<18 | w<<9 | h)
	if width > 511 || height > 511 {
		// BUG?: Does BITMAP_SIZE clears bits set by previous BITMAP_SIZE_H?
		w = uint32(width) >> 9 & 3
		h = uint32(height) >> 9 & 3
		dl.wr32(BITMAP_SIZE_H | w<<2 | h)
	}
}

// BitmapSource sets the source address of bitmap data in graphics memory RAM_G.
func (dl DL) BitmapSource(addr int) {
	dl.wr32(BITMAP_SOURCE | uint32(addr)&0x3FFFFF)
}

// BitmapTransA sets the A coefficient of the bitmap transform matrix.
func (dl DL) BitmapTransformA(s8_8 int) {
	dl.wr32(BITMAP_TRANSFORM_A | uint32(s8_8)&0x1FFFF)
}

// BitmapTransformB sets the B coefficient of the bitmap transform matrix.
func (dl DL) BitmapTransformB(s8_8 int) {
	dl.wr32(BITMAP_TRANSFORM_B | uint32(s8_8)&0x1FFFF)
}

// BitmapTransformC sets the C coefficient of the bitmap transform matrix.
func (dl DL) BitmapTransformC(s8_8 int) {
	dl.wr32(BITMAP_TRANSFORM_C | uint32(s8_8)&0x1FFFF)
}

// BitmapTransformD sets the D coefficient of the bitmap transform matrix.
func (dl DL) BitmapTransformD(s8_8 int) {
	dl.wr32(BITMAP_TRANSFORM_D | uint32(s8_8)&0x1FFFF)
}

// BitmapTransformE sets the E coefficient of the bitmap transform matrix.
func (dl DL) BitmapTransformE(s8_8 int) {
	dl.wr32(BITMAP_TRANSFORM_E | uint32(s8_8)&0x1FFFF)
}

// BitmapTransformF sets the F coefficient of the bitmap transform matrix.
func (dl DL) BitmapTransformF(s8_8 int) {
	dl.wr32(BITMAP_TRANSFORM_F | uint32(s8_8)&0x1FFFF)
}
