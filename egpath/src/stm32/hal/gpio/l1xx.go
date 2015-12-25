// +build l1xx_md l1xx_mdp l1xx_hd l1xx_xl

package gpio

const (
	VeryLow Speed = 0 // 400 kHz (CL = 50 pF, VDD > 2.7 V)
	Low     Speed = 1 //   2 MHz (CL = 50 pF, VDD > 2.7 V)
	Medium  Speed = 2 //  10 MHz (CL = 50 pF, VDD > 2.7 V)
	High    Speed = 3 //  50 MHz (CL = 50 pF, VDD > 2.7 V)
)
