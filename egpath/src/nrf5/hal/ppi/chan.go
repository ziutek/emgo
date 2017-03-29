package ppi

import (
	"unsafe"

	"nrf5/hal/te"
)

// Chan represents PPI channel. There are 31 channels numbered from 0 to 31.
// Channels from 20 to 31 are pre-programmed.
type Chan byte

// Pre-programmed channels.
const (
	TIMER0_COMPARE0__RADIO_TXEN    Chan = 20
	TIMER0_COMPARE0__RADIO_RXEN    Chan = 21
	TIMER0_COMPARE1__RADIO_DISABLE Chan = 22
	RADIO_BCMATCH__AAR_START       Chan = 23
	RADIO_READY__CCM_KSGEN         Chan = 24
	RADIO_ADDRESS__CCM_CRYPT       Chan = 25
	RADIO_ADDRESS__TIMER0_CAPTURE1 Chan = 26
	RADIO_END__TIMER0_CAPTURE2     Chan = 27
	RTC0_COMPARE0__RADIO_TXEN      Chan = 28
	RTC0_COMPARE0__RADIO_RXEN      Chan = 29
	RTC0_COMPARE0__TIMER0_CLEAR    Chan = 30
	RTC0_COMPARE0__TIMER0_START    Chan = 31
)

func (c Chan) Mask() Channels {
	return Channels(1) << c
}

// Enabled reports whether channel c is enabled.
func (c Chan) Enabled() bool {
	return Enabled()&c.Mask() != 0
}

// Enable atomically enables channel c.
func (c Chan) Enable() {
	c.Mask().Enable()
}

// Enable atomically disables channel c.
func (c Chan) Disable() {
	c.Mask().Disable()
}

// EEP returns the value of Event End Point register for channel c.
func (c Chan) EEP() *te.Event {
	return (*te.Event)(unsafe.Pointer(uintptr(r().ch[c].eep.Load())))
}

// SetEEP sets the value of Event End Point register for channel c.
func (c Chan) SetEEP(e *te.Event) {
	r().ch[c].eep.Store(uint32(uintptr(unsafe.Pointer(e))))
}

// TEP returns the value of Task End Point register for channel c.
func (c Chan) TEP() *te.Task {
	return (*te.Task)(unsafe.Pointer(uintptr(r().ch[c].tep.Load())))
}

// SetTEP sets the value of Task End Point register for channel c.
func (c Chan) SetTEP(t *te.Task) {
	r().ch[c].tep.Store(uint32(uintptr(unsafe.Pointer(t))))
}

// FTEP returns the value of Fork Task End Point register for channel c. nRF52.
func (c Chan) FTEP() *te.Task {
	return (*te.Task)(unsafe.Pointer(uintptr(r().forktep[c].Load())))
}

// SetFTEP sets the value of Fork Task End Point register for channel c. nRF52.
func (c Chan) SetFTEP(t *te.Task) {
	r().forktep[c].Store(uint32(uintptr(unsafe.Pointer(t))))
}
