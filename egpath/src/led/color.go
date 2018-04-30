package led

// Color represents 32-bit alpha-premultipled color, having 8 bits for each of
// alpha, red, green and blue (0xAARRGGBB).
type Color uint32

// RGBA returns Color for alpha-premultipled r, g, b (r<=a && g<=a && b<=a). To
// obtain Color from not alpha-premultipled R, G, B use RGB(R, G, B).Mask(a).
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

func (c Color) Scale(alpha byte) Color {
	// Consider replace the original formula (c*alpha + 127) / 255 with cheaper
	// (c*m + 255) >> 8 which avoids division but still (c*255 + 255) >> 8 == c.

	m := uint(alpha)
	r := (uint(c.Red())*m + 127) / 255
	g := (uint(c.Green())*m + 127) / 255
	b := (uint(c.Blue())*m + 127) / 255
	a := (uint(c.Alpha())*m + 127) / 255

	return Color(a<<24 | r<<16 | g<<8 | b)
}

func (c1 Color) Blend(c2 Color, alpha byte) Color {
	r1 := uint(c1.Red())
	g1 := uint(c1.Green())
	b1 := uint(c1.Blue())
	a1 := uint(c1.Alpha())
	r2 := uint(c2.Red())
	g2 := uint(c2.Green())
	b2 := uint(c2.Blue())
	a2 := uint(c2.Alpha())
	m2 := uint(alpha)

	m1 := (255*255 - a2*m2)
	m2 *= 255
	r := (r1*m1 + r2*m2 + 32512) / (255 * 255)
	g := (g1*m1 + g2*m2 + 32512) / (255 * 255)
	b := (b1*m1 + b2*m2 + 32512) / (255 * 255)
	a := (a1*m1 + a2*m2 + 32512) / (255 * 255)

	return Color(a<<24 | r<<16 | g<<8 | b)
}
