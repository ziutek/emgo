package mem

func tomask(a uintptr) uintptr {
	b := a - 1
	if a&b != 0 {
		panic("mem: a isn't power of 2")
	}
	return b
}

// AlignUp returns p aligned up to a. a must be power of 2.
func AlignUp(p, a uintptr) uintptr {
	a = tomask(a)
	return (p + a) &^ a
}

// AlignDown returns p aligned down to a. a must be power of 2.
func AlignDown(p, a uintptr) uintptr {
	tomask(a)
	return p &^ (a - 1)
}
