// +build noos

#define GO(call, lock) do {                   \
	void func() {                             \
		call;                                 \
		builtin$Syscall1(builtin$DELTASK, 0); \
	}                                         \
	newTask(func, lock);                      \
} while(0)

static
void newTask(void (*func)(), bool lock) {
	uintptr$$uintptr r = builtin$Syscall2(
		builtin$NEWTASK, (uintptr)(func), lock
	);
	if (r._1 != 0) {
			panic(INTERFACE(r._1, 0));
	}
}

static inline
void goready() {
	builtin$Syscall0(builtin$TASKUNLOCK);
}