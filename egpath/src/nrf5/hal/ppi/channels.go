package ppi

type Channels uint32

// Enabled returns a bitfield where each bit corresponds to one channel. Zero
// means channel is disabled, one means channel is enabled.
func Enabled() Channels {
	return Channels(r().chen.Load())
}

// Enable atomically enables channels sepcified by c.
func (c Channels) Enable() {
	r().chenset.Store(uint32(c))
}

// Disable atomically disables channels sepcified by c.
func (c Channels) Disable() {
	r().chenclr.Store(uint32(c))
}
