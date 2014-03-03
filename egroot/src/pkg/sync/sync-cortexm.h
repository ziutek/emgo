// +build cortexm3 cortexm4 cortexm4f

__attribute__ ((always_inline))
extern inline void sync_Sync() {
	asm volatile ("dsb":::"memory");
}