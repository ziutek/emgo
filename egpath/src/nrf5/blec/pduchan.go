package blec

import (
	"bluetooth/ble"
)

// PDUChan consist of a channel and a pool of preallocated DataPDUs. It allows
// to send PDUs obtained from the internal poll through channel Ch.
//
// Receiver uses only Ch to receive next PDU (it should not use other methods of
// PDUChan). Receiver owns received PDU as long as it receives next PDU. After
// that it should no use previously received PDU (receiving from Ch returns
// previously received PDU to the internal pool).
//
// Sender obtains free PDU from internal poll using Get method. It owns this PDU
// until it sends it to Ch. After sending, it must call Next to select next free
// PDU, that will be returned by Get method.
type pduChan struct {
	Ch  chan ble.DataPDU
	sli []ble.DataPDU
}

func makePDUChan(maxpay, n int) pduChan {
	ch := make(chan ble.DataPDU, n)
	// Preallocate n+2 PDUs: 1 for receiver, n for channel, 1 for sender.
	sli := make([]ble.DataPDU, n+2)
	for i := range sli {
		sli[i] = ble.MakeDataPDU(maxpay)
	}
	return pduChan{ch, sli[:1]}
}

// Get returns free PDU from internal pool. It returns the same PDU as long as
// Next method selects new one.
func (c *pduChan) Get() ble.DataPDU {
	return c.sli[len(c.sli)-1]
}

// Next selects next free PDU to be returned by Get method.
func (c *pduChan) Next() {
	n := len(c.sli)
	if n < cap(c.sli) {
		c.sli = c.sli[:n+1]
	} else {
		c.sli = c.sli[:1]
	}
}
