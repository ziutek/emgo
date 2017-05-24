package blec

const (
	dcLLID       = 0x03 // Mask for dcL2CAPCont, dcL2CAPStart, dcLLControl.
	dcL2CAPCont  = 0x01 // Continuation fragment of L2CAP message or empty PDU.
	dcL2CAPStart = 0x02 // Start of L2CAP message or complete L2CAP message.
	dcLLControl  = 0x03 // LL Control PDU.
	dcNESN       = 0x04 // Next Expected Sequence Number.
	dcSN         = 0x08 // Sequence Number.
	dcMD         = 0x10 // More Data.
)
