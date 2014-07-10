package noos

// Ticks return current value of system counter that is incremented with some
// period of the order of one to several milliseconds.
func Ticks() uint64 {
	return loadTicks()
}

var tickPeriod int

// TickPeriod returns tick period in milliseconds.
func TickPeriod() int {
	return tickPeriod
}

// SetTickPeriod can be used to set value of tick period.
func SetTickPeriod(ms int) {
	tickPeriod = ms
}

var tickEvent = AssignEvent()

// TickEvent returns event that is send at every tick.
func TickEvent() Event {
	return tickEvent
}
