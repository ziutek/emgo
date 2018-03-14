package tz

import "time"

//emgo:const
var (
	EST = time.Zone{"EST", -5 * 3600}
	EDT = time.DST{
		// Second Sunday in Mar to first Sunday in Nov.
		time.Zone{"EDT", -4 * 3600},
		// Start: March 11, 7:00 UTC, first Sunday this month is March 4.
		(69*24+7)*3600 | 4<<25,
		// End: November 4, 7:00 UTC, first Sunday this month is November 4.
		(307*24+7)*3600 | 4<<25,
	}
)

//emgo:const
var (
	AmericaNewYork = time.Location{"America/New_York", &EST, &EDT}
)
