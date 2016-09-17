void
internal$NewTask(void (*f) (), bool lock) {
	uintptr$$uintptr r = internal$Syscall2(0, (uintptr) (f), (uintptr) (lock));
	uintptr e = r._1;
	if (e != 0) {
		panic(INTERFACE(e, &uintptr$$));
	}
}

