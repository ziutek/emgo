extern byte MaxIRQ;

static inline uintptr runtime_noos_maxIRQ() {
	return (int)&MaxIRQ;
}
