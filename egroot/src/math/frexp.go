package math

func Frexpi(f float64) (int64, int) {
	bits := Float64bits(f)
	frac := int64(bits & (1<<52 - 1))
	exp := int(bits>>52) & (1<<11 - 1)
	if exp == 0 {
		exp = 1 - (1023 + 52)
	} else {
		exp -= 1023 + 52
		frac += 1 << 52
	}
	if f < 0 {
		frac = -frac
	}
	return frac, exp
}
