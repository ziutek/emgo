package tz

import "time"

//emgo:const
var (
	AEST = time.Zone{"AEST", 10 * 3600}
	AEDT = time.DST{
		// First Sunday in October 2:00 to first Sunday in April 3:00.
		time.Zone{"AEDT", 11 * 3600},
		// Start: October 6 16:00 UTC, first Sunday this month is October 7.
		(278*24+16)*3600 | 7<<25,
		// End: March 31 16:00 UTC, first Sunday in April is April 1.
		(89*24+16)*3600 | 1<<25,
	}
)

//emgo:const
var (
	AustraliaSydney = time.Location{"Australia/Sydney", &AEST, &AEDT}
)
