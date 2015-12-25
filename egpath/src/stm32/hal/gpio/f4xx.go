// +build f40_41xxx f411xe

package gpio

const (
	Low      Speed = 0 //   2 MHz (CL = 50 pF, VDD > 2.7 V)
	Medium   Speed = 1 //  25 MHz (CL = 50 pF, VDD > 2.7 V)
	High     Speed = 2 //  50 MHz (CL = 40 pF, VDD > 2.7 V)
	VeryHigh Speed = 3 // 100 MHz (CL = 30 pF, VDD > 2.7 V)
)
