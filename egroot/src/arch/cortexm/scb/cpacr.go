// BaseAddr: 0xe000ed88
//  0: CPACR Coprocessor Access Control Register
package scb

const (
	CP10 CPACR_Field = 3<<siz + 20
	CP11 CPACR_Field = 3<<siz + 22

	AccessDeny = 0
	AccessPriv = 1
	AccessFull = 3
)
