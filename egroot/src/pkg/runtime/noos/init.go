package noos

func init() {
	initCPU()
	if MaxTasks() > 0 {
		initTasker()
	}
}