// +build cortexm0 cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

inline __attribute__((always_inline))
uint16
bits$reverseBytes16(uint16 u) {
	uint16 v;
	asm volatile ("rev16 %0, %1":"=r"(v):"r"(u));
	return v;
}

inline __attribute__((always_inline))
uint32
bits$reverseBytes32(uint32 u) {
	uint32 v;
	asm volatile ("rev %0, %1":"=r"(v):"r"(u));
	return v;
}

inline __attribute__((always_inline))
uint64
bits$reverseBytes64(uint64 u) {
	uint64 v;
	asm volatile (
		"rev %R0, %Q1\n\t"
		"rev %Q0, %R1\n\t"
		:"=&r"(v)
		:"r"(u)
	);
	return v;
}
