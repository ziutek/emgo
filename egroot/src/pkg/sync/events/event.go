// Package events provides low-level communication primitives. They are intended
// for use by low-level library rutines to implement higher level communication
// and synchronization primitives like channels and mutexes.
package events

// Event represents some event that can be send or wait for.
type Event struct {
	bits uint32
}

// Assign returns event from some internal event pool. There is no any guarantee
// that subsequent calls to Assign assigns different events, which means that
// Assign can return Event already assigned by current or another gorutine.
func Assign() Event

// Send sends event that means it waking up all gorutines that wait for e. If
// some gorutine isn't waiting for any event e is saved for this gorutine for
// possible future call to Wait.
func (e Event) Send()

// Wait waits for event. If e == Event{} it returns immediately. Wait clears
// all saved events for current gorutine so the information about sended events
// that Wait hasn't waited for is lost. 
func (e Event) Wait()

// Sum returns logical sum of events. Send sum of events is equal to send all
// that events at once. Wait for sum of events means wait for at least one event
// from sum.
func Sum(el ...Event) Event {
	var sum Event
	for _, e := range el {
		sum.bits |= e.bits
	}
	return sum
}
