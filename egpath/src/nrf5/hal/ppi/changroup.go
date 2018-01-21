package ppi

import (
	"nrf5/hal/te"
)

// Group repersents PPI channel group. There are 4 channel groups (6 in nRF52)
// numbered from 0 to 3 (5 in nRF52).
type Group int8

// Channels returns channels that belongs to the group g.
func (g Group) Channels() Channels {
	return Channels(r().chg[g].Load())
}

// SetChannels sets channels that belongs to the group g.
func (g Group) SetChannels(c Channels) {
	r().chg[g].Store(uint32(c))
}

type Task byte

// EN returns task that can be used to enable channel group g.
func (g Group) EN() Task {
	return Task(g * 2)
}

// DIS returns task that can be used to disable channel group g.
func (g Group) DIS() Task {
	return Task(g*2 + 1)
}

func (t Task) Task() *te.Task { return r().Regs.Task(int(t)) }
