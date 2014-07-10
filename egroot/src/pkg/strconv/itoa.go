package strconv

const digits = "0123456789abcdefghijklmnopqrstuvwxyz"

// Utoa converts u to string and returns offset to most significant digit.
func Utoa(buf []byte, u uint32, base int) int {
	if base < 2 || base > len(digits) {
		panic("strconv: illegal base")
	}
	b := uint32(base)
	n := len(buf)
	for {
		if n == 0 {
			panic("strconv: buffer too short")
		}
		if u == 0 {
			break
		}
		n--
		buf[n] = digits[u%b]
		u /= b
	}
	if n == len(buf) {
		n--
		buf[n] = '0'
	}
	for k := 0; k < n; k++ {
		buf[k] = ' '
	}
	return n
}

// Itoa converts i to string and returns offset to most significant digit or sign.
func Itoa(buf []byte, i int32, base int) int {
	if i >= 0 {
		return Utoa(buf, uint32(i), base)
	}
	if len(buf) == 0 {
		panic("strconv: buffer too short")
	}
	n := Utoa(buf[1:], uint32(-i), base)
	buf[n] = '-'
	return n
}
