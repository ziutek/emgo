// BaseAddr: 0xe000ed88
//  0: CPACR Coprocessor Access Control Register
package scb

const (
	CP10 CPACR_Bits = 3 << 20
	CP11 CPACR_Bits = 3 << 22

	AccessDeny int = 0
	AccessPriv int = 1
	AccessFull int = 3
)
