package time

type Duration int64

type Time struct {
	sec  int64
	nsec int32
	//loc *Location
}

func Now() Time {
	return now()
}

func Set(t Time) {
	set(t)
}
