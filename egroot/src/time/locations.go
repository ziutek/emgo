package time

//emgo:const
var (
	CET  = Zone{"CET", 1 * 3600}
	CEST = DST{
		Zone{"CEST", 2 * 3600},  // Last Sunday in Mar to last Sunday in Oct.
		(83*24+1)*3600 | 1<<25,  // Mar 25, 1:00, first Sun next month is Apr 1.
		(300*24+1)*3600 | 4<<25, // Oct 28, 1:00, first Sun next month is Nov 4.
	}
	EST = Zone{"EST", -5 * 3600}
	EDT = DST{
		Zone{"EDT", -4 * 3600},  // Second Sunday in Mar to first Sunday in Nov.
		(69*24+7)*3600 | 4<<25,  // Mar 11, 7:00, first Sun this month is Mar 4.
		(307*24+7)*3600 | 4<<25, // Nov 4, 7:00, first Sunday this month.
	}
)

//emgo:const
var (
	AmericaNewYork = Location{"America/New_York", &EST, &EDT}
	EuropeBerlin   = Location{"Europe/Berlin", &CET, &CEST}
	EuropeWarsaw   = Location{"Europe/Warsaw", &CET, &CEST}
)
