package ble

// DataPDU represents Data Channel PDU. See PDU for more information.
type DataPDU struct {
	PDU
}

// MakeDataPDU is more readable counterpart of DataPDU{MakePDU(maxpay)}.
func MakeDataPDU(maxpay int) DataPDU {
	return DataPDU{MakePDU(maxpay)}
}

// Header represents first byte of Data Channel PDU header. Second byte of
// header is payload length and can be obtained using PayLen method.
type Header byte

const (
	LLID       Header = 0x03 // Mask for L2CAPCont, L2CAPStart, LLControl.
	L2CAPCont  Header = 0x01 // Continuation of L2CAP message or empty PDU.
	L2CAPStart Header = 0x02 // Start of L2CAP message or complete message.
	LLControl  Header = 0x03 // LL Control PDU.

	NESN Header = 0x04 // Next Expected Sequence Number.
	SN   Header = 0x08 // Sequence Number.
	MD   Header = 0x10 // More Data.
)

// Header returns first byte of header of pdu.
func (pdu DataPDU) Header() Header {
	return Header(pdu.b[0])
}

// SetHeader sets first byte of header of pdu.
func (pdu DataPDU) SetHeader(h Header) {
	 pdu.b[0] = byte(h)
}