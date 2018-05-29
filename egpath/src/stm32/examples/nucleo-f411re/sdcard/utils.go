package main

import (
	"fmt"

	"sdcard"
)

func checkErr(what string, err error, status sdcard.CardStatus) {
	if err != nil {
		fmt.Printf("%s: %v\n", what, err)
		for {
		}
	}
	errFlags := sdcard.OUT_OF_RANGE |
		sdcard.ADDRESS_ERROR |
		sdcard.BLOCK_LEN_ERROR |
		sdcard.ERASE_SEQ_ERROR |
		sdcard.ERASE_PARAM |
		sdcard.WP_VIOLATION |
		sdcard.LOCK_UNLOCK_FAILED |
		sdcard.COM_CRC_ERROR |
		sdcard.ILLEGAL_COMMAND |
		sdcard.CARD_ECC_FAILED |
		sdcard.CC_ERROR |
		sdcard.ERROR |
		sdcard.CSD_OVERWRITE |
		sdcard.WP_ERASE_SKIP |
		sdcard.CARD_ECC_DISABLED |
		sdcard.ERASE_RESET |
		sdcard.AKE_SEQ_ERROR
	if status &= errFlags; status != 0 {
		fmt.Printf("%s: 0x%X", what, status)
		for {
		}
	}
}

func printCID(cid sdcard.CID) {
	y, m := cid.MDT()
	pnm := cid.PNM()
	oid := cid.OID()
	prv := cid.PRV()
	fmt.Printf("Manufacturer ID:       %d\n", cid.MID())
	fmt.Printf("OEM/Application ID:    %s\n", oid[:])
	fmt.Printf("Product name:          %s\n", pnm[:])
	fmt.Printf("Product revision:      %d.%d\n", prv>>4&15, prv&15)
	fmt.Printf("Product serial number: %d\n", cid.PSN())
	fmt.Printf("Manufacturing date:    %04d-%02d\n\n", y, m)
}

func printCSD(csd sdcard.CSD) {
	csdv := csd.Version()
	fmt.Printf("CSD version:        %d\n", csdv)
	fmt.Printf("TAAC:               %d ns\n", csd.TAAC())
	fmt.Printf("NSAC:               %d clk\n", csd.NSAC())
	fmt.Printf("TRAN_SPEED:         %d kbit/s\n", csd.TRAN_SPEED())
	fmt.Printf("CCC:                0b%012b\n", csd.CCC())
	fmt.Printf("READ_BL_LEN:        %d B\n", csd.READ_BL_LEN())
	fmt.Printf("READ_BL_PARTIAL:    %t\n", csd.READ_BL_PARTIAL())
	fmt.Printf("WRITE_BLK_MISALIGN: %t\n", csd.WRITE_BLK_MISALIGN())
	fmt.Printf("READ_BLK_MISALIGN:  %t\n", csd.READ_BLK_MISALIGN())
	fmt.Printf("DSR_IMP:            %t\n", csd.DSR_IMP())
	csize := csd.C_SIZE()
	fmt.Printf("C_SIZE:             %d KiB (%d kB)\n", csize>>1, csize<<9/1000)
	fmt.Printf("ERASE_BLK_EN:       %t\n", csd.ERASE_BLK_EN())
	fmt.Printf("SECTOR_SIZE:        %d * WRITE_BL_LEN\n", csd.SECTOR_SIZE())
	fmt.Printf("WP_GRP_SIZE:        %d * SECTOR_SIZE\n", csd.WP_GRP_SIZE())
	fmt.Printf("WP_GRP_ENABLE:      %t\n", csd.WP_GRP_ENABLE())
	fmt.Printf("R2W_FACTOR:         %d\n", csd.R2W_FACTOR())
	fmt.Printf("WRITE_BL_LEN:       %d B\n", csd.WRITE_BL_LEN())
	fmt.Printf("WRITE_BL_PARTIAL:   %t\n", csd.WRITE_BL_PARTIAL())
	fmt.Printf("FILE_FORMAT:        %d\n", csd.FILE_FORMAT())
	fmt.Printf("COPY:               %t\n", csd.COPY())
	fmt.Printf("PERM_WRITE_PROTECT: %t\n", csd.PERM_WRITE_PROTECT())
	fmt.Printf("TMP_WRITE_PROTECT:  %t\n\n", csd.TMP_WRITE_PROTECT())
}

func printSCR(scr sdcard.SCR) {
	fmt.Printf("SCR_STRUCTURE:         %d\n", scr.SCR_STRUCTURE())
	fmt.Printf("SD_SPEC:               %d\n", scr.SD_SPEC())
	fmt.Printf("DATA_STAT_AFTER_ERASE: %d\n", scr.DATA_STAT_AFTER_ERASE())
	fmt.Printf("SD_SECURITY:           %d\n", scr.SD_SECURITY())
	fmt.Printf("SD_BUS_WIDTHS:         0b%04b\n", scr.SD_BUS_WIDTHS())

}
