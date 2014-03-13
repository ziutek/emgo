package noos

func alignUp(p, a uintptr) uintptr {
	mask := a - 1
	if p&mask != 0 {
		p = (p + a) &^ mask
	}
	return p
}

func alignDown(p, a uintptr) uintptr {
	return p &^ (a - 1)
}

/*func checkAlignment(p, a uintptr) {
	if p & (a-1) != 0 {
		panic("unaligned address")
	}
}*/
