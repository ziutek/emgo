package syscall

//emgo:export
//c:inline
func ftou(func()) uintptr

//emgo:export
//c:inline
func f64tou(func(int64)) uintptr

//emgo:export
//c:inline
func fr64tou(func() int64) uintptr
