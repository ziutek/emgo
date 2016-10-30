package te

import (
	"mmio"
)

// Task represents a peripheral register that is used to trigger events.
type Task struct {
	u32 mmio.U32
}

// Trigger starts action corresponding to task t.
func (r *Task) Trigger() {
	r.u32.Store(1)
}
