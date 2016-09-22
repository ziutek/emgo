// +build noos

void
internal$NewTask(void (*f) (), bool lock) {
	uintptr$$uintptr r = internal$Syscall2(0, (uintptr) (f), (uintptr) (lock));
	uintptr e = r._1;
	if (e != 0) {
		panic(INTERFACE(e, &uintptr$$));
	}
}

#define GO(call, lock) do {                      \
	void func() {                                \
		call;                                    \
		internal$Syscall1(internal$KILLTASK, 0); \
	}                                            \
	internal$NewTask(func, lock);                \
} while(0)

static inline void
goready() {
	internal$Syscall0(internal$TASKUNLOCK);
}
