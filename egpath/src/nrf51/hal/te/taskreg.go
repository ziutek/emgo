package te

import (
	"mmio"
	"unsafe"

	"nrf51/hal/internal"
)

// TaskReg represents peripheral registers that are used to trigger events..
type TaskReg struct {
	u32 mmio.U32
}

// GetTaskReg is for internal use.
func GetTaskReg(ph *internal.Pheader, n int) *TaskReg {
	return (*TaskReg)(unsafe.Pointer(&ph.Tasks[n]))
}

// Trigger starts action corresponding to task t.
func (r *TaskReg) Trigger() {
	r.u32.Store(1)
}
