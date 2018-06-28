package main

import (
	"fmt"

	"sdcard"
)

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
	fmt.Printf("Manufacturing date:    %04d-%02d\n", y, m)
}

//emgo:const
var statusStr = [...]string{
	"?",
	"?",
	"?",
	"AKE_SEQ_ERROR",
	"?",
	"APP_CMD",
	"FX_EVENT",
	"?",
	"READY_FOR_DATA",
	"?",
	"?",
	"?",
	"?",
	"ERASE_RESET",
	"CARD_ECC_DISABLED",
	"WP_ERASE_SKIP",
	"CSD_OVERWRITE",
	"?",
	"?",
	"ERROR",
	"CC_ERROR",
	"CARD_ECC_FAILED",
	"ILLEGAL_COMMAND",
	"COM_CRC_ERROR",
	"LOCK_UNLOCK_FAILED",
	"CARD_IS_LOCKED",
	"WP_VIOLATION",
	"ERASE_PARAM",
	"ERASE_SEQ_ERROR",
	"BLOCK_LEN_ERROR",
	"ADDRESS_ERROR",
	"OUT_OF_RANGE",
}

//emgo:const
var stateStr = [...]string{
	"StateIdle",
	"StateReady",
	"StateIdent",
	"StateStby",
	"StateTran",
	"StateData",
	"StateRcv",
	"StatePrg",
	"StateDis",
	"?",
	"?",
	"?",
	"?",
	"?",
	"?",
	"StateIOOnly",
}

func printStatus(st sdcard.CardStatus) {
	fmt.Printf("\nCard status: ")
	fmt.Printf(stateStr[st&sdcard.CURRENT_STATE>>9])
	for n := uint(0); n < 32; n++ {
		if n == 9 {
			n = 12
			continue
		}
		if st&(1<<n) != 0 {
			fmt.Printf(",")
			fmt.Printf(statusStr[n])
		}
	}
	fmt.Printf("\n")
}

func checkErr(err error, st sdcard.CardStatus) {
	if err == nil {
		return
	}
	printStatus(st)
	fmt.Printf("%v\n", err)
	for {
	}
}
