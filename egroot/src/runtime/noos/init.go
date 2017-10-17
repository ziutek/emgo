package noos

import "internal"

func init() {
	initCPU()
	internal.Panic = panic_
	internal.Alloc = alloc
	internal.MakeChan = makeChan
	internal.Select = selectComm
	internal.TimeChan = makeTimeChan()
	if maxTasks() > 0 {
		tasker.init()
	}
}
