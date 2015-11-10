package syscall

//emgo:export
//c:static inline
func b2u(bool) uintptr

//emgo:export
//c:static inline
func f2u(func()) uintptr
