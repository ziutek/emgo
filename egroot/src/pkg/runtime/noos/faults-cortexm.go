// +build cortexm0 cortexm3 cortexm4 cortexm4f

package noos

import (
	"cortexm"
	"unsafe"
)

func nmiHandler() {
	cortexm.BKPT(0)
}

func hardFaultHandler()

type cfs struct {
	mmfs uint8  `C:"volatile"`
	bfs  uint8  `C:"volatile"`
	ufs  uint16 `C:"volatile"`
}

var cfsr = (*cfs)(unsafe.Pointer(uintptr(0xe000ed28)))

func memFaultHandler() {
	mmfs := cfsr.mmfs
	_ = mmfs
	cortexm.BKPT(2)
}

func busFaultHandler() {
	bfs := cfsr.bfs
	_ = bfs
	cortexm.BKPT(3)
}

func usageFaultHandler() {
	ufs := cfsr.ufs
	pfp := (*stackFrame)(unsafe.Pointer(cortexm.PSP()))
	_, _ = ufs, pfp
	cortexm.BKPT(4)
}
