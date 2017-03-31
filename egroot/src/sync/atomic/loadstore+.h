//#define _MMLD __ATOMIC_ACQUIRE
//#define _MMST __ATOMIC_RELEASE

#define _MMLD   __ATOMIC_RELAXED
#define _MMST   __ATOMIC_RELAXED

extern inline __attribute__((always_inline))
uint32
sync$atomic$LoadUint32(uint32 * addr) {
	return __atomic_load_n(addr, _MMLD);
}

extern inline __attribute__((always_inline))
int_
sync$atomic$LoadInt(int_ * addr) {
	return __atomic_load_n(addr, _MMLD);
}

extern inline __attribute__((always_inline))
uintptr
sync$atomic$LoadUintptr(uintptr * addr) {
	return __atomic_load_n(addr, _MMLD);
}

extern inline __attribute__((always_inline))
unsafe$Pointer
sync$atomic$LoadPointer(unsafe$Pointer * addr) {
	return __atomic_load_n(addr, _MMLD);
}

extern inline __attribute__((always_inline))
void
sync$atomic$StoreUint32(uint32 * addr, uint32 val) {
	return __atomic_store_n(addr, val, _MMST);
}

extern inline __attribute__((always_inline))
void
sync$atomic$StoreInt(int_ * addr, int_ val) {
	return __atomic_store_n(addr, val, _MMST);
}

extern inline __attribute__((always_inline))
void
sync$atomic$StoreUintptr(uintptr * addr, uintptr val) {
	return __atomic_store_n(addr, val, _MMST);
}

extern inline __attribute__((always_inline))
void
sync$atomic$StorePointer(unsafe$Pointer * addr, unsafe$Pointer val) {
	return __atomic_store_n(addr, val, _MMST);
}

#undef _MMLD
#undef _MMST
