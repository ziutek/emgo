package strconv

func panicBuffer() {
	panic("strconv: buffer too short")
}

func panicBase() {
	panic("strconv: unsupported base")
}

func Btoa(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// FormatBool stores text representation of b in buf using format specified by
// base: 
//	|base| == 1: 0 / 1,
//  |base| == 2: fales / true.
// Unused portion of the buffer is filed with spaces.
// If base > 0 then formatted value is right-justified and FormatBool returns
// offset to its first char. If base < 0 then formatted value is left-justified
// and FormatBool returns its length.
func FormatBool(buf []byte, b bool, base int) int {
	left := base < 0
	if left {
		base = -base
	}
	if base != 1 && base != 2 {
		panicBase()
	}
	blen := 4*base - 3
	if len(buf) < blen {
		panicBuffer()
	}
	var str string
	if b {
		str = "1true"[base-1:]
	} else {
		str = "0false"[base-1:]
	}
	str = str[:blen]
	var i, n, end int
	if left {
		i = len(str)
		n = i
		end = len(buf)
		copy(buf, str)
	} else {
		end = len(buf) - len(str)
		n = end
		copy(buf[end:], str)
	}
	for i < end {
		buf[i] = ' '
		i++
	}
	return n
}
