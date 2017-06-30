package ble

type DataRing struct {
	sli []DataPDU
}

func MakeDataRing(maxpay, n int) DataRing {
	sli := make([]DataPDU, n)
	for i := range sli {
		sli[i] = MakeDataPDU(maxpay)
	}
	return DataRing{sli[:1]} // Fail now if n < 1.
}

func (r DataRing) Get() DataPDU {
	return r.sli[len(r.sli)-1]
}

func (r *DataRing) Next() {
	n := len(r.sli)
	if n < cap(r.sli) {
		r.sli = r.sli[:n+1]
	} else {
		r.sli = r.sli[:1]
	}
}
