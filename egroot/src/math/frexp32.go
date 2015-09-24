package math

func Frexpi32(f float32) (int32, int) {
	bits := Float32bits(f)
	frac := int32(bits & (1<<23 - 1))
	exp := int(bits>>23) & (1<<8 - 1)
	if exp == 0 {
		exp = 1 - (127 + 23)
	} else {
		exp -= 127 + 23
		frac += 1 << 23
	}
	if f < 0 {
		frac = -frac
	}
	return frac, exp
}
