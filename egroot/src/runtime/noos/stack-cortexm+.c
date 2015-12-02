// +build cortexm0 cortexm3 cortexm4 cortexm4f

extern byte StackLog2, StackFrac, StackEnd;

static inline
uint runtime$noos$stackLog2() {
	return (uint)&StackLog2;
}

static inline
uint runtime$noos$stackFrac() {
	return (uint)&StackFrac;
}

static inline
uint runtime$noos$stackEnd() {
	return (uintptr)&StackEnd;
}