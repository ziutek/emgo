package ppi

import (
	"nrf5/hal/te"
)

// ChanGroup repersents channel group. There is 6 channel groups numbered from
// 0 to 5.
type ChanGroup byte

// Channels returns channels that belongs to the group c.
func (g ChanGroup) Channels() Channels {
	return Channels(r().chg[g].Load())
}

// SetChannels sets channels that belongs to the group c.
func (g ChanGroup) SetChannels(c Channels) {
	r().chg[g].Store(uint32(c))
}

// EN returns task that can be used to enable channel group g.
func (g ChanGroup) EN() *te.Task {
	return r().Task(int(g) * 2)
}

// DIS returns task that can be used to disable channel group g.
func (g ChanGroup) DIS() *te.Task {
	return r().Task(int(g)*2 + 1)
}
