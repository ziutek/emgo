// Peripheral: RADIO_Periph  2.4 GHz Radio
// Instances:
//	RADIO
// Tasks:
//	0x000  TXEN      Enable RADIO in TX mode.
//	0x004  RXEN      Enable RADIO in RX mode.
//  0x008  START     Start RADIO.
//	0x00C  STOP      Stop RADIO.
//	0x010  DISABLE   Disable RADIO.
//	0x014  RSSISTART Start RSSI measuremen. Take one single sample of RSSI.
//	0x018  RSSISTOP  Stop the RSSI measurement.
//	0x01C  BCSTART   Start the bit counter.
//	0x020  BCSTOP    Stop the bit counter.
// Events:
//	0x100  READY     RADIO has ramped up and is ready to be started.
//	0x104  ADDRESS   Address sent or received.
//	0x108  PAYLOAD   Packet payload sent or received.
//	0x10C  END       Packet sent or received.
//	0x110  DISABLED  RADIO has been disabled.
//	0x114  DEVMATCH  A device address match occurred on the last received pkt.
//	0x118  DEVMISS   No device address match occurred on the last received pkt.
//	0x11C  RSSIEND   Sampling of receive signal strength complete.
//	0x128  BCMATCH   Bit counter reached bit count value.
//	0x130  CRCOK     Packet received with CRC ok.
//	0x134  CRCERROR  Packet received with CRC error.
// Registers:
//  0x400 32b  CRCSTATUS   CRC status.
//  0x408 32i  RXMATCH     Received address.
//  0x40C 32   RXCRC       CRC field of previously received packet.
//	0x410 32i  DAI         Device address match index.
//	0x504 32p  PACKETPTR   Packet pointer.
//	0x508 32   FREQUENCY   Frequency.
//	0x50C 32  TXPOWER     Output power.
//	0x510 32  MODE        Data rate and modulation.
//	0x514 32  PCNF0       Packet configuration register 0.
//	0x518 32  PCNF1       Packet configuration register 1.
//	0x51C 32  BASE[2]     Base address.
//	0x524 32  PREFIX[2]   Prefixes bytes.
//	0x52C 32  TXADDRESS   Transmit address select.
//	0x530 32  RXADDRESSES Receive address select.
//	0x534 32  CRCCNF      CRC configuration.
//	0x538 32  CRCPOLY     CRC polynomial.
//	0x53C 32  CRCINIT     CRC initial value.
// Import
//  nrf5/hal/mmap
package radio

const (
	CH          FREQUENCY_Bits = 0x7f << 0 //+ Radio channel.
	MAP         FREQUENCY_Bits = 0x01 << 8 //+ Channel map selection.
	MAP_DEFAULT FREQUENCY_Bits = 0x00 << 8 //  2400 MHZ .. 2500 MHz.
	MAP_LOW     FREQUENCY_Bits = 0x01 << 8 //  2360 MHZ .. 2460 MHz.
)

const (
	CHn  = 0
	MAPn = 8
)
