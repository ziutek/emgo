.syntax unified

// func HostIO(cmd int, p unsafe.Pointer) int
.global main$HostIO

.thumb_func
main$HostIO:
	bkpt  0xAB
	bx    lr


