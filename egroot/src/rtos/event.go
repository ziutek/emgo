package rtos

// Event represents an event that tasks or ISRs can send and tasks (but not
// ISRs) can wait for.
type Event event

// NewEvent returns new Event.
func NewEvent() *Event {
	return newEvent()
}

// Send sends event.
func (e *Event) Send() {
	eventSend(e)
}

// Wait waits for event until time t. It returns false if timeout occurs and
// true otherwise. Use t < 0 to disable timeout checking. Task can not determine
// how many events was sent before it returned from Wait.
func (e *Event) Wait(t int64) bool {
	return eventWaitUntil(e, t)
}

// WaitEvent waits for at least one from no more than 32 different events until
// time t. It returns bitmask that describes which events occured: the least
// significant bit of returned velue corresponds to events[0]. Use t < 0 to
// disable timeout checking. Task can not determine how many events was sent
// before it returned from Wait.
func WaitEvent(t int64, events ...*Event) uint32 {
	return waitEvent(t, events)
}
