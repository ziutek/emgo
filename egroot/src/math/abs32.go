package math

// Abs32 is like Abs but for float32.
func Abs32(x float32) float32 {
	switch {
	case x < 0:
		return -x
	case x == 0:
		return 0
	}
	return x
}
