// +build cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

uint32
bits$reverse32(uint32 u) {
	uint32 v;
	asm volatile ("rbit %0, %1":"=r"(v):"r"(u));
	return v;
}

uintptr
bits$reversePtr(uintptr u) {
	uintptr v;
	asm volatile ("rbit %0, %1":"=r"(v):"r"(u));
	return v;
}

uint64
bits$reverse64(uint64 u) {
	uint64 v;
	asm volatile (
		"rbit %R0, %Q1\n\t"
		"rbit %Q0, %R1\n\t"
		:"=&r"(v)
		:"r"(u)
	);
	return v;
}
