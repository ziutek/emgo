// +build cortexm0 cortexm3 cortexm4 cortexm4f

extern byte StackExp, StackFrac, StackEnd;

static inline
uint runtime_noos_stackExp() {
	return (uint)&StackExp;
}

static inline
uint runtime_noos_stackFrac() {
	return (uint)&StackFrac;
}

static inline
uint runtime_noos_stackEnd() {
	return (uintptr)&StackEnd;
}