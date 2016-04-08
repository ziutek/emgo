// +build linux

package internal

//c:static inline
func Syscall1(trap, a1 uintptr) uintptr

//c:static inline
func Syscall2(trap, a1, a2 uintptr) uintptr

//c:static inline
func Syscall3(trap, a1, a2, a3 uintptr) uintptr

//c:static inline
func Syscall4(trap, a1, a2, a3, a4 uintptr) uintptr

//c:static inline
func Syscall5(trap, a1, a2, a3, a4, a5 uintptr) uintptr

//c:static inline
func Syscall6(trap, a1, a2, a3, a4, a5, a6 uintptr) uintptr
