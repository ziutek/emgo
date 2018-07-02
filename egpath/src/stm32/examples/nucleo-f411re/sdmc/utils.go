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

func printStatus(st sdcard.CardStatus) {
	fmt.Printf("\nCard status: %s\n", st)
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
