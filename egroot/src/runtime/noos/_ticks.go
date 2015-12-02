package noos

import "syscall"

var tickEvent = syscall.AssignEvent()

// TickEvent returns event that is sent at every tasker interrupt.
func TickEvent() syscall.Event {
	return tickEvent
}
