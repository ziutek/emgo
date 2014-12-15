package sync

import (
	"sync/atomic"
)

type Once struct {
	m    Mutex
	done uint32
}

func (o *Once) Do(f func()) {
	if atomic.LoadUint32(&o.done) != 0 {
		return
	}
	o.m.Lock()
	if o.done == 0 {
		f()
		o.done = 1
	}
	o.m.Unlock()
}
