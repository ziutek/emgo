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
// deadline or true otherwise. Deadline == 0 means no deadline. Wait does not
// clear the flag after return. The common code looks like:
//
//	for {
//		if !flag.Wait(deadline) {
//			handleTimeout()
//			continue
//		}
//
//		// Tell source of events that should not report new ones. Instead it
//		// should save them to report later. Typicall events are reported by
//		// interrupt handlers. Most CPUs implement internal pending flag that
//		// is used to save information about masked or lower priority requests,
//		// so it is safe to mask an interrupt until event will be handled.
//		irq.Disable()
//
//		// Handle the event. Flag do not tell how many times event was reported.
//		// If you need such information use separate variable as a counter.
//		handleEvent()
//
//		// Clear the flag and enable reporting of events.
//		flag.Clear()
//		irq.Enable()
//	}
//
// If it makes sense to accept event reports concurently with handling them, the
// code might look like:
//
//	for {
//		if !flag.Wait(deadline) {
//			handleTimeout()
//			continue
//		}
//		flag.Clear()
//		handleEvent()
//	}
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
