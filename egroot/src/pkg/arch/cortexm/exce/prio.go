package exce

// Prio represents Cortex-M setable interrupt priority.
type Prio byte

const (
	Highest Prio = 0
	Lowest  Prio = 255
)

// Lower resturns true if priority p is lower than o.
func (p Prio) Lower(o Prio) bool {
	return p > o
}

// Higher resturns true if priority p is higher than o.
func (p Prio) Higher(o Prio) bool {
	return p < o
}

// SetPriority sets priority level for exception.
func (e Exce) SetPriority(prio Prio) {
	switch {
	case e >= MemManage && e < IRQ0:
		shp.setByte(e-MemManage, byte(prio))
	case e >= IRQ0:
		ip.setByte(e-IRQ0, byte(prio))
	}
}
