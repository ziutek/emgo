package led

// Color represents 32-bit alpha-premultipled color, having 8 bits for each of
// alpha, red, green and blue (0xAARRGGBB).
type Color uint32

func RGBA(r, g, b, a byte) Color {
	return Color(a)<<24 | Color(r)<<16 | Color(g)<<8 | Color(b)
}

func RGB(r, g, b byte) Color {
	return 0xFF<<24 | Color(r)<<16 | Color(g)<<8 | Color(b)
}

func (c Color) Alpha() byte {
	return byte(c >> 24)
}

func (c Color) Red() byte {
	return byte(c >> 16)
}

func (c Color) Green() byte {
	return byte(c >> 8)
}

func (c Color) Blue() byte {
	return byte(c)
}

func (c Color) Mul(alpha byte) Color {
	a := byte((uint32(c.Alpha())*uint32(alpha) + 255) >> 8)
	r := byte((uint32(c.Red())*uint32(alpha) + 255) >> 8)
	g := byte((uint32(c.Green())*uint32(alpha) + 255) >> 8)
	b := byte((uint32(c.Blue())*uint32(alpha) + 255) >> 8)
	return RGBA(r, g, b, a)
}
