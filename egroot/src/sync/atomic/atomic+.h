// +build !cortexm0

//#define _MMLD   __ATOMIC_ACQUIRE
//#define _MMST   __ATOMIC_RELEASE
//#define _MMLDST __ATOMIC_ACQ_REL

#define _MMLD   __ATOMIC_RELAXED
#define _MMST   __ATOMIC_RELAXED
#define _MMLDST __ATOMIC_RELAXED


extern inline __attribute__((always_inline))
bool
sync$atomic$compareAndSwapUint32(uint32 * addr, uint32 old, uint32 new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, _MMLDST, _MMLD);
}

extern inline __attribute__((always_inline))
bool
sync$atomic$compareAndSwapInt(int_ * addr, int_ old, int_ new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, _MMLDST, _MMLD);
}

extern inline __attribute__((always_inline))
bool
sync$atomic$compareAndSwapUintptr(uintptr * addr, uintptr old, uintptr new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, _MMLDST, _MMLD);
}

extern inline __attribute__((always_inline))
bool
sync$atomic$compareAndSwapPointer(unsafe$Pointer * addr, unsafe$Pointer old,
	unsafe$Pointer new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, _MMLDST, _MMLD);
}


extern inline __attribute__((always_inline))
int_
sync$atomic$swapInt(int_ * addr, int_ new) {
	return __atomic_exchange_n(addr, new, _MMLDST);
}

extern inline __attribute__((always_inline))
unsafe$Pointer
sync$atomic$swapPointer(unsafe$Pointer * addr, unsafe$Pointer new) {
	return __atomic_exchange_n(addr, new, _MMLDST);
}

extern inline __attribute__((always_inline))
uint32
sync$atomic$swapUint32(uint32 * addr, uint32 new) {
	return __atomic_exchange_n(addr, new, _MMLDST);
}

extern inline __attribute__((always_inline))
uintptr
sync$atomic$swapUintptr(uintptr * addr, uintptr new) {
	return __atomic_exchange_n(addr, new, _MMLDST);
}


extern inline __attribute__((always_inline))
uint32
sync$atomic$addUint32(uint32 * addr, uint32 delta) {
	return __atomic_add_fetch(addr, delta, _MMLDST);
}

extern inline __attribute__((always_inline))
int_
sync$atomic$addInt(int_ * addr, int_ delta) {
	return __atomic_add_fetch(addr, delta, _MMLDST);
}

extern inline __attribute__((always_inline))
uintptr
sync$atomic$addUintptr(uintptr * addr, uintptr delta) {
	return __atomic_add_fetch(addr, delta, _MMLDST);
}


extern inline __attribute__((always_inline))
uint32
sync$atomic$orUint32(uint32 * addr, uint32 mask) {
	return __atomic_or_fetch(addr, mask, _MMLDST);
}

extern inline __attribute__((always_inline))
uintptr
sync$atomic$orUintptr(uintptr * addr, uintptr mask) {
	return __atomic_or_fetch(addr, mask, _MMLDST);
}


extern inline __attribute__((always_inline))
uint32
sync$atomic$andUint32(uint32 * addr, uint32 mask) {
	return __atomic_and_fetch(addr, mask, _MMLDST);
}

extern inline __attribute__((always_inline))
uintptr
sync$atomic$andUintptr(uintptr * addr, uintptr mask) {
	return __atomic_and_fetch(addr, mask, _MMLDST);
}


extern inline __attribute__((always_inline))
uint32
sync$atomic$xorUint32(uint32 * addr, uint32 mask) {
	return __atomic_xor_fetch(addr, mask, _MMLDST);
}

extern inline __attribute__((always_inline))
uintptr
sync$atomic$xorUintptr(uintptr * addr, uintptr mask) {
	return __atomic_xor_fetch(addr, mask, _MMLDST);
}


#undef _MMLD
#undef _MMST
#undef _MMLDST
