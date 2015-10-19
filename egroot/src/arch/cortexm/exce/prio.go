package exce

// Prio represents Cortex-M configurable interrupt priority. After reset all
// configurable priorities are set to Highest (Emgo runtime or OS usually changes
// this default values).
type Prio int

const (
	PrioHighest Prio = 0
	PrioLowest  Prio = 255

	PrioRange = PrioHighest - PrioLowest
)

// Lower resturns true if priority p is lower than o.
func (p Prio) Lower(o Prio) bool {
	return p > o
}

// Higher resturns true if priority p is higher than o.
func (p Prio) Higher(o Prio) bool {
	return p < o
}
