// +build noos

#define GO(call, lock) do {                    \
	void func() {                              \
		call;                                  \
		internal$Syscall1(internal$KILLTASK, 0); \
	}                                          \
	internal$NewTask(func, lock);               \
} while(0)

static inline void
goready() {
	internal$Syscall0(internal$TASKUNLOCK);
}

/*
static void
newTask(void (*func)(), bool lock) {
	uintptr$$uintptr r = internal$Syscall2(
		internal$NEWTASK, (uintptr)(func), lock
	);
	if (r._1 != 0) {
			panic(INTERFACE(r._1, nil));
	}
}
*/
