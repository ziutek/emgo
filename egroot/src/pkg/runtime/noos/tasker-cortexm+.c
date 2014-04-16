// +build cortexm0 cortexm3 cortexm4 cortexm4f

extern byte EVTExp;

static inline
uint runtime$noos$evtExp() {
	return (uint)&EVTExp;
}