// +build !cortexm3,!cortexm4,!cortexm4f

static int mem;

// This isn't formally correct, but function call (never inlined), followed by
// hardware full memory barrier, followed by R/W memory access, followed by
// return from function, should be enough on all real platform. 
void sync_Sync() {
	__sync_synchronize();
	mem++
}