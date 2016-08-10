package main

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/ziutek/dvb/ts"
)

func checkErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	sd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	checkErr(err)
	err = syscall.SetsockoptInt(
		sd, syscall.SOL_SOCKET, syscall.SO_RCVBUF, 524288,
	)
	checkErr(err)
	sa := syscall.RawSockaddrInet4{
		Family: syscall.AF_INET,
		Port:   syscall.Htons(1234),
	}
	err = syscall.Bind(sd, &sa)
	checkErr(err)

	r := ts.NewPktPktReader(
		os.NewFile(uintptr(sd), ""), make([]byte, 7*ts.PktLen),
	)
	pkt := new(ts.ArrayPkt)

	for {
		checkErr(r.ReadPkt(pkt))
		checkErr(err)
		t := time.Now()
		fmt.Printf("%d.%09d pkt %d B\n", t.Unix(), t.Nanosecond(), n)
	}

}
