package syscall

//emgo:export
func ftou(func()) uintptr

//emgo:export
func f64btou(func(int64, bool)) uintptr

//emgo:export
func fr64tou(func() int64) uintptr
