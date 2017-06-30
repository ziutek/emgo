package blec

import (
	"bluetooth/ble"
)

type pduChan struct {
	Ch  chan ble.DataPDU
	sli []ble.DataPDU
}

func makePDUChan(maxpay, n int) pduChan {
	ch := make(chan ble.DataPDU, n)
	sli := make([]ble.DataPDU, n+2)
	for i := range sli {
		sli[i] = ble.MakeDataPDU(maxpay)
	}
	return pduChan{ch, sli[:1]}
}

func (c *pduChan) Get() ble.DataPDU {
	return c.sli[len(c.sli)-1]
}

func (c *pduChan) Next() {
	n := len(c.sli)
	if n < cap(c.sli) {
		c.sli = c.sli[:n+1]
	} else {
		c.sli = c.sli[:1]
	}
}
