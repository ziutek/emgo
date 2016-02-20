package dma

type DMA struct {
	*registers
}

func (p DMA) EnableClock(lp bool) {
	enableClock(p, lp)
}

func (p DMA) DisableClock() {
	disableClock(p)
}

// Channel returns channel number n (p.Channel(1) returns first channel).
func (p DMA) Channel(n int) Channel {
	return makeChannel(p)
}

type Channel channel
