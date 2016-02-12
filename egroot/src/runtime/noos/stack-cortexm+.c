// +build cortexm0 cortexm3 cortexm4 cortexm4f

extern byte StackTaskLog2, StackTaskFrac, StackEnd;

static inline
uint runtime$noos$stackTaskLog2() {
	return (uint)&StackTaskLog2;
}

static inline
uint runtime$noos$stackTaskFrac() {
	return (uint)&StackTaskFrac;
}

static inline
uint runtime$noos$stackEnd() {
	return (uintptr)&StackEnd;
}