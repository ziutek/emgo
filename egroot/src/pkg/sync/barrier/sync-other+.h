// +build !cortexm0,!cortexm3,!cortexm4,!cortexm4f

// This isn't formally correct but hardware full memory barrier, followed by
// R/W memory access, should be enough on most platform.
__attribute__ ((always_inline))
extern inline void sync_barrier_Sync() {
	static int mem;
	__sync_synchronize();
	mem++
}