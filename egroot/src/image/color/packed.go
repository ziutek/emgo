package color

// RGB16 represents 16-bit (565) color.
type RGB16 uint16

func (c RGB16) RGBA() (r, g, b, a uint32) {
	v := uint32(c)
	r = v & (0x1f << 11)
	r |= r>>5 | r>>10 | r>>15
	g = v & (0x3f << 5)
	g = g<<5 | g>>1 | g>>7
	b = v & 0x1f
	b = b<<11 | b<<6 | b<<1 | b>>4
	a = 0xffff
	return
}

// RGBA32 represents 32-bit (8888) alpha-premultiplied color.
type RGBA32 uint32

func (c RGBA32) RGBA() (r, g, b, a uint32) {
	v := uint32(c)
	r = v >> 24
	r |= r << 8
	g = v >> 16 & 0xff
	g |= g << 8
	b = v >> 8 & 0xff
	b |= b << 8
	a = v & 0xff
	a |= a << 8
	return
}
