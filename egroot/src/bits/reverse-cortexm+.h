// +build cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

inline __attribute__((always_inline))
uint32
bits$reverse32(uint32 u) {
	uint32 v;
	asm volatile ("rbit %0, %1":"=r"(v):"r"(u));
	return v;
}

inline __attribute__((always_inline))
uintptr
bits$reversePtr(uintptr u) {
	uintptr v;
	asm volatile ("rbit %0, %1":"=r"(v):"r"(u));
	return v;
}

inline __attribute__((always_inline))
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
