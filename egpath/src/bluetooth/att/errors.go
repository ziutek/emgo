package att

type ErrorCode byte

const (
	InvalidHandle                 ErrorCode = 0x01 // The attribute handle given was not valid on this server.
	ReadNotPermitted              ErrorCode = 0x02 // The attribute cannot be read.
	WriteNotPermitted             ErrorCode = 0x03 // The attribute cannot be written.
	InvalidPDU                    ErrorCode = 0x04 // The attribute PDU was invalid.
	InsufficientAuthentication    ErrorCode = 0x05 // The attribute requires authentication before it can be read or written.
	RequestNotSupported           ErrorCode = 0x06 // Attribute server does not support the request received from the client.
	InvalidOffset                 ErrorCode = 0x07 // Offset specified was past the end of the attribute.
	InsufficientAuthorization     ErrorCode = 0x08 // The attribute requires authorization before it can be read or written.
	PrepareQueueFull              ErrorCode = 0x09 // Too many prepare writes have been queued.
	AttributeNotFound             ErrorCode = 0x0A // No attribute found within the given attri-bute handle range.
	AttributeNotLong              ErrorCode = 0x0B // The attribute cannot be read or written using the Read Blob Request
	InsufficientEncryptionKeySize ErrorCode = 0x0C // The Encryption Key Size used for encrypting this link is insufficient.
	InvalidAttributeValueLength   ErrorCode = 0x0D // The attribute value length is invalid for the operation.
	UnlikelyError                 ErrorCode = 0x0E // The attribute request that was requested has encountered an error that was unlikely, and therefore could not be completed as requested.
	InsufficientEncryption        ErrorCode = 0x0F // The attribute requires encryption before it can be read or written.
	UnsupportedGroupType          ErrorCode = 0x10 // The attribute type is not a supported grouping attribute as defined by a higher layer specification.
	InsufficientResources         ErrorCode = 0x11 // Insufficient Resources to complete the request.
)
