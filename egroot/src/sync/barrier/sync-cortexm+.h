// +build cortexm0 cortexm3 cortexm4 cortexm4f

void
sync$barrier$Sync() {
	asm volatile ("dsb":::"memory");
}
