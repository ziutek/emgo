package rtos

// EventFlag can be used by tasks and ISRs to report the occurence of events.
// Any task or ISR can set, clear or check value of the flag. Tasks (but not
// ISRs) can wait for the flag set. Waiting time can be limited by specify
// deadline. The zero value of EventFlag is cleared flag ready to use.
type EventFlag eventFlag

// Set sets the flag.
func (f *EventFlag) Set() {
	flagSet(f)
}

// Clear clears the flag.
func (f *EventFlag) Clear() {
	flagClear(f)
}

// Value returns current value of the flag.
func (f *EventFlag) Val() int {
	return flagVal(f)
}

// Wait waits until flag is set or deadline occured. It returns false in case of
// timeout, true otherwise. Deadline == 0 means no deadline. Wait does not clear
// the flag after return. Typical usage looks like:
//
//	for {
//		if !flag.Wait(deadline) {
//			handleTimeout()
//			continue
//		}
//		// We can miss some events here as long as flag is not cleared.
//		flag.Clear()
//		handleEvent()
//	}
//
// or
//
//	param = 3456  // Prepare some data in memory.
//	flag.Clear()  // This works as memory barrier.
//	start()       // ISR will signal that the work was done by setting flag.
//	done := flag.Wait(deadline)
//
func (f *EventFlag) Wait(deadline int64) bool {
	return flagWait(f, deadline)
}

// WaitEvent waits until at least one from no more than 32 different flags will
// be set or deadline occur. Deadline == 0 means no deadline. WaitEvent returns
// bitmask that represents flags (the least significant bit of returned value
// corresponds to flags[0]). Returned zero value informs that deadline occured.
func WaitEvent(deadline int64, flags ...*EventFlag) uint32 {
	return waitEvent(deadline, flags)
}
