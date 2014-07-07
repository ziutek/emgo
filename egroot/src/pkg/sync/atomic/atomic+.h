#define  MMODEL __ATOMIC_RELAXED

__attribute__ ((always_inline))
extern inline 
bool sync$atomic$compareAndSwapInt32(int32 *addr, int32 old, int32 new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, MMODEL, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
bool sync$atomic$compareAndSwapUint32(uint32 *addr, uint32 old, uint32 new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, MMODEL, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
bool sync$atomic$compareAndSwapUintptr(uintptr *addr, uintptr old, uintptr new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, MMODEL, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
bool sync$atomic$compareAndSwapPointer(unsafe$Pointer *addr, unsafe$Pointer old, unsafe$Pointer new) {
	return __atomic_compare_exchange_n(addr, &old, new, false, MMODEL, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
int32 sync$atomic$addInt32(int32 *addr, int32 delta) {
	return __atomic_add_fetch(addr, delta, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uint32 sync$atomic$addUint32(uint32 *addr, uint32 delta) {
	return __atomic_add_fetch(addr, delta, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uintptr sync$atomic$addUintptr(uintptr *addr, uintptr delta) {
	return __atomic_add_fetch(addr, delta, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
int32 sync$atomic$orInt32(int32 *addr, int32 mask) {
	return __atomic_or_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uint32 sync$atomic$orUint32(uint32 *addr, uint32 mask) {
	return __atomic_or_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uintptr sync$atomic$orUintptr(uintptr *addr, uintptr mask) {
	return __atomic_or_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
int32 sync$atomic$andInt32(int32 *addr, int32 mask) {
	return __atomic_and_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uint32 sync$atomic$andUint32(uint32 *addr, uint32 mask) {
	return __atomic_and_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uintptr sync$atomic$andUintptr(uintptr *addr, uintptr mask) {
	return __atomic_and_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
int32 sync$atomic$xorInt32(int32 *addr, int32 mask) {
	return __atomic_xor_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uint32 sync$atomic$xorUint32(uint32 *addr, uint32 mask) {
	return __atomic_xor_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline 
uintptr sync$atomic$xorUintptr(uintptr *addr, uintptr mask) {
	return __atomic_xor_fetch(addr, mask, MMODEL);
}

__attribute__ ((always_inline))
extern inline
int32 sync$atomic$swapInt32(int32 *addr, int32 new) {
	return __atomic_exchange_n(addr, new, MMODEL);
}

__attribute__ ((always_inline))
extern inline
unsafe$Pointer sync$atomic$swapPointer(unsafe$Pointer *addr, unsafe$Pointer new) {
	return __atomic_exchange_n(addr, new, MMODEL);
}

__attribute__ ((always_inline))
extern inline
uint32 sync$atomic$swapUint32(uint32 *addr, uint32 new) {
	return __atomic_exchange_n(addr, new, MMODEL);
}

__attribute__ ((always_inline))
extern inline
uintptr sync$atomic$swapUintptr(uintptr *addr, uintptr new) {
	return __atomic_exchange_n(addr, new, MMODEL);
}

__attribute__ ((always_inline))
extern inline
int32 sync$atomic$loadInt32(int32 *addr) {
	return __atomic_load_n(addr, MMODEL);
}

__attribute__ ((always_inline))
extern inline
uint32 sync$atomic$loadUint32(uint32 *addr) {
	return __atomic_load_n(addr, MMODEL);
}

__attribute__ ((always_inline))
extern inline
uintptr sync$atomic$loadUintptr(uintptr *addr) {
	return __atomic_load_n(addr, MMODEL);
}

__attribute__ ((always_inline))
extern inline
unsafe$Pointer sync$atomic$loadPointer(unsafe$Pointer *addr) {
	return __atomic_load_n(addr, MMODEL);
}

__attribute__ ((always_inline))
extern inline
void sync$atomic$storeInt32(int32 *addr, int32 val) {
	return __atomic_store_n(addr, val, MMODEL);
}

__attribute__ ((always_inline))
extern inline
void sync$atomic$storeUint32(uint32 *addr, uint32 val) {
	return __atomic_store_n(addr, val, MMODEL);
}

__attribute__ ((always_inline))
extern inline
void sync$atomic$storePointer(unsafe$Pointer *addr, unsafe$Pointer val) {
	return __atomic_store_n(addr, val, MMODEL);
}

__attribute__ ((always_inline))
extern inline
void sync$atomic$storeUintptr(uintptr *addr, uintptr val) {
	return __atomic_store_n(addr, val, MMODEL);
}
