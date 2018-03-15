package tz

import "time"

//emgo:const
var (
	CET  = time.Zone{"CET", 1 * 3600}
	CEST = time.DST{
		// Last Sunday in March 2:00 to last Sunday in October 3:00.
		time.Zone{"CEST", 2 * 3600},
		// Start: March 25 1:00 UTC, first Sunday next month is April 1.
		(83*24+1)*3600 | 1<<25,
		// End: October 28 1:00 UTC, first Sunday next month is November 4.
		(300*24+1)*3600 | 4<<25,
	}
)

//emgo:const
var (
	EuropeBerlin = time.Location{"Europe/Berlin", &CET, &CEST}
	EuropeWarsaw = time.Location{"Europe/Warsaw", &CET, &CEST}
)
