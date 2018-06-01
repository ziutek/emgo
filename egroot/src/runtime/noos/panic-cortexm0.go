// +build cortexm0

package noos

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
	}
	_ = s
	for {
	}
}
