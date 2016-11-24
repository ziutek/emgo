package rtos

// EventFlag can be used by tasks and ISRs to report the occurence of events and
// used by tasks (but not ISRs) to waiting for them. The methods of EventFlag
// are synchronization operations with acquire or release semantic. The zero
// value of EventFlag is cleared flag ready to use.
//
// Typical usage of EventFlag looks like:
//
//	for {
//		if !flag.Wait(1, deadline) {
//			handleTimeout()
//			continue
//		}
//		// We can miss some events here as long as flag is not cleared.
//		flag.Reset(0)
//		handleEvent()
//	}
//
// or
//
//	data = 3456 // Prepare some data for DMA.
//	flag.Reset(0)
//	startDMA() // It must call fence.W() before set DMA MMIO start bit.
//	done := flag.Wait(1, deadline) // ISR will signal DMA transfer complete.
type EventFlag eventFlag

// Signal sets the flag to the least significant bit of val and tries to wake up
// all tasks that wait for this value of f. To ensure that waiter will
// definitely wake up, flag must not been modified before waking him. Signal has
// release semantic.
func (f *EventFlag) Signal(val int) {
	f.signal(val)
}

// Wait waits until flag will be set to the least significant bit of val or will
// deadline occur. It returns false in case of timeout, true otherwise.
// Deadline == 0 means no deadline. Wait does not change the flag after return.
// Wait has acquire semantic.
func (f *EventFlag) Wait(val int, deadline int64) bool {
	return f.wait(val, deadline)
}

// Reset sets the flag to the least significant bit of val. It does not try
// to wake up any task but tasks can wake up spontaneously if the flag has
// proper value. Reset can be thought of as Signal but without trying to waking
// up waiters, so it has release semantic.
func (f *EventFlag) Reset(val int) {
	f.reset(val)
}

// Value returns the current value of flag. It can be thought of as Wait with
// deadline in the past (but more efficient), so Value has acquire semantic.
func (f *EventFlag) Value() int {
	return f.value()
}

// WaitEvent waits until at least one from no more than 32 different flags will
// be set to the least significant bit of val or will deadline occur.
// Deadline == 0 means no deadline. WaitEvent returns bitmask that represents
// values of flags (the least significant bit of returned value corresponds to
// flags[0]). WaitEvent has acquire semantic.
func WaitEvent(val int, deadline int64, flags ...*EventFlag) uint32 {
	return waitEvent(val, deadline, flags)
}
