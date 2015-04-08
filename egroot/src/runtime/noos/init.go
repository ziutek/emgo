package noos

import "builtin"

func init() {
	initCPU()
	builtin.Panic = panic_
	builtin.Alloc = alloc
	builtin.MakeChan = makeChan
	builtin.Select = selectComm
	if MaxTasks() > 0 {
		tasker.init()
	}
}