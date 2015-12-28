// +build f411xe

// Peripheral: ETH_Periph  Ethernet MAC.
// Instances:
//  ETH  mmap.ETH_BASE
// Registers:
//  0x00   32  MACCR
//  0x04   32  MACFFR
//  0x08   32  MACHTHR
//  0x0C   32  MACHTLR
//  0x10   32  MACMIIAR
//  0x14   32  MACMIIDR
//  0x18   32  MACFCR
//  0x1C   32  MACVLANTR   8.
//  0x28   32  MACRWUFFR   11.
//  0x2C   32  MACPMTCSR
//  0x38   32  MACSR       15.
//  0x3C   32  MACIMR
//  0x40   32  MACA0HR
//  0x44   32  MACA0LR
//  0x48   32  MACA1HR
//  0x4C   32  MACA1LR
//  0x50   32  MACA2HR
//  0x54   32  MACA2LR
//  0x58   32  MACA3HR
//  0x5C   32  MACA3LR     24.
//  0x100  32  MMCCR       65.
//  0x104  32  MMCRIR
//  0x108  32  MMCTIR
//  0x10C  32  MMCRIMR
//  0x110  32  MMCTIMR     69.
//  0x14C  32  MMCTGFSCCR  84.
//  0x150  32  MMCTGFMSCCR
//  0x168  32  MMCTGFCR
//  0x194  32  MMCRFCECR
//  0x198  32  MMCRFAECR
//  0x1C4  32  MMCRGUFCR
//  0x700  32  PTPTSCR
//  0x704  32  PTPSSIR
//  0x708  32  PTPTSHR
//  0x70C  32  PTPTSLR
//  0x710  32  PTPTSHUR
//  0x714  32  PTPTSLUR
//  0x718  32  PTPTSAR
//  0x71C  32  PTPTTHR
//  0x720  32  PTPTTLR
//  0x724  32  RESERVED8
//  0x728  32  PTPTSSR
//  0x1000 32  DMABMR
//  0x1004 32  DMATPDR
//  0x1008 32  DMARPDR
//  0x100C 32  DMARDLAR
//  0x1010 32  DMATDLAR
//  0x1014 32  DMASR
//  0x1018 32  DMAOMR
//  0x101C 32  DMAIER
//  0x1020 32  DMAMFBOCR
//  0x1024 32  DMARSWTR
//  0x1048 32  DMACHTDR
//  0x104C 32  DMACHRDR
//  0x1050 32  DMACHTBAR
//  0x1054 32  DMACHRBAR
// Import:
//  stm32/o/f411xe/mmap
package eth

const (
	WD        MACCR_Bits = 0x01 << 23 //+ Watchdog disable.
	JD        MACCR_Bits = 0x01 << 22 //+ Jabber disable.
	IFG       MACCR_Bits = 0x07 << 17 //+ Inter-frame gap.
	IFG_96Bit MACCR_Bits = 0x00 << 17 //  Minimum IFG between frames during transmission is 96Bit.
	IFG_88Bit MACCR_Bits = 0x01 << 17 //  Minimum IFG between frames during transmission is 88Bit.
	IFG_80Bit MACCR_Bits = 0x02 << 17 //  Minimum IFG between frames during transmission is 80Bit.
	IFG_72Bit MACCR_Bits = 0x03 << 17 //  Minimum IFG between frames during transmission is 72Bit.
	IFG_64Bit MACCR_Bits = 0x04 << 17 //  Minimum IFG between frames during transmission is 64Bit.
	IFG_56Bit MACCR_Bits = 0x05 << 17 //  Minimum IFG between frames during transmission is 56Bit.
	IFG_48Bit MACCR_Bits = 0x06 << 17 //  Minimum IFG between frames during transmission is 48Bit.
	IFG_40Bit MACCR_Bits = 0x07 << 17 //  Minimum IFG between frames during transmission is 40Bit.
	CSD       MACCR_Bits = 0x01 << 16 //+ Carrier sense disable (during transmission).
	FES       MACCR_Bits = 0x01 << 14 //+ Fast ethernet speed.
	ROD       MACCR_Bits = 0x01 << 13 //+ Receive own disable.
	LM        MACCR_Bits = 0x01 << 12 //+ loopback mode.
	DM        MACCR_Bits = 0x01 << 11 //+ Duplex mode.
	IPCO      MACCR_Bits = 0x01 << 10 //+ IP Checksum offload.
	RD        MACCR_Bits = 0x01 << 9  //+ Retry disable.
	APCS      MACCR_Bits = 0x01 << 7  //+ Automatic Pad/CRC stripping.
	BL        MACCR_Bits = 0x03 << 5  //+ Back-off limit: random integer number (r) of slot time delays before rescheduling.
	BL_10     MACCR_Bits = 0x00 << 5  //  k = min (n, 10).
	BL_8      MACCR_Bits = 0x01 << 5  //  k = min (n, 8).
	BL_4      MACCR_Bits = 0x02 << 5  //  k = min (n, 4).
	BL_1      MACCR_Bits = 0x03 << 5  //  k = min (n, 1).
	DC        MACCR_Bits = 0x01 << 4  //+ Defferal check.
	TE        MACCR_Bits = 0x01 << 3  //+ Transmitter enable.
	RE        MACCR_Bits = 0x01 << 2  //+ Receiver enable.
)

const (
	RA                          MACFFR_Bits = 0x01 << 31 //+ Receive all.
	HPF                         MACFFR_Bits = 0x01 << 10 //+ Hash or perfect filter.
	SAF                         MACFFR_Bits = 0x01 << 9  //+ Source address filter enable.
	SAIF                        MACFFR_Bits = 0x01 << 8  //+ SA inverse filtering.
	PCF                         MACFFR_Bits = 0x03 << 6  //+ Pass control frames: 3 cases.
	PCF_BlockAll                MACFFR_Bits = 0x01 << 6  //  MAC filters all control frames from reaching the application.
	PCF_ForwardAll              MACFFR_Bits = 0x02 << 6  //  MAC forwards all control frames to application even if they fail the Address Filter.
	PCF_ForwardPassedAddrFilter MACFFR_Bits = 0x03 << 6  //  MAC forwards control frames that pass the Address Filter..
	BFD                         MACFFR_Bits = 0x01 << 5  //+ Broadcast frame disable.
	PAM                         MACFFR_Bits = 0x01 << 4  //+ Pass all mutlicast.
	DAIF                        MACFFR_Bits = 0x01 << 3  //+ DA Inverse filtering.
	HM                          MACFFR_Bits = 0x01 << 2  //+ Hash multicast.
	HU                          MACFFR_Bits = 0x01 << 1  //+ Hash unicast.
	PM                          MACFFR_Bits = 0x01 << 0  //+ Promiscuous mode.
)

const (
	HTH MACHTHR_Bits = 0xFFFFFFFF << 0 //+ Hash table high.
)

const (
	HTL MACHTLR_Bits = 0xFFFFFFFF << 0 //+ Hash table low.
)

const (
	PA        MACMIIAR_Bits = 0x1F << 11 //+ Physical layer address.
	MR        MACMIIAR_Bits = 0x1F << 6  //+ MII register in the selected PHY.
	CR        MACMIIAR_Bits = 0x07 << 2  //+ CR clock range: 6 cases.
	CR_Div42  MACMIIAR_Bits = 0x00 << 2  //  HCLK:60-100 MHz; MDC clock= HCLK/42.
	CR_Div62  MACMIIAR_Bits = 0x01 << 2  //  HCLK:100-150 MHz; MDC clock= HCLK/62.
	CR_Div16  MACMIIAR_Bits = 0x02 << 2  //  HCLK:20-35 MHz; MDC clock= HCLK/16.
	CR_Div26  MACMIIAR_Bits = 0x03 << 2  //  HCLK:35-60 MHz; MDC clock= HCLK/26.
	CR_Div102 MACMIIAR_Bits = 0x04 << 2  //  HCLK:150-168 MHz; MDC clock= HCLK/102.
	MW        MACMIIAR_Bits = 0x01 << 1  //+ MII write.
	MB        MACMIIAR_Bits = 0x01 << 0  //+ MII busy.
)

const (
	MD MACMIIDR_Bits = 0xFFFF << 0 //+ MII data: read/write data from/to PHY.
)

const (
	PT           MACFCR_Bits = 0xFFFF << 16 //+ Pause time.
	ZQPD         MACFCR_Bits = 0x01 << 7    //+ Zero-quanta pause disable.
	PLT          MACFCR_Bits = 0x03 << 4    //+ Pause low threshold: 4 cases.
	PLT_Minus4   MACFCR_Bits = 0x00 << 4    //  Pause time minus 4 slot times.
	PLT_Minus28  MACFCR_Bits = 0x01 << 4    //  Pause time minus 28 slot times.
	PLT_Minus144 MACFCR_Bits = 0x02 << 4    //  Pause time minus 144 slot times.
	PLT_Minus256 MACFCR_Bits = 0x03 << 4    //  Pause time minus 256 slot times.
	UPFD         MACFCR_Bits = 0x01 << 3    //+ Unicast pause frame detect.
	RFCE         MACFCR_Bits = 0x01 << 2    //+ Receive flow control enable.
	TFCE         MACFCR_Bits = 0x01 << 1    //+ Transmit flow control enable.
	FCBBPA       MACFCR_Bits = 0x01 << 0    //+ Flow control busy/backpressure activate.
)

const (
	VLANTC MACVLANTR_Bits = 0x01 << 16  //+ 12-bit VLAN tag comparison.
	VLANTI MACVLANTR_Bits = 0xFFFF << 0 //+ VLAN tag identifier (for receive frames).
)

const (
	D MACRWUFFR_Bits = 0xFFFFFFFF << 0 //+ Wake-up frame filter register data.
)

const (
	WFFRPR MACPMTCSR_Bits = 0x01 << 31 //+ Wake-Up Frame Filter Register Pointer Reset.
	GU     MACPMTCSR_Bits = 0x01 << 9  //+ Global Unicast.
	WFR    MACPMTCSR_Bits = 0x01 << 6  //+ Wake-Up Frame Received.
	MPR    MACPMTCSR_Bits = 0x01 << 5  //+ Magic Packet Received.
	WFE    MACPMTCSR_Bits = 0x01 << 2  //+ Wake-Up Frame Enable.
	MPE    MACPMTCSR_Bits = 0x01 << 1  //+ Magic Packet Enable.
	PD     MACPMTCSR_Bits = 0x01 << 0  //+ Power Down.
)

const (
	TSTS   MACSR_Bits = 0x01 << 9 //+ Time stamp trigger status.
	MMCTS  MACSR_Bits = 0x01 << 6 //+ MMC transmit status.
	MMMCRS MACSR_Bits = 0x01 << 5 //+ MMC receive status.
	MMCS   MACSR_Bits = 0x01 << 4 //+ MMC status.
	PMTS   MACSR_Bits = 0x01 << 3 //+ PMT status.
)

const (
	TSTIM MACIMR_Bits = 0x01 << 9 //+ Time stamp trigger interrupt mask.
	PMTIM MACIMR_Bits = 0x01 << 3 //+ PMT interrupt mask.
)

const (
	MACA0H MACA0HR_Bits = 0xFFFF << 0 //+ MAC address0 high.
)

const (
	MACA0L MACA0LR_Bits = 0xFFFFFFFF << 0 //+ MAC address0 low.
)

const (
	AE             MACA1HR_Bits = 0x01 << 31  //+ Address enable.
	SA             MACA1HR_Bits = 0x01 << 30  //+ Source address.
	MBC            MACA1HR_Bits = 0x3F << 24  //+ Mask byte control: bits to mask for comparison of the MAC Address bytes.
	MBC_HBits15_8  MACA1HR_Bits = 0x20 << 24  //  Mask MAC Address high reg bits [15:8].
	MBC_HBits7_0   MACA1HR_Bits = 0x10 << 24  //  Mask MAC Address high reg bits [7:0].
	MBC_LBits31_24 MACA1HR_Bits = 0x08 << 24  //  Mask MAC Address low reg bits [31:24].
	MBC_LBits23_16 MACA1HR_Bits = 0x04 << 24  //  Mask MAC Address low reg bits [23:16].
	MBC_LBits15_8  MACA1HR_Bits = 0x02 << 24  //  Mask MAC Address low reg bits [15:8].
	MBC_LBits7_0   MACA1HR_Bits = 0x01 << 24  //  Mask MAC Address low reg bits [7:0].
	MACA1H         MACA1HR_Bits = 0xFFFF << 0 //+ MAC address1 high.
)

const (
	MACA1L MACA1LR_Bits = 0xFFFFFFFF << 0 //+ MAC address1 low.
)

const (
	AE             MACA2HR_Bits = 0x01 << 31  //+ Address enable.
	SA             MACA2HR_Bits = 0x01 << 30  //+ Source address.
	MBC            MACA2HR_Bits = 0x3F << 24  //+ Mask byte control.
	MBC_HBits15_8  MACA2HR_Bits = 0x20 << 24  //  Mask MAC Address high reg bits [15:8].
	MBC_HBits7_0   MACA2HR_Bits = 0x10 << 24  //  Mask MAC Address high reg bits [7:0].
	MBC_LBits31_24 MACA2HR_Bits = 0x08 << 24  //  Mask MAC Address low reg bits [31:24].
	MBC_LBits23_16 MACA2HR_Bits = 0x04 << 24  //  Mask MAC Address low reg bits [23:16].
	MBC_LBits15_8  MACA2HR_Bits = 0x02 << 24  //  Mask MAC Address low reg bits [15:8].
	MBC_LBits7_0   MACA2HR_Bits = 0x01 << 24  //  Mask MAC Address low reg bits [70].
	MACA2H         MACA2HR_Bits = 0xFFFF << 0 //+ MAC address1 high.
)

const (
	MACA2L MACA2LR_Bits = 0xFFFFFFFF << 0 //+ MAC address2 low.
)

const (
	AE             MACA3HR_Bits = 0x01 << 31  //+ Address enable.
	SA             MACA3HR_Bits = 0x01 << 30  //+ Source address.
	MBC            MACA3HR_Bits = 0x3F << 24  //+ Mask byte control.
	MBC_HBits15_8  MACA3HR_Bits = 0x20 << 24  //  Mask MAC Address high reg bits [15:8].
	MBC_HBits7_0   MACA3HR_Bits = 0x10 << 24  //  Mask MAC Address high reg bits [7:0].
	MBC_LBits31_24 MACA3HR_Bits = 0x08 << 24  //  Mask MAC Address low reg bits [31:24].
	MBC_LBits23_16 MACA3HR_Bits = 0x04 << 24  //  Mask MAC Address low reg bits [23:16].
	MBC_LBits15_8  MACA3HR_Bits = 0x02 << 24  //  Mask MAC Address low reg bits [15:8].
	MBC_LBits7_0   MACA3HR_Bits = 0x01 << 24  //  Mask MAC Address low reg bits [70].
	MACA3H         MACA3HR_Bits = 0xFFFF << 0 //+ MAC address3 high.
)

const (
	MACA3L MACA3LR_Bits = 0xFFFFFFFF << 0 //+ MAC address3 low.
)

const (
	MCFHP MMCCR_Bits = 0x01 << 5 //+ MMC counter Full-Half preset.
	MCP   MMCCR_Bits = 0x01 << 4 //+ MMC counter preset.
	MCF   MMCCR_Bits = 0x01 << 3 //+ MMC Counter Freeze.
	ROR   MMCCR_Bits = 0x01 << 2 //+ Reset on Read.
	CSR   MMCCR_Bits = 0x01 << 1 //+ Counter Stop Rollover.
	CR    MMCCR_Bits = 0x01 << 0 //+ Counters Reset.
)

const (
	RGUFS MMCRIR_Bits = 0x01 << 17 //+ Set when Rx good unicast frames counter reaches half the maximum value.
	RFAES MMCRIR_Bits = 0x01 << 6  //+ Set when Rx alignment error counter reaches half the maximum value.
	RFCES MMCRIR_Bits = 0x01 << 5  //+ Set when Rx crc error counter reaches half the maximum value.
)

const (
	TGFS    MMCTIR_Bits = 0x01 << 21 //+ Set when Tx good frame count counter reaches half the maximum value.
	TGFMSCS MMCTIR_Bits = 0x01 << 15 //+ Set when Tx good multi col counter reaches half the maximum value.
	TGFSCS  MMCTIR_Bits = 0x01 << 14 //+ Set when Tx good single col counter reaches half the maximum value.
)

const (
	RGUFM MMCRIMR_Bits = 0x01 << 17 //+ Mask the interrupt when Rx good unicast frames counter reaches half the maximum value.
	RFAEM MMCRIMR_Bits = 0x01 << 6  //+ Mask the interrupt when when Rx alignment error counter reaches half the maximum value.
	RFCEM MMCRIMR_Bits = 0x01 << 5  //+ Mask the interrupt when Rx crc error counter reaches half the maximum value.
)

const (
	TGFM    MMCTIMR_Bits = 0x01 << 21 //+ Mask the interrupt when Tx good frame count counter reaches half the maximum value.
	TGFMSCM MMCTIMR_Bits = 0x01 << 15 //+ Mask the interrupt when Tx good multi col counter reaches half the maximum value.
	TGFSCM  MMCTIMR_Bits = 0x01 << 14 //+ Mask the interrupt when Tx good single col counter reaches half the maximum value.
)

const (
	TGFSCC MMCTGFSCCR_Bits = 0xFFFFFFFF << 0 //+ Number of successfully transmitted frames after a single collision in Half-duplex mode..
)

const (
	TGFMSCC MMCTGFMSCCR_Bits = 0xFFFFFFFF << 0 //+ Number of successfully transmitted frames after more than a single collision in Half-duplex mode..
)

const (
	TGFC MMCTGFCR_Bits = 0xFFFFFFFF << 0 //+ Number of good frames transmitted..
)

const (
	RFCEC MMCRFCECR_Bits = 0xFFFFFFFF << 0 //+ Number of frames received with CRC error..
)

const (
	RFAEC MMCRFAECR_Bits = 0xFFFFFFFF << 0 //+ Number of frames received with alignment (dribble) error.
)

const (
	RGUFC MMCRGUFCR_Bits = 0xFFFFFFFF << 0 //+ Number of good unicast frames received..
)

const (
	TSCNT PTPTSCR_Bits = 0x03 << 16 //+ Time stamp clock node type.
	TSARU PTPTSCR_Bits = 0x01 << 5  //+ Addend register update.
	TSITE PTPTSCR_Bits = 0x01 << 4  //+ Time stamp interrupt trigger enable.
	TSSTU PTPTSCR_Bits = 0x01 << 3  //+ Time stamp update.
	TSSTI PTPTSCR_Bits = 0x01 << 2  //+ Time stamp initialize.
	TSFCU PTPTSCR_Bits = 0x01 << 1  //+ Time stamp fine or coarse update.
	TSE   PTPTSCR_Bits = 0x01 << 0  //+ Time stamp enable.
)

const (
	STSSI PTPSSIR_Bits = 0xFF << 0 //+ System time Sub-second increment value.
)

const (
	STS PTPTSHR_Bits = 0xFFFFFFFF << 0 //+ System Time second.
)

const (
	STPNS PTPTSLR_Bits = 0x01 << 31      //+ System Time Positive or negative time.
	STSS  PTPTSLR_Bits = 0x7FFFFFFF << 0 //+ System Time sub-seconds.
)

const (
	TSUS PTPTSHUR_Bits = 0xFFFFFFFF << 0 //+ Time stamp update seconds.
)

const (
	TSUPNS PTPTSLUR_Bits = 0x01 << 31      //+ Time stamp update Positive or negative time.
	TSUSS  PTPTSLUR_Bits = 0x7FFFFFFF << 0 //+ Time stamp update sub-seconds.
)

const (
	TSA PTPTSAR_Bits = 0xFFFFFFFF << 0 //+ Time stamp addend.
)

const (
	TTSH PTPTTHR_Bits = 0xFFFFFFFF << 0 //+ Target time stamp high.
)

const (
	TTSL PTPTTLR_Bits = 0xFFFFFFFF << 0 //+ Target time stamp low.
)

const (
	TSSMRME    PTPTSSR_Bits = 0x01 << 15 //+ Time stamp snapshot for message relevant to master enable.
	TSSEME     PTPTSSR_Bits = 0x01 << 14 //+ Time stamp snapshot for event message enable.
	TSSIPV4FE  PTPTSSR_Bits = 0x01 << 13 //+ Time stamp snapshot for IPv4 frames enable.
	TSSIPV6FE  PTPTSSR_Bits = 0x01 << 12 //+ Time stamp snapshot for IPv6 frames enable.
	TSSPTPOEFE PTPTSSR_Bits = 0x01 << 11 //+ Time stamp snapshot for PTP over ethernet frames enable.
	TSPTPPSV2E PTPTSSR_Bits = 0x01 << 10 //+ Time stamp PTP packet snooping for version2 format enable.
	TSSSR      PTPTSSR_Bits = 0x01 << 9  //+ Time stamp Sub-seconds rollover.
	TSSARFE    PTPTSSR_Bits = 0x01 << 8  //+ Time stamp snapshot for all received frames enable.
	TSTTR      PTPTSSR_Bits = 0x01 << 5  //+ Time stamp target time reached.
	TSSO       PTPTSSR_Bits = 0x01 << 4  //+ Time stamp seconds overflow.
)

const (
	AAB               DMABMR_Bits = 0x01 << 25   //+ Address-Aligned beats.
	FPM               DMABMR_Bits = 0x01 << 24   //+ 4xPBL mode.
	USP               DMABMR_Bits = 0x01 << 23   //+ Use separate PBL.
	RDP               DMABMR_Bits = 0x3F << 17   //+ RxDMA PBL.
	RDP_1Beat         DMABMR_Bits = 0x01 << 17   //  maximum number of beats to be transferred in one RxDMA transaction is 1.
	RDP_2Beat         DMABMR_Bits = 0x02 << 17   //  maximum number of beats to be transferred in one RxDMA transaction is 2.
	RDP_4Beat         DMABMR_Bits = 0x04 << 17   //  maximum number of beats to be transferred in one RxDMA transaction is 4.
	RDP_8Beat         DMABMR_Bits = 0x08 << 17   //  maximum number of beats to be transferred in one RxDMA transaction is 8.
	RDP_16Beat        DMABMR_Bits = 0x10 << 17   //  maximum number of beats to be transferred in one RxDMA transaction is 16.
	RDP_32Beat        DMABMR_Bits = 0x20 << 17   //  maximum number of beats to be transferred in one RxDMA transaction is 32.
	RDP_4xPBL_4Beat   DMABMR_Bits = 0x81 << 17   //  maximum number of beats to be transferred in one RxDMA transaction is 4.
	RDP_4xPBL_8Beat   DMABMR_Bits = 0x82 << 17   //  maximum number of beats to be transferred in one RxDMA transaction is 8.
	RDP_4xPBL_16Beat  DMABMR_Bits = 0x84 << 17   //  maximum number of beats to be transferred in one RxDMA transaction is 16.
	RDP_4xPBL_32Beat  DMABMR_Bits = 0x88 << 17   //  maximum number of beats to be transferred in one RxDMA transaction is 32.
	RDP_4xPBL_64Beat  DMABMR_Bits = 0x90 << 17   //  maximum number of beats to be transferred in one RxDMA transaction is 64.
	RDP_4xPBL_128Beat DMABMR_Bits = 0xA0 << 17   //  maximum number of beats to be transferred in one RxDMA transaction is 128.
	FB                DMABMR_Bits = 0x01 << 16   //+ Fixed Burst.
	RTPR              DMABMR_Bits = 0x03 << 14   //+ Rx Tx priority ratio.
	RTPR_1_1          DMABMR_Bits = 0x00 << 14   //  Rx Tx priority ratio.
	RTPR_2_1          DMABMR_Bits = 0x01 << 14   //  Rx Tx priority ratio.
	RTPR_3_1          DMABMR_Bits = 0x02 << 14   //  Rx Tx priority ratio.
	RTPR_4_1          DMABMR_Bits = 0x03 << 14   //  Rx Tx priority ratio.
	PBL               DMABMR_Bits = 0x3F << 8    //+ Programmable burst length.
	PBL_1Beat         DMABMR_Bits = 0x01 << 8    //  maximum number of beats to be transferred in one TxDMA (or both) transaction is 1.
	PBL_2Beat         DMABMR_Bits = 0x02 << 8    //  maximum number of beats to be transferred in one TxDMA (or both) transaction is 2.
	PBL_4Beat         DMABMR_Bits = 0x04 << 8    //  maximum number of beats to be transferred in one TxDMA (or both) transaction is 4.
	PBL_8Beat         DMABMR_Bits = 0x08 << 8    //  maximum number of beats to be transferred in one TxDMA (or both) transaction is 8.
	PBL_16Beat        DMABMR_Bits = 0x10 << 8    //  maximum number of beats to be transferred in one TxDMA (or both) transaction is 16.
	PBL_32Beat        DMABMR_Bits = 0x20 << 8    //  maximum number of beats to be transferred in one TxDMA (or both) transaction is 32.
	PBL_4xPBL_4Beat   DMABMR_Bits = 0x10001 << 8 //  maximum number of beats to be transferred in one TxDMA (or both) transaction is 4.
	PBL_4xPBL_8Beat   DMABMR_Bits = 0x10002 << 8 //  maximum number of beats to be transferred in one TxDMA (or both) transaction is 8.
	PBL_4xPBL_16Beat  DMABMR_Bits = 0x10004 << 8 //  maximum number of beats to be transferred in one TxDMA (or both) transaction is 16.
	PBL_4xPBL_32Beat  DMABMR_Bits = 0x10008 << 8 //  maximum number of beats to be transferred in one TxDMA (or both) transaction is 32.
	PBL_4xPBL_64Beat  DMABMR_Bits = 0x10010 << 8 //  maximum number of beats to be transferred in one TxDMA (or both) transaction is 64.
	PBL_4xPBL_128Beat DMABMR_Bits = 0x10020 << 8 //  maximum number of beats to be transferred in one TxDMA (or both) transaction is 128.
	EDE               DMABMR_Bits = 0x01 << 7    //+ Enhanced Descriptor Enable.
	DSL               DMABMR_Bits = 0x1F << 2    //+ Descriptor Skip Length.
	DA                DMABMR_Bits = 0x01 << 1    //+ DMA arbitration scheme.
	SR                DMABMR_Bits = 0x01 << 0    //+ Software reset.
)

const (
	TPD DMATPDR_Bits = 0xFFFFFFFF << 0 //+ Transmit poll demand.
)

const (
	RPD DMARPDR_Bits = 0xFFFFFFFF << 0 //+ Receive poll demand.
)

const (
	SRL DMARDLAR_Bits = 0xFFFFFFFF << 0 //+ Start of receive list.
)

const (
	STL DMATDLAR_Bits = 0xFFFFFFFF << 0 //+ Start of transmit list.
)

const (
	TSTS             DMASR_Bits = 0x01 << 29 //+ Time-stamp trigger status.
	PMTS             DMASR_Bits = 0x01 << 28 //+ PMT status.
	MMCS             DMASR_Bits = 0x01 << 27 //+ MMC status.
	EBS              DMASR_Bits = 0x07 << 23 //+ Error bits status.
	EBS_DescAccess   DMASR_Bits = 0x04 << 23 //  Error bits 0-data buffer, 1-desc. access.
	EBS_ReadTransf   DMASR_Bits = 0x02 << 23 //  Error bits 0-write trnsf, 1-read transfr.
	EBS_DataTransfTx DMASR_Bits = 0x01 << 23 //  Error bits 0-Rx DMA, 1-Tx DMA.
	TPS              DMASR_Bits = 0x07 << 20 //+ Transmit process state.
	TPS_Stopped      DMASR_Bits = 0x00 << 20 //  Stopped - Reset or Stop Tx Command issued.
	TPS_Fetching     DMASR_Bits = 0x01 << 20 //  Running - fetching the Tx descriptor.
	TPS_Waiting      DMASR_Bits = 0x02 << 20 //  Running - waiting for status.
	TPS_Reading      DMASR_Bits = 0x03 << 20 //  Running - reading the data from host memory.
	TPS_Suspended    DMASR_Bits = 0x06 << 20 //  Suspended - Tx Descriptor unavailabe.
	TPS_Closing      DMASR_Bits = 0x07 << 20 //  Running - closing Rx descriptor.
	RPS              DMASR_Bits = 0x07 << 17 //+ Receive process state.
	RPS_Stopped      DMASR_Bits = 0x00 << 17 //  Stopped - Reset or Stop Rx Command issued.
	RPS_Fetching     DMASR_Bits = 0x01 << 17 //  Running - fetching the Rx descriptor.
	RPS_Waiting      DMASR_Bits = 0x03 << 17 //  Running - waiting for packet.
	RPS_Suspended    DMASR_Bits = 0x04 << 17 //  Suspended - Rx Descriptor unavailable.
	RPS_Closing      DMASR_Bits = 0x05 << 17 //  Running - closing descriptor.
	RPS_Queuing      DMASR_Bits = 0x07 << 17 //  Running - queuing the recieve frame into host memory.
	NIS              DMASR_Bits = 0x01 << 16 //+ Normal interrupt summary.
	AIS              DMASR_Bits = 0x01 << 15 //+ Abnormal interrupt summary.
	ERS              DMASR_Bits = 0x01 << 14 //+ Early receive status.
	FBES             DMASR_Bits = 0x01 << 13 //+ Fatal bus error status.
	ETS              DMASR_Bits = 0x01 << 10 //+ Early transmit status.
	RWTS             DMASR_Bits = 0x01 << 9  //+ Receive watchdog timeout status.
	RPSS             DMASR_Bits = 0x01 << 8  //+ Receive process stopped status.
	RBUS             DMASR_Bits = 0x01 << 7  //+ Receive buffer unavailable status.
	RS               DMASR_Bits = 0x01 << 6  //+ Receive status.
	TUS              DMASR_Bits = 0x01 << 5  //+ Transmit underflow status.
	ROS              DMASR_Bits = 0x01 << 4  //+ Receive overflow status.
	TJTS             DMASR_Bits = 0x01 << 3  //+ Transmit jabber timeout status.
	TBUS             DMASR_Bits = 0x01 << 2  //+ Transmit buffer unavailable status.
	TPSS             DMASR_Bits = 0x01 << 1  //+ Transmit process stopped status.
	TS               DMASR_Bits = 0x01 << 0  //+ Transmit status.
)

const (
	DTCEFD       DMAOMR_Bits = 0x01 << 26 //+ Disable Dropping of TCP/IP checksum error frames.
	RSF          DMAOMR_Bits = 0x01 << 25 //+ Receive store and forward.
	DFRF         DMAOMR_Bits = 0x01 << 24 //+ Disable flushing of received frames.
	TSF          DMAOMR_Bits = 0x01 << 21 //+ Transmit store and forward.
	FTF          DMAOMR_Bits = 0x01 << 20 //+ Flush transmit FIFO.
	TTC          DMAOMR_Bits = 0x07 << 14 //+ Transmit threshold control.
	TTC_64Bytes  DMAOMR_Bits = 0x00 << 14 //  threshold level of the MTL Transmit FIFO is 64 Bytes.
	TTC_128Bytes DMAOMR_Bits = 0x01 << 14 //  threshold level of the MTL Transmit FIFO is 128 Bytes.
	TTC_192Bytes DMAOMR_Bits = 0x02 << 14 //  threshold level of the MTL Transmit FIFO is 192 Bytes.
	TTC_256Bytes DMAOMR_Bits = 0x03 << 14 //  threshold level of the MTL Transmit FIFO is 256 Bytes.
	TTC_40Bytes  DMAOMR_Bits = 0x04 << 14 //  threshold level of the MTL Transmit FIFO is 40 Bytes.
	TTC_32Bytes  DMAOMR_Bits = 0x05 << 14 //  threshold level of the MTL Transmit FIFO is 32 Bytes.
	TTC_24Bytes  DMAOMR_Bits = 0x06 << 14 //  threshold level of the MTL Transmit FIFO is 24 Bytes.
	TTC_16Bytes  DMAOMR_Bits = 0x07 << 14 //  threshold level of the MTL Transmit FIFO is 16 Bytes.
	ST           DMAOMR_Bits = 0x01 << 13 //+ Start/stop transmission command.
	FEF          DMAOMR_Bits = 0x01 << 7  //+ Forward error frames.
	FUGF         DMAOMR_Bits = 0x01 << 6  //+ Forward undersized good frames.
	RTC          DMAOMR_Bits = 0x03 << 3  //+ receive threshold control.
	RTC_64Bytes  DMAOMR_Bits = 0x00 << 3  //  threshold level of the MTL Receive FIFO is 64 Bytes.
	RTC_32Bytes  DMAOMR_Bits = 0x01 << 3  //  threshold level of the MTL Receive FIFO is 32 Bytes.
	RTC_96Bytes  DMAOMR_Bits = 0x02 << 3  //  threshold level of the MTL Receive FIFO is 96 Bytes.
	RTC_128Bytes DMAOMR_Bits = 0x03 << 3  //  threshold level of the MTL Receive FIFO is 128 Bytes.
	OSF          DMAOMR_Bits = 0x01 << 2  //+ operate on second frame.
	SR           DMAOMR_Bits = 0x01 << 1  //+ Start/stop receive.
)

const (
	NISE  DMAIER_Bits = 0x01 << 16 //+ Normal interrupt summary enable.
	AISE  DMAIER_Bits = 0x01 << 15 //+ Abnormal interrupt summary enable.
	ERIE  DMAIER_Bits = 0x01 << 14 //+ Early receive interrupt enable.
	FBEIE DMAIER_Bits = 0x01 << 13 //+ Fatal bus error interrupt enable.
	ETIE  DMAIER_Bits = 0x01 << 10 //+ Early transmit interrupt enable.
	RWTIE DMAIER_Bits = 0x01 << 9  //+ Receive watchdog timeout interrupt enable.
	RPSIE DMAIER_Bits = 0x01 << 8  //+ Receive process stopped interrupt enable.
	RBUIE DMAIER_Bits = 0x01 << 7  //+ Receive buffer unavailable interrupt enable.
	RIE   DMAIER_Bits = 0x01 << 6  //+ Receive interrupt enable.
	TUIE  DMAIER_Bits = 0x01 << 5  //+ Transmit Underflow interrupt enable.
	ROIE  DMAIER_Bits = 0x01 << 4  //+ Receive Overflow interrupt enable.
	TJTIE DMAIER_Bits = 0x01 << 3  //+ Transmit jabber timeout interrupt enable.
	TBUIE DMAIER_Bits = 0x01 << 2  //+ Transmit buffer unavailable interrupt enable.
	TPSIE DMAIER_Bits = 0x01 << 1  //+ Transmit process stopped interrupt enable.
	TIE   DMAIER_Bits = 0x01 << 0  //+ Transmit interrupt enable.
)

const (
	OFOC DMAMFBOCR_Bits = 0x01 << 28  //+ Overflow bit for FIFO overflow counter.
	MFA  DMAMFBOCR_Bits = 0x7FF << 17 //+ Number of frames missed by the application.
	OMFC DMAMFBOCR_Bits = 0x01 << 16  //+ Overflow bit for missed frame counter.
	MFC  DMAMFBOCR_Bits = 0xFFFF << 0 //+ Number of frames missed by the controller.
)

const (
	HTDAP DMACHTDR_Bits = 0xFFFFFFFF << 0 //+ Host transmit descriptor address pointer.
)

const (
	HRDAP DMACHRDR_Bits = 0xFFFFFFFF << 0 //+ Host receive descriptor address pointer.
)

const (
	HTBAP DMACHTBAR_Bits = 0xFFFFFFFF << 0 //+ Host transmit buffer address pointer.
)

const (
	HRBAP DMACHRBAR_Bits = 0xFFFFFFFF << 0 //+ Host receive buffer address pointer.
)
