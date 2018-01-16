package ppi

import (
	"nrf5/hal/te"
)

// Group repersents PPI channel group. There are 6 channel groups numbered from
// 0 to 5.
type Group byte

// Channels returns channels that belongs to the group g.
func (g Group) Channels() Channels {
	return Channels(r().chg[g].Load())
}

// SetChannels sets channels that belongs to the group g.
func (g Group) SetChannels(c Channels) {
	r().chg[g].Store(uint32(c))
}

// EN returns task that can be used to enable channel group g.
func (g Group) EN() *te.Task {
	return r().Task(int(g) * 2)
}

// DIS returns task that can be used to disable channel group g.
func (g Group) DIS() *te.Task {
	return r().Task(int(g)*2 + 1)
}
