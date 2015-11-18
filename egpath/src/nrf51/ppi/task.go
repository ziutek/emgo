package ppi

import (
	"mmio"

	"nrf51/internal"
)

// Task represents nRF51 task.
type Task struct {
	reg *mmio.U32
}

func GetTask(te *internal.TasksEvents, n int) Task {
	return Task{&te.Tasks[n]}
}

// Trig triggers task.
func (t Task) Trig() {
	t.reg.Store(1)
}
