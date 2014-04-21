#define  MMODEL __ATOMIC_RELAXED

__attribute__ ((always_inline))
extern inline 
bool sync$atomic$CompareAndSwapInt32(int32 *addr, int32 old, int32 new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, MMODEL, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
bool sync$atomic$CompareAndSwapInt64(int64 *addr, int64 old, int64 new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, MMODEL, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
bool sync$atomic$CompareAndSwapPointer(unsafe$Pointer *addr, unsafe$Pointer old, unsafe$Pointer new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, MMODEL, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
bool sync$atomic$CompareAndSwapUint32(uint32 *addr, uint32 old, uint32 new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, MMODEL, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
bool sync$atomic$CompareAndSwapUint64(uint64 *addr, uint64 old, uint64 new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, MMODEL, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
bool sync$atomic$CompareAndSwapUintptr(uintptr *addr, uintptr old, uintptr new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, MMODEL, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
int32 sync$atomic$AddInt32(int32 *addr, int32 delta) {
	return __atomic_add_fetch(addr, delta, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uint32 sync$atomic$AddUint32(uint32 *addr, uint32 delta) {
	return __atomic_add_fetch(addr, delta, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
int64 sync$atomic$AddInt64(int64 *addr, int64 delta) {
	return __atomic_add_fetch(addr, delta, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uint64 sync$atomic$AddUint64(uint64 *addr, uint64 delta) {
	return __atomic_add_fetch(addr, delta, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uintptr sync$atomic$AddUintptr(uintptr *addr, uintptr delta) {
	return __atomic_add_fetch(addr, delta, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
int32 sync$atomic$OrInt32(int32 *addr, int32 mask) {
	return __atomic_or_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uint32 sync$atomic$OrUint32(uint32 *addr, uint32 mask) {
	return __atomic_or_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
int64 sync$atomic$OrInt64(int64 *addr, int64 mask) {
	return __atomic_or_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uint64 sync$atomic$OrUint64(uint64 *addr, uint64 mask) {
	return __atomic_or_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uintptr sync$atomic$OrUintptr(uintptr *addr, uintptr mask) {
	return __atomic_or_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
int32 sync$atomic$AndInt32(int32 *addr, int32 mask) {
	return __atomic_and_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uint32 sync$atomic$AndUint32(uint32 *addr, uint32 mask) {
	return __atomic_and_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
int64 sync$atomic$AndInt64(int64 *addr, int64 mask) {
	return __atomic_and_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uintptr sync$atomic$AndUintptr(uintptr *addr, uintptr mask) {
	return __atomic_and_fetch(addr, mask, MMODEL);
}


__attribute__ ((always_inline))
extern inline 
uint64 sync$atomic$AndUint64(uint64 *addr, uint64 mask) {
	return __atomic_and_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline
int32 sync$atomic$SwapInt32(int32 *addr, int32 new) {
	return __atomic_exchange_n(addr, new, MMODEL);
}

__attribute__ ((always_inline))
extern inline
int64 sync$atomic$SwapInt64(int64 *addr, int64 new) {
	return __atomic_exchange_n(addr, new, MMODEL);
}

__attribute__ ((always_inline))
extern inline
unsafe$Pointer sync$atomic$SwapPointer(unsafe$Pointer *addr, unsafe$Pointer new) {
	return __atomic_exchange_n(addr, new, MMODEL);
}

__attribute__ ((always_inline))
extern inline
uint32 sync$atomic$SwapUint32(uint32 *addr, uint32 new) {
	return __atomic_exchange_n(addr, new, MMODEL);
}

__attribute__ ((always_inline))
extern inline
uint64 sync$atomic$SwapUint64(uint64 *addr, uint64 new) {
	return __atomic_exchange_n(addr, new, MMODEL);
}

__attribute__ ((always_inline))
extern inline
uintptr sync$atomic$SwapUintptr(uintptr *addr, uintptr new) {
	return __atomic_exchange_n(addr, new, MMODEL);
}
