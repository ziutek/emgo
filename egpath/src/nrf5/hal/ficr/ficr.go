// Peripheral: Periph  Factory information configuration register
// Instances:
//  FICR  mmap.FICR_BASE
// Registers:
//  0x010 32  CODEPAGESIZE     Code memory page size.
//  0x014 32  CODESIZE         Code memory size.
//  0x028 32  CLENR0           Length of code region 0 in bytes.
//  0x02C 32  PPFC             Pre-programmed factory code present.
//  0x034 32  NUMRAMBLOCK      Number of individually controllable RAM blocks.
//  0x038 32  SIZERAMBLOCK[4]  Size of RAM block n in bytes.
//  0x05C 32  CONFIGID         Configuration identifier.
//  0x060 32  DEVICEID[2]      Device identifier.
//  0x080 32  ER[4]            Encryption Root.
//  0x090 32  IR[4]            Identity Root.
//  0x0A0 32  DEVICEADDRTYPE   Device address type.
//  0x0A4 32  DEVICEADDR[2]    Device address.
//  0x0AC 32  OVERRIDDEN       Override enable.
//  0x0B0 32  NRF_1MBIT[5]     RADIO.OVERRIDE[n] values for NRF_1MBIT mode.
//  0x0EC 32  BLE_1MBIT[5]     RADIO.OVERRIDE[n] values for BLE_1MBIT mode.
//  0x100 32  INFO_PART        Part code.
//  0x104 32  INFO_VARIANT     Part variant, hardware vers. and production conf.
//  0x108 32  INFO_PACKAGE     Package option.
//  0x10C 32  INFO_RAM         RAM variant.
//  0x110 32  INFO_FLASH       Flash variant.
//  0x404 32  TEMP_A[6]        Slope definition An.
//  0x41C 32  TEMP_B[6]        y-intercept Bn.
//  0x434 32  TEMP_T[5]        Segment end Tn.
//  0x450 32  NFC_TAGHEADER[4] Default header for NFC Tag.
// Import:
//  nrf5/hal/internal/mmap
package ficr

const (
	NRF_1MBIT_OK OVERRIDDEN_Bits = 1 << 0 //+ Use default values for NRF_1MBIT.
	BLE_1MBIT_OK OVERRIDDEN_Bits = 1 << 3 //+ Use default values for BLE_1MBIT.
)
