package noos

import (
	"sync/atomic"
	"sync/fence"
)

type counter64 struct {
	cnt [2]uint64
	aba uintptr
}

func (c *counter64) Add(delta uint64) {
	aba := c.aba
	cnt := c.cnt[aba&1]
	aba++
	c.cnt[aba&1] = cnt + delta
	fence.Memory()
	atomic.StoreUintptr(&c.aba, aba)
}

// AtomicSum returns sum of current value of c and value returned by subcnt.
// subcnt should return value of subcounter (typicaly derived from value of
// hardware timer counter) that is used to update c. There should be guarantee
// that any user of c must observe update of c and subcounter reset as one
// atomic operation. This is usually achieved by assign highest priority to
// subcounter interrupt.
func (c *counter64) AtomicSum(subcnt func() uint64) uint64 {
	var v uint64
	aba := atomic.LoadUintptr(&c.aba)
	for {
		fence.Compiler()
		v = subcnt() + c.cnt[aba&1]
		fence.Compiler()
		aba1 := atomic.LoadUintptr(&c.aba)
		if aba == aba1 {
			break
		}
		aba = aba1
	}
	return v
}

var (
	sysCounter counter64
	loadrtc     func() uint64 = timerCnt
	sysTimerHz uint64
)

func uptime() uint64 {
	return muldiv64(sysCounter.AtomicSum(loadrtc), 1e9, sysTimerHz)
}

func muldiv64(x, m, d uint64) uint64 {
	divx := x / d
	modx := x - divx*d
	divm := m / d
	modm := m - divm*d
	return divx*m + modx*divm + modx*modm/d
}
