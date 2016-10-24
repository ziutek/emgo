package syscall

//emgo:export
//c:static inline
func ftou(func()) uintptr

//emgo:export
//c:static inline
func f64btou(func(int64, bool)) uintptr

//emgo:export
//c:static inline
func fr64tou(func() int64) uintptr
