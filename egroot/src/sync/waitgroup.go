package sync

type WaitGroup waitgroup

func (wg *WaitGroup) Add(delta int) {
	add(wg, delta)
}

func (wg *WaitGroup) Done() {
	add(wg, -1)
}

func (wg *WaitGroup) Wait() {
	wait(wg)
}
