package syscall

//emgo:export
//c:static inline
func btou(bool) uintptr

//emgo:export
//c:static inline
func ftou(func()) uintptr

//emgo:export
//c:static inline
func f64tou(func(int64)) uintptr

//emgo:export
//c:static inline
func fr64tou(func() int64) uintptr
