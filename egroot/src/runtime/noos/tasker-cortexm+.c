// +build cortexm0 cortexm3 cortexm4 cortexm4f

extern byte IRTExp;

static inline
uint runtime$noos$irtExp() {
	return (uint)&IRTExp;
}