// +build cortexm0 cortexm3 cortexm4 cortexm4f

extern byte StackExp, StackFrac, StackEnd;

static inline
uint runtime$noos$stackExp() {
	return (uint)&StackExp;
}

static inline
uint runtime$noos$stackFrac() {
	return (uint)&StackFrac;
}

static inline
uint runtime$noos$stackEnd() {
	return (uintptr)&StackEnd;
}