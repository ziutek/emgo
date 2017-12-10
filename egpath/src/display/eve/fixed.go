package eve

// Fixed is 32-bit binary fixed-point signed number with scaling factor 1/16 =
// 0.0625. It divides 32-bit word to 28-bit integer part and 4-bit fractional
// part. Fixed can store values from -134217728 to 134217727.9375.
type Fixed int32

func F(i int) Fixed {
	return Fixed(i << 4)
}

func (a Fixed) Int() int {
	return int(a) >> 4
}

func (a Fixed) Round() int {
	return int(a+8) >> 4
}

func (a Fixed) NearEven() int {
	b := a >> 4 & 1
	return int(a+b+7) >> 4
}

func (a Fixed) Mul(b Fixed) Fixed {
	return a * b >> 4
}
