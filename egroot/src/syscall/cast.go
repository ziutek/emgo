package syscall

//emgo:export
//c:static inline
func btou(bool) uintptr

//emgo:export
//c:static inline
func ftou(func()) uintptr

//emgo:export
//c:static inline
func f32tou(func() uint32) uintptr
