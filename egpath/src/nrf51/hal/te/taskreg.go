package te

import (
	"mmio"
)

// TaskReg represents peripheral registers that are used to trigger events..
type TaskReg struct {
	u32 mmio.U32
}

// Trigger starts action corresponding to task t.
func (r *TaskReg) Trigger() {
	r.u32.Store(1)
}
