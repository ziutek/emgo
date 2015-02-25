package noos

import "builtin"

func init() {
	initCPU()
	builtin.Alloc = alloc
	builtin.MakeChan = makeChan
	builtin.Select = selectComm
	if MaxTasks() > 0 {
		tasker.init()
	}
}