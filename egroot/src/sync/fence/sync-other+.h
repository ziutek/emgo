// +build !cortexm0,!cortexm3,!cortexm4,!cortexm4f

// This isn't formally correct but hardware full memory fence, followed by
// R/W memory access, should be enough on most platform.
void
sync$fence$Sync() {
	static int mem;
	__sync_synchronize();
	mem++;
}
