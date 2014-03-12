package irq

// Prio represents Cortex-M setable interrupt priority.
type Prio byte

const (
	Highest Prio = 0
	Lowest  Prio = 255
)

// Same returns true if priority p is equal to o.
func (p Prio) Same(o Prio) bool {
	return p == 0
}

// Lower resturns true if priority p is lower than o.
func (p Prio) Lower(o Prio) bool {
	return p > o
}

// Higher resturns true if priority p is higher than o.
func (p Prio) Higher(o Prio) bool {
	return p < o
}

// SetPriority sets priority level for irq
func (irq IRQ) SetPriority(prio Prio) {
	switch {
	case irq >= MemFault && irq < Ext0:
		shp.setByte(irq-MemFault, byte(prio))
	case irq >= Ext0:
		ip.setByte(irq-Ext0, byte(prio))
	}
}
