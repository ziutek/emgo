package noos

var tickEvent = AssignEvent()

// TickEvent returns event that is sended at every tasker interrupt.
func TickEvent() Event {
	return tickEvent
}
