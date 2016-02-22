// +build noos

package internal

//c:inline
func b2u(bool) uintptr

//c:inline
func f2u(func()) uintptr

func NewTask(f func(), lock bool) {
	if _, e := Syscall2(NEWTASK, f2u(f), b2u(lock)); e != 0 {
		panic(e)
	}
}
