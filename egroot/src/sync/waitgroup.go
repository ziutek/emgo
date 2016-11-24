package sync

type WaitGroup waitgroup

func (wg *WaitGroup) Add(delta int) {
	wg.add(delta)
}

func (wg *WaitGroup) Done() {
	wg.add(-1)
}

func (wg *WaitGroup) Wait() {
	wg.wait()
}
