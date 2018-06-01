// +build cortexm3 cortexm4 cortexm4f cortexm7f cortexm7d

package noos

import (
	"arch/cortexm/debug/itm"
)

type stringer interface {
	String() string
}

func panic_(i interface{}) {
	var s string
	switch v := i.(type) {
	case string:
		s = v
	case error:
		s = v.Error()
	case stringer:
		s = v.String()
	default:
		s = "<no text descr>"
	}
	dbg := itm.Port(0)
	dbg.WriteString("\npanic: ")
	dbg.WriteString(s)
	dbg.WriteByte('\n')
	for {
	}
}
